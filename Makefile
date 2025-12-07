.PHONY: all clean help build tidy fmt

# Default target
all: build

# Build site (landing page + blog posts)
build: tidy
	@echo "Building md2html..."
	cd md2html && go build -o md2html ./cmd/md2html
	@echo "Generating site..."
	mkdir -p dist/blog
	md2html/md2html -input md -output dist -css "css/style-blog.css" -addlist -recursive -rss -site-url 'https://krea.to'
	@echo "Copying static assets..."
	cp -r css dist/
	cp -r js dist/
	cp -r fonts dist/
	cp -r assets dist/
	cp CNAME dist/
	@echo "Build complete. Output in dist/"

# Clean all generated files
clean:
	@echo "Cleaning..."
	rm -f md2html/md2html
	rm -rf dist

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
	@echo "  help          - Show this message"
	@echo ""
	@echo "To add blog posts, create markdown files in md/blog/"
	@echo "The landing page is defined in md/index.md" 
