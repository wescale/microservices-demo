import { WebTracerProvider } from '@opentelemetry/sdk-trace-web';
import { SimpleSpanProcessor } from '@opentelemetry/sdk-trace-base';
import { registerInstrumentations } from '@opentelemetry/instrumentation';
import { Resource } from '@opentelemetry/resources';
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http';
import { FetchInstrumentation } from '@opentelemetry/instrumentation-fetch';
import { ZoneContextManager } from '@opentelemetry/context-zone';

import { getEndpoints } from './api';

const endpoints = await getEndpoints();

if (endpoints.otlpEndpoint) {
  const provider = new WebTracerProvider({
    resource: Resource.default().merge(new Resource({
      'service.name': 'front-admin',
    })),
    spanProcessors: [new SimpleSpanProcessor(new OTLPTraceExporter({
      url : endpoints.otlpEndpoint
    }))]

  });

  provider.register({
    contextManager: new ZoneContextManager()
  });

  // Registering instrumentations
  registerInstrumentations({
    instrumentations: [
      new FetchInstrumentation({
        enabled: true,
        ignoreUrls: ['/config/endpoints.json']
      }),
    ],
    tracerProvider: provider
  });

  console.log('Tracing service started');
}
