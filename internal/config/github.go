package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type githubHostConfig struct {
	OAuthToken string `yaml:"oauth_token"`
}

func LoadGitHubToken() (string, error) {
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	tokenPath := home + "/.config/github/cli/hosts.yml"
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return "", fmt.Errorf("github token not found: %w", err)
	}

	var hosts map[string]githubHostConfig
	if err := yaml.Unmarshal(data, &hosts); err != nil {
		return "", fmt.Errorf("failed to parse github config: %w", err)
	}

	for _, host := range hosts {
		if host.OAuthToken != "" {
			return host.OAuthToken, nil
		}
	}

	return "", fmt.Errorf("oauth_token not found in github config")
}
