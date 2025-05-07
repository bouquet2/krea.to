# md2html

A simple tool to convert Markdown files to HTML with Pokemon DP Pro theme styling for Kreato's website.

## Features

- Converts Markdown files to HTML
- Customizable HTML templates with Pokemon DP Pro theme
- Support for metadata in HTML comments
- Automatically handles common Markdown extensions
- Clean, responsive output with light/dark mode support

## Installation

Make sure you have Go installed (version 1.19+), then run:

```bash
# Clone the repository
git clone https://github.com/kreatoo/md2html.git
cd md2html

# Build the tool
go build -o md2html ./cmd/md2html
```

## Usage

```bash
# Convert all markdown files in a directory
./md2html --input=/path/to/markdown/files --output=/path/to/output/directory

# Create a default template file
./md2html --create-template --template=my-template.html

# Use a custom template and CSS
./md2html --input=markdown --output=html --template=my-template.html --css=styles.css

# Set default author and site title
./md2html --input=markdown --output=html --author="Kreato" --title="Kreato's Blog"
```

### Command-line options

- `--input`: Directory containing markdown files (required)
- `--output`: Directory to save HTML files (default: "output")
- `--template`: HTML template file (optional, will create default if not provided)
- `--css`: Path to CSS file (optional)
- `--title`: Site title (default: "Kreato's Website")
- `--author`: Default author name (default: "Kreato")
- `--create-template`: Create a default template and exit

## Markdown Metadata

You can add metadata to your Markdown files using HTML comments at the top of the file:

```markdown
<!--
Title: My Cool Page
Author: Kreato
Description: A page about cool things
-->

# My Cool Page

Content goes here...
```

## Example

```markdown
<!--
Title: Installing Kreato Linux
Author: Kreato
Description: A guide to installing Kreato Linux
-->

# Installing Kreato Linux

Welcome to this guide on installing Kreato Linux on your system.

## System Requirements

- 2GB RAM
- 20GB disk space
- 64-bit processor

## Installation Steps

1. Download the ISO
2. Create bootable USB
3. Boot from USB
4. Follow the installer
```

Will be converted to styled HTML with the Pokemon DP Pro theme.

## License

MIT 