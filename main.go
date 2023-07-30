package main

import (
	"github.com/chain710/copilot/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"os"
)

var (
	configFile string
	logLevel   = zapLogLevel{value: zap.InfoLevel}
)

func initViper() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigName(".dev-agent")
		viper.SetConfigType("yaml")
	}
	viper.AutomaticEnv()
	err := viper.ReadInConfig()

	log.SetLogLevel(logLevel.value)
	if err != nil {
		log.Errorf("viper config error %s", err)
	} else {
		log.Debugf("viper config used: %s", viper.ConfigFileUsed())
	}
}

var rootCmd = &cobra.Command{
	Use:   "dev-agent",
	Short: "AI Agent for developer",
}

func Execute() int {
	cobra.OnInitialize(initViper)
	defer log.Sync()
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default $HOME/.dev-agent)")
	rootCmd.PersistentFlags().VarP(&logLevel, "log-level", "L", "log level")
	rootCmd.PersistentFlags().StringP(flagNameAzureOpenAIKey, "", "", "Azure OpenAI Key")
	rootCmd.PersistentFlags().StringP(flagNameAzureOpenAIEndpoint, "", "", "Azure OpenAI Endpoint")
	rootCmd.PersistentFlags().StringP(flagNameAzureOpenAIModel, "", "", "Azure OpenAI Model")
	if err := rootCmd.Execute(); err != nil {
		return -1
	}
	return 0
}

func main() {
	os.Exit(Execute())
}
