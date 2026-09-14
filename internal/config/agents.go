package config

import (
	"fmt"

	"kotha/internal/llm/models"
)

type AgentName string

const (
	AgentCoder      AgentName = "coder"
	AgentTask       AgentName = "task"
	AgentTitle      AgentName = "title"
	AgentSummarizer AgentName = "summarizer"
)

type Agent struct {
	Model           models.ModelID `json:"model,omitempty"`
	MaxTokens       int            `json:"maxTokens,omitempty"`
	ReasoningEffort string         `json:"reasoningEffort,omitempty"`
}

type ErrUnsupportedModel struct {
	Model models.ModelID
}

func (e ErrUnsupportedModel) Error() string {
	return fmt.Sprintf("model %s is not supported", e.Model)
}

func validateAgent(cfg *Config, name AgentName, agent Agent) error {
	if agent.Model == "" {
		return nil
	}
	if _, ok := models.SupportedModels[agent.Model]; !ok {
		return ErrUnsupportedModel{Model: agent.Model}
	}
	return nil
}
