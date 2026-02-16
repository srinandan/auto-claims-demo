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

# agent_executor.py
import logging
import json
import os
import google.auth
from a2a.server.agent_execution import AgentExecutor, RequestContext
from a2a.server.events import EventQueue
from a2a.utils import new_agent_text_message
from a2a.server.tasks import TaskUpdater
from a2a.types import TaskState, TextPart, UnsupportedOperationError
from a2a.utils.errors import ServerError
from google.adk.artifacts import GcsArtifactService
from google.adk.memory.in_memory_memory_service import InMemoryMemoryService
from google.adk.sessions import InMemorySessionService
from google.adk.agents import LlmAgent
from google.adk import Runner
from pydantic import BaseModel, Field
from google.genai import types
from google.adk.plugins.bigquery_agent_analytics_plugin import (
    BigQueryAgentAnalyticsPlugin,
    BigQueryLoggerConfig,
)
from google.cloud import bigquery
from google.adk.agents import Agent
from google.adk.models import Gemini

logger = logging.getLogger(__name__)
logging.basicConfig(level=logging.INFO)

# --- Agents Definition ---
MODEL_NAME = "gemini-2.5-flash"

_, project_id = google.auth.default()
PROJECT_ID = os.environ.setdefault("GOOGLE_CLOUD_PROJECT", project_id)
# Initialize BigQuery Analytics
_plugins = []
_dataset_id = os.environ.get("BQ_ANALYTICS_DATASET_ID", "adk_agent_analytics")
_location = os.environ.get("GOOGLE_CLOUD_LOCATION", "us-central1")

if project_id:
    try:
        bq = bigquery.Client(project=project_id)
        bq.create_dataset(f"{project_id}.{_dataset_id}", exists_ok=True)

        _plugins.append(
            BigQueryAgentAnalyticsPlugin(
                project_id=project_id,
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

artifact_service = GcsArtifactService(bucket_name=f"gs://{PROJECT_ID}")

AGENT_INSTRUCTION = """
You are an expert insurance adjuster.
Analyze the provided list of vehicle damage findings.
Determine if the overall damage is 'Simple' or 'Complex'.
- Simple: Minor dents, scratches, single panel damage.
- Complex: Structural damage, multiple panels, broken glass, airbag deployment.

Output *only* the word 'Simple' or 'Complex'.
"""

# Conversation schema
class AssessmentOutput(BaseModel):
    severity: str = Field(description="Severity classification: 'Simple' or 'Complex'")

class AssessorAgentExecutor(AgentExecutor):
    """Agent executor that uses the ADK to assess vehicle damage."""

    def __init__(self):
        self.agent = None
        self.runner = None

    def _init_agent(self):
        self.agent = Agent(
    name="AssessorAgent",
    model=Gemini(
        model=MODEL_NAME,
        retry_options=types.HttpRetryOptions(attempts=3),
    ),
    description="Assess damage severity based on findings.",
    instruction="""
    You are an expert insurance adjuster.
    Analyze the provided list of vehicle damage findings.
    Determine if the overall damage is 'Simple' or 'Complex'.
    - Simple: Minor dents, scratches, single panel damage.
    - Complex: Structural damage, multiple panels, broken glass, airbag deployment.

    Output *only* the word 'Simple' or 'Complex'.
    """,
    output_key="severity"
)
        self.runner = Runner(
            app_name=self.agent.name,
            agent=self.agent,
            artifact_service=artifact_service,
            session_service=InMemorySessionService(),
            memory_service=InMemoryMemoryService(),
        )

    async def cancel(self, context: RequestContext, event_queue: EventQueue):
        raise ServerError(error=UnsupportedOperationError())

    async def execute(
        self,
        context: RequestContext,
        event_queue: EventQueue,
    ) -> None:
        if self.agent is None:
            self._init_agent()
        logger.debug(f"Executing agent {self.agent.name}")

        query = context.get_user_input()

        # Extract user_id from headers
        user_id = "fake"
        if context.call_context and context.call_context.state:
            headers = context.call_context.state.get("headers", {})
            user_id = headers.get("x-user-id", "fake")
        logger.info(f"User ID: {user_id}")

        updater = TaskUpdater(event_queue, context.task_id, context.context_id)

        if not context.current_task:
            await updater.submit()

        await updater.start_work()

        content = types.Content(role="user", parts=[types.Part(text=query)])
        
        # Use InMemorySessionService
        session = await self.runner.session_service.get_session(
            app_name=self.runner.app_name,
            user_id=user_id,
            session_id=context.context_id,
        ) or await self.runner.session_service.create_session(
            app_name=self.runner.app_name,
            user_id=user_id,
            session_id=context.context_id,
        )

        async for event in self.runner.run_async(
            session_id=session.id, user_id=user_id, new_message=content
        ):
            logger.debug(f"Event from ADK {event}")
            if event.is_final_response():
                parts = event.content.parts
                text_parts = [TextPart(text=part.text) for part in parts if part.text]
                await updater.add_artifact(
                    text_parts,
                    name="result",
                )
                await updater.complete()
                break
            await updater.update_status(
                TaskState.working, message=new_agent_text_message("Working...")
            )
        else:
            logger.debug("Agent failed to complete")
            await updater.update_status(
                TaskState.failed,
                message=new_agent_text_message("Failed to generate a response."),
            )