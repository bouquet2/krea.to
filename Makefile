.PHONY: all clean help build tidy fmt serve build-bin copy-assets convert-markdown minify-css minify-js

# Include configuration
include config.mk

# Environment variables
DEBUG_FLAG = $(if $(filter 1,$(DEBUG)),--debug,)

# Build configuration
SHOW_COMMIT_INFO_FLAG = $(if $(filter 1,$(SHOW_COMMIT_INFO)),--show-commit-info --git-web-url $(GITHUB_COMMIT_PREFIX),)
MD2HTML_BIN = md2html/md2html
DIST_DIR = dist

ASSETS = fonts assets
CONVERT_FLAGS = --input md --default-theme $(DEFAULT_THEME) --title $(TITLE) --output $(DIST_DIR) --css "css/style.css" --js "js/script.js" --addlist --recursive --rss --site-url $(SITE_URL) $(DEBUG_FLAG) $(SHOW_COMMIT_INFO_FLAG)

# Default target
all: build

# Build md2html binary
build-bin:
	@echo "Building md2html..."
	cd md2html && go build -o md2html ./cmd/md2html

# Copy static assets to dist directory
copy-assets:
	@echo "Copying static assets..."
	mkdir -p $(DIST_DIR)/blog
	$(foreach asset,$(ASSETS),cp -r $(asset) $(DIST_DIR)/;)
	cp CNAME $(DIST_DIR)/

# Minify CSS files
minify-css: build-bin
	@echo "Minifying CSS..."
	$(MD2HTML_BIN) minify --input css --output $(DIST_DIR)/css --type css

# Minify JS files
minify-js: build-bin
	@echo "Minifying JS..."
	$(MD2HTML_BIN) minify --input js --output $(DIST_DIR)/js --type js

# Convert markdown to HTML
convert-markdown:
	@echo "Generating site..."
	$(MD2HTML_BIN) convert $(CONVERT_FLAGS)

# Build site (landing page + blog posts)
build: tidy build-bin copy-assets minify-css minify-js convert-markdown
	@echo "Build complete. Output in $(DIST_DIR)/"

# Clean all generated files
clean:
	@echo "Cleaning..."
	rm -f $(MD2HTML_BIN)
	rm -rf $(DIST_DIR)

# Serve the site locally with a development server
serve: tidy build-bin copy-assets minify-css minify-js
	@echo "Generating site..."
	$(MD2HTML_BIN) convert $(CONVERT_FLAGS) --serve --port '8080'

# Format code
fmt:
	@echo "Formatting code..."
	gofmt -s -w md2html/cmd/md2html md2html/pkg
	cd md2html && go vet ./...

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	cd md2html && go mod tidy

# Help information
help:
	@echo "Available targets:"
	@echo "  all (default) - Build site from markdown files"
	@echo "  build         - Alias for all"
	@echo "  clean         - Remove generated files and binary"
	@echo "  fmt           - Format and lint code"
	@echo "  tidy          - Tidy Go dependencies"
	@echo "  serve         - Build site and start development server on localhost:8080"
	@echo "  help          - Show this message"
	@echo ""
	@echo "Configuration:"
	@echo "  Site options are configured in config.mk file"
	@echo "  Environment variables:"
	@echo "    DEBUG=1     - Enable debug mode for md2html (default: 0)"
	@echo ""
	@echo "Examples:"
	@echo "  make build           # Build without debug"
	@echo "  make DEBUG=1 build   # Build with debug enabled"
	@echo "  make DEBUG=1 serve   # Serve with debug enabled"
	@echo ""
	@echo "To add blog posts, create markdown files in md/blog/"
	@echo "The landing page is defined in md/index.md"
	@echo ""
	@echo "Edit config.mk to customize site title, theme, and other options"

