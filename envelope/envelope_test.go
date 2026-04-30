package envelope

import (
	"encoding/json"
	"math"
	"strings"
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

func TestRoundTripPreservesTypedTransportConfig(t *testing.T) {
	env := &Envelope{
		Transports: []TransportDescriptor{{
			Name: "primary",
			Kind: "grpc",
			Mode: "hot",
			TypedConfig: &TransportConfig{
				TransportRef: "livekit-default",
				ModeReason:   "streaming-default",
			},
		}},
	}

	frame, err := ToBinaryFrame(env)
	if err != nil {
		t.Fatalf("ToBinaryFrame failed: %v", err)
	}
	decoded, err := FromBinaryFrame(frame)
	if err != nil {
		t.Fatalf("FromBinaryFrame failed: %v", err)
	}
	if len(decoded.Transports) != 1 {
		t.Fatalf("expected 1 transport, got %d", len(decoded.Transports))
	}
	typed := decoded.Transports[0].TypedConfig
	if typed == nil {
		t.Fatal("expected typed transport config to round-trip")
	}
	if typed.TransportRef != "livekit-default" || typed.ModeReason != "streaming-default" {
		t.Fatalf("typed config mismatch: %+v", typed)
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

func TestMarshalReturnsErrorWhenEnvelopeIsNil(t *testing.T) {
	_, err := Marshal(nil)
	if err == nil {
		t.Fatal("expected nil envelope error")
	}
}

func TestUnmarshalReturnsErrorWhenPayloadIsEmpty(t *testing.T) {
	_, err := Unmarshal(nil)
	if err == nil {
		t.Fatal("expected empty payload error")
	}
}

func TestUnmarshalReturnsErrorWhenPayloadIsInvalidJSON(t *testing.T) {
	_, err := Unmarshal([]byte(`{"payload":`))
	if err == nil {
		t.Fatal("expected invalid JSON error")
	}
}

func TestFromBinaryFrameReturnsErrorWhenFrameIsNil(t *testing.T) {
	_, err := FromBinaryFrame(nil)
	if err == nil {
		t.Fatal("expected nil frame error")
	}
}

func TestFromBinaryFrameReturnsErrorWhenPayloadIsInvalidJSON(t *testing.T) {
	_, err := FromBinaryFrame(&transportpb.BinaryFrame{
		Payload:  []byte(`{"payload":`),
		MimeType: MIMEType,
	})
	if err == nil {
		t.Fatal("expected invalid JSON error")
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

func TestMarshalDoesNotDefaultKindWhenNoPayloadInputsOrTransports(t *testing.T) {
	data, err := Marshal(&Envelope{})
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	decoded, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if decoded.Kind != "" {
		t.Fatalf("expected empty kind, got %q", decoded.Kind)
	}
}

func TestMarshalDefaultsKindWhenOnlyInputsPresent(t *testing.T) {
	data, err := Marshal(&Envelope{
		Inputs: json.RawMessage(`{"prompt":"hello"}`),
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

func TestMarshalDefaultsKindWhenOnlyTransportsPresent(t *testing.T) {
	data, err := Marshal(&Envelope{
		Transports: []TransportDescriptor{{Name: "primary"}},
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

func TestMarshalPreservesExplicitVersion(t *testing.T) {
	data, err := Marshal(&Envelope{
		Version: "v2",
		Payload: json.RawMessage(`{"ok":true}`),
	})
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	decoded, err := Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if decoded.Version != "v2" {
		t.Fatalf("expected version v2, got %q", decoded.Version)
	}
}

func TestMarshalRejectsInvalidRawJSON(t *testing.T) {
	_, err := Marshal(&Envelope{
		Payload: json.RawMessage(`{"text":`),
	})
	if err == nil {
		t.Fatal("expected invalid payload error")
	}
	if !strings.Contains(err.Error(), "payload") {
		t.Fatalf("expected payload error, got %v", err)
	}
}

func TestMarshalRejectsInvalidInputsJSON(t *testing.T) {
	_, err := Marshal(&Envelope{
		Inputs: json.RawMessage(`{"prompt":`),
	})
	if err == nil {
		t.Fatal("expected invalid inputs error")
	}
	if !strings.Contains(err.Error(), "inputs") {
		t.Fatalf("expected inputs error, got %v", err)
	}
}

func TestMarshalRejectsNegativeTimestamp(t *testing.T) {
	_, err := Marshal(&Envelope{
		TimestampMs: -1,
		Payload:     json.RawMessage(`{"ok":true}`),
	})
	if err == nil {
		t.Fatal("expected negative timestamp error")
	}
}

func TestUnmarshalRejectsOversizedEnvelope(t *testing.T) {
	_, err := Unmarshal(make([]byte, MaxEnvelopeSize+1))
	if err == nil {
		t.Fatal("expected oversized envelope error")
	}
	if !strings.Contains(err.Error(), "exceeds maximum size") {
		t.Fatalf("expected maximum size error, got %v", err)
	}
}

func TestFromBinaryFrameRejectsTimestampOverflow(t *testing.T) {
	frame := &transportpb.BinaryFrame{
		Payload:     []byte(`{"version":"v1","kind":"data","payload":{"ok":true}}`),
		MimeType:    MIMEType,
		TimestampMs: uint64(math.MaxInt64) + 1,
	}

	_, err := FromBinaryFrame(frame)
	if err == nil {
		t.Fatal("expected timestamp overflow error")
	}
	if !strings.Contains(err.Error(), "timestamp") {
		t.Fatalf("expected timestamp error, got %v", err)
	}
}

func TestHookPayloadMarshalUnmarshalRoundTrip(t *testing.T) {
	payload := HookPayload{
		Version: "v1",
		Event:   "storyrun.ready",
		Source:  "bobrapet",
		Data: map[string]any{
			"storyRun": "storyrun-1",
		},
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	var decoded HookPayload
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if decoded.Event != payload.Event || decoded.Source != payload.Source {
		t.Fatalf("hook payload mismatch: got %+v want %+v", decoded, payload)
	}
	if decoded.Data["storyRun"] != "storyrun-1" {
		t.Fatalf("hook payload data mismatch: %+v", decoded.Data)
	}
}

func TestHookPayloadWithNilData(t *testing.T) {
	payload := HookPayload{
		Event: "session.start",
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}
	var decoded HookPayload
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}
	if decoded.Event != "session.start" {
		t.Fatalf("hook event mismatch: %q", decoded.Event)
	}
	if decoded.Data != nil {
		t.Fatalf("expected nil hook data, got %+v", decoded.Data)
	}
}
