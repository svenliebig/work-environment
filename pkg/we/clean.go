package we

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils/bytes"
	"github.com/svenliebig/work-environment/pkg/utils/git"
	"github.com/svenliebig/work-environment/pkg/utils/tablewriter"
)

type CleanOptions struct {
	Preview bool
}

func Clean(ctx context.BaseContext, o *CleanOptions) error {
	projects := ctx.GetProjectsInPath()

	tw := tablewriter.New()

	totalSize := int64(0)
	totalCleanableSize := int64(0)

	tw.Write([]byte("📁 Project   \t| 📦 Total Size \t| 🧹 Cleanable Size \t| 💾 Ratio"))
	tw.Line()

	type projectResult struct {
		project       *core.Project
		size          int64
		cleanableSize int64
		err           error
	}

	resultChan := make(chan projectResult)

	// Launch goroutines for each project
	for _, p := range projects {
		go func(proj *core.Project) {
			result := projectResult{project: proj}
			
			// Use waitgroup to synchronize the two goroutines
			var wg sync.WaitGroup
			wg.Add(2)
			
			// Get total size concurrently
			go func() {
				defer wg.Done()
				size, err := getTotalSize(proj)
				if err != nil {
					result.err = err
					return
				}
				result.size = size
			}()
			
			// Get cleanable size concurrently
			go func() {
				defer wg.Done()
				cleanSize, err := getCleanableSize(proj)
				if err != nil {
					result.err = err
					return
				}
				result.cleanableSize = cleanSize
			}()
			
			wg.Wait()
			resultChan <- result
		}(p)
	}

	formatOptions := &bytes.FormatOptions{
		Colorize:  true,
		Precision: 2,
	}

	// Collect and process results
	for range projects {
		result := <-resultChan
		if result.err != nil {
			return result.err
		}

		totalSize += result.size
		totalCleanableSize += result.cleanableSize

		ratio := float64(result.cleanableSize) / float64(result.size)

		tw.Write([]byte(result.project.Identifier + " \t| " + bytes.Format(result.size, formatOptions) + " \t| " + bytes.Format(result.cleanableSize, formatOptions) + " \t| " + fmt.Sprintf("%.2f%%", ratio * 100)))
		tw.Print()
	}

	tw.Line()
	tw.Write([]byte("🔍 Total   \t| " + bytes.Format(totalSize, formatOptions) + " \t| " + bytes.Format(totalCleanableSize, formatOptions) + " \t| " + fmt.Sprintf("%.2f%%", float64(totalCleanableSize) / float64(totalSize) * 100)))
	tw.Print()

	return nil
}

// returns the total size of the project in bytes
func getTotalSize(p *core.Project) (int64, error) {
	totalSize, err := getTotalSizeOfPath(p.Path)
	if err != nil {
		return 0, err
	}
	
	return totalSize, nil
}

func getTotalSizeOfPath(path string) (int64, error) {
	var totalSize int64 = 0

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})
	
	if err != nil {
		return 0, err
	}
	
	return totalSize, nil
}


func getCleanableSize(p *core.Project) (int64, error) {
	var cleanableSize int64 = 0

	ignoredDirectories, err := getIgnoredFilesOrDirectories(p)

	// fmt.Println("ignoredDirectories", ignoredDirectories)

	if err != nil {
		return 0, err
	}

	for _, fileOrDirectory := range ignoredDirectories {
		_, err := os.Stat(fileOrDirectory)
		if os.IsNotExist(err) {
			continue
		}
		
		size, err := getTotalSizeOfPath(fileOrDirectory)
		if err != nil {
			return 0, err
		}
		cleanableSize += size
	}

	return cleanableSize, nil
}

const (
	NodeModules = "node_modules"
	Dist = "dist"
	Target = "target"
	Build = "build"
)

// checks the .gitignore against some common big files and directories
// and returns the paths of these files and directories.
//
// - node_modules
// - dist
// - target
// - build
// - .next
// - .env
// - .env.local
// - .env.development.local
func getIgnoredFilesOrDirectories(p *core.Project) ([]string, error) {
	var ignoredFiles []string

	ignoredDirectories, err := git.FindAllIgnoredExistingDirectories(p.Path)
	if err != nil {
		return nil, err
	}

	for _, ignoredDirectory := range ignoredDirectories {
		_, err := os.OpenFile(ignoredDirectory, os.O_RDONLY, 0644)
		if os.IsNotExist(err) {
			return []string{}, nil
		} else if err != nil {
			return nil, err
		}

		switch {
			case matchFileOrDirectory(ignoredDirectory, NodeModules):
				fallthrough
			case matchFileOrDirectory(ignoredDirectory, Dist):
				fallthrough
			case matchFileOrDirectory(ignoredDirectory, Target):
				fallthrough
			case matchFileOrDirectory(ignoredDirectory, Build):
				ignoredFiles = append(ignoredFiles, ignoredDirectory)
		}
	}

	return ignoredFiles, nil
}

func matchFileOrDirectory(line string, folder string) bool {
	// fmt.Println("line", line, "folder", folder, "match", path.Base(line) == folder)
	return path.Base(line) == folder
}