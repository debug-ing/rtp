const dgram = require('dgram');

// ساخت کلاینت UDP
const client = dgram.createSocket('udp4');

// تنظیمات بسته RTP
function createRTPPacket(sequenceNumber, timestamp, payload, includeOptionalField = false) {
    const version = 2; 
    const padding = 0;
    const extension = includeOptionalField ? 1 : 0;
    const csrcCount = 0;
    const marker = 0;
    const payloadType = 96; 
    const ssrc = 12345678;

    const header = Buffer.alloc(12);
    header[0] = (version << 6) | (padding << 5) | (extension << 4) | csrcCount;
    header[1] = (marker << 7) | payloadType;
    header.writeUInt16BE(sequenceNumber, 2);
    header.writeUInt32BE(timestamp, 4);
    header.writeUInt32BE(ssrc, 8);

    let optionalField = Buffer.alloc(0);
    if (includeOptionalField) {
        optionalField = Buffer.alloc(4);
        optionalField.writeUInt32BE(98765432); 
    }

    return Buffer.concat([header, optionalField, Buffer.from(payload)]);
}

const host = '127.0.0.1';
const port = 5004;

let sequenceNumber = 0;

function sendPacket() {
    const timestamp = Math.floor(Date.now() / 1000);
    const payload = `Hello RTP! Packet #${sequenceNumber}`;

    const packet = createRTPPacket(sequenceNumber, timestamp, payload);

    client.send(packet, port, host, (err) => {
        if (err) {
            console.error(`Error sending packet: ${err.message}`);
            client.close();
        } else {
            console.log(`Sent packet #${sequenceNumber}`);
        }
    });

    sequenceNumber++;
}

client.on('message', (msg, rinfo) => {
    console.log(`Received response from ${rinfo.address}:${rinfo.port} - ${msg}`);
});

setInterval(sendPacket, 2000);