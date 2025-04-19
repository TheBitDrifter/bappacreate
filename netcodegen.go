package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings" // Add missing imports if necessary (e.g., binaryExtensions map definition)
)

// handleNetcodeTemplate handles the platformer-netcode template using direct filesystem copy.
func handleNetcodeTemplate(projectName, username, moduleName string) error {
	// Calculate the full module path with username
	projectImportPath := fmt.Sprintf("github.com/%s/%s", username, moduleName)

	// Determine source path relative to executable
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)
	// Templates directory is located next to the executable
	sourceBaseDir := filepath.Join(exeDir, "templates", "platformer-netcode")

	// Verify source directory exists
	sourceInfo, err := os.Stat(sourceBaseDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("source template directory not found relative to executable (%s): %s", exeDir, sourceBaseDir)
	} else if err != nil {
		return fmt.Errorf("error stating source template directory %s: %w", sourceBaseDir, err)
	} else if !sourceInfo.IsDir() {
		return fmt.Errorf("source template path is not a directory: %s", sourceBaseDir)
	}

	// Create project directory
	err = os.MkdirAll(projectName, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("error creating project directory %s: %w", projectName, err)
	}

	// Define subdirectories for the netcode template
	subdirs := []string{"client", "server", "shared", "sharedclient", "bot", "standalone"}

	// Create subdirectories
	for _, subdir := range subdirs {
		targetSubdirPath := filepath.Join(projectName, subdir)
		err := os.MkdirAll(targetSubdirPath, 0755)
		if err != nil {
			return fmt.Errorf("error creating target subdirectory %s: %w", targetSubdirPath, err)
		}
	}

	// Create asset directories
	sharedClientAssets := filepath.Join(projectName, "sharedclient", "assets")
	os.MkdirAll(filepath.Join(sharedClientAssets, "images"), 0755)
	os.MkdirAll(filepath.Join(sharedClientAssets, "sounds"), 0755)

	// Define template import path for replacement
	templateImportPath := "github.com/TheBitDrifter/bappacreate/templates/platformer-netcode"

	// Process all template files
	var filesProcessedCount int
	var walkStarted bool

	err = filepath.WalkDir(sourceBaseDir, func(sourcePath string, d os.DirEntry, walkErr error) error {
		if !walkStarted {
			walkStarted = true
		}

		// Handle walk errors
		if walkErr != nil {
			if d != nil && d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Get relative path from source dir
		relPath, err := filepath.Rel(sourceBaseDir, sourcePath)
		if err != nil {
			return fmt.Errorf("internal error: failed to get relative path for %q from %q: %w", sourcePath, sourceBaseDir, err)
		}

		// Skip root directory
		if relPath == "." {
			return nil
		}

		// Create target path
		targetPath := filepath.Join(projectName, relPath)

		// Create directories recursively
		if d.IsDir() {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create target directory %q: %w", targetPath, err)
			}
			return nil
		}

		// Skip hidden files
		if strings.HasPrefix(d.Name(), ".") || strings.HasPrefix(d.Name(), "_") {
			return nil
		}

		filesProcessedCount++
		// fmt.Println("Creating:", targetPath)

		// Process file based on type
		ext := strings.ToLower(filepath.Ext(sourcePath))

		if binaryExtensions[ext] {
			// Copy binary files directly
			sourceFile, err := os.Open(sourcePath)
			if err != nil {
				return nil
			}
			defer sourceFile.Close()

			destFile, err := os.Create(targetPath)
			if err != nil {
				return nil
			}
			defer destFile.Close()

			_, err = io.Copy(destFile, sourceFile)
			if err != nil {
				return nil
			}
		} else {
			// Process text files with replacements
			contentBytes, err := os.ReadFile(sourcePath)
			if err != nil {
				return nil
			}
			fileContent := string(contentBytes)

			// Replace import paths
			fileContent = strings.ReplaceAll(fileContent, templateImportPath, projectImportPath)

			// Special handling for go.mod
			if filepath.Base(sourcePath) == "go.mod" {
				pathParts := strings.Split(relPath, string(filepath.Separator))
				if len(pathParts) > 1 {
					currentSubdir := pathParts[0]
					lines := strings.Split(fileContent, "\n")
					for i, line := range lines {
						if strings.HasPrefix(line, "module ") {
							subModulePath := fmt.Sprintf("%s/%s", projectImportPath, currentSubdir)
							lines[i] = "module " + subModulePath
							break
						}
					}
					fileContent = strings.Join(lines, "\n")
				}
			}

			// Write processed content
			err = os.WriteFile(targetPath, []byte(fileContent), 0644)
			if err != nil {
				return nil
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error during filesystem walk/copy: %w", err)
	}

	if filesProcessedCount == 0 {
		return fmt.Errorf("no template files processed from filesystem")
	}
	logNetSuc(projectName)
	return nil
}

func logNetSuc(projName string) {
	log.Println("-------------------------------------------------")
	log.Println("Created bappa netcode project!")
	log.Println("-------------------------------------------------")
	log.Println("To run networked:")

	log.Println(fmt.Sprintf("cd %s/server", projName))
	log.Println("go mod tidy")
	log.Println("go run .")
	log.Println("&&")
	log.Println(fmt.Sprintf("cd %s/client", projName))
	log.Println("go mod tidy")
	log.Println("go run .")
	log.Println("-------------------------------------------------")
	log.Println("To run single player/standalone:")
	log.Println(fmt.Sprintf("cd %s/standalone", projName))
	log.Println("go mod tidy")
	log.Println("go run .")
	log.Println("-------------------------------------------------")
}
