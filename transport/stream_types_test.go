package transport

import (
	"regexp"
	"testing"

	transportv1 "github.com/bubustack/tractatus/gen/go/proto/transport/v1"
)

func TestStreamTypeConstantsValuesAreStable(t *testing.T) {
	tests := map[string]string{
		"StreamTypeSpeechTranscript":      StreamTypeSpeechTranscript,
		"StreamTypeSpeechTranslation":     StreamTypeSpeechTranslation,
		"StreamTypeSpeechTranscriptDelta": StreamTypeSpeechTranscriptDelta,
		"StreamTypeSpeechTranscriptDone":  StreamTypeSpeechTranscriptDone,
		"StreamTypeSpeechVADActive":       StreamTypeSpeechVADActive,
		"StreamTypeSpeechVADInactive":     StreamTypeSpeechVADInactive,
		"StreamTypeSpeechAudio":           StreamTypeSpeechAudio,
		"StreamTypeSpeechAudioDelta":      StreamTypeSpeechAudioDelta,
		"StreamTypeSpeechAudioDone":       StreamTypeSpeechAudioDone,
		"StreamTypeSpeechAudioSummary":    StreamTypeSpeechAudioSummary,
		"StreamTypeSpeechTurn":            StreamTypeSpeechTurn,
		"StreamTypeChatMessage":           StreamTypeChatMessage,
		"StreamTypeChatResponse":          StreamTypeChatResponse,
		"StreamTypeOpenAIChat":            StreamTypeOpenAIChat,
	}
	expected := map[string]string{
		"StreamTypeSpeechTranscript":      "speech.transcript.v1",
		"StreamTypeSpeechTranslation":     "speech.translation.v1",
		"StreamTypeSpeechTranscriptDelta": "speech.transcript.delta",
		"StreamTypeSpeechTranscriptDone":  "speech.transcript.done",
		"StreamTypeSpeechVADActive":       "speech.vad.active",
		"StreamTypeSpeechVADInactive":     "speech.vad.inactive",
		"StreamTypeSpeechAudio":           "speech.audio.v1",
		"StreamTypeSpeechAudioDelta":      "speech.audio.delta",
		"StreamTypeSpeechAudioDone":       "speech.audio.done",
		"StreamTypeSpeechAudioSummary":    "speech.audio.summary",
		"StreamTypeSpeechTurn":            "speech.turn.v1",
		"StreamTypeChatMessage":           "chat.message.v1",
		"StreamTypeChatResponse":          "chat.response.v1",
		"StreamTypeOpenAIChat":            "openai.chat.v1",
	}
	for name, got := range tests {
		if got != expected[name] {
			t.Fatalf("%s changed wire value: got %q, want %q", name, got, expected[name])
		}
	}
}

func TestStreamTypeConstantsFollowNamingConvention(t *testing.T) {
	pattern := regexp.MustCompile(`^[a-z0-9]+(\.[a-z0-9]+)+(\.v[0-9]+|\.delta|\.done|\.active|\.inactive|\.summary)$`)
	values := []string{
		StreamTypeSpeechTranscript,
		StreamTypeSpeechTranslation,
		StreamTypeSpeechTranscriptDelta,
		StreamTypeSpeechTranscriptDone,
		StreamTypeSpeechVADActive,
		StreamTypeSpeechVADInactive,
		StreamTypeSpeechAudio,
		StreamTypeSpeechAudioDelta,
		StreamTypeSpeechAudioDone,
		StreamTypeSpeechAudioSummary,
		StreamTypeSpeechTurn,
		StreamTypeChatMessage,
		StreamTypeChatResponse,
		StreamTypeOpenAIChat,
	}
	for _, value := range values {
		if !pattern.MatchString(value) {
			t.Fatalf("stream type %q does not follow category.kind.version convention", value)
		}
	}
}

func TestStreamTypeProtoEnumMapsToCanonicalNames(t *testing.T) {
	tests := map[transportv1.StreamType]string{
		transportv1.StreamType_STREAM_TYPE_SPEECH_TRANSCRIPT:       StreamTypeSpeechTranscript,
		transportv1.StreamType_STREAM_TYPE_SPEECH_TRANSLATION:      StreamTypeSpeechTranslation,
		transportv1.StreamType_STREAM_TYPE_SPEECH_TRANSCRIPT_DELTA: StreamTypeSpeechTranscriptDelta,
		transportv1.StreamType_STREAM_TYPE_SPEECH_TRANSCRIPT_DONE:  StreamTypeSpeechTranscriptDone,
		transportv1.StreamType_STREAM_TYPE_SPEECH_VAD_ACTIVE:       StreamTypeSpeechVADActive,
		transportv1.StreamType_STREAM_TYPE_SPEECH_VAD_INACTIVE:     StreamTypeSpeechVADInactive,
		transportv1.StreamType_STREAM_TYPE_SPEECH_AUDIO:            StreamTypeSpeechAudio,
		transportv1.StreamType_STREAM_TYPE_SPEECH_AUDIO_DELTA:      StreamTypeSpeechAudioDelta,
		transportv1.StreamType_STREAM_TYPE_SPEECH_AUDIO_DONE:       StreamTypeSpeechAudioDone,
		transportv1.StreamType_STREAM_TYPE_SPEECH_AUDIO_SUMMARY:    StreamTypeSpeechAudioSummary,
		transportv1.StreamType_STREAM_TYPE_SPEECH_TURN:             StreamTypeSpeechTurn,
		transportv1.StreamType_STREAM_TYPE_CHAT_MESSAGE:            StreamTypeChatMessage,
		transportv1.StreamType_STREAM_TYPE_CHAT_RESPONSE:           StreamTypeChatResponse,
		transportv1.StreamType_STREAM_TYPE_OPENAI_CHAT:             StreamTypeOpenAIChat,
	}
	for typ, expected := range tests {
		got, ok := CanonicalStreamTypeName(typ)
		if !ok {
			t.Fatalf("expected mapping for %s", typ)
		}
		if got != expected {
			t.Fatalf("stream type mapping mismatch for %s: got %q want %q", typ, got, expected)
		}
	}
}
