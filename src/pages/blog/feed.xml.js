import rss from '@astrojs/rss';
import { getCollection } from 'astro:content';
import { getPostTitle } from '../../utils/content';
import { getPostComputedMetadataById } from '../../utils/postMetadata';
import siteConfig from '../../../site.config.mjs';

export async function GET(context) {
  const posts = await getCollection('blog');
  posts.sort((a, b) => (b.data.date?.valueOf() || 0) - (a.data.date?.valueOf() || 0));
  
  return rss({
    title: siteConfig.TITLE,
    description: "Blog Posts and Articles by Kreato",
    site: context.site,
    items: posts.map(post => {
      const computed = getPostComputedMetadataById(post.id);
      const commitDate = computed?.commitDate ? new Date(computed.commitDate) : undefined;

      return {
        title: getPostTitle(post),
        pubDate: post.data.date || commitDate,
        description: post.data.description,
        link: `/blog/${post.id.replace(/\.md$/, '')}/`,
        author: post.data.author,
      };
    }),
  });
}
