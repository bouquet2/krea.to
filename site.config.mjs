// Site Configuration Options
// These settings control the appearance and behavior of the website.

const siteConfig = {
  // The title used for the blog pages and RSS feed.
  TITLE: "Kreato's Blog",

  // Base URL for the site (used for RSS feeds, sitemap, and absolute links).
  SITE_URL: 'https://krea.to',

  // Default theme for the website.
  // Available themes: nord, latte, frappe, mocha, macchiato, gruvbox,
  // tokyonight, monokai, onedark, solarized, kanagawa, pinkie
  DEFAULT_THEME: 'pinkie',

  // Show commit info links on blog pages.
  // true to enable, false to disable
  SHOW_COMMIT_INFO: true,

  // Debug mode for metadata generation.
  // true to enable, false to disable
  DEBUG: false,
};

export default siteConfig;
