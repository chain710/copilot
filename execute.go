package main

import (
	"github.com/chain710/dev_agent/log"
	"github.com/chain710/dev_agent/plan"
	"github.com/chain710/dev_agent/util"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// executeCmd executes a multi step prompt plan
	executeCmd = &cobra.Command{
		Use:   "execute <plan>",
		Short: "Execute a multi step prompt plan",
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			util.BindRequiredFlags(cmd,
				flagNameAzureOpenAIKey,
				flagNameAzureOpenAIEndpoint,
				flagNameAzureOpenAIModel)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			planPath := args[0]
			testingPlan, err := plan.New(plan.FromFile(planPath))
			if err != nil {
				return err
			}
			key := viper.GetString(flagNameAzureOpenAIKey)
			endpoint := viper.GetString(flagNameAzureOpenAIEndpoint)
			model := viper.GetString(flagNameAzureOpenAIModel)
			config := openai.DefaultAzureConfig(key, endpoint)
			client := openai.NewClientWithConfig(config)
			var data struct {
				CodeContent string
			}

			data.CodeContent = `
func Execute(ctx context.Context, p *Plan, data any, model string, client *openai.Client) (*Result, error) {
	for _, step := range p.Steps {
		fmt.Printf("executing step %s\n", step.Name)
		messages, err := step.ToMessage(data)
		if err != nil {
			return nil, err
		}
		resp, err := client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:    model,
				Messages: messages,
			},
		)

		if err != nil {
			return nil, fmt.Errorf("error executing step %s: %w", step.Name, err)
		}

		for _, choice := range resp.Choices {
			fmt.Println(choice.Message)
		}
	}

	return nil, nil
}
`
			result, err := plan.Execute(cmd.Context(), testingPlan, data, model, client)
			if err != nil {
				log.Errorf("error executing plan: %v", err)
				return err
			}

			cmd.Println("Result:\n")
			cmd.Println(result.Content)
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(executeCmd)
}
