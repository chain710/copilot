package main

import (
	"github.com/chain710/copilot/log"
	"github.com/chain710/copilot/plan"
	"github.com/chain710/copilot/util"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"text/template"
)

func mustReadFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(data)
}

var (
	planArguments     PlanArguments
	templateFunctions = template.FuncMap{
		"ReadFile": mustReadFile,
	}

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
			templateArgs, err := planArguments.AssignList.Evaluate(templateFunctions)
			if err != nil {
				log.Errorf("error evaluating template arguments: %v", err)
				return err
			}

			result, err := plan.Execute(cmd.Context(), testingPlan, templateArgs, model, client)
			if err != nil {
				log.Errorf("error executing plan: %v", err)
				return err
			}

			cmd.Println("Result:")
			cmd.Println(result.Content)
			return nil
		},
	}
)

func init() {
	executeCmd.Flags().VarP(&planArguments, "plan-args", "p", "plan arguments. example: 'Content=ReadFile(`source.go`) Function=`GetPriority`'")
	rootCmd.AddCommand(executeCmd)
}
