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
from unittest.mock import MagicMock, patch, mock_open
import sys
import os
import io

# --- MOCKING START ---
class DummyType:
    pass

class DummyNPArray(DummyType):
    pass

class DummyPILImage(DummyType):
    pass

# We must mock sys.modules BEFORE importing the detector to handle heavy dependencies
mock_cv2 = MagicMock()
sys.modules['cv2'] = mock_cv2

mock_np = MagicMock()
mock_np.ndarray = DummyNPArray
mock_np.uint8 = DummyType
sys.modules['numpy'] = mock_np

mock_pil = MagicMock()
mock_pil.Image = MagicMock()
mock_pil.Image.Image = DummyPILImage
sys.modules['PIL'] = mock_pil

mock_torch = MagicMock()
sys.modules['torch'] = mock_torch
mock_yolo = MagicMock()
sys.modules['ultralytics'] = mock_yolo
mock_hf = MagicMock()
sys.modules['huggingface_hub'] = mock_hf
mock_requests = MagicMock()
sys.modules['requests'] = mock_requests
# --- MOCKING END ---

# Add ai-service to path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))
from car_damage_detector import CarDamageDetector

class TestCarDamageDetector(unittest.TestCase):
    @patch('os.path.exists', return_value=True)
    @patch('car_damage_detector.hf_hub_download', return_value="dummy.pt")
    @patch('car_damage_detector.YOLO')
    def setUp(self, mock_yolo_cls, mock_hf_hub, mock_exists):
        self.detector = CarDamageDetector(model_path="dummy.pt")

    @patch('os.path.exists', return_value=False)
    def test_load_model_file_not_found(self, mock_exists):
        """
        Test that _load_model correctly handles a missing model file and
        raises the appropriate RuntimeError.
        """
        with self.assertRaises(RuntimeError) as context:
            CarDamageDetector(model_path="dummy.pt")

        self.assertIn("Failed to load model: Model file not found", str(context.exception))

    @patch('os.path.exists', return_value=True)
    def test_load_model_yolo_exception(self, mock_exists):
        """
        Test that _load_model correctly handles an exception during the YOLO
        constructor call and raises a RuntimeError.
        """
        with patch('car_damage_detector.YOLO', side_effect=Exception("YOLO init error")):
            with self.assertRaises(RuntimeError) as context:
                CarDamageDetector(model_path="dummy.pt")

            self.assertIn("Failed to load model: YOLO init error", str(context.exception))

    # --- _preprocess_image Tests ---

    def test_preprocess_image_invalid_type(self):
        """Test that _preprocess_image raises TypeError for unsupported types."""
        with self.assertRaises(TypeError):
            self.detector._preprocess_image(123)

    def test_preprocess_image_invalid_file_path(self):
        """Test that _preprocess_image raises ValueError for non-existent file paths."""
        mock_cv2.imread.return_value = None
        with self.assertRaises(ValueError) as context:
            self.detector._preprocess_image("non_existent.jpg")
        self.assertIn("Could not load image from non_existent.jpg", str(context.exception))

    def test_preprocess_image_valid_file_path(self):
        """Test that _preprocess_image correctly processes a valid file path."""
        mock_img = MagicMock()
        mock_cv2.imread.return_value = mock_img
        mock_cv2.cvtColor.return_value = "processed_img"

        result = self.detector._preprocess_image("valid.jpg")

        mock_cv2.imread.assert_called_with("valid.jpg")
        self.assertEqual(result, "processed_img")

    def test_preprocess_image_valid_url(self):
        """Test that _preprocess_image correctly processes a valid URL."""
        mock_requests.get.side_effect = None
        mock_response = MagicMock()
        mock_response.content = b"fake_image_data"
        mock_requests.get.return_value = mock_response

        mock_pil_img = MagicMock()
        mock_pil_img.convert.return_value = mock_pil_img
        mock_pil.Image.open.return_value = mock_pil_img  # pylint: disable=no-member

        with patch('car_damage_detector.np.array', return_value="np_array_img"):
            result = self.detector._preprocess_image("http://example.com/image.jpg")

        mock_requests.get.assert_called_with("http://example.com/image.jpg", timeout=10)
        self.assertEqual(result, "np_array_img")

    def test_preprocess_image_invalid_url(self):
        """Test that _preprocess_image handles URL request failures."""
        mock_requests.get.side_effect = Exception("Connection error")

        with self.assertRaises(ValueError) as context:
            self.detector._preprocess_image("http://example.com/image.jpg")
        self.assertIn("Error preprocessing image: Connection error", str(context.exception))

    def test_preprocess_image_pil_image(self):
        """Test that _preprocess_image correctly processes a PIL Image."""
        mock_pil_img = DummyPILImage()
        mock_pil_img.convert = MagicMock(return_value=mock_pil_img)

        with patch('car_damage_detector.np.array', return_value="np_array_from_pil"):
            result = self.detector._preprocess_image(mock_pil_img)

        self.assertEqual(result, "np_array_from_pil")

    def test_preprocess_image_numpy_array(self):
        """Test that _preprocess_image correctly processes a numpy array."""
        fake_arr = DummyNPArray()
        fake_arr.copy = MagicMock(return_value="copied_arr")

        result = self.detector._preprocess_image(fake_arr)

        self.assertEqual(result, "copied_arr")

    def test_preprocess_image_valid_bytes(self):
        """Test that _preprocess_image correctly processes valid bytes."""
        image_bytes = b"fake_bytes"
        mock_np.frombuffer.return_value = "nparr"
        mock_cv2.imdecode.return_value = "decoded_img"
        mock_cv2.cvtColor.return_value = "rgb_img"

        result = self.detector._preprocess_image(image_bytes)

        self.assertEqual(result, "rgb_img")

    def test_preprocess_image_invalid_bytes(self):
        """Test that _preprocess_image handles invalid bytes."""
        image_bytes = b"invalid_bytes"
        mock_cv2.imdecode.return_value = None

        with self.assertRaises(ValueError) as context:
            self.detector._preprocess_image(image_bytes)
        self.assertIn("Could not decode image bytes", str(context.exception))

if __name__ == '__main__':
    unittest.main()
