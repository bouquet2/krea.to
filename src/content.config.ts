import { defineCollection, z } from 'astro:content';
import { glob } from 'astro/loaders';

const blog = defineCollection({
  loader: glob({ pattern: '**/*.md', base: './src/content/blog' }),
  schema: z.object({
    author: z.string().default('Kreato'),
    tags: z.array(z.string()).default([]),
    date: z.coerce.date().optional(),
    title: z.string().optional(),
    description: z.string().optional(),
    commitHash: z.string().optional(),
    commitDate: z.string().optional(),
    commitAuthor: z.string().optional(),
    readTime: z.string().optional(),
  }),
});

const landing = defineCollection({
  loader: glob({ pattern: '**/*.md', base: './src/content/landing' }),
  schema: z.object({
    title: z.string(),
    description: z.string(),
    author: z.string().optional(),
    template: z.string().optional(),
    settings: z.object({
      'hide-topbar': z.boolean().optional(),
      'hide-shell': z.boolean().optional(),
      'fullscreen': z.boolean().optional(),
    }).optional(),
  }),
});

export const collections = { blog, landing };
