import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import { DocumentLoadInstrumentation } from '@opentelemetry/instrumentation-document-load';
import { XMLHttpRequestInstrumentation } from '@opentelemetry/instrumentation-xml-http-request';
import { FetchInstrumentation } from '@opentelemetry/instrumentation-fetch';
import { ZoneContextManager } from '@opentelemetry/context-zone';
import { ATTR_SERVICE_NAME } from '@opentelemetry/semantic-conventions';
import { resourceFromAttributes } from '@opentelemetry/resources';

export function initTelemetry() {
  const exporter = new OTLPTraceExporter({
    url: import.meta.env.VITE_OTLP_EXPORTER_URL,
    headers: {},
  });

  const resource = resourceFromAttributes({
    [ATTR_SERVICE_NAME]: 'auto-claims-frontend',
  });

  const provider = new WebTracerProvider({
    resource: resource,
    spanProcessors: [new BatchSpanProcessor(exporter)],
  });

  provider.register({
    contextManager: new ZoneContextManager(),
  });

  registerInstrumentations({
    instrumentations: [
      new DocumentLoadInstrumentation(),
      new XMLHttpRequestInstrumentation({
        propagateTraceHeaderCorsUrls: [/.+/g], // Propagate context to all domains
      }),
      new FetchInstrumentation({
        propagateTraceHeaderCorsUrls: [/.+/g],
      }),
    ],
  });

  console.log('OpenTelemetry initialized for frontend');
}
