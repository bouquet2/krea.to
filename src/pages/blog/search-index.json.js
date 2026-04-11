import { getCollection } from 'astro:content';
import { getPostTitle } from '../../utils/content';
import { marked } from 'marked';

function stripHtml(html) {
  return html.replace(/<[^>]*>/g, ' ').replace(/\s+/g, ' ').trim();
}

function truncateText(text, maxLength = 5000) {
  if (text.length <= maxLength) return text;
  return text.substring(0, maxLength) + '...';
}

export async function GET() {
  const posts = await getCollection('blog');
  
  const searchIndex = posts.map(post => {
    const content = truncateText(stripHtml(marked(post.body)), 5000);
    
    return {
      title: getPostTitle(post),
      description: post.data.description || '',
      slug: post.id.replace(/\.md$/, ''),
      url: `/blog/${post.id.replace(/\.md$/, '')}/`,
      id: post.id,
      tags: post.data.tags || [],
      date: post.data.date?.toISOString() || '',
      content
    };
  });
  
  return new Response(JSON.stringify(searchIndex), {
    headers: {
      'Content-Type': 'application/json'
    }
  });
}
