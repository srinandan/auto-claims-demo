import unittest
from unittest.mock import MagicMock, AsyncMock, patch
import sys
import os

# Add ai-service to path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

# --- MOCKING START ---
# We must mock sys.modules BEFORE importing main or anything else that imports google

# Create a mock for google package
mock_google = MagicMock()
sys.modules['google'] = mock_google

# Set up google.auth.default to return (credentials, project_id)
mock_auth = MagicMock()
mock_auth.default.return_value = (MagicMock(), "test-project")
sys.modules['google.auth'] = mock_auth
sys.modules['google.auth.transport'] = MagicMock()
sys.modules['google.auth.transport.requests'] = MagicMock()

# Mock other google submodules
sys.modules['google.adk'] = MagicMock()
sys.modules['google.adk.agents'] = MagicMock()
sys.modules['google.adk.agents.remote_a2a_agent'] = MagicMock()
sys.modules['google.adk.sessions'] = MagicMock()
sys.modules['google.adk.runners'] = MagicMock()
sys.modules['google.adk.plugins'] = MagicMock()
sys.modules['google.adk.cli'] = MagicMock()
sys.modules['google.adk.cli.adk_web_server'] = MagicMock()
sys.modules['google.adk.telemetry'] = MagicMock()
sys.modules['google.adk.telemetry.google_cloud'] = MagicMock()
sys.modules['google.adk.telemetry.setup'] = MagicMock()
sys.modules['google.adk.tools'] = MagicMock()
sys.modules['google.cloud'] = MagicMock()
sys.modules['google.cloud.storage'] = MagicMock()
sys.modules['google.genai'] = MagicMock()
sys.modules['google.genai.types'] = MagicMock()

sys.modules['ultralytics'] = MagicMock()
sys.modules['cv2'] = MagicMock()
sys.modules['PIL'] = MagicMock()
sys.modules['numpy'] = MagicMock()
sys.modules['fastapi.concurrency'] = MagicMock()

# Mock a2a for claims_agent import
# Need to mock the package properly so submodules can be imported
mock_a2a = MagicMock()
sys.modules['a2a'] = mock_a2a
sys.modules['a2a.client'] = MagicMock()
sys.modules['a2a.types'] = MagicMock()

# Mock run_in_threadpool
sys.modules['fastapi.concurrency'].run_in_threadpool = AsyncMock(side_effect=lambda func, *args, **kwargs: func(*args, **kwargs))

# --- MOCKING END ---

# Import main and patch its dependencies
with patch.dict(os.environ, {"MOCK_MODE": "false"}):
    with patch('telemetry.setup_telemetry'):
         from main import process_claims, ClaimsRequest

class TestProcessClaimsLogic(unittest.IsolatedAsyncioTestCase):

    async def test_process_claims_yolo_integration(self):
        # Setup mocks
        mock_detector = MagicMock()
        mock_detector.detect_damage.return_value = {
            "image_shape": (1000, 1000), # Height, Width
            "damages": [
                {
                    "type": "dent",
                    "bbox": [100, 100, 200, 200], # x1, y1, x2, y2 (pixels)
                    "confidence": 0.95,
                    "severity": "moderate"
                },
                {
                    "type": "scratch",
                    "bbox": [500, 500, 600, 600],
                    "confidence": 0.85,
                    "severity": "light"
                }
            ],
            "highest_severity": "moderate"
        }

        # Patch the global detector in main
        with patch('main.detector', mock_detector):
            # Patch read_image_from_gcs
            with patch('main.read_image_from_gcs', return_value=b'dummy_image_bytes'):
                # Patch run_claims_agent
                with patch('main.run_claims_agent', new_callable=AsyncMock) as mock_agent:
                    mock_agent.return_value = {"decision": "Approved"}

                    request = ClaimsRequest(file_uris=["gs://bucket/image.jpg"])
                    response = await process_claims(request)

                    # Verify detector call
                    mock_detector.detect_damage.assert_called_once_with(b'dummy_image_bytes')

                    # Verify findings
                    self.assertIn("dent", response.findings)
                    self.assertIn("scratch", response.findings)

                    # Verify photo analyses structure
                    detections = response.photo_analyses["gs://bucket/image.jpg"]
                    self.assertEqual(len(detections), 2)

                    # Check first detection (dent)
                    d1 = detections[0]
                    self.assertEqual(d1.label, "dent")
                    self.assertEqual(d1.score, 0.95)
                    # Expected normalized box: [100/1000, 100/1000, 200/1000, 200/1000] = [0.1, 0.1, 0.2, 0.2]
                    self.assertEqual(d1.box, [0.1, 0.1, 0.2, 0.2])

                    # Check second detection (scratch)
                    d2 = detections[1]
                    self.assertEqual(d2.label, "scratch")
                    self.assertEqual(d2.score, 0.85)
                    self.assertEqual(d2.box, [0.5, 0.5, 0.6, 0.6])

if __name__ == '__main__':
    unittest.main()
