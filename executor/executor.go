package executor

import (
	"context"
	"fmt"
	"github.com/chain710/copilot/log"
	"github.com/chain710/copilot/plan"
	"github.com/sashabaranov/go-openai"
)

func NewExecutor(client *openai.Client, model string) Executor {
	return &executor{
		client: client,
		model:  model,
	}
}

type Executor interface {
	Do(ctx context.Context, p *plan.Plan, options ...Option) (*Result, error)
}

type executor struct {
	client *openai.Client
	model  string
}

func (e executor) Do(ctx context.Context, p *plan.Plan, options ...Option) (*Result, error) {
	var opt Options
	for _, f := range options {
		f(&opt)
	}
	provider, err := plan.NewPromptProvider(p.Messages, opt.data)
	if err != nil {
		return nil, err
	}

	var result Result
	for i, step := range p.Steps {
		var callbackData any
		lastStep := i == len(p.Steps)-1
		if opt.callback.Before != nil {
			callbackData = opt.callback.Before(&p.Steps[i], lastStep)
		}
		messages, err := step.ToMessage(provider)
		if err != nil {
			return nil, err
		}

		log.Debugf("executing step %s", step.Name)
		l := log.With("step", step.Name)
		for j := range messages {
			l.Debugw("prompt", "index", j, "message", messages)
		}
		resp, err := e.client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:       step.Model,
				Messages:    messages,
				Temperature: p.Temperature,
				MaxTokens:   p.MaxTokens,
				TopP:        p.TopP,
			},
		)

		if err != nil {
			return nil, fmt.Errorf("error executing step %s: %w", step.Name, err)
		}

		if len(resp.Choices) == 0 {
			return nil, fmt.Errorf("no choices returned for step %s", step.Name)
		}

		if step.ResultAs != nil {
			provider.Add(step.ResultAs.Name, &plan.Message{
				Name:    step.ResultAs.Name,
				Role:    step.ResultAs.Role,
				Content: resp.Choices[0].Message.Content,
			})
		}

		for _, choice := range resp.Choices {
			log.Debugf("choice[%d] finish reason: %s, message: %s", choice.Index, choice.FinishReason, choice.Message)
		}

		if lastStep {
			result.Content = resp.Choices[0].Message.Content
			result.Finished = true // TODO: if finish_reason = length, should collect all
		}

		if opt.callback.After != nil {
			opt.callback.After(&p.Steps[i], lastStep, callbackData, resp.Choices[0].Message.Content)
		}
	}
	return &result, nil
}
