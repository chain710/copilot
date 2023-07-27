package plan

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Option func(p *Plan) error

func FromFile(path string) Option {
	return func(p *Plan) error {
		f, err := os.Open(path)
		if err != nil {
			return err
		}

		defer f.Close()
		decoder := yaml.NewDecoder(f)
		if err := decoder.Decode(p); err != nil {
			return err
		}

		return nil
	}
}
