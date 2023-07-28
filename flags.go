package main

import "go.uber.org/zap/zapcore"

const (
	flagNameAzureOpenAIKey      = "azure-openai-key"
	flagNameAzureOpenAIEndpoint = "azure-openai-endpoint"
	flagNameAzureOpenAIModel    = "azure-openai-model"
)

type zapLogLevel struct {
	value zapcore.Level
}

func (l *zapLogLevel) String() string {
	return l.value.String()
}

func (l *zapLogLevel) Set(s string) error {
	level, err := zapcore.ParseLevel(s)
	if err != nil {
		return err
	}

	l.value = level
	return nil
}

func (l *zapLogLevel) Type() string {
	return "string"
}
