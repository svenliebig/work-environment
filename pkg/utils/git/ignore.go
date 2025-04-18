package git

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// returns all directories that are ignored by the .gitignore files
// and exist.
func FindAllIgnoredExistingDirectories(path string) ([]string, error) {
	gitignoreFiles, err := findAllGitignoreFiles(path)
	if err != nil {
		return nil, err
	}
	
	ignoredDirs := make(map[string]bool)
	ignoredPatterns := make(map[string][]string)
	result := []string{}
	
	// Read all gitignore files and collect their patterns
	for _, gitignorePath := range gitignoreFiles {
		dir := filepath.Dir(gitignorePath)
		file, err := os.Open(gitignorePath)
		if err != nil {
			return nil, err
		}
		
		scanner := bufio.NewScanner(file)
		var patterns []string
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			// Skip empty lines and comments
			if line != "" && !strings.HasPrefix(line, "#") {
				patterns = append(patterns, line)
			}
		}
		file.Close()
		
		ignoredPatterns[dir] = patterns
	}
	
	// First pass: Find all ignored top-level directories
	dirsToProcess := []string{path}
	
	for len(dirsToProcess) > 0 {
		currentDir := dirsToProcess[0]
		dirsToProcess = dirsToProcess[1:]
		
		entries, err := os.ReadDir(currentDir)
		if err != nil {
			continue // Skip directories we can't read
		}
		
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			
			entryPath := filepath.Join(currentDir, entry.Name())
			
			// Skip .git directory
			if entry.Name() == ".git" {
				continue
			}
			
			// Check if any parent directory is already marked as ignored
			// If so, skip this subdirectory
			isSubdirOfIgnored := false
			tmpDir := filepath.Dir(entryPath)
			for tmpDir != "." && tmpDir != "/" && tmpDir != "" {
				if ignoredDirs[tmpDir] {
					isSubdirOfIgnored = true
					break
				}
				tmpDir = filepath.Dir(tmpDir)
			}
			
			if isSubdirOfIgnored {
				continue
			}
			
			// Check if current directory matches node_modules pattern directly
			if entry.Name() == "node_modules" {
				ignoredDirs[entryPath] = true
				result = append(result, entryPath)
				continue
			}
			
			// Check if this directory is ignored by any gitignore pattern
			isIgnored := false
			for gitignoreDir, patterns := range ignoredPatterns {
				// Check if gitignore path is an ancestor of the current path or in the same directory
				if !isAncestorDir(gitignoreDir, entryPath) {
					continue
				}
				
				// Make path relative to the gitignore directory
				relPath, err := filepath.Rel(gitignoreDir, entryPath)
				if err != nil {
					continue
				}
				
				// Process regular patterns first, then negation patterns
				var regularPatterns []string
				var negationPatterns []string
				
				for _, pattern := range patterns {
					if strings.HasPrefix(pattern, "!") {
						negationPatterns = append(negationPatterns, pattern)
					} else {
						regularPatterns = append(regularPatterns, pattern)
					}
				}
				
				// Check against regular patterns
				matchedAnyPattern := false
				for _, pattern := range regularPatterns {
					if matchesGitignorePattern(relPath, pattern) || 
					// Handle simple patterns without slashes (match at any level)
					(!strings.Contains(pattern, "/") && strings.HasSuffix(entryPath, "/"+pattern)) ||
					(!strings.Contains(pattern, "/") && entry.Name() == pattern) {
						matchedAnyPattern = true
						break
					}
				}
				
				// If matched, check negation patterns
				if matchedAnyPattern {
					for _, pattern := range negationPatterns {
						patternWithoutNegation := pattern[1:]
						if matchesGitignorePattern(relPath, patternWithoutNegation) ||
						// Handle simple patterns without slashes (match at any level)
						(!strings.Contains(patternWithoutNegation, "/") && strings.HasSuffix(entryPath, "/"+patternWithoutNegation)) ||
						(!strings.Contains(patternWithoutNegation, "/") && entry.Name() == patternWithoutNegation) {
							matchedAnyPattern = false // Un-ignore
							break
						}
					}
					
					if matchedAnyPattern {
						isIgnored = true
						break
					}
				}
			}
			
			// If the directory is ignored, add it to results and don't process its subdirectories
			if isIgnored {
				ignoredDirs[entryPath] = true
				result = append(result, entryPath)
			} else {
				// Only process subdirectories of non-ignored directories
				dirsToProcess = append(dirsToProcess, entryPath)
			}
		}
	}
	
	return result, nil
}

func findAllGitignoreFiles(p string) ([]string, error) {
	var gitignoreFiles []string
	ignoredPatterns := make(map[string][]string)
	
	// First pass: collect all gitignore files and their patterns
	err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			return nil
		}

		if info.Name() == ".gitignore" {
			dir := filepath.Dir(path)
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			
			scanner := bufio.NewScanner(file)
			var patterns []string
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				// Skip empty lines and comments
				if line != "" && !strings.HasPrefix(line, "#") {
					patterns = append(patterns, line)
				}
			}
			
			ignoredPatterns[dir] = patterns
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	
	// Second pass: check each gitignore file against patterns from parent directories
	err = filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			return nil
		}

		if info.Name() == ".gitignore" {
			// Check if this gitignore file is ignored by any parent gitignore
			_, err := filepath.Rel(p, path)
			if err != nil {
				return err
			}
			
			isIgnored := false
			dir := filepath.Dir(path)
			
			// Check against patterns from all parent directories
			for patternDir, patterns := range ignoredPatterns {
				// Only check patterns from parent directories
				if !isParentDir(patternDir, dir) {
					continue
				}
				
				// Make the path relative to the pattern directory
				filePathRelToPatternDir, err := filepath.Rel(patternDir, path)
				if err != nil {
					continue
				}
				
				// First collect all regular patterns and negation patterns separately
				var regularPatterns []string
				var negationPatterns []string
				
				for _, pattern := range patterns {
					if strings.HasPrefix(pattern, "!") {
						negationPatterns = append(negationPatterns, pattern)
					} else {
						regularPatterns = append(regularPatterns, pattern)
					}
				}
				
				// Check against regular patterns first
				matchedAnyPattern := false
				for _, pattern := range regularPatterns {
					// Skip negation patterns for this initial check
					if strings.HasPrefix(pattern, "!") {
						continue
					}
					
					// Check if the file matches this pattern
					if matchesGitignorePattern(filePathRelToPatternDir, pattern) {
						matchedAnyPattern = true
						break
					}
				}
				
				// If matched by a regular pattern, check if it's negated by any negation pattern
				if matchedAnyPattern {
					// Apply negation patterns - they can "un-ignore" a file
					for _, pattern := range negationPatterns {
						// Negation patterns start with !, which was already checked above
						// but we'll verify again to be safe
						if !strings.HasPrefix(pattern, "!") {
							continue
						}
						
						// Remove the ! prefix for matching
						patternWithoutNegation := pattern[1:]
						
						// If the file matches the negation pattern, it should NOT be ignored
						// The ! in the pattern says "don't ignore this, even if other patterns would ignore it"
						if matchesGitignorePattern(filePathRelToPatternDir, patternWithoutNegation) {
							matchedAnyPattern = false // Un-ignore the file
							break
						}
					}
					
					if matchedAnyPattern {
						isIgnored = true
						break
					}
				}
			}
			
			if !isIgnored {
				gitignoreFiles = append(gitignoreFiles, path)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return gitignoreFiles, nil
}


// Checks if dir1 is a parent directory of dir2
func isParentDir(dir1, dir2 string) bool {
	rel, err := filepath.Rel(dir1, dir2)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..") && rel != "."
}

// Checks if dir1 is an ancestor directory of dir2 or the same directory
func isAncestorDir(dir1, dir2 string) bool {
	rel, err := filepath.Rel(dir1, dir2)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..") || rel == "."
}

// Basic implementation of gitignore pattern matching
// This is a simplified version and doesn't handle all gitignore pattern features
func matchesGitignorePattern(path, pattern string) bool {
	// Ignore empty patterns
	if pattern == "" {
		return false
	}

	// Handle negation (patterns starting with !)
	negate := false
	if strings.HasPrefix(pattern, "!") {
		negate = true
		pattern = pattern[1:]
		// If after removing the ! the pattern is empty, return false
		if pattern == "" {
			return false
		}
	}

	// Remove leading slash if present - it matches files only in the root
	if strings.HasPrefix(pattern, "/") {
		// This is a pattern that only matches at the root level
		pattern = pattern[1:]
	}

	// Trim trailing spaces
	pattern = strings.TrimSpace(pattern)
	
	// Skip comments
	if strings.HasPrefix(pattern, "#") {
		return false
	}
	
	// Check if it's a directory pattern (ending with /)
	isDirPattern := strings.HasSuffix(pattern, "/")
	if isDirPattern {
		pattern = pattern[:len(pattern)-1]
	}
	
	// Convert gitignore pattern to filepath.Match pattern
	// * matches anything except /
	// ** matches anything including /
	
	// Check for exact match first
	exactMatch := path == pattern || path+"/" == pattern || path == pattern+"/"
	
	// If not an exact match, try wildcard matching
	var matched bool
	if !exactMatch {
		// Simple case: Check if the path is exactly this pattern or starts with pattern/
		if strings.HasPrefix(path, pattern+"/") || path == pattern {
			matched = true
		} else {
			// Try to use filepath.Match for simple wildcard patterns
			// This won't handle ** properly but works for simple cases
			matched, _ = filepath.Match(pattern, path)
		}
	} else {
		matched = true
	}
	
	// For directory patterns, check if the path is or is inside a directory that matches
	if isDirPattern && !matched {
		// Check if the path is within a directory matching this pattern
		pathComponents := strings.Split(path, "/")
		for i := range pathComponents {
			partialPath := strings.Join(pathComponents[:i+1], "/")
			if partialPath == pattern {
				matched = true
				break
			}
		}
	}
	
	// Apply negation correctly
	if negate {
		return !matched
	}
	
	return matched
}
