package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// Database connection
	PgConnectionString string `envconfig:"PG_CONNECTION_STRING" required:"true"`

	// Background event handler configuration
	EventBatchInterval time.Duration `envconfig:"EVENT_BATCH_INTERVAL" default:"2s"`
	EventHandlerSleep  time.Duration `envconfig:"EVENT_HANDLER_SLEEP" default:"100ms"`

	// Statistics computation configuration
	StatisticsInterval time.Duration `envconfig:"STATISTICS_INTERVAL" default:"120s"`

	// Minimap generation configuration
	MinimapInitialInterval time.Duration `envconfig:"MINIMAP_INITIAL_INTERVAL" default:"1s"`
	MinimapIdleInterval    time.Duration `envconfig:"MINIMAP_IDLE_INTERVAL" default:"10m"`
	MinimapLockTimeout     time.Duration `envconfig:"MINIMAP_LOCK_TIMEOUT" default:"10m"`

	// Channel and buffer size configuration
	ButtonEventChannelSize int `envconfig:"BUTTON_EVENT_CHANNEL_SIZE" default:"2000"`
	MinimapChannelSize     int `envconfig:"MINIMAP_CHANNEL_SIZE" default:"10000"`
	EventBatchCapacity     int `envconfig:"EVENT_BATCH_CAPACITY" default:"1000"`
	SignalChannelSize      int `envconfig:"SIGNAL_CHANNEL_SIZE" default:"2"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
