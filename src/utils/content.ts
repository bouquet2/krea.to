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

export async function getRecentPosts(limit: number = 5): Promise<Array<{ title: string; link: string; commitHash?: string; commitURL?: string }>> {
  const posts = await getCollection('blog');
  posts.sort((a, b) => (b.data.date?.valueOf() || 0) - (a.data.date?.valueOf() || 0));
  
  return posts.slice(0, limit).map(post => {
    const metadata = getPostComputedMetadataById(post.id);
    return {
      title: getPostTitle(post),
      link: `/blog/${post.id.replace(/\.md$/, '')}/`,
      commitHash: metadata?.commitHash,
      commitURL: metadata?.commitURL
    };
  });
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

export interface LandingSection {
  label: string;
  html: string;
}

export interface LandingLink {
  name: string;
  url: string;
  external: boolean;
}

/**
 * Parse landing-page markdown into a header, labelled sections, and links.
 * Sections are delimited by `<!-- Section: <label> -->`; links by `<!-- Links -->`.
 */
export function parseLandingContent(content: string): {
  headerHtml: string;
  sections: LandingSection[];
  links: LandingLink[];
} {
  const sections: LandingSection[] = [];
  const links: LandingLink[] = [];

  const lines = content.split('\n');
  const headerLines: string[] = [];
  let currentSection: { label: string; contentLines: string[] } | null = null;
  let inLinks = false;
  let sawSection = false;

  const flushSection = () => {
    if (!currentSection) return;
    const raw = currentSection.contentLines.join('\n').trim();
    if (raw) sections.push({ label: currentSection.label, html: marked(raw) });
  };

  for (const line of lines) {
    const sectionMatch = line.match(/<!--\s*Section:\s*(.+?)\s*-->/);
    if (sectionMatch) {
      flushSection();
      currentSection = { label: sectionMatch[1].trim(), contentLines: [] };
      sawSection = true;
      inLinks = false;
      continue;
    }

    if (/<!--\s*Links\s*-->/.test(line)) {
      flushSection();
      currentSection = null;
      inLinks = true;
      continue;
    }

    if (inLinks) {
      const linkMatch = line.match(/-\s*\[([^\]]+)\]\(([^)]+)\)/);
      if (linkMatch) {
        const url = linkMatch[2];
        links.push({ name: linkMatch[1], url, external: url.startsWith('http://') || url.startsWith('https://') });
      }
      continue;
    }

    if (currentSection) {
      currentSection.contentLines.push(line);
    } else if (!sawSection) {
      headerLines.push(line);
    }
  }

  flushSection();

  return { headerHtml: marked(headerLines.join('\n').trim()), sections, links };
}
