package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/Praqma/Praqmatic-Automated-Changelog/src/model"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Generate configuration file",
	Long:  `Generate a PAC configuration file in YAML or JSON format with current settings`,
	Run: func(cmd *cobra.Command, args []string) {
		// Build settings from flags using the shared builder
		builder := NewSettingsBuilder(cmd)
		settings := builder.BuildFromFlags()
		
		// Output in requested format
		switch OutputFormat {
		case "json":
			outputJSON(settings)
		case "yaml", "yml":
			outputYAML(settings)
		default:
			fmt.Fprintf(os.Stderr, "Invalid output format: %s. Use 'yaml' or 'json'\n", OutputFormat)
			os.Exit(1)
		}
	},
}

func outputJSON(settings *model.Settings) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(settings); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}
}

func outputYAML(settings *model.Settings) {
	encoder := yaml.NewEncoder(os.Stdout)
	if err := encoder.Encode(settings); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding YAML: %v\n", err)
		os.Exit(1)
	}
}
