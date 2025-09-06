package main

import (
	"crypto/rand"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed templates/*
var templatesFS embed.FS

var UserConfigDirectory = os.UserConfigDir

var verbosity int

type VaultConfig struct {
	Path string `json:"path"`
	Ts   int64  `json:"ts"`
	Open bool   `json:"open"`
}

type ObsidianConfig struct {
	Vaults      map[string]VaultConfig `json:"vaults"`
	OpenSchemes map[string]bool        `json:"openSchemes"`
}

var rootCmd = &cobra.Command{
	Use:   "petiteplatypus",
	Short: "Generate Obsidian vault scaffolding",
	Long:  "A CLI tool to generate Obsidian vault directory structure and configuration",
}

var generateCmd = &cobra.Command{
	Use:   "generate [vault-path]",
	Short: "Generate a new Obsidian vault",
	Args:  cobra.ExactArgs(1),
	RunE:  generateVault,
}

func init() {
	rootCmd.PersistentFlags().CountVarP(&verbosity, "verbose", "v", "increase verbosity level (use multiple times for more verbose output)")
	rootCmd.AddCommand(generateCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func debugLog(level int, format string, args ...interface{}) {
	if verbosity >= level {
		prefix := fmt.Sprintf("[DEBUG%d] ", level)
		log.Printf(prefix+format, args...)
	}
}

func generateVault(cmd *cobra.Command, args []string) error {
	vaultPath := args[0]
	debugLog(1, "Starting vault generation for path: %s", vaultPath)
	debugLog(2, "Verbosity level: %d", verbosity)

	// Create absolute path
	debugLog(1, "Converting to absolute path")
	absVaultPath, err := filepath.Abs(vaultPath)
	if err != nil {
		debugLog(1, "Failed to get absolute path: %v", err)
		return fmt.Errorf("failed to get absolute path: %w", err)
	}
	debugLog(1, "Absolute vault path: %s", absVaultPath)

	// Create vault directory
	debugLog(1, "Creating vault directory: %s", absVaultPath)
	if err := os.MkdirAll(absVaultPath, 0o755); err != nil {
		debugLog(1, "Failed to create vault directory: %v", err)
		return fmt.Errorf("failed to create vault directory: %w", err)
	}
	debugLog(2, "Successfully created vault directory")

	// Create .obsidian directory
	obsidianDir := filepath.Join(absVaultPath, ".obsidian")
	debugLog(1, "Creating .obsidian directory: %s", obsidianDir)
	if err := os.MkdirAll(obsidianDir, 0o755); err != nil {
		debugLog(1, "Failed to create .obsidian directory: %v", err)
		return fmt.Errorf("failed to create .obsidian directory: %w", err)
	}
	debugLog(2, "Successfully created .obsidian directory")

	// Generate vault ID
	debugLog(1, "Generating vault ID")
	vaultID, err := generateVaultID()
	if err != nil {
		debugLog(1, "Failed to generate vault ID: %v", err)
		return fmt.Errorf("failed to generate vault ID: %w", err)
	}
	debugLog(1, "Generated vault ID: %s", vaultID)

	// Create obsidian config files
	debugLog(1, "Creating Obsidian config files")
	if err := createObsidianFiles(obsidianDir); err != nil {
		debugLog(1, "Failed to create obsidian config files: %v", err)
		return fmt.Errorf("failed to create obsidian config files: %w", err)
	}
	debugLog(2, "Successfully created Obsidian config files")

	// Create initial markdown files
	debugLog(1, "Creating initial markdown files")
	if err := createInitialFiles(absVaultPath); err != nil {
		debugLog(1, "Failed to create initial files: %v", err)
		return fmt.Errorf("failed to create initial files: %w", err)
	}
	debugLog(2, "Successfully created initial markdown files")

	// Update global obsidian.json
	debugLog(1, "Updating global Obsidian configuration")
	if err := updateGlobalConfig(absVaultPath, vaultID); err != nil {
		debugLog(1, "Failed to update global config: %v", err)
		return fmt.Errorf("failed to update global config: %w", err)
	}
	debugLog(2, "Successfully updated global configuration")

	fmt.Printf("Successfully created vault at: %s\n", absVaultPath)
	fmt.Printf("Vault ID: %s\n", vaultID)
	debugLog(1, "Vault generation completed successfully")

	return nil
}

func generateVaultID() (string, error) {
	debugLog(2, "Generating cryptographically secure random vault ID")
	// Generate 8 random bytes and convert to hex
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		debugLog(2, "Failed to generate random bytes: %v", err)
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	vaultID := fmt.Sprintf("%x", bytes)
	debugLog(2, "Generated vault ID: %s", vaultID)
	return vaultID, nil
}

func loadEmbeddedTemplate(filename string) (string, error) {
	debugLog(3, "Loading embedded template file: %s", filename)
	content, err := templatesFS.ReadFile(filename)
	if err != nil {
		debugLog(2, "Failed to read embedded template file %s: %v", filename, err)
		return "", fmt.Errorf("failed to read embedded template file %s: %w", filename, err)
	}
	debugLog(3, "Successfully loaded embedded template file: %s", filename)
	return string(content), nil
}

func createObsidianFiles(obsidianDir string) error {
	debugLog(2, "Defining Obsidian config files from embedded templates")
	templateFiles := map[string]string{
		"app.json":          "templates/app.json",
		"appearance.json":   "templates/appearance.json",
		"core-plugins.json": "templates/core-plugins.json",
		"graph.json":        "templates/graph.json",
		"workspace.json":    "templates/workspace.json",
	}

	debugLog(2, "Writing %d Obsidian config files", len(templateFiles))
	for filename, templatePath := range templateFiles {
		debugLog(3, "Loading embedded template for: %s", filename)
		content, err := loadEmbeddedTemplate(templatePath)
		if err != nil {
			return err
		}

		filePath := filepath.Join(obsidianDir, filename)
		debugLog(3, "Writing file: %s", filePath)
		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			debugLog(2, "Failed to write %s: %v", filename, err)
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}
		debugLog(3, "Successfully wrote file: %s", filename)
	}

	return nil
}

func createInitialFiles(vaultPath string) error {
	debugLog(2, "Creating initial markdown files")
	templateFiles := map[string]string{
		"Welcome.md": "templates/Welcome.md",
	}

	debugLog(2, "Writing %d initial files", len(templateFiles))
	for filename, templatePath := range templateFiles {
		debugLog(3, "Loading embedded template for: %s", filename)
		content, err := loadEmbeddedTemplate(templatePath)
		if err != nil {
			return err
		}

		filePath := filepath.Join(vaultPath, filename)
		debugLog(3, "Writing initial file: %s", filePath)
		if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
			debugLog(2, "Failed to write %s: %v", filename, err)
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}
		debugLog(3, "Successfully wrote initial file: %s", filename)
	}

	return nil
}

func updateGlobalConfig(vaultPath, vaultID string) error {
	debugLog(2, "Getting user config directory")
	configDir, err := UserConfigDirectory()
	if err != nil {
		debugLog(2, "Failed to get user config directory: %v", err)
		return fmt.Errorf("failed to get user config directory: %w", err)
	}
	debugLog(2, "User config directory: %s", configDir)

	obsidianConfigPath := filepath.Join(configDir, "obsidian", "obsidian.json")
	debugLog(1, "Global Obsidian config path: %s", obsidianConfigPath)

	// Create obsidian config directory if it doesn't exist
	configDirPath := filepath.Dir(obsidianConfigPath)
	debugLog(2, "Ensuring config directory exists: %s", configDirPath)
	if err := os.MkdirAll(configDirPath, 0o755); err != nil {
		debugLog(2, "Failed to create obsidian config directory: %v", err)
		return fmt.Errorf("failed to create obsidian config directory: %w", err)
	}
	debugLog(3, "Config directory created/verified")

	// Read existing config or create new one
	debugLog(2, "Reading existing global config")
	var config ObsidianConfig
	if data, err := os.ReadFile(obsidianConfigPath); err == nil {
		debugLog(3, "Found existing config file, parsing JSON")
		if err := json.Unmarshal(data, &config); err != nil {
			debugLog(2, "Failed to parse existing config: %v", err)
			return fmt.Errorf("failed to parse existing config: %w", err)
		}
		debugLog(3, "Successfully parsed existing config with %d vaults", len(config.Vaults))
	} else if !errors.Is(err, os.ErrNotExist) {
		debugLog(2, "Failed to read config file: %v", err)
		return fmt.Errorf("failed to read config file: %w", err)
	} else {
		debugLog(3, "No existing config file found, creating new one")
	}

	// Initialize maps if they're nil
	if config.Vaults == nil {
		debugLog(3, "Initializing empty vaults map")
		config.Vaults = make(map[string]VaultConfig)
	}
	if config.OpenSchemes == nil {
		debugLog(3, "Initializing default open schemes")
		config.OpenSchemes = map[string]bool{
			"vscode":           true,
			"chrome-extension": true,
		}
	}

	// Add new vault
	debugLog(2, "Adding new vault to config: %s -> %s", vaultID, vaultPath)
	config.Vaults[vaultID] = VaultConfig{
		Path: vaultPath,
		Ts:   1757173820641, // Using timestamp from example
		Open: true,
	}
	debugLog(3, "Vault added to config")

	// Write updated config
	debugLog(2, "Marshaling updated config to JSON")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		debugLog(2, "Failed to marshal config: %v", err)
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	debugLog(3, "Config marshaled successfully")

	debugLog(2, "Writing updated config to: %s", obsidianConfigPath)
	if err := os.WriteFile(obsidianConfigPath, data, 0o644); err != nil {
		debugLog(2, "Failed to write config file: %v", err)
		return fmt.Errorf("failed to write config file: %w", err)
	}
	debugLog(3, "Global config file written successfully")

	return nil
}
