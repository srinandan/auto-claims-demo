from google.adk.agents import LlmAgent
from google.adk.runners import InMemoryRunner
from google.adk.tools import google_search
from google.genai.types import Content, Part
import json
import uuid
import re

# Repair Shop Agent Definition
MODEL_NAME = "gemini-2.5-flash"

repair_shop_agent = LlmAgent(
    name="RepairShopAgent",
    model=MODEL_NAME,
    description="Finds local auto repair shops using Google Search.",
    instruction="""
    You are an expert assistant helping users find auto repair shops.
    The user will provide their location (Zip Code, City, State), Vehicle details (Make, Model, Year), and Damage description.

    Your task:
    1. Use Google Search to find 3-5 highly-rated auto body repair shops near the provided location.
    2. Look for shops that specialize in the specific vehicle make or type of damage if possible.
    3. Return a JSON list of objects. Each object must have the following keys:
       - "name": The name of the shop.
       - "address": The full address of the shop.
       - "rating": The average rating (e.g., 4.5) as a number, or null if not found.
       - "phone": The phone number, or null if not found.
       - "reasoning": A brief explanation of why this shop is a good fit (e.g., "Specializes in BMW repair", "High rating for dent removal").

    IMPORTANT:
    - Output ONLY valid JSON.
    - Do not include markdown formatting (like ```json ... ```).
    - If you cannot find any shops, return an empty list [].
    """,
    tools=[google_search],
    output_key="shops"
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

    # Create Content object
    prompt_content = Content(parts=[Part(text=prompt_text)])

    runner = InMemoryRunner(agent=repair_shop_agent)
    runner.auto_create_session = True
    session_id = str(uuid.uuid4())
    user_id = "system"

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
