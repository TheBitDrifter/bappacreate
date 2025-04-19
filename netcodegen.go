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

	// Find the netcode template directory
	sourceBaseDir, err := findNetcodeTemplate()
	if err != nil {
		// If template wasn't found, show a helpful error message
		fmt.Println("Error: The netcode template was not found.")
		fmt.Println("To use the netcode template, you need to clone the repository:")
		fmt.Println("  git clone https://github.com/TheBitDrifter/bappacreate")
		fmt.Println("  cd bappacreate")
		fmt.Println("  go build .")
		fmt.Println("  ./bappacreate username/project --template platformer-netcode")
		return err
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

func findNetcodeTemplate() (string, error) {
	// Try multiple locations for the netcode template
	possibleLocations := []string{}

	// Check relative to executable first
	exePath, err := os.Executable()
	if err == nil {
		exeDir := filepath.Dir(exePath)
		possibleLocations = append(possibleLocations, filepath.Join(exeDir, "templates", "platformer-netcode"))
	}

	// Check common installation directories
	homeDir, err := os.UserHomeDir()
	if err == nil {
		// Unix/Linux locations
		possibleLocations = append(possibleLocations, filepath.Join(homeDir, ".local", "share", "bappacreate", "templates", "platformer-netcode"))
		// macOS location
		possibleLocations = append(possibleLocations, filepath.Join(homeDir, "Library", "Application Support", "bappacreate", "templates", "platformer-netcode"))
		// Windows location
		possibleLocations = append(possibleLocations, filepath.Join(homeDir, "AppData", "Local", "bappacreate", "templates", "platformer-netcode"))
	}

	// Check relative to current directory
	currentDir, err := os.Getwd()
	if err == nil {
		possibleLocations = append(possibleLocations, filepath.Join(currentDir, "templates", "platformer-netcode"))
		// Check if we're in the project dir
		possibleLocations = append(possibleLocations, filepath.Join(currentDir, "..", "templates", "platformer-netcode"))
	}

	// Check GOPATH location
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		possibleLocations = append(possibleLocations, filepath.Join(gopath, "src", "github.com", "TheBitDrifter", "bappacreate", "templates", "platformer-netcode"))
	}

	// Try each location
	for _, location := range possibleLocations {
		info, err := os.Stat(location)
		if err == nil && info.IsDir() {
			return location, nil
		}
	}

	return "", fmt.Errorf("netcode template not found in any of the expected locations")
}
