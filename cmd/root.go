package cmd

import (
	"os"

	"github.com/Vypolor/fixturec/pkg/fixturec"
	"github.com/spf13/cobra"
)

const (
	flagType = "type"
)

var rootCmd = &cobra.Command{
	Use:   "fixturec",
	Short: "Tool for generate test fixtures",
	RunE: func(cmd *cobra.Command, args []string) error {
		typeName, err := cmd.Flags().GetString(flagType)
		if err != nil {
			return err
		}

		cfg := fixturec.Config{
			TypeName: typeName,
		}

		return fixturec.GenerateFixture(cfg)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP(flagType, "t", "", "The type of the structure for which the fixture should be generated")
	rootCmd.Flags().BoolP("generate", "g", true, "Execute go:generate for mocks")
	rootCmd.Flags().BoolP("external", "e", false, "Generate mocks for external dependencies fields")
}
