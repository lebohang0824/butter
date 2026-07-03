package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"butter/pkg/formatter"
	"butter/pkg/lexer"
	"butter/pkg/output"
	_ "butter/pkg/output/htmltree"
	_ "butter/pkg/output/json"
	_ "butter/pkg/output/yaml"
	"butter/pkg/parser"
	"butter/pkg/semantic"

	"github.com/spf13/cobra"
)

func formatList() string {
	return strings.Join(output.Names(), ", ")
}

const Version = "1.10.0"

var outputFile string
var outputFormat string
var checkMode bool
var showVersion bool
var fmtCheckMode bool

var rootCmd = &cobra.Command{
	Use:   "butter",
	Short: "Butter is a high-performance, indentation-aware specification compiler.",
	Long:  `A clean compiler framework that turns minimalist indentation-based .butter specifications into beautifully formatted JSON or YAML structures.`,
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
	Short: "Compile a .butter specification file to JSON (default) or YAML",
	Long:  `Compile a .butter file to JSON or YAML. Use --check to validate syntax without generating output.`,
	Args:  cobra.ExactArgs(1),
	SilenceUsage: true,
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

		diags := semantic.Analyze(appAST)
		hasErrors := false
		for _, d := range diags {
			fmt.Fprintf(os.Stderr, "%s\n", d)
			if d.Severity == semantic.SemError {
				hasErrors = true
			}
		}

		if checkMode {
			if hasErrors {
				return fmt.Errorf("semantic analysis found errors")
			}
			fmt.Println("OK")
			return nil
		}

		if hasErrors {
			return fmt.Errorf("semantic analysis failed — output not generated")
		}

		extIface, ok := output.Get(outputFormat)
		if !ok {
			return fmt.Errorf("unsupported output format %q — supported: %s", outputFormat, formatList())
		}
		outputData, err := extIface.Serialize(appAST)
		if err != nil {
			return fmt.Errorf("%s serialization failed: %w", outputFormat, err)
		}
		ext := extIface.FileExtension()

		if outputFile == "" {
			outputFile = strings.TrimSuffix(inputFile, filepath.Ext(inputFile)) + ext
		}

		if err := os.WriteFile(outputFile, outputData, 0644); err != nil {
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
	compileCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Specify custom path for output file destination (defaults to input name + .json for json, .yaml for yaml)")
	compileCmd.Flags().StringVarP(&outputFormat, "format", "f", "json", fmt.Sprintf("Output format (default: json). Supported: %s", formatList()))
	compileCmd.Flags().BoolVar(&checkMode, "check", false, "Check syntax without generating output")
	fmtCmd.Flags().BoolVar(&fmtCheckMode, "check", false, "Check formatting without modifying")
	rootCmd.AddCommand(compileCmd)
	rootCmd.AddCommand(fmtCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
