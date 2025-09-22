const net = require('net');

if (process.argv.length < 6) {
  console.error('Usage: node test-rpc-client.js <host> <port> <method> <paramsJson>');
  process.exit(2);
}

const host = process.argv[2];
const port = parseInt(process.argv[3], 10);
const method = process.argv[4];
let params = {};
try {
  params = JSON.parse(process.argv[5]);
} catch (e) {
  console.error('Invalid JSON for params:', e.message);
  process.exit(2);
}

const payload = JSON.stringify({ method, params }) + '\n';
const client = net.createConnection({ host, port }, () => {
  console.log(`connected -> ${host}:${port}`);
  client.write(payload);
});

let resp = '';
client.on('data', (chunk) => { resp += chunk.toString(); });
client.on('end', () => {
  console.log('socket ended');
  try { console.log('response:', JSON.parse(resp)); }
  catch (e) { console.log('raw response:', resp); }
  process.exit(0);
});
client.on('error', (err) => {
  console.error('socket error:', err.message);
  process.exit(3);
});

// ensure we close after 5s
setTimeout(() => {
  if (!client.destroyed) {
    client.end();
  }
}, 5000);
