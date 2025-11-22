#!/bin/bash

google-chrome \
  --origin-to-force-quic-on=localhost:17000 \
  --ignore-certificate-errors \
  --ignore-ssl-errors \
  --ignore-certificate-errors-spki-list \
  --disable-web-security \
  --allow-running-insecure-content
