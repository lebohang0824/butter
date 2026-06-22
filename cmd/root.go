package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"butter/pkg/lexer"
	"butter/pkg/parser"

	"github.com/spf13/cobra"
)

var outputFile string

var rootCmd = &cobra.Command{
	Use:   "butter",
	Short: "Butter is a high-performance, indentation-aware specification compiler.",
	Long:  `A clean compiler framework that turns minimalist indentation-based .butter specifications into beautifully formatted JSON structures.`,
}

var compileCmd = &cobra.Command{
	Use:   "compile [input file]",
	Short: "Compile a .butter specification file down to pretty JSON",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile := args[0]
		if !strings.HasSuffix(inputFile, ".butter") {
			return fmt.Errorf("invalid file context: source files must end explicitly with the '.butter' extension")
		}

		content, err := os.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read source file: %w", err)
		}

		l := lexer.NewLexer(string(content))
		p := parser.NewParser(l)
		appAST, err := p.Parse()
		if err != nil {
			return fmt.Errorf("compilation syntax compilation error:\n%w", err)
		}

		jsonOutput, err := parser.GenerateJSONSpec(appAST)
		if err != nil {
			return fmt.Errorf("json packaging generation failed: %w", err)
		}

		if outputFile == "" {
			outputFile = strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + ".json"
		}

		if err := os.WriteFile(outputFile, jsonOutput, 0644); err != nil {
			return fmt.Errorf("failed to write compiled asset to target destination disk: %w", err)
		}

		fmt.Printf("Successfully compiled: %s ──> %s\n", inputFile, outputFile)
		return nil
	},
}

func init() {
	compileCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Specify custom path for output file destination (defaults to input name + .json)")
	rootCmd.AddCommand(compileCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
