package regolithe

import (
	"strings"

	"github.com/aporeto-inc/regolithe/spec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewCommand generates a new CLI for regolith
func NewCommand(
	name string,
	description string,
	nameConvertFunc spec.AttributeNameConverterFunc,
	typeConvertFunc spec.AttributeTypeConverterFunc,
	typeMappingName string,
	generatorFunc func(*spec.SpecificationSet) error,
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

	var cmdGen = &cobra.Command{
		Use:   "folder",
		Short: "Generate the model using a local directory containing the specs.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			specSet, err := spec.NewSpecificationSet(
				viper.GetString("dir"),
				nameConvertFunc,
				typeConvertFunc,
				typeMappingName,
			)
			if err != nil {
				return err
			}

			return generatorFunc(specSet)
		},
	}
	cmdGen.Flags().StringP("dir", "d", "", "Path of the specifications folder.")

	var githubGen = &cobra.Command{
		Use:   "github",
		Short: "Generate the model using a remote github repository.",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO
			return generatorFunc(nil)
		},
	}
	githubGen.Flags().StringP("api", "a", "https://api.gituhub.com", "Endpoint for the github api.")
	githubGen.Flags().StringP("repo", "r", "", "Repository where to pull the specification from.")
	githubGen.Flags().StringP("path", "p", "", "Internal path to a directory in the repo if not in the root.")
	githubGen.Flags().StringP("ref", "R", "master", "Branch or tag to use.")

	rootCmd.AddCommand(
		cmdGen,
		githubGen,
	)

	return rootCmd
}
