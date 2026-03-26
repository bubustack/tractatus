package transport

// Stream event type identifiers shared across connectors and engrams.
// They intentionally remain provider-agnostic so multiple STT/VAD engines can
// interoperate on the same transport without per-vendor routing rules.
const (
	// StreamTypeSpeechTranscript is a finalized speech-to-text transcript payload.
	StreamTypeSpeechTranscript = "speech.transcript.v1"
	// StreamTypeSpeechTranslation is a finalized translated transcript payload.
	StreamTypeSpeechTranslation = "speech.translation.v1"
	// StreamTypeSpeechTranscriptDelta carries incremental transcript text.
	StreamTypeSpeechTranscriptDelta = "speech.transcript.delta"
	// StreamTypeSpeechTranscriptDone finalizes a transcript stream (usage, logprobs, etc.).
	StreamTypeSpeechTranscriptDone = "speech.transcript.done"
	// StreamTypeSpeechVADActive indicates VAD detected speech onset.
	StreamTypeSpeechVADActive = "speech.vad.active"
	// StreamTypeSpeechVADInactive indicates VAD detected silence/end of speech.
	StreamTypeSpeechVADInactive = "speech.vad.inactive"

	// StreamTypeSpeechAudio represents synthesized speech/audio payloads (e.g., TTS results).
	StreamTypeSpeechAudio = "speech.audio.v1"
	// StreamTypeSpeechAudioDelta carries incremental synthesized audio chunks.
	StreamTypeSpeechAudioDelta = "speech.audio.delta"
	// StreamTypeSpeechAudioDone closes a streaming audio sequence or provides a summary.
	StreamTypeSpeechAudioDone = "speech.audio.done"
	// StreamTypeSpeechAudioSummary carries synthesized-audio summary metadata.
	StreamTypeSpeechAudioSummary = "speech.audio.summary"
	// StreamTypeSpeechTurn marks a turn-detection event emitted by realtime pipelines.
	StreamTypeSpeechTurn = "speech.turn.v1"

	// StreamTypeChatMessage carries a text chat message (e.g., from LiveKit data channels).
	StreamTypeChatMessage = "chat.message.v1"
	// StreamTypeChatResponse carries an assistant/bot reply to a chat message.
	StreamTypeChatResponse = "chat.response.v1"
	// StreamTypeOpenAIChat carries OpenAI chat payloads emitted by the chat engram.
	StreamTypeOpenAIChat = "openai.chat.v1"
)
