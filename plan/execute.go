package plan

import (
	"context"
	"errors"
	"github.com/sashabaranov/go-openai"
)

func Execute(ctx context.Context, p *Plan, data any, model string, client *openai.Client) (*Result, error) {
	return nil, errors.New("TODO")
	//assistant := openai.ChatCompletionMessage{
	//	Role: openai.ChatMessageRoleAssistant,
	//}
	//for _, step := range p.Steps {
	//	messages, err := step.ToMessage(data)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	fmt.Printf("executing step %s\nmessages: %+v\n", step.Name, messages)
	//	resp, err := client.CreateChatCompletion(
	//		ctx,
	//		openai.ChatCompletionRequest{
	//			Model:       model,
	//			Messages:    messages,
	//			Temperature: p.Temperature,
	//			MaxTokens:   p.MaxTokens,
	//			TopP:        p.TopP,
	//		},
	//	)
	//
	//	if err != nil {
	//		return nil, fmt.Errorf("error executing step %s: %w", step.Name, err)
	//	}
	//
	//	if len(resp.Choices) == 0 {
	//		return nil, fmt.Errorf("no choices returned for step %s", step.Name)
	//	}
	//
	//	assistant.Content = resp.Choices[0].Message.Content
	//
	//	for _, choice := range resp.Choices {
	//		fmt.Println(choice.Message)
	//	}
	//}
	//return nil, nil
}
