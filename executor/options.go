package executor

import "github.com/chain710/copilot/plan"

type BeforeCallbackFunc func(step *plan.Step, last bool) any
type AfterCallbackFunc func(step *plan.Step, last bool, beforeData any, content string)

type Callback struct {
	Before BeforeCallbackFunc
	After  AfterCallbackFunc
}

type Options struct {
	data     any
	callback Callback
}

type Option func(*Options)

func BeforeCallback(f BeforeCallbackFunc) Option {
	return func(options *Options) {
		options.callback.Before = f
	}
}

func AfterCallback(f AfterCallbackFunc) Option {
	return func(options *Options) {
		options.callback.After = f
	}
}

func PlanData(d any) Option {
	return func(options *Options) {
		options.data = d
	}
}
