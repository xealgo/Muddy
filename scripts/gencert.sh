#!/bin/bash

CWD="$(pwd)"
cd ./certs

# Install mkcert if needed (uncomment if not installed)
# sudo apt-get install mkcert
# sudo apt install libnss3-tools

# Install the local CA (only needed once)
# mkcert -install

# Create the certificate (THIS NEEDS TO RUN!)
mkcert localhost 127.0.0.1 ::1

mv localhost+2.pem muddy.crt
mv localhost+2-key.pem muddy-server.key

chmod 644 muddy.crt
chmod 600 muddy-server.key

cd "$CWD"