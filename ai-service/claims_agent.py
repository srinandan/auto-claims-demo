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

import os
from google.adk.agents.remote_a2a_agent import AGENT_CARD_WELL_KNOWN_PATH
from google.adk.agents.remote_a2a_agent import RemoteA2aAgent
from google.adk.agents import LlmAgent, SequentialAgent
from google.adk.runners import InMemoryRunner
from google.genai.types import Content, Part
import json
import uuid
import re

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

ASSESSOR_AGENT_URL = os.getenv("ASSESSOR_AGENT_URL")

# Assessor Agent
assessor_agent = RemoteA2aAgent(
        name="assessor_agent_test",
        description=(
            "Assess damage severity based on findings."
        ),
        agent_card=f"{ASSESSOR_AGENT_URL}{AGENT_CARD_WELL_KNOWN_PATH}",
        # a2a_client_factory=factory,
)

# Processor Agent
processor_agent = LlmAgent(
    name="ProcessorAgent",
    model=MODEL_NAME,
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
    tools=[generate_repair_cost],
    output_key="final_result"
)

# Sequential Agent
claims_sequential_agent = SequentialAgent(
    name="ClaimsSequentialAgent",
    sub_agents=[assessor_agent, processor_agent],
    description="Assess damage and process claim sequentially."
)

# Helper to run the agent
async def run_claims_agent(findings: list[str]) -> dict:
    """
    Runs the sequential agent with the provided findings.
    Returns the final result dictionary.
    """

    # Prepare input text
    findings_text = "\n".join([f"- {f}" for f in findings])
    prompt_text = f"Here are the damage findings from the vehicle images:\n{findings_text}\n\nPlease assess and process this claim."

    # Create Content object
    prompt_content = Content(parts=[Part(text=prompt_text)])

    runner = InMemoryRunner(agent=claims_sequential_agent)
    runner.auto_create_session = True
    session_id = str(uuid.uuid4())
    user_id = "system" # internal usage

    final_text = ""

    # Run asynchronously
    async for event in runner.run_async(user_id=user_id, session_id=session_id, new_message=prompt_content):
        # Accumulate text from events
        if hasattr(event, 'text') and event.text:
            final_text += event.text
        elif isinstance(event, str):
            final_text += event
        elif hasattr(event, 'content') and event.content and event.content.parts:
            for part in event.content.parts:
                if part.text:
                   final_text += part.text

    if not final_text:
         return {
            "decision": "Error",
            "reasoning": "No response text received from agent.",
            "estimate": {}
        }

    try:
        # Improved JSON extraction
        cleaned_text = final_text.strip()

        # Try to find JSON block if mixed with other text
        json_match = re.search(r'\{.*\}', cleaned_text, re.DOTALL)
        if json_match:
            potential_json = json_match.group(0)
            # Try parsing the finding
            try:
                return json.loads(potential_json)
            except json.JSONDecodeError:
                # If that failed, maybe it has markdown blocks
                pass

        # Markdown cleanup fallback
        if "```json" in cleaned_text:
            cleaned_text = cleaned_text.split("```json")[1].split("```")[0]
        elif "```" in cleaned_text:
            cleaned_text = cleaned_text.split("```")[1].split("```")[0]

        return json.loads(cleaned_text.strip())
    except Exception as e:
        print(f"Error parsing agent response: {e}. Raw: {final_text}")
        return {
            "decision": "Error",
            "reasoning": f"Failed to parse agent response: {str(e)}",
            "estimate": {},
            "raw": final_text
        }
