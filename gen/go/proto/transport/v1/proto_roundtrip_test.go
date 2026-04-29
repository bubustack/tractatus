package transportv1

import (
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestDataRequestRoundTripAllFrameTypes(t *testing.T) {
	tests := map[string]*DataRequest{
		"audio": {
			Frame: &DataRequest_Audio{Audio: &AudioFrame{
				Pcm:          []byte("audio-data"),
				SampleRateHz: 48000,
				Channels:     2,
				Codec:        "opus",
				TimestampMs:  1234567890,
			}},
		},
		"video": {
			Frame: &DataRequest_Video{Video: &VideoFrame{
				Payload:     []byte("video-data"),
				Codec:       "h264",
				Width:       1920,
				Height:      1080,
				TimestampMs: 1234567890,
			}},
		},
		"binary": {
			Frame: &DataRequest_Binary{Binary: &BinaryFrame{
				Payload:     []byte(`{"ok":true}`),
				MimeType:    "application/json",
				TimestampMs: 1234567890,
			}},
		},
	}

	for name, msg := range tests {
		t.Run(name, func(t *testing.T) {
			roundTripDataRequest(t, msg)
		})
	}
}

func TestDataRequestRoundTripComplexStructuredContext(t *testing.T) {
	payload, err := structpb.NewStruct(map[string]any{
		"text": "hello",
		"nested": map[string]any{
			"score": 0.99,
			"tags":  []any{"chat", "demo"},
		},
	})
	if err != nil {
		t.Fatalf("payload struct: %v", err)
	}
	inputs, err := structpb.NewStruct(map[string]any{
		"prompt": "translate",
		"limits": map[string]any{"maxTokens": 128},
	})
	if err != nil {
		t.Fatalf("inputs struct: %v", err)
	}
	cfg, err := structpb.NewStruct(map[string]any{
		"endpoint": "unix:///tmp/bubu.sock",
		"tls":      true,
	})
	if err != nil {
		t.Fatalf("transport config struct: %v", err)
	}

	roundTripDataRequest(t, &DataRequest{
		Frame: &DataRequest_Binary{Binary: &BinaryFrame{
			Payload:  []byte(`{"ok":true}`),
			MimeType: "application/json",
		}},
		Metadata: map[string]string{
			"storyrun": "sr-1",
			"step":     "ingress",
		},
		Payload: payload,
		Inputs:  inputs,
		Transports: []*TransportDescriptor{{
			Name:   "primary",
			Kind:   "grpc",
			Mode:   "hot",
			Config: cfg,
		}},
		Envelope: &StreamEnvelope{
			StreamId:   "stream-1",
			Sequence:   42,
			Partition:  "speaker-1",
			ChunkId:    "chunk-1",
			ChunkIndex: 1,
			ChunkCount: 3,
			ChunkBytes: 512,
			TotalBytes: 1536,
		},
	})
}

func TestControlRequestRoundTripAllActions(t *testing.T) {
	for name, value := range ControlAction_value {
		action := ControlAction(value)
		msg := &ControlRequest{
			Action: action,
			Metadata: map[string]string{
				"name": name,
			},
			Flow: &FlowControl{
				Ack:             12,
				CreditsMessages: 4,
				CreditsBytes:    2048,
				Signal:          FlowControlSignal_FLOW_CONTROL_SIGNAL_RESUME,
			},
		}
		data, err := proto.Marshal(msg)
		if err != nil {
			t.Fatalf("marshal %s: %v", name, err)
		}
		var decoded ControlRequest
		if err := proto.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("unmarshal %s: %v", name, err)
		}
		if !proto.Equal(msg, &decoded) {
			t.Fatalf("round trip mismatch for %s: got %+v want %+v", name, &decoded, msg)
		}
	}
}

func TestControlResponseRoundTripStructuredError(t *testing.T) {
	msg := &ControlResponse{
		Action: ControlAction_CONTROL_ACTION_ERROR,
		Error: &ControlError{
			Code:    42,
			Message: "downstream rejected frame",
			Details: map[string]string{
				"component": "connector",
				"reason":    "invalid_schema",
			},
		},
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded ControlResponse
	if err := proto.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !proto.Equal(msg, &decoded) {
		t.Fatalf("round trip mismatch: got %+v want %+v", &decoded, msg)
	}
}

func TestStreamTypeEnumNamesAreStable(t *testing.T) {
	expected := map[StreamType]string{
		StreamType_STREAM_TYPE_SPEECH_TRANSCRIPT:       "STREAM_TYPE_SPEECH_TRANSCRIPT",
		StreamType_STREAM_TYPE_SPEECH_TRANSLATION:      "STREAM_TYPE_SPEECH_TRANSLATION",
		StreamType_STREAM_TYPE_SPEECH_TRANSCRIPT_DELTA: "STREAM_TYPE_SPEECH_TRANSCRIPT_DELTA",
		StreamType_STREAM_TYPE_SPEECH_TRANSCRIPT_DONE:  "STREAM_TYPE_SPEECH_TRANSCRIPT_DONE",
		StreamType_STREAM_TYPE_SPEECH_VAD_ACTIVE:       "STREAM_TYPE_SPEECH_VAD_ACTIVE",
		StreamType_STREAM_TYPE_SPEECH_VAD_INACTIVE:     "STREAM_TYPE_SPEECH_VAD_INACTIVE",
		StreamType_STREAM_TYPE_SPEECH_AUDIO:            "STREAM_TYPE_SPEECH_AUDIO",
		StreamType_STREAM_TYPE_SPEECH_AUDIO_DELTA:      "STREAM_TYPE_SPEECH_AUDIO_DELTA",
		StreamType_STREAM_TYPE_SPEECH_AUDIO_DONE:       "STREAM_TYPE_SPEECH_AUDIO_DONE",
		StreamType_STREAM_TYPE_SPEECH_AUDIO_SUMMARY:    "STREAM_TYPE_SPEECH_AUDIO_SUMMARY",
		StreamType_STREAM_TYPE_SPEECH_TURN:             "STREAM_TYPE_SPEECH_TURN",
		StreamType_STREAM_TYPE_CHAT_MESSAGE:            "STREAM_TYPE_CHAT_MESSAGE",
		StreamType_STREAM_TYPE_CHAT_RESPONSE:           "STREAM_TYPE_CHAT_RESPONSE",
		StreamType_STREAM_TYPE_OPENAI_CHAT:             "STREAM_TYPE_OPENAI_CHAT",
	}
	for value, name := range expected {
		if value.String() != name {
			t.Fatalf("stream type enum changed: got %q want %q", value.String(), name)
		}
	}
}

func roundTripDataRequest(t *testing.T, msg *DataRequest) {
	t.Helper()
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded DataRequest
	if err := proto.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !proto.Equal(msg, &decoded) {
		t.Fatalf("round trip mismatch: got %+v want %+v", &decoded, msg)
	}
}
