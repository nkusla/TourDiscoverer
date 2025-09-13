const { NodeSDK } = require('@opentelemetry/sdk-node');
const { Resource } = require('@opentelemetry/resources');
const { ATTR_SERVICE_NAME } = require('@opentelemetry/semantic-conventions');
const { OTLPTraceExporter } = require('@opentelemetry/exporter-trace-otlp-http');
const { getNodeAutoInstrumentations } = require('@opentelemetry/auto-instrumentations-node');
const { trace } = require('@opentelemetry/api');
const { JAEGER_SERVICE_URL } = require('./constants');

let tp;

function initTracer() {
	const url = JAEGER_SERVICE_URL;

	if (!url || url.length === 0) {
		console.error('JAEGER_SERVICE_URL environment variable is required for tracing');
		process.exit(1);
	}

	return initJaegerTracer(url);
}

function initJaegerTracer(url) {
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
	tp = trace.getTracerProvider();
	return tp;
}

function getTracerProvider() {
	return tp;
}

module.exports = {
	initTracer,
	getTracerProvider
};