package mongodb

import (
	"testing"
)

func TestNewClient_InvalidURI(t *testing.T) {
	// Attempting to connect with an invalid scheme or malformed URI string
	client, err := NewClient("invalid-scheme://invalid-host:999999")
	if err == nil {
		t.Error("expected error for invalid URI scheme, got nil")
	}
	if client != nil {
		t.Error("expected nil client on connection error, got instance")
	}
}
