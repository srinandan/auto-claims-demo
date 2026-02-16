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
