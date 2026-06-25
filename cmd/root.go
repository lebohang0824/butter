package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"butter/pkg/formatter"
	"butter/pkg/lexer"
	"butter/pkg/parser"

	"github.com/spf13/cobra"
)

const Version = "1.3.0"

var outputFile string
var checkMode bool
var showVersion bool
var fmtCheckMode bool

var rootCmd = &cobra.Command{
	Use:   "butter",
	Short: "Butter is a high-performance, indentation-aware specification compiler.",
	Long:  `A clean compiler framework that turns minimalist indentation-based .butter specifications into beautifully formatted JSON structures.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Printf("butter v%s\n", Version)
			return nil
		}
		return cmd.Help()
	},
}

var compileCmd = &cobra.Command{
	Use:   "compile [input file]",
	Short: "Compile a .butter specification file down to pretty JSON",
	Long:  `Compile a .butter file to JSON. Use --check to validate syntax without generating output.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile := args[0]
		if !strings.HasSuffix(inputFile, ".butter") {
			return fmt.Errorf("input file must have a .butter extension")
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

		if checkMode {
			fmt.Println("OK")
			return nil
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

var fmtCmd = &cobra.Command{
	Use:   "fmt [input file]",
	Short: "Format a .butter specification file",
	Long:  `Format a .butter file according to standard conventions. Use --check to validate formatting without modifying.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile := args[0]
		if !strings.HasSuffix(inputFile, ".butter") {
			return fmt.Errorf("input file must have a .butter extension")
		}

		content, err := os.ReadFile(inputFile)
		if err != nil {
			return fmt.Errorf("failed to read source file: %w", err)
		}

		formatted, err := formatter.Format(content)
		if err != nil {
			return fmt.Errorf("formatting error: %w", err)
		}

		if fmtCheckMode {
			if string(content) != string(formatted) {
				return fmt.Errorf("file is not formatted")
			}
			fmt.Println("OK")
			return nil
		}

		if err := os.WriteFile(inputFile, formatted, 0644); err != nil {
			return fmt.Errorf("failed to write formatted file: %w", err)
		}

		fmt.Printf("Formatted: %s\n", inputFile)
		return nil
	},
}

func init() {
	rootCmd.Flags().BoolVar(&showVersion, "version", false, "Print the version number")
	compileCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Specify custom path for output file destination (defaults to input name + .json)")
	compileCmd.Flags().BoolVar(&checkMode, "check", false, "Check syntax without generating output")
	fmtCmd.Flags().BoolVar(&fmtCheckMode, "check", false, "Check formatting without modifying")
	rootCmd.AddCommand(compileCmd)
	rootCmd.AddCommand(fmtCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
