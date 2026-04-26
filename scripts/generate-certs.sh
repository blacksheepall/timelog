#!/bin/bash

# Generate self-signed certificate for HTTPS
# Usage: ./scripts/generate-certs.sh [hostname] [days_valid]

HOSTNAME=${1:-localhost}
DAYS=${2:-365}
CERT_DIR="./certs"

# Create certs directory if it doesn't exist
mkdir -p "$CERT_DIR"

echo "Generating self-signed certificate for $HOSTNAME..."
echo "Certificate will be valid for $DAYS days"

# Generate private key and certificate
openssl req -x509 -newkey rsa:2048 -keyout "$CERT_DIR/key.pem" -out "$CERT_DIR/cert.pem" \
  -days "$DAYS" -nodes \
  -subj "/C=CN/ST=State/L=City/O=Organization/CN=$HOSTNAME"

if [ $? -eq 0 ]; then
  echo "✓ Certificate generated successfully"
  echo "  Certificate: $CERT_DIR/cert.pem"
  echo "  Private key: $CERT_DIR/key.pem"
  echo ""
  echo "To use these certificates, set in config.yml:"
  echo "  server:"
  echo "    https_enabled: true"
  echo "    cert_file: ./certs/cert.pem"
  echo "    key_file: ./certs/key.pem"
else
  echo "✗ Failed to generate certificate"
  exit 1
fi
