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

			srcdir := viper.GetString("src")
			dstdir := viper.GetString("dst")
			file := viper.GetString("file")

			if file != "" && (srcdir != "" || dstdir != "") {
				return fmt.Errorf("If you pass --file, you cannot pass anything else")
			}

			if file != "" || (file == "" && srcdir == "" && dstdir == "") {
				return beautifyOne(file)
			}

			if srcdir == "" {
				return fmt.Errorf("You must provide --src")
			}

			if dstdir == "" {
				dstdir = srcdir
			}

			if err := os.MkdirAll(dstdir, 0755); err != nil {
				return fmt.Errorf("Unable to create dest dir: %s", err)
			}

			return beautifyAll(srcdir, dstdir)
		},
	}
	beautifyCmd.Flags().StringP("src", "s", "", "Source directory containing the specs")
	beautifyCmd.Flags().StringP("dst", "d", "", "Destination directory where to write the specs")
	beautifyCmd.Flags().StringP("file", "f", "", "File to beautify. The result will be printed on stdout")

	rootCmd.AddCommand(
		beautifyCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
