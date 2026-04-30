package transport

import transportv1 "github.com/bubustack/tractatus/gen/go/proto/transport/v1"

var streamTypeNames = map[transportv1.StreamType]string{
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

// CanonicalStreamTypeName returns the dot-delimited wire name for a proto
// StreamType enum value.
func CanonicalStreamTypeName(typ transportv1.StreamType) (string, bool) {
	name, ok := streamTypeNames[typ]
	return name, ok
}
