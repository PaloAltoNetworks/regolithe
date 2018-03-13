package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	name        = "rego"
	description = "Tool to manipulate regolithe specifications"
)

func main() {

	cobra.OnInitialize(func() {
		viper.SetEnvPrefix(name)
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	})

	var rootCmd = &cobra.Command{
		Use:   name,
		Short: description,
	}

	var beautifyCmd = &cobra.Command{
		Use:   "beautify",
		Short: "Beautify a set of specifications.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return beautify(
				viper.GetString("src"),
				viper.GetString("dst"),
			)
		},
	}
	beautifyCmd.Flags().StringP("src", "s", "", "Source directory containing the specs")
	beautifyCmd.Flags().StringP("dst", "d", "", "Destination directory where to write the specs")

	rootCmd.AddCommand(
		beautifyCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
