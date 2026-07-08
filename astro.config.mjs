import { defineConfig } from 'astro/config';
import sitemap from '@astrojs/sitemap';
import siteConfig from './site.config.mjs';

export default defineConfig({
  site: siteConfig.SITE_URL,
  integrations: [sitemap()],
});
