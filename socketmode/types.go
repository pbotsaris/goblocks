package socketmode

import "encoding/json"

// Envelope wraps all messages received from Slack via Socket Mode.
type Envelope struct {
	EnvelopeID             string          `json:"envelope_id"`
	Type                   string          `json:"type"`
	Payload                json.RawMessage `json:"payload"`
	AcceptsResponsePayload bool            `json:"accepts_response_payload"`
	RetryAttempt           int             `json:"retry_attempt,omitempty"`
	RetryReason            string          `json:"retry_reason,omitempty"`
}

// Envelope types sent by Slack.
const (
	EnvelopeTypeHello         = "hello"
	EnvelopeTypeDisconnect    = "disconnect"
	EnvelopeTypeEventsAPI     = "events_api"
	EnvelopeTypeInteractive   = "interactive"
	EnvelopeTypeSlashCommands = "slash_commands"
)

// HelloMessage is sent by Slack upon successful WebSocket connection.
type HelloMessage struct {
	Type           string         `json:"type"`
	ConnectionInfo ConnectionInfo `json:"connection_info"`
	NumConnections int            `json:"num_connections"`
	DebugInfo      DebugInfo      `json:"debug_info"`
}

// ConnectionInfo contains app identification from the hello message.
type ConnectionInfo struct {
	AppID string `json:"app_id"`
}

// DebugInfo contains connection debugging information from the hello message.
type DebugInfo struct {
	Host                      string `json:"host"`
	Started                   string `json:"started"`
	BuildNumber               int    `json:"build_number"`
	ApproximateConnectionTime int    `json:"approximate_connection_time"`
}

// DisconnectMessage is sent by Slack when requesting disconnection.
type DisconnectMessage struct {
	Type      string    `json:"type"`
	Reason    string    `json:"reason"`
	DebugInfo DebugInfo `json:"debug_info"`
}

// Disconnect reasons sent by Slack.
const (
	DisconnectReasonLinkDisabled     = "link_disabled"
	DisconnectReasonWarning          = "warning"
	DisconnectReasonRefreshRequested = "refresh_requested"
)

// Ack is the acknowledgment sent back to Slack for each envelope.
type Ack struct {
	EnvelopeID string `json:"envelope_id"`
	Payload    any    `json:"payload,omitempty"`
}

// ConnectionOpenResponse is the response from apps.connections.open API.
type ConnectionOpenResponse struct {
	OK    bool   `json:"ok"`
	URL   string `json:"url"`
	Error string `json:"error,omitempty"`
}
