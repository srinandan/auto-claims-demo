# Copyright © 2025 Google LLC.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

source ./lib/utils.sh

source ENV

# ====================================================================

check_shell_variables CLOUDRUN_SERVICE_NAME CLOUDRUN_PROJECT_ID CLOUDRUN_REGION

check_required_commands gcloud 

# A function to echo a command and then execute it
log_and_run() {
  # Print the command. The '$*' joins all arguments with a space.
  echo "▶️ Running: $*"
  # Execute the command. The '"$@"' treats each argument as a separate
  # quoted string, which correctly handles spaces and special characters.
  "$@"
  
}

# Important to use   --no-cpu-throttling  because this runs as a daemon.
# It needs to run always, even when there is no pending request.

log_and_run gcloud run deploy "${CLOUDRUN_SERVICE_NAME}" \
  --source "./api-load-generator" \
  --project "${CLOUDRUN_PROJECT_ID}" \
  --concurrency 5 \
  --no-cpu-throttling \
  --cpu 1 \
  --memory '512Mi' \
  --min-instances 1 \
  --max-instances 1 \
  --allow-unauthenticated \
  --region "${CLOUDRUN_REGION}" \
  --set-env-vars="OTEL_SERVICE_NAME=${CLOUDRUN_SERVICE_NAME}" \
  --timeout 180

printf "\nOK.\n"
