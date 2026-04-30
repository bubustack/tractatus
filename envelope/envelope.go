package envelope

import (
	"encoding/json"
	"fmt"
	"math"

	transportpb "github.com/bubustack/tractatus/gen/go/proto/transport/v1"
)

const (
	// MIMEType is the canonical identifier for envelope-based BinaryFrame payloads.
	MIMEType = "application/vnd.bubu.packet+json"
	// MaxEnvelopeSize is the maximum JSON envelope payload accepted for decoding.
	MaxEnvelopeSize = 10 * 1024 * 1024
	// LatestVersion is the default schema version applied when one is not provided.
	LatestVersion = "v1"

	// KindData is the default value for regular business payload packets.
	KindData = "data"
	// KindHeartbeat marks transport keepalive packets.
	KindHeartbeat = "heartbeat"
	// KindNoop marks packets that should be ignored by downstream business logic.
	KindNoop = "noop"
	// KindHook marks event-style packets emitted by runtime sources (session.start, storyrun.ready, etc.).
	KindHook = "hook"
)

// Envelope mirrors the historical DataPacket shape that engrams exchange.
// It allows connectors/SDKs in any language to embed structured payloads inside
// BinaryFrame messages without baking provider-specific dependencies into the
// transport contract.
type Envelope struct {
	// Version is the envelope schema version carried inside the BinaryFrame payload.
	Version string `json:"version,omitempty"`
	// Kind annotates the semantic intent of the packet (e.g., "data", "heartbeat", "telemetry").
	Kind string `json:"kind,omitempty"`
	// MessageID carries an optional sender-defined identifier useful for deduplication or tracing.
	MessageID string `json:"messageId,omitempty"`
	// TimestampMs encodes when the packet was produced, expressed in Unix milliseconds.
	TimestampMs int64 `json:"timestampMs,omitempty"`
	// Metadata carries additional packet-level routing or tracing attributes.
	Metadata map[string]string `json:"metadata,omitempty"`
	// Payload carries the primary structured business payload.
	Payload json.RawMessage `json:"payload,omitempty"`
	// Inputs carries evaluated step inputs associated with the packet.
	Inputs json.RawMessage `json:"inputs,omitempty"`
	// Transports enumerates transport descriptors relevant to the packet.
	Transports []TransportDescriptor `json:"transports,omitempty"`
}

// HookPayload is the recommended event payload shape when Envelope.Kind == KindHook.
// It keeps hook packets generic across AI and non-AI workflows.
type HookPayload struct {
	// Version is the hook payload schema version.
	Version string `json:"version,omitempty"`
	// Event names the emitted hook event (for example, "session.start").
	Event string `json:"event"`
	// Source identifies the component or runtime that emitted the hook.
	Source string `json:"source,omitempty"`
	// Data carries event-specific structured details.
	Data map[string]any `json:"data,omitempty"`
}

// TransportDescriptor mirrors Story transport declarations so downstream steps
// can infer hot vs cold paths without reloading pod env.
type TransportDescriptor struct {
	// Name is the logical transport identifier.
	Name string `json:"name"`
	// Kind identifies the transport implementation or class.
	Kind string `json:"kind,omitempty"`
	// Mode identifies how the transport is used (for example, hot vs cold path).
	Mode string `json:"mode,omitempty"`
	// Config carries transport-specific structured configuration.
	Config map[string]any `json:"config,omitempty"`
	// TypedConfig carries bounded runtime descriptor metadata.
	TypedConfig *TransportConfig `json:"typedConfig,omitempty"`
}

// TransportConfig carries safe runtime descriptor metadata.
type TransportConfig struct {
	// TransportRef identifies the Transport resource that produced this descriptor.
	TransportRef string `json:"transportRef,omitempty"`
	// ModeReason explains why the runtime selected the descriptor mode.
	ModeReason string `json:"modeReason,omitempty"`
}

// Marshal encodes the envelope into JSON suitable for a BinaryFrame payload.
func Marshal(env *Envelope) ([]byte, error) {
	if env == nil {
		return nil, fmt.Errorf("envelope is nil")
	}
	out := *env
	if out.TimestampMs < 0 {
		return nil, fmt.Errorf("envelope timestamp_ms must not be negative")
	}
	if len(out.Payload) > 0 && !json.Valid(out.Payload) {
		return nil, fmt.Errorf("envelope payload is not valid JSON")
	}
	if len(out.Inputs) > 0 && !json.Valid(out.Inputs) {
		return nil, fmt.Errorf("envelope inputs is not valid JSON")
	}
	if out.Version == "" {
		out.Version = LatestVersion
	}
	if out.Kind == "" && (len(out.Payload) > 0 || len(out.Inputs) > 0 || len(out.Transports) > 0) {
		out.Kind = KindData
	}
	return json.Marshal(&out)
}

// Unmarshal decodes JSON payloads into an Envelope structure.
func Unmarshal(data []byte) (*Envelope, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("envelope payload empty")
	}
	if len(data) > MaxEnvelopeSize {
		return nil, fmt.Errorf("envelope payload exceeds maximum size of %d bytes", MaxEnvelopeSize)
	}
	var env Envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, err
	}
	if env.Version == "" {
		env.Version = LatestVersion
	}
	return &env, nil
}

// ToBinaryFrame converts the envelope into a BinaryFrame message.
func ToBinaryFrame(env *Envelope) (*transportpb.BinaryFrame, error) {
	payload, err := Marshal(env)
	if err != nil {
		return nil, err
	}
	frame := &transportpb.BinaryFrame{
		Payload:  payload,
		MimeType: MIMEType,
	}
	if env != nil && env.TimestampMs > 0 {
		frame.TimestampMs = uint64(env.TimestampMs)
	}
	return frame, nil
}

// FromBinaryFrame decodes a BinaryFrame into an envelope as long as the MIME
// type matches the contract.
func FromBinaryFrame(frame *transportpb.BinaryFrame) (*Envelope, error) {
	if frame == nil {
		return nil, fmt.Errorf("binary frame is nil")
	}
	if frame.GetMimeType() != MIMEType {
		return nil, fmt.Errorf("unexpected mime type %q", frame.GetMimeType())
	}
	env, err := Unmarshal(frame.GetPayload())
	if err != nil {
		return nil, err
	}
	if env.TimestampMs == 0 && frame.GetTimestampMs() > 0 {
		if frame.GetTimestampMs() > uint64(math.MaxInt64) {
			return nil, fmt.Errorf("binary frame timestamp_ms %d exceeds maximum supported timestamp", frame.GetTimestampMs())
		}
		env.TimestampMs = int64(frame.GetTimestampMs())
	}
	return env, nil
}
