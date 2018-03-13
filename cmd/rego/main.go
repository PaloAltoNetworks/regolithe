package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/aporeto-inc/regolithe/spec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {

	const name = "rego"
	const description = "Tool to manipulate regolithe specifications"

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

func beautify(srcdir, dstdir string) error {

	if srcdir == "" {
		return fmt.Errorf("You must provide --src")
	}

	if dstdir == "" {
		dstdir = srcdir
	}

	if err := os.MkdirAll(dstdir, 0755); err != nil {
		return fmt.Errorf("Unable to create dest dir: %s", err)
	}

	filesInfo, err := ioutil.ReadDir(srcdir)
	if err != nil {
		return fmt.Errorf("Unable to read dir: %s", err)
	}

	for _, info := range filesInfo {

		if path.Ext(info.Name()) != ".spec" {
			continue
		}

		s, err := spec.LoadSpecification(path.Join(srcdir, info.Name()))
		if err != nil {
			return fmt.Errorf("Unable to read file '%s': %s", info.Name(), err)
		}

		if err = s.Write(dstdir); err != nil {
			return fmt.Errorf("Unable to write file '%s': %s", info.Name(), err)
		}
	}

	return nil
}
