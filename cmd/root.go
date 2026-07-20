package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	profile string
	region  string
	tag     string
)

var RootCmd = &cobra.Command{
	Use:   "aws-kill",
	Short: "AWS Kill Switch is a dependency-aware infrastructure cleanup tool",
	Long: `AWS Kill Switch automatically scans, plans, destroys, and verifies 
AWS infrastructure created during development. 

The commands must be executed in the exact sequence:
  scan -> plan -> kill -> verify`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.aws-kill.yaml)")
	RootCmd.PersistentFlags().StringVar(&profile, "profile", "", "AWS CLI profile to use")
	RootCmd.PersistentFlags().StringVar(&region, "region", "", "AWS region to target")
	RootCmd.PersistentFlags().StringVar(&tag, "tag", "", "Filter resources by tag (e.g. Environment=dev)")

	viper.BindPFlag("profile", RootCmd.PersistentFlags().Lookup("profile"))
	viper.BindPFlag("region", RootCmd.PersistentFlags().Lookup("region"))
	viper.BindPFlag("tag", RootCmd.PersistentFlags().Lookup("tag"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err == nil {
			viper.AddConfigPath(home)
			viper.SetConfigType("yaml")
			viper.SetConfigName(".aws-kill")
		}
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
