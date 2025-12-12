package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kreatoo/md2html/pkg/converter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	// Global flags
	verbose bool
	debug   bool

	// Convert command flags
	inputDir            string
	outputDir           string
	convertTemplateFile string
	cssPath             string
	jsPath              string
	siteTitle           string
	author              string
	addList             bool
	recursive           bool
	generateRSS         bool
	siteURL             string
	serve               bool
	port                string
	gitWebURL           string
	showCommitInfo      bool

	// Create-template command flags
	createTemplateFile string

	// Serve command flags
	serveOutputDir string
	servePort      string

	// Minify command flags
	minifyInput  string
	minifyOutput string
	minifyType   string

	// Root command
	rootCmd = &cobra.Command{
		Use:   "md2html",
		Short: "Convert Markdown files to HTML",
		Long:  "A tool to convert Markdown files to HTML with template support, RSS generation, and development server",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Setup logging (verbose by default, can be made less verbose with --quiet)
			if debug {
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			} else if verbose {
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			} else {
				// Default to InfoLevel (verbose)
				zerolog.SetGlobalLevel(zerolog.InfoLevel)
			}
		},
	}

	// Convert command
	convertCmd = &cobra.Command{
		Use:   "convert",
		Short: "Convert markdown files to HTML",
		Long:  "Convert markdown files in a directory to HTML using specified templates and options",
		RunE:  runConvert,
	}

	// Create-template command
	createTemplateCmd = &cobra.Command{
		Use:   "create-template",
		Short: "Create a default HTML template",
		Long:  "Generate a default HTML template file that can be customized",
		RunE:  runCreateTemplate,
	}

	// Serve command
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start a development server",
		Long:  "Start an HTTP development server to serve the converted HTML files",
		RunE:  runServe,
	}

	// Minify command
	minifyCmd = &cobra.Command{
		Use:   "minify",
		Short: "Minify CSS and JS files",
		Long:  "Minify CSS and/or JavaScript files or directories to reduce file size",
		RunE:  runMinify,
	}
)

func init() {
	// Add subcommands to root
	rootCmd.AddCommand(convertCmd)
	rootCmd.AddCommand(createTemplateCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(minifyCmd)

	// Global flags
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")

	// Convert command flags
	convertCmd.Flags().StringVarP(&inputDir, "input", "i", "", "Directory containing markdown files (required)")
	convertCmd.Flags().StringVarP(&outputDir, "output", "o", "output", "Directory to save HTML files")
	convertCmd.Flags().StringVarP(&convertTemplateFile, "template", "t", "", "HTML template file (optional)")
	convertCmd.Flags().StringVar(&cssPath, "css", "", "Path to CSS file")
	convertCmd.Flags().StringVar(&jsPath, "js", "", "Path to JavaScript file")
	convertCmd.Flags().StringVar(&siteTitle, "title", "", "Site title")
	convertCmd.Flags().StringVar(&author, "author", "", "Default author name")
	convertCmd.Flags().BoolVar(&addList, "addlist", false, "Generate index.html with a list of all blog posts")
	convertCmd.Flags().BoolVar(&recursive, "recursive", false, "Process subdirectories recursively")
	convertCmd.Flags().BoolVar(&generateRSS, "rss", false, "Generate RSS feed (feed.xml)")
	convertCmd.Flags().StringVar(&siteURL, "site-url", "", "Site URL for RSS feed")
	convertCmd.Flags().StringVar(&gitWebURL, "git-web-url", "", "Base URL for git web interface (e.g., https://github.com/user/repo/commit/)")
	convertCmd.Flags().BoolVar(&showCommitInfo, "show-commit-info", false, "Display commit information in templates")
	convertCmd.Flags().BoolVar(&serve, "serve", false, "Start development server after building")
	convertCmd.Flags().StringVar(&port, "port", "8080", "Port to serve on")

	// Create-template command flags
	createTemplateCmd.Flags().StringVarP(&createTemplateFile, "template", "t", "templates/default.html", "Path for template file")

	// Serve command flags
	serveCmd.Flags().StringVarP(&serveOutputDir, "output", "o", "output", "Directory to serve")
	serveCmd.Flags().StringVar(&servePort, "port", "8080", "Port to serve on")

	// Minify command flags
	minifyCmd.Flags().StringVarP(&minifyInput, "input", "i", "", "Input file or directory (required)")
	minifyCmd.Flags().StringVarP(&minifyOutput, "output", "o", "", "Output file or directory (required)")
	minifyCmd.Flags().StringVarP(&minifyType, "type", "T", "css", "Type of files to minify: css, js, or all")
}

func runConvert(cmd *cobra.Command, args []string) error {
	logger := log.With().Str("command", "convert").Logger()

	// Validate required fields
	if inputDir == "" {
		logger.Error().Msg("input directory is required")
		return cmd.Usage()
	}

	logger.Info().Str("input", inputDir).Str("output", outputDir).Msg("Starting conversion")

	// Create default template if none provided
	if convertTemplateFile == "" {
		tmpDir := filepath.Join(outputDir, "tmp")
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			logger.Error().Err(err).Str("dir", tmpDir).Msg("Failed to create temp directory")
			return err
		}

		defaultTemplate := filepath.Join(tmpDir, "default_template.html")
		if err := converter.GenerateTemplateFile(defaultTemplate); err != nil {
			logger.Error().Err(err).Str("path", defaultTemplate).Msg("Failed to create default template")
			return err
		}

		convertTemplateFile = defaultTemplate
		logger.Debug().Str("path", defaultTemplate).Msg("Generated default template")

		// Clean up temp directory at exit
		defer os.RemoveAll(tmpDir)
	}

	// Configure the converter
	config := converter.Config{
		TemplateFile:   convertTemplateFile,
		OutputDir:      outputDir,
		CSSPath:        cssPath,
		JSPath:         jsPath,
		SiteTitle:      siteTitle,
		DefaultAuthor:  author,
		GenerateList:   addList,
		Recursive:      recursive,
		GenerateRSS:    generateRSS,
		SiteURL:        siteURL,
		GitWebURL:      gitWebURL,
		ShowCommitInfo: showCommitInfo,
	}

	// Convert files
	if err := converter.ConvertDirectory(inputDir, config); err != nil {
		logger.Error().Err(err).Msg("Conversion failed")
		return err
	}

	logger.Info().Str("output", outputDir).Msg("Conversion complete")

	// Start server if requested
	if serve {
		logger.Info().Str("port", port).Msg("Starting development server")
		if err := converter.StartServer(outputDir, port); err != nil {
			logger.Error().Err(err).Msg("Server failed")
			return err
		}
	}

	return nil
}

func runCreateTemplate(cmd *cobra.Command, args []string) error {
	logger := log.With().Str("command", "create-template").Logger()

	if err := converter.GenerateTemplateFile(createTemplateFile); err != nil {
		logger.Error().Err(err).Str("path", createTemplateFile).Msg("Failed to create template")
		return err
	}

	logger.Info().Str("path", createTemplateFile).Msg("Template created successfully")
	return nil
}

func runServe(cmd *cobra.Command, args []string) error {
	logger := log.With().Str("command", "serve").Logger()

	logger.Info().Str("dir", serveOutputDir).Str("port", servePort).Msg("Starting development server")
	if err := converter.StartServer(serveOutputDir, servePort); err != nil {
		logger.Error().Err(err).Msg("Server failed")
		return err
	}

	return nil
}

func runMinify(cmd *cobra.Command, args []string) error {
	logger := log.With().Str("command", "minify").Logger()

	// Validate required fields
	if minifyInput == "" {
		logger.Error().Msg("input is required")
		return cmd.Usage()
	}
	if minifyOutput == "" {
		logger.Error().Msg("output is required")
		return cmd.Usage()
	}

	// Validate type
	minifyType = strings.ToLower(minifyType)
	if minifyType != "css" && minifyType != "js" && minifyType != "all" {
		logger.Error().Str("type", minifyType).Msg("invalid type: must be css, js, or all")
		return cmd.Usage()
	}

	logger.Info().Str("input", minifyInput).Str("output", minifyOutput).Str("type", minifyType).Msg("Starting minification")

	// Check if input is a file or directory
	info, err := os.Stat(minifyInput)
	if err != nil {
		logger.Error().Err(err).Str("input", minifyInput).Msg("Failed to stat input")
		return err
	}

	if info.IsDir() {
		// Minify directory based on type
		switch minifyType {
		case "css":
			if err := converter.MinifyCSSDir(minifyInput, minifyOutput); err != nil {
				logger.Error().Err(err).Msg("CSS minification failed")
				return err
			}
		case "js":
			if err := converter.MinifyJSDir(minifyInput, minifyOutput); err != nil {
				logger.Error().Err(err).Msg("JS minification failed")
				return err
			}
		case "all":
			if err := converter.MinifyCSSDir(minifyInput, minifyOutput); err != nil {
				logger.Error().Err(err).Msg("CSS minification failed")
				return err
			}
			if err := converter.MinifyJSDir(minifyInput, minifyOutput); err != nil {
				logger.Error().Err(err).Msg("JS minification failed")
				return err
			}
		}
	} else {
		// Minify single file - detect type from extension if not specified
		ext := strings.ToLower(filepath.Ext(minifyInput))
		if minifyType == "all" {
			// Auto-detect from extension
			if ext == ".css" {
				minifyType = "css"
			} else if ext == ".js" {
				minifyType = "js"
			} else {
				logger.Error().Str("extension", ext).Msg("Cannot auto-detect file type")
				return fmt.Errorf("cannot auto-detect file type for extension %s", ext)
			}
		}

		switch minifyType {
		case "css":
			if err := converter.MinifyCSS(minifyInput, minifyOutput); err != nil {
				logger.Error().Err(err).Msg("CSS minification failed")
				return err
			}
		case "js":
			if err := converter.MinifyJS(minifyInput, minifyOutput); err != nil {
				logger.Error().Err(err).Msg("JS minification failed")
				return err
			}
		}
	}

	logger.Info().Msg("Minification complete")
	return nil
}

func main() {
	// Setup zerolog with pretty console output
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
