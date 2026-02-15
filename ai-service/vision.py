import os
import json
from typing import List, Optional
from pydantic import BaseModel
from google import genai
from google.genai.types import (
    GenerateContentConfig,
    HarmBlockThreshold,
    HarmCategory,
    HttpOptions,
    Part,
    SafetySetting,
)

# Configuration
MOCK_MODE = os.getenv("MOCK_MODE", "true").lower() == "true"
PROJECT_ID = os.getenv("PROJECT_ID", "your-project-id")
LOCATION = os.getenv("LOCATION", "us-central1")

class BoundingBox(BaseModel):
    box_2d: List[int] # [y_min, x_min, y_max, x_max]
    label: str

def detect_bounding_boxes(image_bytes: bytes) -> List[BoundingBox]:
    if MOCK_MODE:
        return [
            BoundingBox(box_2d=[200, 300, 400, 500], label="dent on bumper"),
            BoundingBox(box_2d=[600, 600, 800, 800], label="scratch on door")
        ]

    try:
        client = genai.Client(http_options=HttpOptions(api_version="v1"))

        config = GenerateContentConfig(
            system_instruction="""
            Return bounding boxes as an array with labels.
            Never return masks. Limit to 25 objects.
            Identify damage to the vehicle such as dents, scratches, broken glass, etc.
            If an object is present multiple times, give each object a unique label.
            """,
            temperature=0.5,
            safety_settings=[
                SafetySetting(
                    category=HarmCategory.HARM_CATEGORY_DANGEROUS_CONTENT,
                    threshold=HarmBlockThreshold.BLOCK_ONLY_HIGH,
                ),
            ],
            response_mime_type="application/json",
            response_schema=list[BoundingBox],
        )

        response = client.models.generate_content(
            model="gemini-2.0-flash", # Using 2.0 Flash as a proxy for "Gemini 3" or latest fast model
            contents=[
                Part.from_bytes(
                    data=image_bytes,
                    mime_type="image/jpeg", # Assuming JPEG for now, ideally detect mime type
                ),
                "Identify and localize all vehicle damage in this image.",
            ],
            config=config,
        )

        # Parse response
        # response.parsed is expected to be a list of BoundingBox objects due to response_schema
        if response.parsed:
            return response.parsed

        # Fallback if parsed is None but text exists
        if response.text:
            data = json.loads(response.text)
            return [BoundingBox(**item) for item in data]

        return []

    except Exception as e:
        print(f"Error in detect_bounding_boxes: {e}")
        # Fallback to mock/empty in case of error to keep service running?
        # Or re-raise. For now, let's return empty list or re-raise.
        # Re-raising allows the caller to handle it.
        raise e
