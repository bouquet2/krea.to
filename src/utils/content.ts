import { getCollection } from 'astro:content';
import type { CollectionEntry } from 'astro:content';
import { marked } from 'marked';
import { getPostComputedMetadataById } from './postMetadata';

export async function getLandingPage(): Promise<CollectionEntry<'landing'>> {
  const landing = await getCollection('landing');
  if (!landing || landing.length === 0) {
    throw new Error('No landing page found in content collection');
  }
  return landing[0];
}

export async function getRecentPosts(limit: number = 5): Promise<Array<{ title: string; link: string }>> {
  const posts = await getCollection('blog');
  posts.sort((a, b) => (b.data.date?.valueOf() || 0) - (a.data.date?.valueOf() || 0));
  
  return posts.slice(0, limit).map(post => ({
    title: getPostTitle(post),
    link: `/blog/${post.id.replace(/\.md$/, '')}/`
  }));
}

// Get title from slug (which is the filename without extension)
// e.g., slug="Linux/Nix on macOS using nix-darwin, and my initial experiences" -> "Nix on macOS using nix-darwin, and my initial experiences"
export function getTitleFromSlug(slug: string): string {
  const parts = slug.split('/');
  return parts[parts.length - 1];
}

// Get post title: use frontmatter title if provided, otherwise derive from originalFilename
export function getPostTitle(entry: CollectionEntry<'blog'>): string {
  if (entry.data.title) {
    return entry.data.title;
  }
  if (entry.data.originalFilename) {
    return entry.data.originalFilename.replace(/\.md$/, '');
  }
  const metadata = getPostComputedMetadataById(entry.id);
  if (metadata?.title) {
    return metadata.title;
  }

  // Final fallback: convert slug to title (replace hyphens with spaces and capitalize)
  const basename = entry.id.split('/').pop()?.replace(/\.md$/, '') || 'Untitled';
  return basename
    .replace(/-/g, ' ')
    .replace(/\b\w/g, (char) => char.toUpperCase());
}

export interface TerminalSection {
  command: string;
  html: string;
}

export interface LandingLink {
  name: string;
  url: string;
  external: boolean;
}

export function parseLandingContent(content: string): {
  sections: TerminalSection[];
  links: LandingLink[];
} {
  const sections: TerminalSection[] = [];
  const links: LandingLink[] = [];
  
  const lines = content.split('\n');
  let currentSection: { command: string; contentLines: string[] } | null = null;
  let inLinks = false;
  
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    
    const sectionMatch = line.match(/<!--\s*Section:\s*(.+?)\s*-->/);
    if (sectionMatch) {
      if (currentSection) {
        const rawContent = currentSection.contentLines.join('\n').trim();
        if (rawContent) {
          sections.push({
            command: currentSection.command,
            html: marked(rawContent)
          });
        }
      }
      currentSection = {
        command: sectionMatch[1].trim(),
        contentLines: []
      };
      inLinks = false;
      continue;
    }
    
    if (line.includes('<!-- Links -->')) {
      if (currentSection) {
        const rawContent = currentSection.contentLines.join('\n').trim();
        if (rawContent) {
          sections.push({
            command: currentSection.command,
            html: marked(rawContent)
          });
        }
        currentSection = null;
      }
      inLinks = true;
      continue;
    }
    
    if (currentSection && !inLinks) {
      currentSection.contentLines.push(line);
    }
    
    if (inLinks) {
      const linkMatch = line.match(/-\s*\[([^\]]+)\]\(([^)]+)\)/);
      if (linkMatch) {
        const name = linkMatch[1];
        const url = linkMatch[2];
        const external = url.startsWith('http://') || url.startsWith('https://');
        links.push({ name, url, external });
      }
    }
  }
  
  if (currentSection) {
    const rawContent = currentSection.contentLines.join('\n').trim();
    if (rawContent) {
      sections.push({
        command: currentSection.command,
        html: marked(rawContent)
      });
    }
  }
  
  return { sections, links };
}
