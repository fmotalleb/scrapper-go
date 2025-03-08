package config

import "github.com/playwright-community/playwright-go"

type ExecutionConfig struct {
	Pipeline Pipeline `mapstructure:"pipeline"`
}

type Pipeline struct {
	KeepRunning   string                              `mapstructure:"keep_running"`
	Browser       string                              `mapstructure:"browser"`
	BrowserParams playwright.BrowserTypeLaunchOptions `mapstructure:"browser_params"`
	Vars          []Variable                          `mapstructure:"vars"`
	Steps         []Step                              `mapstructure:"steps"`
}

type Variable struct {
	Name         string `mapstructure:"name"`
	Value        string `mapstructure:"value"`
	Random       string `mapstructure:"random"`
	RandomChars  string `mapstructure:"random_chars"`
	RandomLength int    `mapstructure:"random_length"`
	Postfix      string `mapstructure:"postfix,omitempty"`
	Prefix       string `mapstructure:"prefix,omitempty"`
}

type Step map[string]interface{}
