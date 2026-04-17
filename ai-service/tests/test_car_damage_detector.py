# Copyright 2026 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import unittest
from unittest.mock import MagicMock, patch
import sys
import os

# Add ai-service to path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

# --- MOCKING START ---
# We must mock sys.modules BEFORE importing the detector to handle heavy dependencies
mock_cv2 = MagicMock()
sys.modules['cv2'] = mock_cv2
mock_np = MagicMock()
sys.modules['numpy'] = mock_np
mock_pil = MagicMock()
sys.modules['PIL'] = mock_pil
mock_torch = MagicMock()
sys.modules['torch'] = mock_torch
mock_yolo = MagicMock()
sys.modules['ultralytics'] = mock_yolo
mock_hf = MagicMock()
sys.modules['huggingface_hub'] = mock_hf
# --- MOCKING END ---

from car_damage_detector import CarDamageDetector

class TestCarDamageDetector(unittest.TestCase):
    @patch('os.path.exists', return_value=False)
    def test_load_model_file_not_found(self, mock_exists):
        """
        Test that _load_model correctly handles a missing model file and
        raises the appropriate RuntimeError.
        """
        with self.assertRaises(RuntimeError) as context:
            detector = CarDamageDetector(model_path="dummy.pt")

        self.assertIn("Failed to load model: Model file not found", str(context.exception))

    @patch('os.path.exists', return_value=True)
    def test_load_model_yolo_exception(self, mock_exists):
        """
        Test that _load_model correctly handles an exception during the YOLO
        constructor call and raises a RuntimeError.
        """
        with patch('car_damage_detector.YOLO', side_effect=Exception("YOLO init error")):
            with self.assertRaises(RuntimeError) as context:
                detector = CarDamageDetector(model_path="dummy.pt")

            self.assertIn("Failed to load model: YOLO init error", str(context.exception))

if __name__ == '__main__':
    unittest.main()
