package main

//go:generate go-bindata -pkg static -o static/bindata.go templates specset/...

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/aporeto-inc/regolithe/cmd/rego/doc"
	"github.com/aporeto-inc/regolithe/cmd/rego/static"
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

			if err := s.Read(os.Stdin, true); err != nil {
				return err
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

			s, err := spec.LoadSpecificationSet(
				viper.GetString("dir"),
				nil,
				nil,
				viper.GetString("category"),
			)
			if err != nil {
				return err
			}

			return doc.Write(s, viper.GetString("format"))
		},
	}

	var initCmd = &cobra.Command{
		Use:           "init <dest>",
		Short:         "Generate a new set of specification",
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) != 1 {
				return fmt.Errorf("usage: init <dest>")
			}

			dir := args[0]
			if err := os.MkdirAll(path.Base(dir), 0744); err != nil {
				return err
			}

			tmp, err := ioutil.TempDir(os.TempDir(), "rego")
			if err != nil {
				return err
			}

			if err := static.RestoreAssets(tmp, "specset"); err != nil {
				return err
			}

			return os.Rename(path.Join(tmp, "specset"), dir)
		},
	}

	rootCmd.AddCommand(
		formatCmd,
		docCmd,
		initCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
