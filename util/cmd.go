package util

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func BindRequiredFlags(cmd *cobra.Command, names ...string) {
	nameSet := NewSet(names...)
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if !nameSet.Contains(f.Name) {
			return
		}
		_ = viper.BindPFlag(f.Name, f)
		_ = cmd.MarkFlagRequired(f.Name)
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			_ = cmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}
