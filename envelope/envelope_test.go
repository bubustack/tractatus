package envelope

import (
	"encoding/json"
	"testing"

	transportpb "github.com/bubustack/tractatus/gen/go/proto/transport/v1"
)

func TestRoundTrip(t *testing.T) {
	env := &Envelope{
		Kind:        "data",
		MessageID:   "abc123",
		TimestampMs: 42,
		Metadata:    map[string]string{"storyrun": "sr-123", "step": "ingress"},
		Payload:     json.RawMessage(`{"text":"hello"}`),
		Inputs:      json.RawMessage(`{"foo":"bar"}`),
		Transports: []TransportDescriptor{
			{Name: "primary", Kind: "grpc", Mode: "hot"},
		},
	}

	frame, err := ToBinaryFrame(env)
	if err != nil {
		t.Fatalf("ToBinaryFrame failed: %v", err)
	}
	if frame.GetMimeType() != MIMEType {
		t.Fatalf("expected mime type %s, got %s", MIMEType, frame.GetMimeType())
	}
	if frame.GetTimestampMs() != 42 {
		t.Fatalf("expected frame timestamp 42, got %d", frame.GetTimestampMs())
	}

	decoded, err := FromBinaryFrame(frame)
	if err != nil {
		t.Fatalf("FromBinaryFrame failed: %v", err)
	}
	if decoded.Version != LatestVersion {
		t.Fatalf("expected version %s, got %s", LatestVersion, decoded.Version)
	}
	if decoded.Metadata["storyrun"] != "sr-123" {
		t.Fatalf("metadata mismatch: %#v", decoded.Metadata)
	}
	if decoded.Kind != "data" || decoded.MessageID != "abc123" || decoded.TimestampMs != 42 {
		t.Fatalf("header fields not round-tripped: %+v", decoded)
	}
	if string(decoded.Payload) != `{"text":"hello"}` {
		t.Fatalf("payload mismatch: %s", decoded.Payload)
	}
}

func TestFromBinaryFrameRejectsMIME(t *testing.T) {
	frame, err := ToBinaryFrame(&Envelope{})
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	}
	frame.MimeType = "application/octet-stream"
	if _, err := FromBinaryFrame(frame); err == nil {
		t.Fatalf("expected mime type error")
	}
}

func TestFromBinaryFrameUsesFrameTimestampFallback(t *testing.T) {
	frame := &transportpb.BinaryFrame{
		Payload:     []byte(`{"version":"v1","kind":"data","payload":{"ok":true}}`),
		MimeType:    MIMEType,
		TimestampMs: 777,
	}

	decoded, err := FromBinaryFrame(frame)
	if err != nil {
		t.Fatalf("FromBinaryFrame failed: %v", err)
	}
	if decoded.TimestampMs != 777 {
		t.Fatalf("expected frame timestamp fallback, got %d", decoded.TimestampMs)
	}
}

func TestMarshalDefaultsKindForStructuredPayload(t *testing.T) {
	data, err := Marshal(&Envelope{
		Payload: json.RawMessage(`{"text":"hello"}`),
	})
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	decoded, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if decoded.Kind != KindData {
		t.Fatalf("expected default kind %q, got %q", KindData, decoded.Kind)
	}
}
