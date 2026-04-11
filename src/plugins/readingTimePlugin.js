import { calculateReadingTime } from '../utils/readingTime.js';
import { toString } from 'mdast-util-to-string';

export function readingTimePlugin() {
  return (tree, file) => {
    const textOnPage = toString(tree);
    const readTime = calculateReadingTime(textOnPage);
    
    if (file.data.astro?.frontmatter) {
      file.data.astro.frontmatter.readTime = readTime;
    }
  };
}
