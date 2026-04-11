import { execSync } from 'child_process';
import { readdirSync, statSync } from 'fs';
import fs from 'fs';
import * as git from 'isomorphic-git';
import { join, relative, sep } from 'path';
import siteConfig from '../../site.config.mjs';

interface PostComputedMetadata {
  title: string;
  originalDirectory?: string;
  commitHash?: string;
  commitDate?: string;
  commitAuthor?: string;
  commitURL?: string;
}

const BLOG_ROOT = join(process.cwd(), 'src/content/blog');

let cache: Map<string, PostComputedMetadata> | null = null;
const REPOSITORY_URL = await resolveRepositoryURL();

function slugifySegment(segment: string): string {
  return segment
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '');
}

function toEntryId(relativePath: string): string {
  const parts = relativePath.split('/');
  const fileName = parts.pop() || '';
  const baseName = fileName.replace(/\.md$/, '');
  const slugParts = parts.map(slugifySegment);
  slugParts.push(slugifySegment(baseName));
  return slugParts.join('/');
}

function queryGitInfo(repoRelativePath: string): { hash: string; date: string; author: string } | null {
  const command = `git log -1 --format=%H%n%ai%n%an -- "${repoRelativePath}"`;
  const output = execSync(command, { encoding: 'utf-8' }).trim();

  if (siteConfig.DEBUG) {
    console.log(`[postMetadata] git query path=${repoRelativePath} hasOutput=${output.length > 0}`);
  }

  if (!output) return null;

  const [hash, date, author] = output.split('\n');
  if (!hash || !date || !author) return null;

  return { hash, date, author };
}

function normalizeRemoteURL(remoteURL: string): string | undefined {
  const trimmed = remoteURL.trim();
  if (!trimmed) return undefined;

  const scpLikeMatch = trimmed.match(/^git@([^:]+):(.+)$/);
  if (scpLikeMatch) {
    const host = scpLikeMatch[1];
    const repoPath = scpLikeMatch[2].replace(/\.git$/, '');
    return `https://${host}/${repoPath}`;
  }

  const sshLikeMatch = trimmed.match(/^ssh:\/\/git@([^/]+)\/(.+)$/);
  if (sshLikeMatch) {
    const host = sshLikeMatch[1];
    const repoPath = sshLikeMatch[2].replace(/\.git$/, '');
    return `https://${host}/${repoPath}`;
  }

  const httpLikeMatch = trimmed.match(/^https?:\/\/(.+)$/);
  if (httpLikeMatch) {
    return `https://${httpLikeMatch[1]}`.replace(/\.git$/, '');
  }

  return undefined;
}

async function resolveRepositoryURL(): Promise<string | undefined> {
  try {
    const dir = process.cwd();
    let remoteName: string | undefined;

    const currentBranch = await git.currentBranch({ fs, dir, fullname: false });
    if (currentBranch) {
      const upstream = await git.getConfig({ fs, dir, path: `branch.${currentBranch}.remote` });
      if (upstream) remoteName = upstream;
    }

    const remotes = await git.listRemotes({ fs, dir });
    const names = remotes.map((r) => r.remote);

    if (!remoteName || !names.includes(remoteName)) {
      remoteName = names.includes('origin') ? 'origin' : names[0];
    }

    if (!remoteName) return undefined;

    const selected = remotes.find((r) => r.remote === remoteName);
    return selected?.url ? normalizeRemoteURL(selected.url) : undefined;
  } catch {
    return undefined;
  }
}

function readGitInfo(repoRelativePaths: string[]): Omit<PostComputedMetadata, 'title' | 'originalDirectory'> {
  try {
    let gitInfo: { hash: string; date: string; author: string } | null = null;
    for (const path of repoRelativePaths) {
      gitInfo = queryGitInfo(path);
      if (gitInfo) break;
    }

    if (!gitInfo) return {};

    const repoURL = REPOSITORY_URL;

    if (siteConfig.DEBUG) {
      console.log(`[postMetadata] resolved commit ${gitInfo.hash.slice(0, 7)} repoURL=${repoURL || 'none'}`);
    }

    return {
      commitHash: gitInfo.hash.slice(0, 7),
      commitDate: gitInfo.date,
      commitAuthor: gitInfo.author,
      commitURL: repoURL ? `${repoURL}/commit/${gitInfo.hash}` : undefined,
    };
  } catch {
    return {};
  }
}

function getGitInfo(repoRelativePath: string, legacyPath: string): Omit<PostComputedMetadata, 'title' | 'originalDirectory'> {
  return readGitInfo([repoRelativePath, legacyPath]);
}

function walk(dir: string, files: string[]): void {
  const entries = readdirSync(dir);
  for (const entry of entries) {
    const fullPath = join(dir, entry);
    const stat = statSync(fullPath);
    if (stat.isDirectory()) {
      walk(fullPath, files);
    } else if (entry.endsWith('.md')) {
      files.push(fullPath);
    }
  }
}

function buildCache(): Map<string, PostComputedMetadata> {
  const map = new Map<string, PostComputedMetadata>();
  const files: string[] = [];
  walk(BLOG_ROOT, files);

  for (const filePath of files) {
    const rel = relative(BLOG_ROOT, filePath).split(sep).join('/');
    const repoRel = `src/content/blog/${rel}`;
    const legacyRel = `md/blog/${rel}`;
    const id = toEntryId(rel);
    const pathParts = rel.split('/');
    const fileName = pathParts[pathParts.length - 1] || '';
    const title = fileName.replace(/\.md$/, '');
    const originalDirectory = pathParts.length > 1 ? pathParts[pathParts.length - 2] : undefined;

    map.set(id, {
      title,
      originalDirectory,
      ...getGitInfo(repoRel, legacyRel),
    });
  }

  return map;
}

function getCache(): Map<string, PostComputedMetadata> {
  if (!cache) {
    cache = buildCache();
  }
  return cache;
}

export function getPostComputedMetadataById(id: string): PostComputedMetadata | undefined {
  return getCache().get(id);
}
