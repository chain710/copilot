package plan

import (
	"bytes"
	"text/template"
)

func newPromptProvider(tplMessages []Message, data any) (*promptProvider, error) {
	provider := promptProvider{
		messages: make(map[string]*Message),
	}
	for _, m := range tplMessages {
		tpl, err := template.New(m.Name).Option("missingkey=error").Parse(m.Content)
		if err != nil {
			return nil, err
		}
		var buf bytes.Buffer
		if err = tpl.Execute(&buf, data); err != nil {
			return nil, err
		}

		provider.Add(m.Name, &Message{
			Name:    m.Name,
			Role:    m.Role,
			Content: buf.String(),
		})
	}

	return &provider, nil
}

type promptProvider struct {
	messages map[string]*Message // name => message mapping
}

func (p *promptProvider) Get(name string) (*Message, bool) {
	value, ok := p.messages[name]
	return value, ok
}

func (p *promptProvider) Add(name string, value *Message) {
	p.messages[name] = value
}
