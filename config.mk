# Site Configuration Options
# These settings control the appearance and behavior of the generated website

# The title. Currently only used for the blog part of the website.
TITLE = "Kreato's Blog"

# Base URL for the site (used for RSS feeds and absolute links)
SITE_URL = "https://krea.to"

# Default theme for the website
# Available themes: nord, latte, frappe, mocha, macchiato, gruvbox, tokyonight, monokai, onedark, solarized, kanagawa
DEFAULT_THEME = "gruvbox"

# Show commit info to the user on the blog side of the website
# Set to 1 to enable, 0 to disable
# GITHUB_COMMIT_PREFIX is needed if this is enabled
SHOW_COMMIT_INFO ?= 1

# Debug mode on md2html
# Is required when you want to create a GitHub issue.
DEBUG ?= 0

# Base URL for git web interface (used when SHOW_COMMIT_INFO is enabled)
# This should point to your git repository's commit view
# Examples:
#   GitHub:  "https://github.com/user/repo/commit/"
#   GitLab:  "https://gitlab.com/user/repo/-/commit/"
#   Gitea:   "https://gitea.example.com/user/repo/commit/"
GITHUB_COMMIT_PREFIX = "https://github.com/bouquet2/krea.to/commit/"
