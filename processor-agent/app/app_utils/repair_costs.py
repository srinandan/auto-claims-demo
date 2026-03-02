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

def generate_repair_cost(severity: str) -> dict:
    """Generates itemized repair costs based on severity."""
    severity = severity.lower()
    if "simple" in severity:
        return {
            "items": [
                {"part": "Bumper Repair", "cost": 350.00},
                {"part": "Labor (2 hours)", "cost": 200.00},
                {"part": "Paint Touch-up", "cost": 150.00}
            ],
            "total_labor": 200.00,
            "total_parts": 500.00,
            "total_cost": 700.00
        }
    else: # Complex
        return {
            "items": [
                {"part": "Fender Replacement", "cost": 1200.00},
                {"part": "Door Panel Repair", "cost": 800.00},
                {"part": "Labor (10 hours)", "cost": 1000.00},
                {"part": "Painting & Blending", "cost": 1500.00}
            ],
            "total_labor": 1000.00,
            "total_parts": 3500.00,
            "total_cost": 4500.00
        }
