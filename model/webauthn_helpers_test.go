package model

import (
	"testing"

	"github.com/go-webauthn/webauthn/protocol"
)

func TestSerializeCredentialTransportEmpty(t *testing.T) {
	if got := serializeCredentialTransport(nil); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}

func TestSerializeCredentialTransportMultiple(t *testing.T) {
	input := []protocol.AuthenticatorTransport{protocol.USB, protocol.NFC}
	if got := serializeCredentialTransport(input); got != "usb,nfc" {
		t.Fatalf("unexpected: %q", got)
	}
}

func TestParseCredentialTransportEmpty(t *testing.T) {
	if got := parseCredentialTransport(""); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestParseCredentialTransportWithEmptyPart(t *testing.T) {
	got := parseCredentialTransport("usb, ,nfc")
	if len(got) != 2 {
		t.Fatalf("expected 2 transports, got %d", len(got))
	}
}

func TestParseAuthenticatorAttachment(t *testing.T) {
	if got := parseAuthenticatorAttachment(string(protocol.Platform)); got != protocol.Platform {
		t.Fatalf("unexpected attachment: %v", got)
	}
	if got := parseAuthenticatorAttachment(string(protocol.CrossPlatform)); got != protocol.CrossPlatform {
		t.Fatalf("unexpected attachment: %v", got)
	}
	if got := parseAuthenticatorAttachment("unknown"); string(got) != "unknown" {
		t.Fatalf("unexpected attachment: %v", got)
	}
}
