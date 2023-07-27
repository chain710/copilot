package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	configFile string
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
	_ = viper.ReadInConfig()

	//logLevel := zapcore.ErrorLevel
	//_ = logLevel.Set(viper.GetString("logLevel"))
	//bld := log.DefaultBuilder()
	//bld.Config.Level = logLevel
	//log.InitLog(bld)
	//
	//if err != nil {
	//	log.Errorf("viper config error %s", err)
	//} else {
	//	log.Debugf("viper config used: %s", viper.ConfigFileUsed())
	//}
}

var rootCmd = &cobra.Command{
	Use:   "dev-agent",
	Short: "AI Agent for developer",
}

func Execute() int {
	cobra.OnInitialize(initViper)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default $HOME/.dev-agent)")
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
