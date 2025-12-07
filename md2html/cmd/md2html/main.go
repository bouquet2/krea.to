package main

import (
	"os"
	"path/filepath"

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

	// Create-template command flags
	createTemplateFile string

	// Serve command flags
	serveOutputDir string
	servePort      string

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
)

func init() {
	// Add subcommands to root
	rootCmd.AddCommand(convertCmd)
	rootCmd.AddCommand(createTemplateCmd)
	rootCmd.AddCommand(serveCmd)

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
	convertCmd.Flags().BoolVar(&serve, "serve", false, "Start development server after building")
	convertCmd.Flags().StringVar(&port, "port", "8080", "Port to serve on")

	// Create-template command flags
	createTemplateCmd.Flags().StringVarP(&createTemplateFile, "template", "t", "templates/default.html", "Path for template file")

	// Serve command flags
	serveCmd.Flags().StringVarP(&serveOutputDir, "output", "o", "output", "Directory to serve")
	serveCmd.Flags().StringVar(&servePort, "port", "8080", "Port to serve on")
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
		TemplateFile:  convertTemplateFile,
		OutputDir:     outputDir,
		CSSPath:       cssPath,
		JSPath:        jsPath,
		SiteTitle:     siteTitle,
		DefaultAuthor: author,
		GenerateList:  addList,
		Recursive:     recursive,
		GenerateRSS:   generateRSS,
		SiteURL:       siteURL,
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

func main() {
	// Setup zerolog with pretty console output
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
