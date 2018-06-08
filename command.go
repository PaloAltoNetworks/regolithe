package regolithe

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/regolithe/spec"
)

// NewCommand generates a new CLI for regolith
func NewCommand(
	name string,
	description string,
	version string,
	nameConvertFunc spec.AttributeNameConverterFunc,
	typeConvertFunc spec.AttributeTypeConverterFunc,
	typeMappingName string,
	generatorFunc func([]spec.SpecificationSet, string) error,
) *cobra.Command {

	cobra.OnInitialize(func() {
		viper.SetEnvPrefix(name)
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	})

	var rootCmd = &cobra.Command{
		Use:   name,
		Short: description,
	}

	rootCmd.PersistentFlags().StringP("out", "o", "codegen", "Default output path.")

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the version and exit.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	}

	var cmdFolderGen = &cobra.Command{
		Use:           "folder",
		Short:         "Generate the model using a local directory containing the specs.",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(viper.GetStringSlice("dir")) == 0 {
				return errors.New("--dir is required")
			}

			var specSets []spec.SpecificationSet

			for _, dir := range viper.GetStringSlice("dir") {
				set, err := spec.LoadSpecificationSet(
					dir,
					nameConvertFunc,
					typeConvertFunc,
					typeMappingName,
				)
				if err != nil {
					return err
				}

				specSets = append(specSets, set)
			}

			return generatorFunc(specSets, viper.GetString("out"))
		},
	}
	cmdFolderGen.Flags().StringSliceP("dir", "d", nil, "Path of the specifications folder.")

	var githubGen = &cobra.Command{
		Use:           "github",
		Short:         "Generate the model using a remote github repository.",
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			specSet, err := spec.LoadSpecificationSetFromGithub(
				viper.GetString("token"),
				viper.GetString("repo"),
				viper.GetString("ref"),
				viper.GetString("path"),
				nameConvertFunc,
				typeConvertFunc,
				typeMappingName,
			)
			if err != nil {
				return err
			}

			return generatorFunc([]spec.SpecificationSet{specSet}, viper.GetString("out"))
		},
	}
	githubGen.Flags().StringP("repo", "r", "", "Endpoint for the github api.")
	githubGen.Flags().StringP("path", "p", "", "Internal path to a directory in the repo if not in the root.")
	githubGen.Flags().StringP("ref", "R", "master", "Branch or tag to use.")
	githubGen.Flags().StringP("token", "t", "", "The api token to use.")

	rootCmd.AddCommand(
		versionCmd,
		cmdFolderGen,
		githubGen,
	)

	return rootCmd
}
