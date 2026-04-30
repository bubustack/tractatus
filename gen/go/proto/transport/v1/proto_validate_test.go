package transportv1

import (
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

func TestProtovalidateRejectsInvalidTransportMessages(t *testing.T) {
	oversizedFrame := strings.Repeat("a", 10*1024*1024+1)
	oversizedMetadataValue := strings.Repeat("v", 1025)

	tests := map[string]proto.Message{
		"audio pcm too large": &AudioFrame{
			Pcm:          []byte(oversizedFrame),
			SampleRateHz: 48000,
			Channels:     2,
		},
		"audio sample rate too low": &AudioFrame{
			Pcm:          []byte("pcm"),
			SampleRateHz: 7999,
			Channels:     2,
		},
		"audio channels too high": &AudioFrame{
			Pcm:          []byte("pcm"),
			SampleRateHz: 48000,
			Channels:     9,
		},
		"video payload too large": &VideoFrame{
			Payload: []byte(oversizedFrame),
			Codec:   "h264",
			Width:   1920,
			Height:  1080,
		},
		"binary mime type too long": &BinaryFrame{
			Payload:  []byte("payload"),
			MimeType: strings.Repeat("m", 256),
		},
		"publish metadata too many pairs": &PublishRequest{
			Metadata: manyMetadataPairs(101),
		},
		"publish metadata empty key": &PublishRequest{
			Metadata: map[string]string{"": "value"},
		},
		"publish metadata value too long": &PublishRequest{
			Metadata: map[string]string{"key": oversizedMetadataValue},
		},
		"publish transports too many": &PublishRequest{
			Transports: manyTransportDescriptors(11),
		},
		"data packet transports too many": &DataPacket{
			Transports: manyTransportDescriptors(11),
		},
		"transport typed config transport ref too long": &TransportDescriptor{
			Name: "primary",
			TypedConfig: &TransportConfig{
				TransportRef: strings.Repeat("t", 254),
			},
		},
		"transport typed config mode reason too long": &TransportDescriptor{
			Name: "primary",
			TypedConfig: &TransportConfig{
				TransportRef: "livekit",
				ModeReason:   strings.Repeat("m", 513),
			},
		},
		"partition acks too many": &FlowControl{
			PartitionAcks: manyPartitionAcks(257),
		},
		"control error message too long": &ControlResponse{
			Action: ControlAction_CONTROL_ACTION_ERROR,
			Error: &ControlError{
				Code:    1,
				Message: strings.Repeat("e", 2049),
			},
		},
		"binding info endpoint too long": &BindingInfo{
			Endpoint: strings.Repeat("e", 2049),
		},
		"binding info payload too large": &BindingInfo{
			Payload: []byte(strings.Repeat("p", 1024*1024+1)),
		},
		"chunk index out of range": &StreamEnvelope{
			StreamId:   "stream-1",
			ChunkId:    "chunk-1",
			ChunkIndex: 2,
			ChunkCount: 2,
		},
		"chunk missing stream id": &StreamEnvelope{
			ChunkId:    "chunk-1",
			ChunkIndex: 0,
			ChunkCount: 2,
		},
		"chunk fields set without chunk count": &StreamEnvelope{
			ChunkId:    "chunk-1",
			ChunkIndex: 1,
		},
	}

	for name, msg := range tests {
		t.Run(name, func(t *testing.T) {
			if err := protovalidate.Validate(msg); err == nil {
				t.Fatalf("expected protovalidate to reject %T", msg)
			}
		})
	}
}

func manyMetadataPairs(count int) map[string]string {
	metadata := make(map[string]string, count)
	for i := 0; i < count; i++ {
		metadata["key-"+strings.Repeat("x", i%8)+string(rune('a'+i%26))+string(rune('A'+i%26))] = "value"
	}
	return metadata
}

func manyTransportDescriptors(count int) []*TransportDescriptor {
	transports := make([]*TransportDescriptor, 0, count)
	for i := 0; i < count; i++ {
		transports = append(transports, &TransportDescriptor{Name: "transport"})
	}
	return transports
}

func manyPartitionAcks(count int) []*PartitionAck {
	acks := make([]*PartitionAck, 0, count)
	for i := 0; i < count; i++ {
		acks = append(acks, &PartitionAck{Partition: "partition"})
	}
	return acks
}
