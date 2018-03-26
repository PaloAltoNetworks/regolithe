package main

//go:generate go-bindata -pkg static -o static/bindata.go templates

import (
	"fmt"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
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
		Use:           "format",
		Short:         "Reads a specification from stdin and prints it formatted on std out.",
		SilenceErrors: true,
		SilenceUsage:  true,
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

	var docCmd = &cobra.Command{
		Use:           "doc",
		Short:         "Generate a documentation for the given specification set",
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			s, err := spec.NewSpecificationSet(
				viper.GetString("dir"),
				nil,
				nil,
				viper.GetString("category"),
			)
			if err != nil {
				return err
			}

			return writeDoc(s, viper.GetString("format"))
		},
	}

	docCmd.Flags().StringP("dir", "d", "", "Directory containing the specification")
	docCmd.Flags().String("format", "markdown", "Format of the documentation")
	docCmd.Flags().String("category", "", "Category of the type mapping to look for")

	rootCmd.AddCommand(
		formatCmd,
		docCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Fatal("Error")
		os.Exit(1)
	}
}
