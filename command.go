package regolithe

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/aporeto-inc/regolithe/spec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	git "gopkg.in/src-d/go-git.v4"
)

// NewCommand generates a new CLI for regolith
func NewCommand(
	name string,
	description string,
	version string,
	nameConvertFunc spec.AttributeNameConverterFunc,
	typeConvertFunc spec.AttributeTypeConverterFunc,
	typeMappingName string,
	generatorFunc func([]*spec.SpecificationSet, string) error,
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

			var specSets []*spec.SpecificationSet

			for _, dir := range viper.GetStringSlice("dir") {
				set, err := spec.NewSpecificationSet(
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

			var auth transport.AuthMethod
			if viper.GetString("token") != "" {
				auth = &http.BasicAuth{
					Username: "Bearer",
					Password: viper.GetString("token"),
				}
			}

			tmpFolder, err := ioutil.TempDir("", "regolithe-refs-head")
			if err != nil {
				return err
			}
			defer func(f string) { os.RemoveAll(f) }(tmpFolder) // nolint: errcheck

			var (
				ref           plumbing.ReferenceName
				needsCheckout bool
			)

			givenHash := plumbing.NewHash(viper.GetString("ref"))
			if !givenHash.IsZero() {
				ref = plumbing.NewReferenceFromStrings("refs/heads/master", "").Name()
				needsCheckout = true
			} else {
				ref = plumbing.NewReferenceFromStrings("refs/heads/"+viper.GetString("ref"), "").Name()
			}

			logrus.WithFields(logrus.Fields{
				"ref":  viper.GetString("ref"),
				"repo": viper.GetString("repo"),
				"path": viper.GetString("path"),
			}).Info("Retrieving repository")

			cloneFunc := func(folder string, ref plumbing.ReferenceName) (*git.Repository, error) {
				return git.PlainClone(
					folder,
					false,
					&git.CloneOptions{
						URL:           viper.GetString("repo"),
						Progress:      nil,
						ReferenceName: ref,
						Auth:          auth,
					})
			}

			repo, err := cloneFunc(tmpFolder, ref)

			if err != nil {
				if err == plumbing.ErrReferenceNotFound {
					logrus.WithFields(logrus.Fields{
						"err":  err,
						"ref":  viper.GetString("ref"),
						"repo": viper.GetString("repo"),
						"path": viper.GetString("path"),
					}).Warn("Trying to clone with refs/tags - failed to clone with refs/heads")

					// Need to recreate a folder, get error repository already created otherwise
					// Happened even if old tmp folder is deleted...
					tmpFolder, err = ioutil.TempDir("", "regolithe-refs-tags")
					if err != nil {
						return err
					}
					defer func(f string) { os.RemoveAll(f) }(tmpFolder) // nolint: errcheck

					ref = plumbing.NewReferenceFromStrings("refs/tags/"+viper.GetString("ref"), "").Name()
					repo, err = cloneFunc(tmpFolder, ref)

					if err != nil {
						return err
					}
				} else {
					return err
				}
			}

			if needsCheckout {
				wt, e := repo.Worktree()
				if e != nil {
					return e
				}

				if err = wt.Checkout(
					&git.CheckoutOptions{
						Hash: givenHash,
					}); err != nil {
					return err
				}
			}

			specSet, err := spec.NewSpecificationSet(
				path.Join(tmpFolder, viper.GetString("path")),
				nameConvertFunc,
				typeConvertFunc,
				typeMappingName,
			)
			if err != nil {
				return err
			}

			return generatorFunc([]*spec.SpecificationSet{specSet}, viper.GetString("out"))
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
