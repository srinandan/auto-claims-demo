# ruff: noqa
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

import datetime
from zoneinfo import ZoneInfo

from google.adk.agents import Agent
from google.adk.apps import App
from google.adk.models import Gemini
from google.adk.tools import LongRunningFunctionTool
from google.genai import types
import logging
from google.adk.plugins.bigquery_agent_analytics_plugin import (
    BigQueryAgentAnalyticsPlugin,
    BigQueryLoggerConfig,
)
from google.cloud import bigquery

import os
import google.auth

_, project_id = google.auth.default()
os.environ["GOOGLE_CLOUD_PROJECT"] = project_id
os.environ["GOOGLE_CLOUD_LOCATION"] = "global"
os.environ["GOOGLE_GENAI_USE_VERTEXAI"] = "True"

# Tool for generating mock repair costs
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

# --- Agents Definition ---
MODEL_NAME = "gemini-2.5-flash"

root_agent = Agent(
    name="ProcessorAgent",
    model=Gemini(
        model=MODEL_NAME,
        retry_options=types.HttpRetryOptions(attempts=3),
    ),
    description="Generate repair cost and decision.",
    instruction="""
    You are a claims processor.
    Based on the 'severity' determined by the AssessorAgent:
    1. Call the 'generate_repair_cost' tool with the severity.
    2. Determine the decision:
       - If Simple: 'Approved'.
       - If Complex: 'Review Required'.

    Output a valid JSON object with the following structure:
    {
        "decision": "Approved" or "Review Required",
        "estimate": {result from generate_repair_cost tool},
        "reasoning": "Brief explanation of the decision."
    }
    IMPORTANT: Ensure the 'estimate' field directly contains the object returned by the tool (it has keys "items", "total_labor", "total_parts", "total_cost"). Do NOT wrap it in another key like 'generate_repair_cost_response'.
    
    Do not output markdown code blocks. Just the raw JSON string.
    """,
    tools=[
        generate_repair_cost,
    ],
    output_key="final_result"
)

# Initialize BigQuery Analytics
_plugins = []
_project_id = os.environ.get("GOOGLE_CLOUD_PROJECT")
_dataset_id = os.environ.get("BQ_ANALYTICS_DATASET_ID", "adk_agent_analytics")
_location = os.environ.get("GOOGLE_CLOUD_LOCATION", "us-central1")

if _project_id:
    try:
        bq = bigquery.Client(project=_project_id)
        bq.create_dataset(f"{_project_id}.{_dataset_id}", exists_ok=True)

        _plugins.append(
            BigQueryAgentAnalyticsPlugin(
                project_id=_project_id,
                dataset_id=_dataset_id,
                location=_location,
                config=BigQueryLoggerConfig(
                    gcs_bucket_name=os.environ.get("BQ_ANALYTICS_GCS_BUCKET"),
                    connection_id=os.environ.get("BQ_ANALYTICS_CONNECTION_ID"),
                ),
            )
        )
    except Exception as e:
        logging.warning(f"Failed to initialize BigQuery Analytics: {e}")

app = App(
    root_agent=root_agent,
    name="app",
    plugins=_plugins,
)
