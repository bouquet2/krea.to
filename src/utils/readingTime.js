export function calculateReadingTime(content) {
  const wordsPerMinute = 200;
  const words = content.trim().split(/\s+/).filter(Boolean).length;
  const minutes = Math.ceil(words / wordsPerMinute);
  return `~${minutes} min read`;
}

function stripMarkdown(body = '') {
  return body
    .replace(/```[\s\S]*?```/g, ' ')
    .replace(/~~~[\s\S]*?~~~/g, ' ')
    .replace(/`[^`]*`/g, ' ')
    .replace(/!\[[^\]]*\]\([^)]*\)/g, ' ')
    .replace(/\[([^\]]*)\]\([^)]*\)/g, '$1')
    .replace(/^>{1,}\s+/gm, ' ')
    .replace(/^#{1,6}\s+/gm, ' ')
    .replace(/^\s*[-*+]\s+/gm, ' ')
    .replace(/[*_~]/g, ' ')
    .replace(/^\s*[-*_]{3,}\s*$/gm, ' ')
    .replace(/\s+/g, ' ')
    .trim();
}

export function getReadingTime(body) {
  return calculateReadingTime(stripMarkdown(body));
}
