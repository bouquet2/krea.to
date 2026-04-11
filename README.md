# krea.to

My personal website and blog, featuring a terminal-style interface and an Astro-powered static site.

🌐 **Live Site**: [krea.to](https://krea.to)

## Features

### Main Site
- **Interactive Terminal UI** - Browse the site through a simulated terminal interface
- **11 Color Schemes** - Choose from Catppuccin variants, Gruvbox, Nord, Tokyo Night, Monokai, One Dark, Solarized, and Kanagawa while having custom selected backgrounds for each
- **Theme Toggle** - Quick switch between light and dark modes
- **Accessibility Mode** - Enhanced readability with Comic Sans font and increased spacing
- **Customizable Settings** - Adjust font size, background, and transparency (including the intensity) with ease and slick animations
- **Automatic Generation** - Pages are generated from Markdown content collections using Astro

### Blog
- **Clean Reading Experience** - Distraction-free blog post layout that uses the most of the current device
- **Syntax Highlighting** - Code blocks with syntax highlighting
- **Responsive Design** - Works great on all screen sizes, including but not limited to phone, tablet, and laptop

## Building

Install dependencies and build with Bun:

```bash
bun install
bun run build
```

For local development:

```bash
bun run dev
```

## Adding New Blog Posts

1. Create a new markdown file in `src/content/blog/<Category>/`.
2. Add frontmatter metadata at the top of the file:
   ```markdown
   ---
   title: "Your Post Title"
   author: "Kreato"
   description: "A brief description"
   tags: ["tag1", "tag2"]
   date: 2026-01-01
   ---
   ```
3. Run `bun run build` to generate the site.
4. Commit and push the changes.

## Landing Page Settings

The landing page template supports additional settings to customize its appearance. Add a `Settings` field to the metadata with comma-separated options:

```markdown
<!--
Title: My Site
Description: Site description
Template: landing
Settings: hide-topbar, fullscreen, hide-shell
-->
```

### Available Settings

| Setting | Description |
|---------|-------------|
| `hide-topbar` | Removes the terminal header bar (window buttons and title) |
| `fullscreen` | Makes the terminal take up the full viewport without margins or borders |
| `hide-shell` | Hides the shell prompts (`kreato@akiri:~$`) before each section |

Settings can be combined as needed. For example, `Settings: fullscreen, hide-topbar` creates a clean fullscreen terminal without the window chrome.

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
