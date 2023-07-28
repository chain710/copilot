package plan

import (
	"context"
	"fmt"
	"github.com/chain710/dev_agent/log"

	"github.com/sashabaranov/go-openai"
)

func Execute(ctx context.Context, p *Plan, data any, model string, client *openai.Client) (*Result, error) {
	provider, err := newPromptProvider(p.Messages, data)
	if err != nil {
		return nil, err
	}

	var result Result
	for i, step := range p.Steps {
		lastStep := i == len(p.Steps)-1
		messages, err := step.ToMessage(provider)
		if err != nil {
			return nil, err
		}

		log.Infof("executing step %s", step.Name)
		l := log.With("step", step.Name)
		for j := range messages {
			l.Debugw("prompt", "index", j, "message", messages)
		}
		resp, err := client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:       model,
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
			provider.Add(step.ResultAs.Name, &Message{
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
			result.Finished = true // TODO
		}
	}
	return &result, nil
}
