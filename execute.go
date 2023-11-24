package main

import (
	"context"
	"fmt"
	"github.com/chain710/copilot/executor"
	"github.com/chain710/copilot/log"
	"github.com/chain710/copilot/plan"
	"github.com/chain710/copilot/util"
	"github.com/sashabaranov/go-openai"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"text/template"
	"time"
)

func mustReadStdin() string {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func mustReadFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(data)
}

var (
	planArguments     PlanArguments
	printEveryStep    bool
	templateFunctions = template.FuncMap{
		"ReadFile": mustReadFile,
		"Stdin":    mustReadStdin,
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

			bar := progressbar.NewOptions(-1,
				progressbar.OptionSetItsString("step"),
				progressbar.OptionSpinnerType(11),
				progressbar.OptionClearOnFinish())
			defer bar.Close()
			bgIncrementProgress(cmd.Context(), bar)

			exec := executor.NewExecutor(client, model)
			totalSteps := len(testingPlan.Steps)
			currentStep := 0
			_, err = exec.Do(cmd.Context(), testingPlan,
				executor.PlanData(templateArgs),
				executor.BeforeCallback(func(step *plan.Step, last bool) any {
					currentStep++
					bar.Describe(fmt.Sprintf("%d/%d %s", currentStep, totalSteps, step.Name))
					return nil
				}),
				executor.AfterCallback(func(step *plan.Step, last bool, beforeData any, content string) {
					if last || printEveryStep {
						_ = bar.Clear()
						cmd.Printf("Step[%s] Result:\n", step.Name)
						cmd.Println(content)
					}
				}))
			if err != nil {
				log.Errorf("error executing plan: %v", err)
				return err
			}

			return nil
		},
	}
)

func bgIncrementProgress(ctx context.Context, bar *progressbar.ProgressBar) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				_ = bar.Add(1)
				time.Sleep(time.Millisecond * 200)
			}
		}
	}()
}

func init() {
	executeCmd.Flags().VarP(&planArguments, "plan-args", "p", "plan arguments. example: 'Content=ReadFile(`source.go`) Function=`GetPriority`'")
	executeCmd.Flags().BoolVar(&printEveryStep, "print-every-step", false, "print every step's result")
	rootCmd.AddCommand(executeCmd)
}
