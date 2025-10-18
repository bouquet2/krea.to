# krea.to

Mu personal website and blog, featuring a terminal-style interface and a custom markdown-to-HTML blog system.

🌐 **Live Site**: [krea.to](https://krea.to)

## Features

### Main Site
- **Interactive Terminal UI** - Browse the site through a simulated terminal interface
- **11 Color Schemes** - Choose from Catppuccin variants, Gruvbox, Nord, Tokyo Night, Monokai, One Dark, Solarized, and Kanagawa while having custom selected backgrounds for each
- **Theme Toggle** - Quick switch between light and dark modes
- **Accessibility Mode** - Enhanced readability with Comic Sans font and increased spacing
- **Customizable Settings** - Adjust font size, background, and transparency (including the intensity) with ease and slick animations
- **Terminal Commands** - Interactive commands like `help`, `ls`, `cd`, `cat`, and more to tinker with and explore

### Blog
- **Clean Reading Experience** - Distraction-free blog post layout
- **Automatic Generation** - Blog posts are generated from Markdown files
- **Syntax Highlighting** - Code blocks with syntax highlighting
- **Responsive Design** - Works great on all screen sizes, including but not limited to phone, tablet, and laptop

## Building the Blog

The blog uses a custom Go-based markdown converter. To generate blog posts:

```bash
# Build the converter and generate all blog pages
make blog

# Clean generated files
make clean

# View all available commands
make help
```

## Adding New Blog Posts

1. Create a new markdown file in the `md/` directory (organized by category)
2. Add metadata at the top of the file:
   ```markdown
   <!--
   Title: Your Post Title
   Author: Kreato
   Description: A brief description
   -->
   
   # Your Post Title
   
   Your content here...
   ```
3. Run `make blog` to generate the HTML
4. Commit and push the changes

## Development

The site uses vanilla HTML, CSS, and JavaScript - no build process required for the main site.

- Edit `index.html` for main site content
- Edit `css/` files for styling changes
- Edit `js/script.js` for interactive features
- Blog posts are in `md/` directory

## Color Schemes

The site includes 11 carefully selected color schemes:
- **Catppuccin**: Mocha (default), Frappe, Latte, Macchiato
- **Gruvbox**: Retro groove colors
- **Nord**: Arctic, north-bluish color palette
- **Tokyo Night**: Dark theme inspired by Tokyo at night
- **Monokai**: Classic code editor theme
- **One Dark**: Atom's iconic One Dark theme
- **Solarized**: Precision colors for machines and people
- **Kanagawa**: Dark theme inspired by Kanagawa paintings

## Font

The main site uses the **Scientifica** font for a clean, retro terminal aesthetic, while the blog uses **Pokemon DP Pro** for that nostalgic feel, with Arial as a fallback.

## License

See [LICENSE](LICENSE) file for details.
