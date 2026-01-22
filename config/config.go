package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the complete application configuration
type Config struct {
	Platforms PlatformsConfig `json:"platforms"`
	OpenAI    OpenAIConfig    `json:"openai"`
	Database  DatabaseConfig  `json:"database"`
	YouTube   YouTubeConfig   `json:"youtube"`
	Community CommunityConfig `json:"community"`
	DryRun    bool            `json:"dry_run"`
}

// PlatformsConfig contains all platform-specific credentials
type PlatformsConfig struct {
	Twitter  []TwitterCredentials  `json:"twitter"`
	LinkedIn []LinkedInCredentials `json:"linkedin"`
	Discord  []DiscordCredentials  `json:"discord"`
	Bluesky  []BlueskyCredentials  `json:"bluesky"`
}

// TwitterCredentials represents OAuth1 credentials for Twitter
type TwitterCredentials struct {
	Handle         string `json:"handle"`
	ConsumerKey    string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
	AccessToken    string `json:"access_token"`
	AccessSecret   string `json:"access_secret"`
}

// LinkedInCredentials represents OAuth2 credentials for LinkedIn
type LinkedInCredentials struct {
	Handle      string `json:"handle"`
	AccessToken string `json:"access_token"`
	// Optional: for token refresh
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresAt    string `json:"expires_at,omitempty"`
}

// DiscordCredentials represents Discord bot credentials
type DiscordCredentials struct {
	Handle    string `json:"handle"`
	BotToken  string `json:"bot_token"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`
}

// BlueskyCredentials represents Bluesky AT Protocol credentials
type BlueskyCredentials struct {
	Handle      string `json:"handle"`
	AppPassword string `json:"app_password"`
	// Optional: custom PDS (Personal Data Server) endpoint
	PDSEndpoint string `json:"pds_endpoint,omitempty"`
}

// OpenAIConfig represents OpenAI-compatible API configuration
type OpenAIConfig struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"` // e.g., "https://api.openai.com/v1"
	Model   string `json:"model"`    // e.g., "gpt-4"
	// Optional: for custom parameters
	Temperature      float64 `json:"temperature,omitempty"`
	MaxTokens        int     `json:"max_tokens,omitempty"`
	TopP             float64 `json:"top_p,omitempty"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64 `json:"presence_penalty,omitempty"`
}

// DatabaseConfig represents Supabase database configuration
type DatabaseConfig struct {
	SupabaseURL        string `json:"supabase_url"`         // e.g., "https://xxx.supabase.co"
	SupabaseKey        string `json:"supabase_key"`         // Supabase anon/service key
	SupabaseDBPassword string `json:"supabase_db_password"` // Direct PostgreSQL password
	SupabaseProjectRef string `json:"supabase_project_ref"` // Project reference ID
}

// YouTubeConfig represents YouTube RSS feed configuration
type YouTubeConfig struct {
	Channels []string `json:"channels"` // YouTube channel IDs
}

// CommunityConfig represents community-specific context
type CommunityConfig struct {
	PrimarySite string `json:"primary_site"` // e.g., "forgeutah.tech"
	EventsSite  string `json:"events_site"`  // e.g., "utahtechevents.dev"
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config JSON: %w", err)
	}

	// Override with environment variables if present
	config.applyEnvironmentOverrides()

	return &config, nil
}

// applyEnvironmentOverrides allows environment variables to override config file values
func (c *Config) applyEnvironmentOverrides() {
	// Database overrides
	if url := os.Getenv("SUPABASE_URL"); url != "" {
		c.Database.SupabaseURL = url
	}
	if key := os.Getenv("SUPABASE_KEY"); key != "" {
		c.Database.SupabaseKey = key
	}
	if password := os.Getenv("SUPABASE_DB_PASSWORD"); password != "" {
		c.Database.SupabaseDBPassword = password
	}
	if projectRef := os.Getenv("SUPABASE_PROJECT_REF"); projectRef != "" {
		c.Database.SupabaseProjectRef = projectRef
	}

	// OpenAI overrides
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		c.OpenAI.APIKey = apiKey
	}
	if baseURL := os.Getenv("OPENAI_BASE_URL"); baseURL != "" {
		c.OpenAI.BaseURL = baseURL
	}
	if model := os.Getenv("OPENAI_MODEL"); model != "" {
		c.OpenAI.Model = model
	}

	// Dry run mode
	if dryRun := os.Getenv("DRY_RUN"); dryRun == "true" || dryRun == "1" {
		c.DryRun = true
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Database validation
	if c.Database.SupabaseDBPassword == "" || c.Database.SupabaseProjectRef == "" {
		return fmt.Errorf("database configuration incomplete: SUPABASE_DB_PASSWORD and SUPABASE_PROJECT_REF required")
	}

	// OpenAI validation
	if c.OpenAI.APIKey == "" {
		return fmt.Errorf("OpenAI API key is required")
	}
	if c.OpenAI.BaseURL == "" {
		c.OpenAI.BaseURL = "https://api.openai.com/v1" // Default
	}
	if c.OpenAI.Model == "" {
		c.OpenAI.Model = "gpt-4" // Default
	}

	// Platform validation (at least one platform must be configured)
	totalPlatforms := len(c.Platforms.Twitter) + len(c.Platforms.LinkedIn) +
		len(c.Platforms.Discord) + len(c.Platforms.Bluesky)

	if totalPlatforms == 0 {
		return fmt.Errorf("at least one platform must be configured")
	}

	// YouTube validation
	if len(c.YouTube.Channels) == 0 {
		return fmt.Errorf("at least one YouTube channel must be configured")
	}

	return nil
}

// GetTwitterCredentials retrieves Twitter credentials by handle
func (c *Config) GetTwitterCredentials(handle string) (*TwitterCredentials, error) {
	for _, cred := range c.Platforms.Twitter {
		if cred.Handle == handle {
			return &cred, nil
		}
	}
	return nil, fmt.Errorf("twitter credentials not found for handle: %s", handle)
}

// GetLinkedInCredentials retrieves LinkedIn credentials by handle
func (c *Config) GetLinkedInCredentials(handle string) (*LinkedInCredentials, error) {
	for _, cred := range c.Platforms.LinkedIn {
		if cred.Handle == handle {
			return &cred, nil
		}
	}
	return nil, fmt.Errorf("linkedin credentials not found for handle: %s", handle)
}

// GetDiscordCredentials retrieves Discord credentials by handle
func (c *Config) GetDiscordCredentials(handle string) (*DiscordCredentials, error) {
	for _, cred := range c.Platforms.Discord {
		if cred.Handle == handle {
			return &cred, nil
		}
	}
	return nil, fmt.Errorf("discord credentials not found for handle: %s", handle)
}

// GetBlueskyCredentials retrieves Bluesky credentials by handle
func (c *Config) GetBlueskyCredentials(handle string) (*BlueskyCredentials, error) {
	for _, cred := range c.Platforms.Bluesky {
		if cred.Handle == handle {
			return &cred, nil
		}
	}
	return nil, fmt.Errorf("bluesky credentials not found for handle: %s", handle)
}

// GetAllTwitterHandles returns all configured Twitter handles
func (c *Config) GetAllTwitterHandles() []string {
	handles := make([]string, len(c.Platforms.Twitter))
	for i, cred := range c.Platforms.Twitter {
		handles[i] = cred.Handle
	}
	return handles
}

// GetAllLinkedInHandles returns all configured LinkedIn handles
func (c *Config) GetAllLinkedInHandles() []string {
	handles := make([]string, len(c.Platforms.LinkedIn))
	for i, cred := range c.Platforms.LinkedIn {
		handles[i] = cred.Handle
	}
	return handles
}

// GetAllDiscordHandles returns all configured Discord handles
func (c *Config) GetAllDiscordHandles() []string {
	handles := make([]string, len(c.Platforms.Discord))
	for i, cred := range c.Platforms.Discord {
		handles[i] = cred.Handle
	}
	return handles
}

// GetAllBlueskyHandles returns all configured Bluesky handles
func (c *Config) GetAllBlueskyHandles() []string {
	handles := make([]string, len(c.Platforms.Bluesky))
	for i, cred := range c.Platforms.Bluesky {
		handles[i] = cred.Handle
	}
	return handles
}
