package config

const (
	MCPStdio = "stdio"
	MCPSse   = "sse"
)

type LSPServerConfig struct {
	Command  string                 `json:"command,omitempty"`
	Args     []string               `json:"args,omitempty"`
	Disabled bool                   `json:"disabled,omitempty"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

type MCPConfig struct {
	Command string            `json:"command,omitempty"`
	Args    []string          `json:"args,omitempty"`
	Env     []string          `json:"env,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
	Type    string            `json:"type,omitempty"`
	URL     string            `json:"url,omitempty"`
}

type ProviderConfig struct {
	APIKey   string `json:"apiKey,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
	Provider string `json:"provider,omitempty"`
}
