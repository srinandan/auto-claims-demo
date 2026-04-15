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
import subprocess
import sys

def run_command(command, ignore_errors=False):
    """Run a shell command and print its output."""
    print(f"Running: {command}")
    result = subprocess.run(command, shell=True, text=True, capture_output=True)
    if result.returncode != 0 and not ignore_errors:
        print(f"Error executing: {command}")
        print(result.stderr)
        sys.exit(1)
    return result.stdout.strip()

def setup_infrastructure():
    # Configuration
    project_id = os.getenv("GOOGLE_CLOUD_PROJECT")
    if not project_id:
        project_id = run_command("gcloud config get-value project", ignore_errors=True)
    if not project_id:
        print("Error: Could not determine GOOGLE_CLOUD_PROJECT. Set it as an environment variable or via gcloud config.")
        sys.exit(1)

    location = os.getenv("GOOGLE_CLOUD_LOCATION")
    if not location:
        location = run_command("gcloud config get-value compute/region", ignore_errors=True)
    if not location:
        location = "us-central1"

    service_account_name = "auto-claims-sa"
    service_account_email = f"{service_account_name}@{project_id}.iam.gserviceaccount.com"

    # Define resource names
    buckets = [
        f"gs://{project_id}-auto-claims"
    ]

    print(f"Using Project ID: {project_id}")
    print(f"Using Location: {location}")

    # 1. Enable APIs
    print("\n--- Enabling necessary APIs ---")
    apis = [
        "aiplatform.googleapis.com",
        "iam.googleapis.com",
        "serviceusage.googleapis.com",
        "storage.googleapis.com",
        "cloudresourcemanager.googleapis.com",
        "telemetry.googleapis.com",
        "artifactregistry.googleapis.com",
        "run.googleapis.com",
        "cloudbuild.googleapis.com",
        "cloudtrace.googleapis.com",
        "bigquery.googleapis.com",
        "secretmanager.googleapis.com",
        "compute.googleapis.com",
        "apphub.googleapis.com",
        "mapstools.googleapis.com",
    ]
    print(f"Enabling {len(apis)} APIs...")
    api_list = ' '.join(apis)
    run_command(f"gcloud services enable {api_list} --project {project_id}")

    # 2. Create Service Account
    print("\n--- Creating Service Account ---")
    sa_check = run_command(f"gcloud iam service-accounts describe {service_account_email} --project {project_id}", ignore_errors=True)
    if not sa_check:
        print(f"Creating service account: {service_account_name}...")
        run_command(f'gcloud iam service-accounts create {service_account_name} --display-name="auto-claims-sa" --project {project_id}')
    else:
        print(f"Service account {service_account_name} already exists.")

    # 3. Assign Roles
    print("\n--- Assigning Roles to Service Account ---")
    roles = [
        "roles/aiplatform.user",
        "roles/aiplatform.admin",
        "roles/storage.admin",
        "roles/storage.objectViewer",
        "roles/storage.objectCreator",
        "roles/telemetry.writer",
        "roles/telemetry.tracesWriter",
        "roles/cloudtrace.agent",
        "roles/telemetry.metricsWriter",
        "roles/monitoring.metricWriter",
        "roles/logging.logWriter",
        "roles/serviceusage.serviceUsageConsumer",
        "roles/bigquery.admin",
        "roles/secretmanager.secretAccessor",
        "roles/apphub.editor",
        "roles/artifactregistry.admin",
        "roles/artifactregistry.reader",
        "roles/run.admin",
        "roles/iam.serviceAccountTokenCreator",
        "roles/iam.serviceAccountUser",
        "roles/discoveryengine.agentspaceAdmin",
        "roles/agentregistry.viewer"
    ]
    for role in roles:
        print(f"Adding role {role}...")
        run_command(f"gcloud projects add-iam-policy-binding {project_id} --member=serviceAccount:{service_account_email} --role={role} > /dev/null", ignore_errors=True)

    # 4. Create GCS Buckets
    print("\n--- Creating GCS Buckets ---")
    for bucket in buckets:
        bucket_check = run_command(f"gcloud storage buckets describe {bucket}", ignore_errors=True)
        if not bucket_check:
            print(f"Creating GCS bucket: {bucket} in {location}...")
            run_command(f"gcloud storage buckets create {bucket} --location={location} --uniform-bucket-level-access --public-access-prevention --default-storage-class=STANDARD")
        else:
            print(f"GCS bucket {bucket} already exists.")

    # 5. Create Artifact Registry Repository
    print("\n--- Creating Artifact Registry ---")
    repository_name = "auto-claims"
    repo_check = run_command(f"gcloud artifacts repositories describe {repository_name} --location={location} --project={project_id}", ignore_errors=True)
    if not repo_check:
        print(f"Creating Artifact Registry repository: {repository_name} in {location}...")
        run_command(f'gcloud artifacts repositories create {repository_name} --repository-format=docker --location={location} --project={project_id} --description="Docker repository for auto-claims images"')
    else:
        print(f"Artifact Registry repository {repository_name} already exists.")

    print(f"Granting service account permissions to repository {repository_name}...")
    run_command(f'gcloud artifacts repositories add-iam-policy-binding {repository_name} --location={location} --project={project_id} --member="serviceAccount:{service_account_email}" --role="roles/artifactregistry.repoAdmin"')

    # 6. Create Secrets
    print("\n--- Creating Secrets ---")
    secrets = ["hf-token"]
    for secret in secrets:
        secret_check = run_command(f"gcloud secrets describe {secret} --project {project_id}", ignore_errors=True)
        if not secret_check:
            print(f"Creating Secret Manager secret: {secret}...")
            run_command(f"gcloud secrets create {secret} --replication-policy=automatic --project {project_id}")
            print(f"Adding empty placeholder to {secret}...")
            run_command(f"echo 'REPLACEME' | gcloud secrets versions add {secret} --data-file=- --project {project_id}")
        else:
            print(f"Secret {secret} already exists.")

    # 6a. Create, restrict, and store Maps API Key
    print("\n--- Creating and Storing Maps API Key ---")
    # Check if a key with display name "Maps API Key" already exists
    # Note: gcloud services api-keys list does not support filtering by display name directly
    # So, we list all and filter manually
    existing_keys_json = run_command(f"gcloud services api-keys list --project={project_id} --format=json", ignore_errors=True)
    existing_keys = json.loads(existing_keys_json) if existing_keys_json else []
    
    key_exists = False
    key_name = None # To store the resource name of the existing key if found
    for key in existing_keys:
        if key.get("displayName") == "Maps API Key":
            key_exists = True
            key_name = key.get("name") # e.g., projects/PROJECT_NUMBER/locations/global/keys/KEY_ID
            print(f"Maps API Key with display name 'Maps API Key' already exists: {key_name}. Skipping creation.")
            break

    if not key_exists:
        print("Creating Maps API Key with restriction to mapstools.googleapis.com...")
        # Create the key with display name and API target restriction
        # Output: Created new key: projects/PROJECT_NUMBER/locations/global/keys/KEY_ID
        create_output = run_command(f"gcloud services api-keys create --display-name=\"Maps API Key\" --api-target=service=mapstools.googleapis.com --project={project_id}")

        match = re.search(r"projects/.*/locations/.*/keys/.*", create_output)
        if match:
            key_name = match.group(0) # This will be the full resource name of the key
            print(f"Successfully created key: {key_name}")

            # Get the cleartext key string
            print("Retrieving key string...")
            key_string = run_command(f"gcloud services api-keys get-key-string {key_name} --project={project_id}")

            if key_string:
                # Store the key in Secret Manager
                print("Storing key in Secret Manager...")
                run_command(f"echo -n '{key_string}' | gcloud secrets versions add maps-api-key --data-file=- --project={project_id}")
                print("Successfully stored Maps API Key in Secret Manager.")
            else:
                print("Error: Failed to retrieve API key string after creation.")
        else:
            print("Error: Failed to parse key name from creation output.")
    else: # If key already exists, still try to store its string in Secret Manager
        if key_name:
            print(f"Retrieving key string for existing key: {key_name}...")
            key_string = run_command(f"gcloud services api-keys get-key-string {key_name} --project={project_id}")
            if key_string:
                print("Storing existing key in Secret Manager...")
                run_command(f"echo -n '{key_string}' | gcloud secrets versions add maps-api-key --data-file=- --project={project_id}")
                print("Successfully stored existing Maps API Key in Secret Manager.")
            else:
                print("Error: Failed to retrieve API key string for existing key.")

    # 7. Create BigQuery Dataset
    print("\n--- Creating BigQuery Dataset ---")
    bq_dataset = "adk_agent_analytics"
    bq_check = run_command(f"bq show {project_id}:{bq_dataset}", ignore_errors=True)
    if "Not found" in bq_check or not bq_check:
        print(f"Creating BigQuery dataset: {bq_dataset}...")
        run_command(f"bq mk --location={location} -d {project_id}:{bq_dataset}")
    else:
        print(f"BigQuery dataset {bq_dataset} already exists.")

    # 8. Create Cloud Build Worker Pool
    print("\n--- Creating Cloud Build Worker Pool ---")
    pool_name = "auto-claims"
    pool_check = run_command(f"gcloud builds worker-pools describe {pool_name} --region={location} --project={project_id}", ignore_errors=True)
    if not pool_check:
        print(f"Creating Cloud Build worker pool: {pool_name} in {location}...")
        run_command(f"gcloud builds worker-pools create {pool_name} --region={location} --project={project_id} --worker-machine-type=c3-standard-4 --worker-disk-size=500GB")
    else:
        print(f"Cloud Build worker pool {pool_name} already exists.")

    # 9. Create Base Docker Images
    print("\n--- Creating Base Docker Images ---")
    run_command("make base")
    print("\n==================================================================================")
    print("IMPORTANT: Don't forget to update the base image path in ai-service/Dockerfile")
    print("to use the correct project ID before running 'make ai-service'!")
    print("==================================================================================\n")

    print("\nInfrastructure setup script completed successfully.")

if __name__ == "__main__":
    setup_infrastructure()
