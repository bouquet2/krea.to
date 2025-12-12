package converter

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
)

// GitInfo holds git-related information for a file
type GitInfo struct {
	LastModified time.Time
	CommitHash   string
	Author       string
}

// GitRepository wraps go-git repository for file history lookups
type GitRepository struct {
	repo     *git.Repository
	repoRoot string
}

// OpenGitRepository opens a git repository at the given path
// It walks up the directory tree to find the repository root
func OpenGitRepository(path string) (*GitRepository, error) {
	// Find repository root by walking up
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpenWithOptions(absPath, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, err
	}

	// Get the worktree to find the root path
	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}

	return &GitRepository{
		repo:     repo,
		repoRoot: wt.Filesystem.Root(),
	}, nil
}

// GetFileLastModified returns the last commit date for a file
func (gr *GitRepository) GetFileLastModified(filePath string) (*GitInfo, error) {
	logger := log.With().Str("file", filePath).Logger()

	// Convert to path relative to repository root
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, err
	}

	relPath, err := filepath.Rel(gr.repoRoot, absPath)
	if err != nil {
		return nil, err
	}

	// Normalize path separators for git (always use forward slashes)
	relPath = strings.ReplaceAll(relPath, string(filepath.Separator), "/")

	logger.Debug().Str("rel_path", relPath).Msg("Looking up git history for file")

	// Get the commit iterator for this specific file
	cIter, err := gr.repo.Log(&git.LogOptions{
		PathFilter: func(path string) bool {
			return path == relPath
		},
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, err
	}
	defer cIter.Close()

	// Get the first (most recent) commit
	commit, err := cIter.Next()
	if err != nil {
		return nil, err
	}

	logger.Debug().
		Str("commit", commit.Hash.String()[:8]).
		Time("date", commit.Committer.When).
		Str("author", commit.Author.Name).
		Msg("Found last commit for file")

	return &GitInfo{
		LastModified: commit.Committer.When,
		CommitHash:   commit.Hash.String(),
		Author:       commit.Author.Name,
	}, nil
}

// GetFilesLastModified returns the last commit date for multiple files
// Returns a map of file path -> GitInfo
func (gr *GitRepository) GetFilesLastModified(filePaths []string) map[string]*GitInfo {
	result := make(map[string]*GitInfo)

	for _, filePath := range filePaths {
		info, err := gr.GetFileLastModified(filePath)
		if err != nil {
			log.Debug().Str("file", filePath).Err(err).Msg("Could not get git info for file")
			continue
		}
		result[filePath] = info
	}

	return result
}

// GetLastModifiedTime is a convenience function to get just the time
func GetLastModifiedTime(filePath string) (time.Time, error) {
	repo, err := OpenGitRepository(filepath.Dir(filePath))
	if err != nil {
		return time.Time{}, err
	}

	info, err := repo.GetFileLastModified(filePath)
	if err != nil {
		return time.Time{}, err
	}

	return info.LastModified, nil
}

// GetFileGitInfo is a convenience function to get full git info for a file
func GetFileGitInfo(filePath string) (*GitInfo, error) {
	repo, err := OpenGitRepository(filepath.Dir(filePath))
	if err != nil {
		return nil, err
	}

	return repo.GetFileLastModified(filePath)
}

// GetCommitURL creates a URL to the commit in a git web interface
func GetCommitURL(gitWebURL, commitHash string) string {
	if gitWebURL == "" || commitHash == "" {
		return ""
	}

	// Ensure the base URL ends with a slash
	if !strings.HasSuffix(gitWebURL, "/") {
		gitWebURL += "/"
	}

	return gitWebURL + commitHash
}
