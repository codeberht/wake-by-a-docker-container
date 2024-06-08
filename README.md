# wake-by-a-docker-container
Wake-on-LAN from a docker container

## Installation
It must be started at boot time (systemd service).

## Example call with Node.js
For example in a container with a volume on the .sock :
```
const net = require('net');

const socket = net.createConnection("/tmp/wake-by-a-docker-container.sock")
    .on('connect', () => {
        socket.write("00:00:5e:00:53:01");
    })
    .on('error', function(data) {
        console.error(data);
    });
```
