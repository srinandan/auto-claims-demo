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

# Official Node.js runtime as a parent image
FROM node:22-alpine

# working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json to the working directory
COPY api-load-generator/package*.json ./

# Install packages
RUN npm install

# Bundle app source
COPY api-load-generator/api-load-generator.js ./
COPY api-load-generator/tracing.js ./
COPY api-load-generator/lib/ ./lib/
COPY api-load-generator/config/ ./config/

# Make port 5950 available
EXPOSE 5950

# Default OTEL config (can be overridden by Cloud Run env vars)
ENV OTEL_SERVICE_NAME=auto-claims-loadgen
ENV OTEL_TRACES_EXPORTER=otlp
ENV OTEL_EXPORTER_OTLP_ENDPOINT=https://telemetry.googleapis.com

# Run the app when the container launches
CMD ["npm", "start"]
