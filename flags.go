package main

import (
	"github.com/chain710/dev_agent/args"
	"go.uber.org/zap/zapcore"
)

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

type PlanArguments struct {
	AssignList args.AssignStatementList
}

func (p *PlanArguments) String() string {
	return p.AssignList.String()
}

func (p *PlanArguments) Set(s string) error {
	list, err := args.ParseAssignStatementList(s)
	if err != nil {
		return err
	}

	p.AssignList = list
	return nil
}

func (p *PlanArguments) Type() string {
	return "string"
}
