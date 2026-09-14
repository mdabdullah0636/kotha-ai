package prompt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"kotha/internal/config"
)

func TestGetContextFromPaths(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfg, err := config.Load(tmpDir, false)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	cfg.WorkingDir = tmpDir
	cfg.ContextPaths = []string{
		"file.txt",
		"directory/",
	}
	testFiles := []string{
		"file.txt",
		"directory/file_a.txt",
		"directory/file_b.txt",
		"directory/file_c.txt",
	}

	createTestFiles(t, tmpDir, testFiles)

	context := getContextFromPaths()
	assert.Contains(t, context, "# From:"+tmpDir+"/file.txt")
	assert.Contains(t, context, "file.txt: test content")
	assert.Contains(t, context, "# From:"+tmpDir+"/directory/file_a.txt")
	assert.Contains(t, context, "directory/file_a.txt: test content")
	assert.Contains(t, context, "# From:"+tmpDir+"/directory/file_b.txt")
	assert.Contains(t, context, "directory/file_b.txt: test content")
	assert.Contains(t, context, "# From:"+tmpDir+"/directory/file_c.txt")
	assert.Contains(t, context, "directory/file_c.txt: test content")
}

func createTestFiles(t *testing.T, tmpDir string, testFiles []string) {
	t.Helper()
	for _, path := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		if path[len(path)-1] == '/' {
			err := os.MkdirAll(fullPath, 0755)
			require.NoError(t, err)
		} else {
			dir := filepath.Dir(fullPath)
			err := os.MkdirAll(dir, 0755)
			require.NoError(t, err)
			err = os.WriteFile(fullPath, []byte(path+": test content"), 0644)
			require.NoError(t, err)
		}
	}
}
