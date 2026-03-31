/**
 * Copyright 2026 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import opentelemetry from '@opentelemetry/sdk-node';
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node';
import pkgResource from '@opentelemetry/resources';
const { resourceFromAttributes } = pkgResource;
import pkgSemConv from '@opentelemetry/semantic-conventions';
const { ATTR_SERVICE_NAME } = pkgSemConv;

// Use Google Cloud Exporter
import { TraceExporter } from '@google-cloud/opentelemetry-cloud-trace-exporter';
const traceExporter = new TraceExporter();

const sdk = new opentelemetry.NodeSDK({
  resource: resourceFromAttributes({
    [ATTR_SERVICE_NAME]: process.env.OTEL_SERVICE_NAME || 'auto-claims-loadgen',
  }),
  traceExporter,
  instrumentations: [
    getNodeAutoInstrumentations({
      // Disable noisy FS instrumentations
      '@opentelemetry/instrumentation-fs': { enabled: false },
    }),
  ],
});

sdk.start();

process.on('SIGTERM', () => {
  sdk.shutdown()
    .then(() => console.log('Tracing terminated'))
    .catch((error) => console.log('Error terminating tracing', error))
    .finally(() => process.exit(0));
});
