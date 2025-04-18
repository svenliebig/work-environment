package we

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/svenliebig/work-environment/pkg/context"
	"github.com/svenliebig/work-environment/pkg/core"
	"github.com/svenliebig/work-environment/pkg/utils/bytes"
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
	for _, p := range projects {
		size, err := getTotalSize(p)
		if err != nil {
			return err
		}

		cleanableSize, err := getCleanableSize(p)
		if err != nil {
			return err
		}

		totalSize += size
		totalCleanableSize += cleanableSize

		ratio := float64(cleanableSize) / float64(size)

		tw.Write([]byte(p.Identifier + " \t| " + bytes.Format(size, &bytes.FormatOptions{
			Colorize: true,
		}) + " \t| " + bytes.Format(cleanableSize, &bytes.FormatOptions{
			Colorize: true,
		}) + " \t| " + fmt.Sprintf("%.2f%%", ratio * 100)))
		tw.Print()
	}

	tw.Line()
	tw.Write([]byte("🔍 Total   \t| " + bytes.Format(totalSize, &bytes.FormatOptions{
		Colorize: true,
	}) + " \t| " + bytes.Format(totalCleanableSize, &bytes.FormatOptions{
		Colorize: true,
	}) + " \t| " + fmt.Sprintf("%.2f%%", float64(totalCleanableSize) / float64(totalSize) * 100)))
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

	ignoredFilesOrDirectories, err := getIgnoredFilesOrDirectories(p)
	if err != nil {
		return 0, err
	}

	for _, fileOrDirectory := range ignoredFilesOrDirectories {
		_, err := os.Stat(filepath.Join(p.Path, fileOrDirectory))
		if os.IsNotExist(err) {
			continue
		}
		
		size, err := getTotalSizeOfPath(filepath.Join(p.Path, fileOrDirectory))
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

	gitignorePath := filepath.Join(p.Path, ".gitignore")

	file, err := os.OpenFile(gitignorePath, os.O_RDONLY, 0644)
	if os.IsNotExist(err) {
		return []string{}, nil
	} else if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case matchFileOrDirectory(line, NodeModules):
			fallthrough
		case matchFileOrDirectory(line, Dist):
			fallthrough
		case matchFileOrDirectory(line, Target):
			fallthrough
		case matchFileOrDirectory(line, Build):
			ignoredFiles = append(ignoredFiles, line)
		}
	}

	return ignoredFiles, nil
}

func matchFileOrDirectory(line string, folder string) bool {
	return strings.HasPrefix(line, folder) || strings.HasPrefix(line, "/" + folder) || strings.HasPrefix(line, "./" + folder)
}
