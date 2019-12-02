package cmd

import (
	"fmt"
	"os"

	"github.com/apex/log"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/springload/ssm-parent/ssm/transformations"
)

var (
	config              string
	debug               bool
	paths               []string
	names               []string
	plainPaths          []string
	plainNames          []string
	recursive           bool
	strict              bool
	expand              bool
	transformationsList []transformations.Transformation
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ssm-parent",
	Short: "Docker entrypoint that get parameters from AWS SSM Parameter Store",
	Long: `SSM-Parent is a docker entrypoint.
	
It gets specified parameters (possibly secret) from AWS SSM Parameter Store,
then exports them to the underlying process. Or creates a .env file to be consumed by an application.

It reads parameters in the following order: path->name->plain-path->plain-name.
So that every rightmost parameter overrides the previous one.
`,
}

// Execute is the entrypoint for cmd/ module
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initSettings() {
	if config != "" {
		viper.SetConfigFile(config)
		if err := viper.ReadInConfig(); err == nil {
			log.Infof("Using config file: %s", viper.ConfigFileUsed())
		} else {
			log.WithError(err).Fatal("Had some errors while parsing the config")
		}
	}
	// parse to an array first
	var transformationsInterfaceArray = []interface{}{}
	if err := viper.UnmarshalKey("transformations", &transformationsInterfaceArray); err != nil {
		log.WithError(err).Fatal("can't decode config")
	}
	// unmarshal to the tiny struct first to see what the action is
	for n, t := range transformationsInterfaceArray {
		var hint = struct{ Action string }{}

		if err := mapstructure.Decode(t, &hint); err != nil {
			log.WithError(err).Fatal("can't decode config")
		}
		switch hint.Action {

		case "delete":
			tr := new(transformations.DeleteTransformation)
			if err := mapstructure.Decode(t, tr); err != nil {
				log.WithFields(log.Fields{
					"transformation_number": n,
					"transformation_action": hint.Action,
				}).WithError(err).Fatal("can't decode config")
			}
			transformationsList = append(transformationsList, tr)
		case "rename":
			tr := new(transformations.RenameTransformation)
			if err := mapstructure.Decode(t, tr); err != nil {
				log.WithFields(log.Fields{
					"transformation_number": n,
					"transformation_action": hint.Action,
				}).WithError(err).Fatal("can't decode config")
			}
			transformationsList = append(transformationsList, tr)
		case "template":
			tr := new(transformations.TemplateTransformation)
			if err := mapstructure.Decode(t, tr); err != nil {
				log.WithFields(log.Fields{
					"transformation_number": n,
					"transformation_action": hint.Action,
				}).WithError(err).Fatal("can't decode config")
			}
			transformationsList = append(transformationsList, tr)

		default:
			log.Warnf("Got unparsed action: %s", hint.Action)
		}
	}

	//	}
	if viper.GetBool("debug") {
		log.SetLevel(log.DebugLevel)
	}
}
func init() {
	cobra.OnInitialize(initSettings)
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "", "Path to the config file (optional). Allows to set transformations")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Turn on debug logging")
	rootCmd.PersistentFlags().BoolVarP(&expand, "expand", "e", false, "Expand arguments and values using shell-style syntax")
	rootCmd.PersistentFlags().StringArrayVarP(&paths, "path", "p", []string{}, "Path to a SSM parameter. Expects JSON in the value. Can be specified multiple times.")
	rootCmd.PersistentFlags().StringArrayVarP(&names, "name", "n", []string{}, "Name of the SSM parameter to retrieve. Expects JSON in the value. Can be specified multiple times.")
	rootCmd.PersistentFlags().StringArrayVarP(&plainPaths, "plain-path", "", []string{}, "Path to a SSM parameter. Expects actual parameter in the value. Can be specified multiple times.")
	rootCmd.PersistentFlags().StringArrayVarP(&plainNames, "plain-name", "", []string{}, "Name of the SSM parameter to retrieve. Expects actual parameter in the value. Can be specified multiple times.")
	rootCmd.PersistentFlags().BoolVarP(&recursive, "recursive", "r", false, "Walk through the provided SSM paths recursively.")
	rootCmd.PersistentFlags().BoolVarP(&strict, "strict", "s", false, "Strict mode. Fail if found less parameters than number of names.")

	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("expand", rootCmd.PersistentFlags().Lookup("expand"))
	viper.BindPFlag("path", rootCmd.PersistentFlags().Lookup("path"))
	viper.BindPFlag("name", rootCmd.PersistentFlags().Lookup("name"))
	viper.BindPFlag("plain-path", rootCmd.PersistentFlags().Lookup("plain-path"))
	viper.BindPFlag("plain-name", rootCmd.PersistentFlags().Lookup("plain-name"))
	viper.BindPFlag("recursive", rootCmd.PersistentFlags().Lookup("recursive"))
	viper.BindPFlag("strict", rootCmd.PersistentFlags().Lookup("strict"))
}
