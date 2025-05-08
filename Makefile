.PHONY: all blog clean help build-md2html clean-md2html fmt-md2html tidy-md2html

# Default target
all: blog

# Blog generation target
blog: build-md2html
	@echo "Converting markdown files..."
	mkdir -p blog
	md2html/md2html -input md -output blog -css "css/style-blog.css" -addlist -recursive

# Build md2html binary
build-md2html: tidy-md2html
	@echo "Building md2html..."
	cd md2html && go build -o md2html ./cmd/md2html

# Clean output files
clean: clean-md2html
	@echo "Cleaning output files..."
	rm -rf blog

# Clean md2html binary
clean-md2html:
	@echo "Cleaning md2html..."
	rm -f md2html/md2html

# Format md2html code
fmt-md2html:
	@echo "Formatting md2html code..."
	gofmt -s -w md2html/cmd/md2html md2html/pkg
	cd md2html && go vet ./...

# Tidy md2html dependencies
tidy-md2html:
	@echo "Tidying md2html dependencies..."
	cd md2html && go mod tidy

# Help information
help:
	@echo "Available targets:"
	@echo "  all          - Generate all blog pages (default)"
	@echo "  blog         - Generate blog pages from markdown files"
	@echo "  clean        - Remove all generated files"
	@echo "  build-md2html - Build the md2html binary"
	@echo "  clean-md2html - Clean the md2html binary"
	@echo "  fmt-md2html  - Format md2html code"
	@echo "  tidy-md2html - Tidy md2html dependencies"
	@echo "  help         - Show this help message"
	@echo ""
	@echo "To add new blog posts, create markdown files in the md2html directory"
	@echo "Generated output will be in the blog directory" 
