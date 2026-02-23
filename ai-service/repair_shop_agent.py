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

from google.adk.agents import LlmAgent
from google.adk.runners import InMemoryRunner
from google.adk.tools import google_search
from google.genai.types import Content, Part
import json
import uuid
import re
import httpx
import os
from google.adk.sessions import InMemorySessionService
from google.adk.runners import Runner
from google.adk.agents.remote_a2a_agent import AGENT_CARD_WELL_KNOWN_PATH
from google.adk.agents.remote_a2a_agent import RemoteA2aAgent
from a2a.client import ClientConfig, ClientFactory
from a2a.types import TransportProtocol
from google.auth.transport.requests import Request
from google.auth import default

factory = None
reasoningEngineId = ""

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

def get_last_element(url_string):
    return url_string.split('/')[-1]

# Repair Shop Agent Definition
MODEL_NAME = "gemini-2.5-flash"

REPAIR_SHOP_AGENT_URL = os.getenv("REPAIR_SHOP_AGENT_URL")
AGENT_ENGINE_ENABLED = os.getenv("AGENT_ENGINE_ENABLED", "false")

if AGENT_ENGINE_ENABLED == "true":
    REPAIR_SHOP_AGENT_CARD_URL = f"{REPAIR_SHOP_AGENT_URL}/a2a/v1/card"
    factory = ClientFactory(
        ClientConfig(
            # Specify supported transport mechanisms
            supported_transports=[TransportProtocol.http_json],
            # Use client preferences for protocol negotiation
            use_client_preference=True,
            # Configure HTTP client with authentication            
            httpx_client=httpx.AsyncClient(
                auth=GoogleAuth(),
                timeout=60.0
            ),
        )
    )
    reasoningEngineId = get_last_element(REPAIR_SHOP_AGENT_URL)
else:
    REPAIR_SHOP_AGENT_CARD_URL = f"{REPAIR_SHOP_AGENT_URL}/a2a/app{AGENT_CARD_WELL_KNOWN_PATH}"

repair_shop_agent = RemoteA2aAgent(
        name="repair_shop_agent",
        description=(
            "Finds local auto repair shops using Google Search."
        ),
        agent_card=REPAIR_SHOP_AGENT_CARD_URL,
        a2a_client_factory=factory,
)

async def run_repair_shop_agent(zip_code: str, state: str, make: str, model: str, damage_type: str) -> list:
    """
    Runs the repair shop agent to find shops.
    Returns a list of shop dictionaries.
    """

    prompt_text = f"""
    Please find auto repair shops for a customer.
    Location: Zip Code {zip_code}, State {state}
    Vehicle: {make} {model}
    Damage: {damage_type}

    Find the best repair shops near this location.
    """

    # session_service = VertexAiSessionService(
    #     project=os.getenv("GOOGLE_CLOUD_PROJECT"),
    #     location=os.getenv("GOOGLE_CLOUD_LOCATION"),
    #     agent_engine_id=os.getenv("REASONING_ENGINE_ID")
    # )
    session_service = InMemorySessionService()    

    # Create Content object
    prompt_content = Content(parts=[Part(text=prompt_text)])

    runner = Runner(
        agent=repair_shop_agent,
        app_name=reasoningEngineId,
        session_service=session_service
    )
    user_id = "system" # internal usage

    final_text = ""

    session = await session_service.create_session(
       app_name=reasoningEngineId,
       user_id=user_id
    )

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
        return []

    try:
        # Improved JSON extraction (similar to claims_agent.py)
        cleaned_text = final_text.strip()

        # Try to find JSON block if mixed with other text
        json_match = re.search(r'\[.*\]', cleaned_text, re.DOTALL)
        if json_match:
            potential_json = json_match.group(0)
            try:
                return json.loads(potential_json)
            except json.JSONDecodeError:
                pass

        # Markdown cleanup fallback
        if "```json" in cleaned_text:
            cleaned_text = cleaned_text.split("```json")[1].split("```")[0]
        elif "```" in cleaned_text:
            cleaned_text = cleaned_text.split("```")[1].split("```")[0]

        return json.loads(cleaned_text.strip())
    except Exception as e:
        print(f"Error parsing repair shop agent response: {e}. Raw: {final_text}")
        return []
