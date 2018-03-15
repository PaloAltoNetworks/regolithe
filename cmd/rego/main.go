package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aporeto-inc/regolithe/spec"
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

	var formatCmd = &cobra.Command{
		Use:   "format",
		Short: "Reads a specification from stdin and prints it formatted on std out.",
		RunE: func(cmd *cobra.Command, args []string) error {

			s := spec.NewSpecification()

			if err := s.Read(os.Stdin); err != nil {
				return fmt.Errorf("Unable to load specs: %s", err)
			}

			if err := s.Write(os.Stdout); err != nil {
				return err
			}

			return nil
		},
	}

	rootCmd.AddCommand(
		formatCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
