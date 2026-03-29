package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type discordClient interface {
	JoinChannel(roleID string) error
	GetJoinChannelLink(roleID string) string
	CreateChannel(channelName string) (string, string, error)
}

type DiscordHttpClient struct {
	BaseURL    string
	HttpClient *http.Client
	GuildID    string
}

type DiscordCreatedResponse struct {
	RoleID    string `json:"roleID"`
	ChannelID string `json:"channelID"`
}

type DiscordJoinResponse struct {
	authResponse string `json:`
}

func NewDiscordClient(baseURL string, guildID string) *DiscordHttpClient {
	return &DiscordHttpClient{
		BaseURL: baseURL,
		HttpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		GuildID: guildID,
	}
}

func (c *DiscordHttpClient) CreateChannel(channelName string) (string, string, error) {
	requestData := map[string]string{
		"name":    channelName,
		"guildID": c.GuildID,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/create-channel", c.BaseURL+"/bot-api"), "application/json", bytes.NewBuffer(jsonData))
	if err != nil || resp.StatusCode != http.StatusCreated {
		return "", "", fmt.Errorf("bot return error status: %v", err)
	}
	defer resp.Body.Close()

	var data DiscordCreatedResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", "", fmt.Errorf("failed to decode body response: %v", err)
	}

	return data.RoleID, data.ChannelID, err

}

func (c *DiscordHttpClient) JoinChannel(roleID string) (string, error) {
	// We don't use http.Get here.
	// We construct the URL that the STUDENT'S browser needs to visit.
	if c.BaseURL == "" {
		return "", fmt.Errorf("discord base URL is not configured")
	}

	joinURL := fmt.Sprintf("%s/join-course/%s", c.BaseURL+"/bot-api", roleID)

	// Return the URL string to the service/controller
	return joinURL, nil
}

func (c *DiscordHttpClient) GetJoinChannelLink(roleID string) string {
	return ""
}
