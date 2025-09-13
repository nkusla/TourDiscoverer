const { NodeSDK } = require('@opentelemetry/sdk-node');
const { Resource } = require('@opentelemetry/resources');
const { ATTR_SERVICE_NAME } = require('@opentelemetry/semantic-conventions');
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-http');
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');
const { trace } = require('@opentelemetry/api');
const { JAEGER_SERVICE_URL } = require('./constants');

class TracingManager {
	static tp;

	static initTracer() {
		const url = JAEGER_SERVICE_URL;

		if (!url || url.length === 0) {
			console.error('JAEGER_SERVICE_URL environment variable is required for tracing');
			process.exit(1);
		}

		return this._initJaegerTracer(url);
	}

	static _initJaegerTracer(url) {
		console.log(`Initializing tracing to jaeger at ${url}`);

		const sdk = new NodeSDK({
			resource: new Resource({
				[ATTR_SERVICE_NAME]: 'api-gateway',
			}),
			traceExporter: new OTLPTraceExporter({
				url: url,
			}),
			instrumentations: [getNodeAutoInstrumentations()],
		});

		sdk.start();
		this.tp = trace.getTracerProvider();
		return this.tp;
	}

	static getTracerProvider() {
		return this.tp;
	}
}

module.exports = TracingManager;