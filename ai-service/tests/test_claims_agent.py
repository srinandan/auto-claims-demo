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
from unittest.mock import MagicMock
import sys
import os

# Mocking external dependencies that might not be available
sys.modules['google'] = MagicMock()
sys.modules['google.adk'] = MagicMock()
sys.modules['google.adk.agents'] = MagicMock()
sys.modules['google.adk.runners'] = MagicMock()
sys.modules['google.genai'] = MagicMock()
sys.modules['google.genai.types'] = MagicMock()

# Add the parent directory to sys.path to allow importing from claims_agent
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from claims_agent import generate_repair_cost

class TestClaimsAgent(unittest.TestCase):
    def test_generate_repair_cost_simple_lowercase(self):
        result = generate_repair_cost("simple")
        self.assertEqual(result["total_cost"], 700.00)
        self.assertEqual(len(result["items"]), 3)
        self.assertEqual(result["total_labor"], 200.00)
        self.assertEqual(result["total_parts"], 500.00)

    def test_generate_repair_cost_simple_capitalized(self):
        result = generate_repair_cost("Simple")
        self.assertEqual(result["total_cost"], 700.00)
        self.assertEqual(len(result["items"]), 3)

    def test_generate_repair_cost_simple_substring(self):
        result = generate_repair_cost("very simple damage")
        self.assertEqual(result["total_cost"], 700.00)
        self.assertEqual(len(result["items"]), 3)

    def test_generate_repair_cost_complex(self):
        result = generate_repair_cost("complex")
        self.assertEqual(result["total_cost"], 4500.00)
        self.assertEqual(len(result["items"]), 4)
        self.assertEqual(result["total_labor"], 1000.00)
        self.assertEqual(result["total_parts"], 3500.00)

    def test_generate_repair_cost_other(self):
        # Any string not containing "simple" should trigger the complex branch
        result = generate_repair_cost("totaled")
        self.assertEqual(result["total_cost"], 4500.00)
        self.assertEqual(len(result["items"]), 4)

if __name__ == "__main__":
    unittest.main()
