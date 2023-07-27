package plan

import (
	"bytes"
	"fmt"
	"github.com/chain710/dev_agent/util"
	"github.com/go-playground/validator/v10"
	"github.com/sashabaranov/go-openai"
	"text/template"
)

type Plan struct {
	Steps       []Step    `yaml:"steps" validate:"min=1"`
	Messages    []Message `yaml:"messages"`
	MaxTokens   int       `yaml:"max_tokens,omitempty" validate:"min=1"`
	Temperature float32   `yaml:"temperature,omitempty" validate:"min=0,max=1"`
	TopP        float32   `yaml:"top_p,omitempty" validate:"min=0,max=1"`
}

type ResultAs struct {
	Name string `json:"name" validate:"required"`
	Role string `json:"role" validate:"required"`
}

type Step struct {
	Name     string       `yaml:"name"`
	Messages []MessageRef `yaml:"messages"`
	ResultAs *ResultAs    `yaml:"result_as,omitempty"`
}

type MessageRef struct {
	Name string `yaml:"name"`
}

type Message struct {
	Name    string `yaml:"name"`
	Role    string `yaml:"role"`
	Content string `yaml:"content"`
}

var validate = validator.New()

func New(options ...Option) (*Plan, error) {
	var p Plan
	for _, option := range options {
		if err := option(&p); err != nil {
			return nil, err
		}
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}
	return &p, nil
}

func (p *Plan) Validate() error {
	if err := validate.Struct(p); err != nil {
		return err
	}

	availableMessages := util.NewSet[string]()
	for _, message := range p.Messages {
		availableMessages.Add(message.Name)
	}
	for _, step := range p.Steps {
		for _, message := range step.Messages {
			if !availableMessages.Contains(message.Name) {
				return fmt.Errorf("step `%s` message `%s` not found", step.Name, message.Name)
			}
		}

		if step.ResultAs != nil {
			if availableMessages.Contains(step.ResultAs.Name) {
				return fmt.Errorf("step `%s` result_as conflict with message `%s`", step.Name, step.ResultAs.Name)
			}
			availableMessages.Add(step.ResultAs.Name)
		}
	}
	return nil
}

func (s *Step) ToMessage(arg any) (messages []openai.ChatCompletionMessage, err error) {
	tpl := template.New(s.Name)
	for _, message := range s.Messages {
		_ = message.Name // TODO
		_, err = tpl.Parse( /*message.Content*/ "")
		if err != nil {
			return nil, err
		}

		var buf bytes.Buffer
		err = tpl.Execute(&buf, arg)
		if err != nil {
			return nil, err
		}
		messages = append(messages, openai.ChatCompletionMessage{
			// Role:    message.Role,
			Content: buf.String(),
		})
	}
	return
}