package validation

import (
	"testing"

	transportv1 "github.com/bubustack/tractatus/gen/go/proto/transport/v1"
)

func TestValidateRejectsInvalidTransportMessage(t *testing.T) {
	err := Validate(&transportv1.AudioFrame{
		Pcm:          []byte("pcm"),
		SampleRateHz: 7999,
		Channels:     2,
	})
	if err == nil {
		t.Fatal("expected invalid transport message to be rejected")
	}
}

func TestValidateAcceptsValidTransportMessage(t *testing.T) {
	err := Validate(&transportv1.AudioFrame{
		Pcm:          []byte("pcm"),
		SampleRateHz: 48000,
		Channels:     2,
		Codec:        "pcm16",
	})
	if err != nil {
		t.Fatalf("expected valid transport message to pass: %v", err)
	}
}
