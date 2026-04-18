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
import asyncio
import logging
import os
from typing import Any

import google.auth
import nest_asyncio
import vertexai
from a2a.types import AgentCapabilities, AgentCard, TransportProtocol
from dotenv import load_dotenv
from google.adk.a2a.converters.request_converter import (
    AgentRunRequest,
    convert_a2a_request_to_agent_run_request,
)
from google.adk.a2a.executor.a2a_agent_executor import A2aAgentExecutor
from google.adk.a2a.executor.config import A2aAgentExecutorConfig
from google.adk.a2a.utils.agent_card_builder import AgentCardBuilder
from google.adk.apps import App
from google.adk.artifacts import GcsArtifactService, InMemoryArtifactService
from google.adk.memory import VertexAiMemoryBankService
from google.adk.runners import Runner
from google.adk.sessions import VertexAiSessionService
from google.cloud import logging as google_cloud_logging
from vertexai.preview.reasoning_engines import A2aAgent

from app.agent import app as adk_app
from app.app_utils.custom_types import Feedback
from app.app_utils.telemetry import setup_telemetry

# Load environment variables from .env file at runtime
load_dotenv()


class AgentEngineApp(A2aAgent):
    @staticmethod
    def create(
        app: App | None = None,
        artifact_service: Any = None,
        session_service: Any = None,
        memory_service: Any = None,
    ) -> Any:
        """Create an AgentEngineApp instance.

        This method detects whether it's being called in an async context (like notebooks
        or Agent Engine) and handles agent card creation appropriately.
        """
        if app is None:
            app = adk_app

        def create_runner() -> Runner:
            """Create a Runner for the AgentEngineApp."""
            return Runner(
                app=app,
                session_service=session_service,
                artifact_service=artifact_service,
                memory_service=memory_service,
            )

        # Build agent card in an async context if needed
        try:
            asyncio.get_running_loop()
            # Running event loop detected - enable nested asyncio.run()
            nest_asyncio.apply()
        except RuntimeError:
            pass

        agent_card = asyncio.run(AgentEngineApp.build_agent_card(app=app))

        def custom_request_converter(request, part_converter) -> AgentRunRequest:
            run_request = convert_a2a_request_to_agent_run_request(request, part_converter)
            user_id = None
            if request.call_context and request.call_context.metadata:
                user_id = request.call_context.metadata.get("x-adk-user-id") or request.call_context.metadata.get("user_id")

            if not user_id:
                # The context_id has format ADK$app_name$user_id$session_id. We can try splitting.
                # However, it's safer to just extract from context_id parts directly if we know it.
                parts = request.context_id.split("$") if request.context_id else []
                if len(parts) >= 3:
                    # e.g., parts = ["ADK", "app", "12345", "sess_xyz"]
                    # If prefix is ADK, user_id is the 3rd element
                    if parts[0] == "ADK":
                        user_id = parts[2]

            if not user_id:
                user_id = f'A2A_USER_{request.context_id}'

            return run_request.model_copy(update={"user_id": user_id})

        config = A2aAgentExecutorConfig(request_converter=custom_request_converter)

        return AgentEngineApp(
            agent_executor_builder=lambda: A2aAgentExecutor(runner=create_runner(), config=config),
            agent_card=agent_card,
        )

    @staticmethod
    async def build_agent_card(app: App) -> AgentCard:
        """Builds the Agent Card dynamically from the app."""
        agent_card_builder = AgentCardBuilder(
            agent=app.root_agent,
            # Agent Engine does not support streaming yet
            capabilities=AgentCapabilities(streaming=False),
            rpc_url="http://localhost:9999",
            agent_version=os.getenv("AGENT_VERSION", "0.1.0"),
        )
        agent_card = await agent_card_builder.build()
        agent_card.preferred_transport = TransportProtocol.http_json  # Http Only.
        agent_card.supports_authenticated_extended_card = True
        return agent_card

    def set_up(self) -> None:
        """Initialize the agent engine app with logging and telemetry."""
        vertexai.init()
        setup_telemetry()
        super().set_up()
        logging.basicConfig(level=logging.INFO)
        logging_client = google_cloud_logging.Client()
        self.logger = logging_client.logger(__name__)
        if gemini_location:
            os.environ["GOOGLE_CLOUD_LOCATION"] = gemini_location

    def register_feedback(self, feedback: dict[str, Any]) -> None:
        """Collect and log feedback."""
        feedback_obj = Feedback.model_validate(feedback)
        self.logger.log_struct(feedback_obj.model_dump(), severity="INFO")

    def register_operations(self) -> dict[str, list[str]]:
        """Registers the operations of the Agent."""
        operations = super().register_operations()
        operations[""] = operations.get("", []) + ["register_feedback"]
        return operations

    def clone(self) -> "AgentEngineApp":
        """Returns a clone of the Agent Engine application."""
        return self


_, project_id = google.auth.default()
os.environ["GOOGLE_CLOUD_PROJECT"] = project_id
if not os.environ.get("GOOGLE_CLOUD_LOCATION"):
    os.environ["GOOGLE_CLOUD_LOCATION"] = "us-central1"
gemini_location = os.environ.get("GOOGLE_CLOUD_LOCATION")
engine_id = os.environ.get("SHARED_AGENT_ENGINE_ID")
print(f"DEBUG: Found engine_id: {engine_id}")
logs_bucket_name = os.environ.get("LOGS_BUCKET_NAME")
agent_engine = AgentEngineApp.create(
    app=adk_app,
    artifact_service=(
        GcsArtifactService(bucket_name=logs_bucket_name)
        if logs_bucket_name
        else InMemoryArtifactService()
    ),
    session_service=VertexAiSessionService(
        project=project_id, location=gemini_location, agent_engine_id=engine_id
    ),
    # session_service=InMemorySessionService(),
    memory_service=VertexAiMemoryBankService(
        project=project_id, location=gemini_location, agent_engine_id=engine_id
    ),
)
