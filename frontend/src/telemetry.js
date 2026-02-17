import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { BatchSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import { DocumentLoadInstrumentation } from '@opentelemetry/instrumentation-document-load';
import { XMLHttpRequestInstrumentation } from '@opentelemetry/instrumentation-xml-http-request';
import { FetchInstrumentation } from '@opentelemetry/instrumentation-fetch';
import { resourceFromAttributes } from '@opentelemetry/resources';
import { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions';
import { ZoneContextManager } from '@opentelemetry/context-zone';

export function initTelemetry() {
  try {
    const exporter = new OTLPTraceExporter({
      url: import.meta.env.VITE_OTEL_EXPORTER_OTLP_ENDPOINT || 'https://telemetry.googleapis.com/v1/traces',
    });

    const provider = new WebTracerProvider({
      resource: resourceFromAttributes({
        [SemanticResourceAttributes.SERVICE_NAME]: 'claims-frontend',
      }),
    });

    if (typeof provider.addSpanProcessor === 'function') {
      provider.addSpanProcessor(new BatchSpanProcessor(exporter));
    } else {
      console.warn('WebTracerProvider missing addSpanProcessor, skipping telemetry setup');
      return;
    }

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
  } catch (e) {
    console.error('Failed to initialize telemetry:', e);
  }
}
