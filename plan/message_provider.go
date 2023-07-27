package plan

import (
	"bytes"
	"text/template"
)

func newMessageProvider(tplMessages []Message, data any) (*messageProvider, error) {
	provider := messageProvider{
		messages: make(map[string]*Message),
	}
	for _, m := range tplMessages {
		tpl, err := template.New(m.Name).Parse(m.Content)
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

type messageProvider struct {
	messages map[string]*Message // name => message mapping
}

func (p *messageProvider) Get(name string) (*Message, bool) {
	value, ok := p.messages[name]
	return value, ok
}

func (p *messageProvider) Add(name string, value *Message) {
	p.messages[name] = value
}
