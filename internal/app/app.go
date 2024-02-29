package app

import (
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

type App struct {
	cfg *Config
}

func (a *App) Run() {

}

func (a *App) init(configPath string) error {
	cfg, err := newConfig(configPath)
	if err != nil {
		return fmt.Errorf("cannot create config: %w", err)
	}

	a.cfg = cfg

	return nil
}

func newConfig(configPath string) (*Config, error) {
	cfg := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open config file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("cannot decode config file: %w", err)
	}

	return cfg, nil
}

func NewApp(configPath string) (*App, error) {
	app := &App{}

	err := app.init(configPath)
	if err != nil {
		return nil, fmt.Errorf("cannot init app: %w", err)
	}

	return app, nil
}
