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

from fastapi import FastAPI, UploadFile, File, Form, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import os
import json
import base64
from typing import List, Dict, Any, Optional
from google.cloud import storage
import io

# Local imports
from vision import detect_bounding_boxes
from claims_agent import run_claims_agent
from repair_shop_agent import run_repair_shop_agent
from appointment_agent import run_appointment_agent
from car_damage_detector import CarDamageDetector
from fastapi.concurrency import run_in_threadpool

import google.auth

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Configuration
MOCK_MODE = os.environ.get("MOCK_MODE", "false").lower() == "true"

try:
    _, project_id = google.auth.default()
except Exception:
    project_id = "mock-project-id"

PROJECT_ID = os.environ.setdefault("GOOGLE_CLOUD_PROJECT", project_id)
REGION = os.environ.setdefault("GOOGLE_CLOUD_LOCATION", "us-central1")
os.environ.setdefault("GOOGLE_GENAI_USE_VERTEXAI", "True")
os.environ.setdefault("GOOGLE_CLOUD_QUOTA_PROJECT", project_id)

# Initialize storage client once
storage_client = None
if not MOCK_MODE:
    try:
        storage_client = storage.Client()
    except Exception as e:
        print(f"Failed to initialize storage client: {e}")

# Initialize Car Damage Detector
detector = None
try:
    detector = CarDamageDetector()
except Exception as e:
    print(f"Failed to initialize CarDamageDetector: {e}")

class Detection(BaseModel):
    label: str
    box: List[float] # [x_min, y_min, x_max, y_max] normalized 0-1
    score: float = 1.0 # Default score if not provided by genai

class AnalysisResponse(BaseModel):
    photo_id: str
    quality_score: str # "Good", "Blurry", "Dark"
    detections: List[Detection]
    parts_detected: List[str]
    severity: str # "Low", "Medium", "High"
    total_loss_probability: float # 0.0 to 1.0

class ClaimsRequest(BaseModel):
    file_uris: List[str]

class ClaimsProcessResponse(BaseModel):
    findings: List[str]
    agent_result: Dict[str, Any]
    photo_analyses: Dict[str, List[Detection]] # Key is URI or filename

class RepairShopRequest(BaseModel):
    zip_code: str
    state: str
    make: str
    model: str
    damage_type: str

class RepairShopResponse(BaseModel):
    shops: List[Dict[str, Any]]

class AppointmentRequest(BaseModel):
    session_id: str
    message: str
    context: Optional[Dict[str, Any]] = None

class AppointmentResponse(BaseModel):
    agent_message: str

def read_image_from_gcs(uri: str) -> bytes:
    """Reads image content from GCS URI gs://bucket/path"""
    if not uri.startswith("gs://"):
        raise ValueError(f"Invalid GCS URI: {uri}")

    parts = uri[5:].split("/", 1)
    if len(parts) != 2:
        raise ValueError(f"Invalid GCS URI format: {uri}")

    bucket_name, blob_name = parts

    if storage_client is None:
        raise RuntimeError("Storage client not initialized")

    bucket = storage_client.bucket(bucket_name)
    blob = bucket.blob(blob_name)
    return blob.download_as_bytes()

@app.get("/ping")
async def ping():
    return {"message": "AI Service is running"}

@app.post("/analyze", response_model=AnalysisResponse)
async def analyze_image(
    file: UploadFile = File(...),
    photo_id: str = Form(...)
):
    print(f"Analyzing photo {photo_id} in {'MOCK' if MOCK_MODE else 'REAL'} mode")

    if MOCK_MODE:
        return mock_analysis(photo_id)
    else:
        # Read file content for real analysis
        content = await file.read()
        return await real_analysis(content, photo_id)

@app.post("/process-claims", response_model=ClaimsProcessResponse)
async def process_claims(request: ClaimsRequest):
    print(f"Processing {len(request.file_uris)} URIs for claims agent.")

    aggregated_findings = []
    photo_analyses = {}

    for uri in request.file_uris:
        try:
            if MOCK_MODE:
                content = b"mock content"
            else:
                content = read_image_from_gcs(uri)

            # Use Gemini to detect bounding boxes
            boxes = detect_bounding_boxes(content)

            detections = []
            file_findings = []

            for box in boxes:
                # Convert 0-1000 [y_min, x_min, y_max, x_max] to 0-1 [x_min, y_min, x_max, y_max]
                y_min, x_min, y_max, x_max = box.box_2d
                box_normalized = [
                    x_min / 1000.0,
                    y_min / 1000.0,
                    x_max / 1000.0,
                    y_max / 1000.0
                ]

                detections.append(Detection(
                    label=box.label,
                    box=box_normalized,
                    score=1.0
                ))

                # Add to findings (e.g. "dent on bumper")
                file_findings.append(box.label)
                aggregated_findings.append(box.label)

            photo_analyses[uri] = detections

        except Exception as e:
            print(f"Error processing file {uri}: {e}")
            # Continue with other files or fail?
            # For now, log and continue, maybe empty detection
            photo_analyses[uri] = []

    # Run Sequential Agent
    if not aggregated_findings:
        aggregated_findings = ["No visible damage detected."]

    if MOCK_MODE:
        agent_response = mock_claims_agent_response(aggregated_findings)
    else:
        agent_response = await run_claims_agent(aggregated_findings)

    return ClaimsProcessResponse(
        findings=aggregated_findings,
        agent_result=agent_response,
        photo_analyses=photo_analyses
    )

@app.post("/find-repair-shops", response_model=RepairShopResponse)
async def find_repair_shops(request: RepairShopRequest):
    print(f"Finding repair shops for {request.make} {request.model} near {request.zip_code}")

    if MOCK_MODE:
        return mock_repair_shop_search()

    shops = await run_repair_shop_agent(
        request.zip_code,
        request.state,
        request.make,
        request.model,
        request.damage_type
    )
    return RepairShopResponse(shops=shops)

@app.post("/book-appointment", response_model=AppointmentResponse)
async def book_appointment(request: AppointmentRequest):
    print(f"Booking appointment session {request.session_id}")

    response = await run_appointment_agent(
        request.session_id,
        request.message,
        request.context
    )

    return AppointmentResponse(agent_message=response)

def mock_repair_shop_search():
    return {
        "shops": [
            {
                "name": "Joe's Auto Body",
                "address": "123 Main St, Springfield, IL",
                "rating": 4.5,
                "phone": "555-123-4567",
                "reasoning": "High rating and close to your location."
            },
            {
                "name": "Expert Collision Center",
                "address": "456 Elm St, Springfield, IL",
                "rating": 4.8,
                "phone": "555-987-6543",
                "reasoning": "Specializes in collision repair."
            }
        ]
    }

def mock_claims_agent_response(findings):
    return {
        "decision": "Approved",
        "reasoning": "Damage is minor and consistent with description (MOCK).",
        "estimate": {
            "items": [
                {"part": "Bumper Repair", "cost": 350.00},
                {"part": "Labor (2 hours)", "cost": 200.00},
                {"part": "Paint Touch-up", "cost": 150.00}
            ],
            "total_cost": 700.00
        }
    }

def mock_analysis(photo_id):
    # Return hardcoded sample data
    # We'll return a generic 'damage found' response
    return {
        "photo_id": photo_id,
        "quality_score": "Good",
        "detections": [
            {"label": "dent", "box": [0.2, 0.3, 0.4, 0.5], "score": 0.95},
            {"label": "scratch", "box": [0.6, 0.6, 0.8, 0.65], "score": 0.88}
        ],
        "parts_detected": ["front_bumper", "hood"],
        "severity": "Medium",
        "total_loss_probability": 0.15
    }

async def real_analysis(content, photo_id):
    if detector is None:
        print("Detector not initialized, falling back to mock.")
        return mock_analysis(photo_id)

    try:
        # Run detection in threadpool since it's CPU bound
        result = await run_in_threadpool(detector.detect_damage, content)

        detections = []
        height, width = result['image_shape']

        for d in result['damages']:
            x1, y1, x2, y2 = d['bbox']
            # Normalize to 0-1
            detections.append(Detection(
                label=d['type'],
                box=[x1/width, y1/height, x2/width, y2/height],
                score=d['confidence']
            ))

        # Map severity
        severity_map = {
            "light": "Low",
            "moderate": "Medium",
            "severe": "High",
            "none": "Low"
        }
        severity = severity_map.get(result['highest_severity'], "Medium")

        # Estimate total loss probability based on severity
        loss_prob_map = {
            "Low": 0.1,
            "Medium": 0.4,
            "High": 0.8
        }
        loss_prob = loss_prob_map.get(severity, 0.1)

        return AnalysisResponse(
            photo_id=photo_id,
            quality_score="Good",
            detections=detections,
            parts_detected=[], # Part detection not available in this model
            severity=severity,
            total_loss_probability=loss_prob
        )

    except Exception as e:
        print(f"Error in real_analysis: {e}")
        return mock_analysis(photo_id)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
