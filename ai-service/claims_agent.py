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
from google.adk.sessions import VertexAiSessionService, InMemorySessionService
from google.adk.agents import LlmAgent, SequentialAgent
from google.adk.runners import InMemoryRunner, Runner
from google.genai.types import Content, Part
import json
import uuid
import re
import httpx
from a2a.client import ClientConfig, ClientFactory
from a2a.types import TransportProtocol
from google.auth.transport.requests import Request
from google.auth import default

factory = None

class GoogleAuth(httpx.Auth):
    def __init__(self):
        try:
            self.credentials, _ = default(scopes=['https://www.googleapis.com/auth/cloud-platform'])
            self.request = Request()
        except Exception as e:
            print(f"Error getting credentials: {e}")
            print("Please ensure you have authenticated with 'gcloud auth application-default login'.")
            self.credentials = None
            self.request = None

    def auth_flow(self, request):
        if self.credentials:
            if not self.credentials.valid:
                try:
                    self.credentials.refresh(self.request)
                except Exception as e:
                    print(f"Error refreshing credentials: {e}")
            if self.credentials.token:
                request.headers["Authorization"] = f"Bearer {self.credentials.token}"
        yield request

ASSESSOR_AGENT_URL = os.getenv("ASSESSOR_AGENT_URL")
PROCESSOR_AGENT_URL = os.getenv("PROCESSOR_AGENT_URL")
AGENT_ENGINE_ENABLED = os.getenv("AGENT_ENGINE_ENABLED", "false")

if AGENT_ENGINE_ENABLED == "true":
    ASSESSOR_AGENT_CARD_URL = f"{ASSESSOR_AGENT_URL}/a2a/v1/card"
    PROCESSOR_AGENT_CARD_URL = f"{PROCESSOR_AGENT_URL}/a2a/v1/card"
    factory = ClientFactory(
        ClientConfig(
            # Specify supported transport mechanisms
            supported_transports=[TransportProtocol.http_json],
            # Use client preferences for protocol negotiation
            use_client_preference=True,
            # Configure HTTP client with authentication            
            httpx_client=httpx.AsyncClient(
                auth=GoogleAuth(),
            ),
        )
    )    
else:
    ASSESSOR_AGENT_CARD_URL = f"{ASSESSOR_AGENT_URL}/a2a/app{AGENT_CARD_WELL_KNOWN_PATH}"
    PROCESSOR_AGENT_CARD_URL = f"{PROCESSOR_AGENT_URL}/a2a/app{AGENT_CARD_WELL_KNOWN_PATH}"

# Assessor Agent
assessor_agent = RemoteA2aAgent(
        name="assessor_agent",
        description=(
            "Assess damage severity based on findings."
        ),
        agent_card=ASSESSOR_AGENT_CARD_URL,
        a2a_client_factory=factory,
)

# Processor Agent
processor_agent = RemoteA2aAgent(
        name="processor_agent",
        description=(
            "Generate repair cost and decision."
        ),
        agent_card=PROCESSOR_AGENT_CARD_URL,
        a2a_client_factory=factory,
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

    # session_service = VertexAiSessionService(
    #     project=os.getenv("GOOGLE_CLOUD_PROJECT"),
    #     location=os.getenv("GOOGLE_CLOUD_LOCATION"),
    #     agent_engine_id=os.getenv("REASONING_ENGINE_ID")
    # )
    session_service = InMemorySessionService()

    # Prepare input text
    findings_text = "\n".join([f"- {f}" for f in findings])
    prompt_text = f"Here are the damage findings from the vehicle images:\n{findings_text}\n\nPlease assess and process this claim."

    # Create Content object
    prompt_content = Content(parts=[Part(text=prompt_text)])

    runner = Runner(
        agent=claims_sequential_agent,
        app_name=os.getenv("REASONING_ENGINE_ID"),
        session_service=session_service
    )
    user_id = "system" # internal usage

    session = await session_service.create_session(
       app_name=os.getenv("REASONING_ENGINE_ID"),
       user_id=user_id
    )

    final_text = ""

    # Run asynchronously
    async for event in runner.run_async(user_id=user_id, session_id=session.id, new_message=prompt_content):
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
