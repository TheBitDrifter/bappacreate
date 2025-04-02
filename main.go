package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

//go:embed templates
var templateFS embed.FS

// Define file types that should be copied as-is (not processed for replacements)
var binaryExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".wav":  true,
	".mp3":  true,
	".ogg":  true,
}

// Structure to define common files and where they should be copied to in each template
type CommonFileMapping struct {
	SourcePath      string
	DestinationPath map[string]string // Map of template name → destination path within project
}

// Define the import paths for each template
var templateImportPaths = map[string]string{
	"platformer":            "github.com/TheBitDrifter/bappacreate/templates/platformer",
	"platformer-split":      "github.com/TheBitDrifter/bappacreate/templates/platformer-split",
	"platformer-ldtk":       "github.com/TheBitDrifter/bappacreate/templates/platformer-ldtk",
	"platformer-split-ldtk": "github.com/TheBitDrifter/bappacreate/templates/platformer-split-ldtk",
}

// Common import path that will be replaced in all files
var commonImportPattern = "github.com/TheBitDrifter/bappacreate/templates/common"

// Define all the shared files that should be moved from common to each template
var commonFiles = []CommonFileMapping{
	// Core Systems
	{
		SourcePath: "templates/common/coresystems/gravitysystem.go",
		DestinationPath: map[string]string{
			"platformer":            "coresystems/gravitysystem.go",
			"platformer-split":      "coresystems/gravitysystem.go",
			"platformer-ldtk":       "coresystems/gravitysystem.go",
			"platformer-split-ldtk": "coresystems/gravitysystem.go",
		},
	},
	{
		SourcePath: "templates/common/coresystems/frictionsystem.go",
		DestinationPath: map[string]string{
			"platformer":            "coresystems/frictionsystem.go",
			"platformer-split":      "coresystems/frictionsystem.go",
			"platformer-ldtk":       "coresystems/frictionsystem.go",
			"platformer-split-ldtk": "coresystems/frictionsystem.go",
		},
	},
	{
		SourcePath: "templates/common/coresystems/player_movement_system.go",
		DestinationPath: map[string]string{
			"platformer":            "coresystems/player_movement_system.go",
			"platformer-split":      "coresystems/player_movement_system.go",
			"platformer-ldtk":       "coresystems/player_movement_system.go",
			"platformer-split-ldtk": "coresystems/player_movement_system.go",
		},
	},
	{
		SourcePath: "templates/common/coresystems/player_block_collision_system.go",
		DestinationPath: map[string]string{
			"platformer":            "coresystems/player_block_collision_system.go",
			"platformer-split":      "coresystems/player_block_collision_system.go",
			"platformer-ldtk":       "coresystems/player_block_collision_system.go",
			"platformer-split-ldtk": "coresystems/player_block_collision_system.go",
		},
	},
	{
		SourcePath: "templates/common/coresystems/player_platform_collision_system.go",
		DestinationPath: map[string]string{
			"platformer":            "coresystems/player_platform_collision_system.go",
			"platformer-split":      "coresystems/player_platform_collision_system.go",
			"platformer-ldtk":       "coresystems/player_platform_collision_system.go",
			"platformer-split-ldtk": "coresystems/player_platform_collision_system.go",
		},
	},
	{
		SourcePath: "templates/common/coresystems/on_ground_clearing_system.go",
		DestinationPath: map[string]string{
			"platformer":            "coresystems/on_ground_clearing_system.go",
			"platformer-split":      "coresystems/on_ground_clearing_system.go",
			"platformer-ldtk":       "coresystems/on_ground_clearing_system.go",
			"platformer-split-ldtk": "coresystems/on_ground_clearing_system.go",
		},
	},
	{
		SourcePath: "templates/common/coresystems/ignore_platform_clearing_system.go",
		DestinationPath: map[string]string{
			"platformer":            "coresystems/ignore_platform_clearing_system.go",
			"platformer-split":      "coresystems/ignore_platform_clearing_system.go",
			"platformer-ldtk":       "coresystems/ignore_platform_clearing_system.go",
			"platformer-split-ldtk": "coresystems/ignore_platform_clearing_system.go",
		},
	},
	{
		SourcePath: "templates/common/coresystems/common.go",
		DestinationPath: map[string]string{
			"platformer":            "coresystems/common.go",
			"platformer-split":      "coresystems/common.go",
			"platformer-ldtk":       "coresystems/common.go",
			"platformer-split-ldtk": "coresystems/common.go",
		},
	},
	// Client Systems
	{
		SourcePath: "templates/common/clientsystems/camera_follower_system.go",
		DestinationPath: map[string]string{
			"platformer":            "clientsystems/camera_follower_system.go",
			"platformer-split":      "clientsystems/camera_follower_system.go",
			"platformer-ldtk":       "clientsystems/camera_follower_system.go",
			"platformer-split-ldtk": "clientsystems/camera_follower_system.go",
		},
	},
	{
		SourcePath: "templates/common/clientsystems/player_animation_system.go",
		DestinationPath: map[string]string{
			"platformer":            "clientsystems/player_animation_system.go",
			"platformer-split":      "clientsystems/player_animation_system.go",
			"platformer-ldtk":       "clientsystems/player_animation_system.go",
			"platformer-split-ldtk": "clientsystems/player_animation_system.go",
		},
	},
	{
		SourcePath: "templates/common/clientsystems/player_sound_system.go",
		DestinationPath: map[string]string{
			"platformer":            "clientsystems/player_sound_system.go",
			"platformer-split":      "clientsystems/player_sound_system.go",
			"platformer-ldtk":       "clientsystems/player_sound_system.go",
			"platformer-split-ldtk": "clientsystems/player_sound_system.go",
		},
	},
	{
		SourcePath: "templates/common/clientsystems/musicsystem.go",
		DestinationPath: map[string]string{
			"platformer":            "clientsystems/musicsystem.go",
			"platformer-split":      "clientsystems/musicsystem.go",
			"platformer-ldtk":       "clientsystems/musicsystem.go",
			"platformer-split-ldtk": "clientsystems/musicsystem.go",
		},
	},
	{
		SourcePath: "templates/common/clientsystems/collision_player_transfer_system.go",
		DestinationPath: map[string]string{
			"platformer":            "clientsystems/collision_player_transfer_system.go",
			"platformer-split":      "clientsystems/collision_player_transfer_system.go",
			"platformer-ldtk":       "clientsystems/collision_player_transfer_system.go",
			"platformer-split-ldtk": "clientsystems/collision_player_transfer_system.go",
		},
	},
	{
		SourcePath: "templates/common/clientsystems/scene_deactivation_system.go",
		DestinationPath: map[string]string{
			"platformer-split":      "clientsystems/scene_deactivation_system.go",
			"platformer-split-ldtk": "clientsystems/scene_deactivation_system.go",
		},
	},
	{
		SourcePath: "templates/common/clientsystems/common.go",
		DestinationPath: map[string]string{
			"platformer":            "clientsystems/common.go",
			"platformer-split":      "clientsystems/common.go",
			"platformer-ldtk":       "clientsystems/common.go",
			"platformer-split-ldtk": "clientsystems/common.go",
		},
	},
	// Components
	{
		SourcePath: "templates/common/components/components.go",
		DestinationPath: map[string]string{
			"platformer":            "components/components.go",
			"platformer-split":      "components/components.go",
			"platformer-ldtk":       "components/components.go",
			"platformer-split-ldtk": "components/components.go",
		},
	},
	{
		SourcePath: "templates/common/components/tags.go",
		DestinationPath: map[string]string{
			"platformer":            "components/tags.go",
			"platformer-split":      "components/tags.go",
			"platformer-ldtk":       "components/tags.go",
			"platformer-split-ldtk": "components/tags.go",
		},
	},
	// Other shared code
	{
		SourcePath: "templates/common/animations/animations.go",
		DestinationPath: map[string]string{
			"platformer":            "animations/animations.go",
			"platformer-split":      "animations/animations.go",
			"platformer-ldtk":       "animations/animations.go",
			"platformer-split-ldtk": "animations/animations.go",
		},
	},
	{
		SourcePath: "templates/common/sounds/sounds.go",
		DestinationPath: map[string]string{
			"platformer":            "sounds/sounds.go",
			"platformer-split":      "sounds/sounds.go",
			"platformer-ldtk":       "sounds/sounds.go",
			"platformer-split-ldtk": "sounds/sounds.go",
		},
	},
	{
		SourcePath: "templates/common/actions/actions.go",
		DestinationPath: map[string]string{
			"platformer":            "actions/actions.go",
			"platformer-split":      "actions/actions.go",
			"platformer-ldtk":       "actions/actions.go",
			"platformer-split-ldtk": "actions/actions.go",
		},
	},
	// Render Systems
	{
		SourcePath: "templates/common/rendersystems/common.go",
		DestinationPath: map[string]string{
			"platformer":            "rendersystems/common.go",
			"platformer-split":      "rendersystems/common.go",
			"platformer-ldtk":       "rendersystems/common.go",
			"platformer-split-ldtk": "rendersystems/common.go",
		},
	},
	{
		SourcePath: "templates/common/rendersystems/player_camera_prio_system.go",
		DestinationPath: map[string]string{
			"platformer-split":      "rendersystems/player_camera_prio_system.go",
			"platformer-split-ldtk": "rendersystems/player_camera_prio_system.go",
		},
	},
}

// Import path mappings - defines how to transform component imports for each template
var importPathMappings = map[string]map[string]string{
	"platformer": {
		"github.com/TheBitDrifter/bappacreate/templates/common/components":    "github.com/TheBitDrifter/bappacreate/templates/platformer/components",
		"github.com/TheBitDrifter/bappacreate/templates/common/actions":       "github.com/TheBitDrifter/bappacreate/templates/platformer/actions",
		"github.com/TheBitDrifter/bappacreate/templates/common/animations":    "github.com/TheBitDrifter/bappacreate/templates/platformer/animations",
		"github.com/TheBitDrifter/bappacreate/templates/common/sounds":        "github.com/TheBitDrifter/bappacreate/templates/platformer/sounds",
		"github.com/TheBitDrifter/bappacreate/templates/common/coresystems":   "github.com/TheBitDrifter/bappacreate/templates/platformer/coresystems",
		"github.com/TheBitDrifter/bappacreate/templates/common/clientsystems": "github.com/TheBitDrifter/bappacreate/templates/platformer/clientsystems",
		"github.com/TheBitDrifter/bappacreate/templates/common/rendersystems": "github.com/TheBitDrifter/bappacreate/templates/platformer/rendersystems",
	},
	"platformer-split": {
		"github.com/TheBitDrifter/bappacreate/templates/common/components":    "github.com/TheBitDrifter/bappacreate/templates/platformer-split/components",
		"github.com/TheBitDrifter/bappacreate/templates/common/actions":       "github.com/TheBitDrifter/bappacreate/templates/platformer-split/actions",
		"github.com/TheBitDrifter/bappacreate/templates/common/animations":    "github.com/TheBitDrifter/bappacreate/templates/platformer-split/animations",
		"github.com/TheBitDrifter/bappacreate/templates/common/sounds":        "github.com/TheBitDrifter/bappacreate/templates/platformer-split/sounds",
		"github.com/TheBitDrifter/bappacreate/templates/common/coresystems":   "github.com/TheBitDrifter/bappacreate/templates/platformer-split/coresystems",
		"github.com/TheBitDrifter/bappacreate/templates/common/clientsystems": "github.com/TheBitDrifter/bappacreate/templates/platformer-split/clientsystems",
		"github.com/TheBitDrifter/bappacreate/templates/common/rendersystems": "github.com/TheBitDrifter/bappacreate/templates/platformer-split/rendersystems",
	},
	"platformer-ldtk": {
		"github.com/TheBitDrifter/bappacreate/templates/common/components":    "github.com/TheBitDrifter/bappacreate/templates/platformer-ldtk/components",
		"github.com/TheBitDrifter/bappacreate/templates/common/actions":       "github.com/TheBitDrifter/bappacreate/templates/platformer-ldtk/actions",
		"github.com/TheBitDrifter/bappacreate/templates/common/animations":    "github.com/TheBitDrifter/bappacreate/templates/platformer-ldtk/animations",
		"github.com/TheBitDrifter/bappacreate/templates/common/sounds":        "github.com/TheBitDrifter/bappacreate/templates/platformer-ldtk/sounds",
		"github.com/TheBitDrifter/bappacreate/templates/common/coresystems":   "github.com/TheBitDrifter/bappacreate/templates/platformer-ldtk/coresystems",
		"github.com/TheBitDrifter/bappacreate/templates/common/clientsystems": "github.com/TheBitDrifter/bappacreate/templates/platformer-ldtk/clientsystems",
		"github.com/TheBitDrifter/bappacreate/templates/common/rendersystems": "github.com/TheBitDrifter/bappacreate/templates/platformer-ldtk/rendersystems",
	},
	"platformer-split-ldtk": {
		"github.com/TheBitDrifter/bappacreate/templates/common/components":    "github.com/TheBitDrifter/bappacreate/templates/platformer-split-ldtk/components",
		"github.com/TheBitDrifter/bappacreate/templates/common/actions":       "github.com/TheBitDrifter/bappacreate/templates/platformer-split-ldtk/actions",
		"github.com/TheBitDrifter/bappacreate/templates/common/animations":    "github.com/TheBitDrifter/bappacreate/templates/platformer-split-ldtk/animations",
		"github.com/TheBitDrifter/bappacreate/templates/common/sounds":        "github.com/TheBitDrifter/bappacreate/templates/platformer-split-ldtk/sounds",
		"github.com/TheBitDrifter/bappacreate/templates/common/coresystems":   "github.com/TheBitDrifter/bappacreate/templates/platformer-split-ldtk/coresystems",
		"github.com/TheBitDrifter/bappacreate/templates/common/clientsystems": "github.com/TheBitDrifter/bappacreate/templates/platformer-split-ldtk/clientsystems",
		"github.com/TheBitDrifter/bappacreate/templates/common/rendersystems": "github.com/TheBitDrifter/bappacreate/templates/platformer-split-ldtk/rendersystems",
	},
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	projectName := os.Args[1]
	template := "topdown" // Default template

	// Check if user specified a template
	if len(os.Args) >= 4 && os.Args[2] == "--template" {
		template = os.Args[3]
	}

	createProject(projectName, template)
}

func printUsage() {
	fmt.Println("Bappa Game Template Generator")
	fmt.Println("===============================")
	fmt.Println("Usage: bappacreate username/project-name [--template <template-name>]")
	fmt.Println()
	fmt.Println("This tool creates a new Bappa game project with the specified name.")
	fmt.Println("The username/ prefix is used to create the proper Go module path.")
	fmt.Println("Example: bappacreate johndoe/my-awesome-game --template platformer")
	fmt.Println()
	fmt.Println("Available templates:")
	fmt.Println("  platformer      - A simple platformer game")
	fmt.Println("  platformer-split - A platformer game with split-screen co-op support")
	fmt.Println("  platformer-ldtk - A platformer game with LDtk level editor support")
	fmt.Println("  platformer-split-ldtk - A platformer game with split-screen co-op and LDtk support")
	fmt.Println("  topdown         - A top-down perspective game (default)")
	fmt.Println("  topdown-split   - A top-down game with split-screen co-op support")
	fmt.Println("  sandbox         - An open sandbox game environment")
}

func createProject(projectName, templateName string) {
	// Check if the project name contains a username
	parts := strings.Split(projectName, "/")
	if len(parts) != 2 {
		fmt.Println("Error: Project name must be in the format 'username/project-name'")
		os.Exit(1)
	}

	username := parts[0]
	projectNameOnly := parts[1]

	// Clean the project name for Go module naming
	moduleName := strings.ToLower(strings.ReplaceAll(projectNameOnly, " ", "-"))

	// Create project directory (just the project part, not with username)
	err := os.MkdirAll(projectNameOnly, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Creating new Bappa game project: %s (using %s template)\n", projectNameOnly, templateName)

	// The template import path that will be replaced
	templateImportPath := fmt.Sprintf("github.com/TheBitDrifter/bappacreate/templates/%s", templateName)

	// Make sure the template exists
	templatePath := filepath.Join("templates", templateName)
	_, err = fs.Stat(templateFS, templatePath)
	if err != nil {
		fmt.Printf("Error: Template '%s' not found\n", templateName)
		os.Exit(1)
	}

	// Create additional directories that might not be in the templates
	os.MkdirAll(filepath.Join(projectNameOnly, "assets", "images"), 0755)
	os.MkdirAll(filepath.Join(projectNameOnly, "assets", "sounds"), 0755)

	// Create additional directories for split-screen co-op templates
	if strings.HasSuffix(templateName, "-split") {
		os.MkdirAll(filepath.Join(projectNameOnly, "players"), 0755)
		os.MkdirAll(filepath.Join(projectNameOnly, "splitscreen"), 0755)
	}

	// Process all template files
	err = fs.WalkDir(templateFS, templatePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root template directory itself
		if path == templatePath {
			return nil
		}

		// Create relative path from the template directory
		relPath, err := filepath.Rel(templatePath, path)
		if err != nil {
			return err
		}

		// Create the target path in the new project using projectNameOnly
		targetPath := filepath.Join(projectNameOnly, relPath)

		// Create directories
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// Check if this file should be skipped because it will be created from common template
		skipFile := false
		for _, mapping := range commonFiles {
			destPath, exists := mapping.DestinationPath[templateName]
			if exists && relPath == destPath {
				skipFile = true
				break
			}
		}

		if skipFile {
			fmt.Printf("Skipping: %s (will be created from common template)\n", targetPath)
			return nil
		}

		// Handle the file appropriately
		ext := strings.ToLower(filepath.Ext(path))

		if binaryExtensions[ext] {
			// Copy binary files directly
			return copyBinaryFile(path, targetPath)
		} else {
			// Process text files with replacement
			return processTextFile(path, targetPath, username, moduleName, templateImportPath, templateName)
		}
	})
	if err != nil {
		fmt.Printf("Error creating project: %v\n", err)
		os.Exit(1)
	}

	// Copy common files that should replace template-specific ones
	fmt.Println("\nProcessing common template files...")
	err = copyCommonFiles(projectNameOnly, templateName, username, moduleName)
	if err != nil {
		fmt.Printf("Error processing common files: %v\n", err)
		os.Exit(1)
	}

	// Calculate the full module path with username
	modulePath := fmt.Sprintf("github.com/%s/%s", username, moduleName)

	// Initialize Go module and dependencies
	initGoModule(projectNameOnly, modulePath)

	fmt.Printf("\nSuccessfully created Bappa game project: %s\n", projectNameOnly)
	fmt.Printf("\nTo run your game:\n")
	fmt.Printf("  cd %s\n", projectNameOnly)
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  go run .\n")
}

// Copy common template files to the project with appropriate processing
func copyCommonFiles(projectDir, templateName, username, moduleName string) error {
	// Full project import path
	projectImportPath := fmt.Sprintf("github.com/%s/%s", username, moduleName)

	for _, mapping := range commonFiles {
		// Check if this template needs this file
		destPath, exists := mapping.DestinationPath[templateName]
		if !exists {
			continue // Skip files not needed for this template
		}

		// Full path in the project
		fullDestPath := filepath.Join(projectDir, destPath)

		// Ensure the directory exists
		if err := os.MkdirAll(filepath.Dir(fullDestPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", filepath.Dir(fullDestPath), err)
		}

		// Read the source file from the common template
		content, err := templateFS.ReadFile(mapping.SourcePath)
		if err != nil {
			// Check if file exists in template-specific location instead
			templateSpecificPath := strings.Replace(mapping.SourcePath, "templates/common", "templates/"+templateName, 1)
			fmt.Printf("Common file not found, trying template-specific file: %s\n", templateSpecificPath)
			content, err = templateFS.ReadFile(templateSpecificPath)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %v", mapping.SourcePath, err)
			}
		}

		// Process the content
		processedContent := processCommonFile(string(content), mapping.SourcePath, destPath, templateName, projectImportPath)

		// Write the processed content
		if err := os.WriteFile(fullDestPath, []byte(processedContent), 0644); err != nil {
			return fmt.Errorf("failed to write file %s: %v", fullDestPath, err)
		}

		fmt.Printf("  Created: %s\n", fullDestPath)
	}

	return nil
}

// Process a common file's content to make it template-specific
func processCommonFile(content, sourcePath, destPath, templateName, projectImportPath string) string {
	// 1. Determine target package name from destination
	destDir := filepath.Dir(destPath)
	packageName := filepath.Base(destDir) // E.g., "coresystems" from "coresystems/gravitysystem.go"

	// 2. Replace package declaration
	packagePattern := regexp.MustCompile(`(?m)^package\s+\w+`)
	content = packagePattern.ReplaceAllString(content, "package "+packageName)

	// 3. Replace imports for this specific template
	templateImportPath := templateImportPaths[templateName]

	// First replace any imports from common to template-specific paths
	templateSpecificMappings, hasMappings := importPathMappings[templateName]
	if hasMappings {
		for commonPath, tempPath := range templateSpecificMappings {
			content = strings.ReplaceAll(content, commonPath, tempPath)
		}
	}

	// Common path replacement fallback
	content = strings.ReplaceAll(content, commonImportPattern, templateImportPath)

	// 4. Replace template imports with final project imports
	content = strings.ReplaceAll(content, templateImportPath, projectImportPath)

	return content
}

func copyBinaryFile(sourcePath, targetPath string) error {
	fmt.Println("Copying asset:", targetPath)

	// Read binary file
	data, err := templateFS.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	// Write to target path
	return os.WriteFile(targetPath, data, 0644)
}

func processTextFile(sourcePath, targetPath, username, moduleName, templateImportPath, templateName string) error {
	fmt.Println("Creating:", targetPath)

	// Calculate the new import path including username
	projectImportPath := fmt.Sprintf("github.com/%s/%s", username, moduleName)

	// Read file content
	content, err := templateFS.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	fileContent := string(content)

	// First, handle any imports from common to template-specific paths
	templateSpecificMappings, hasMappings := importPathMappings[templateName]
	if hasMappings {
		for commonPath, tempPath := range templateSpecificMappings {
			fileContent = strings.ReplaceAll(fileContent, commonPath, tempPath)
		}
	}

	// Common import path replacement
	fileContent = strings.ReplaceAll(fileContent, commonImportPattern, templateImportPath)

	// Then replace all template imports with project imports
	fileContent = strings.ReplaceAll(fileContent, templateImportPath, projectImportPath)

	// Special handling for go.mod file
	if filepath.Base(targetPath) == "go.mod" {
		lines := strings.Split(fileContent, "\n")
		for i, line := range lines {
			if strings.HasPrefix(line, "module ") {
				lines[i] = "module " + projectImportPath
				break
			}
		}
		fileContent = strings.Join(lines, "\n")
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	// Write processed content
	return os.WriteFile(targetPath, []byte(fileContent), 0644)
}

func initGoModule(projectName, modulePath string) {
	fmt.Println("\nSetting up Go module...")

	// Change to project directory
	currentDir, _ := os.Getwd()
	projectDir := filepath.Join(currentDir, projectName)
	err := os.Chdir(projectDir)
	if err != nil {
		fmt.Printf("Error changing to project directory: %v\n", err)
		return
	}
	defer os.Chdir(currentDir)

	// Initialize go.mod with the correct module name
	runCommand("go", "mod", "init", modulePath)

	// Get required dependencies
	runCommand("go", "get", "github.com/TheBitDrifter/bappa/coldbrew@latest")
	runCommand("go", "get", "github.com/TheBitDrifter/bappa/blueprint@latest")
	runCommand("go", "get", "github.com/TheBitDrifter/bappa/warehouse@latest")
	runCommand("go", "get", "github.com/TheBitDrifter/bappa/tteokbokki@latest")
	runCommand("go", "get", "github.com/TheBitDrifter/bappa/table@latest")
}

func runCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Running: %s %s\n", command, strings.Join(args, " "))

	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: Command failed: %v\n", err)
	}
}
