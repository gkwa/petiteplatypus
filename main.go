package main

import (
   "encoding/json"
   "errors"
   "fmt"
   "os"
   "path/filepath"

   "github.com/spf13/cobra"
)

var UserConfigDirectory = os.UserConfigDir

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
   Use:   "obsidian-vault-generator",
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
   rootCmd.AddCommand(generateCmd)
}

func main() {
   if err := rootCmd.Execute(); err != nil {
   	fmt.Println(err)
   	os.Exit(1)
   }
}

func generateVault(cmd *cobra.Command, args []string) error {
   vaultPath := args[0]

   // Create absolute path
   absVaultPath, err := filepath.Abs(vaultPath)
   if err != nil {
   	return fmt.Errorf("failed to get absolute path: %w", err)
   }

   // Create vault directory
   if err := os.MkdirAll(absVaultPath, 0755); err != nil {
   	return fmt.Errorf("failed to create vault directory: %w", err)
   }

   // Create .obsidian directory
   obsidianDir := filepath.Join(absVaultPath, ".obsidian")
   if err := os.MkdirAll(obsidianDir, 0755); err != nil {
   	return fmt.Errorf("failed to create .obsidian directory: %w", err)
   }

   // Generate vault ID
   vaultID, err := generateVaultID()
   if err != nil {
   	return fmt.Errorf("failed to generate vault ID: %w", err)
   }

   // Create obsidian config files
   if err := createObsidianFiles(obsidianDir); err != nil {
   	return fmt.Errorf("failed to create obsidian config files: %w", err)
   }

   // Create initial markdown files
   if err := createInitialFiles(absVaultPath); err != nil {
   	return fmt.Errorf("failed to create initial files: %w", err)
   }

   // Update global obsidian.json
   if err := updateGlobalConfig(absVaultPath, vaultID); err != nil {
   	return fmt.Errorf("failed to update global config: %w", err)
   }

   fmt.Printf("Successfully created vault at: %s\n", absVaultPath)
   fmt.Printf("Vault ID: %s\n", vaultID)

   return nil
}

func generateVaultID() (string, error) {
   // Generate 8 random bytes and convert to hex
   bytes := make([]byte, 8)
   file, err := os.Open("/dev/urandom")
   if err != nil {
   	return "", err
   }
   defer file.Close()

   _, err = file.Read(bytes)
   if err != nil {
   	return "", err
   }

   return fmt.Sprintf("%x", bytes), nil
}

func createObsidianFiles(obsidianDir string) error {
   files := map[string]string{
   	"app.json": "{}",
   	"appearance.json": "{}",
   	"core-plugins.json": `{
 "file-explorer": true,
 "global-search": true,
 "switcher": true,
 "graph": true,
 "backlink": true,
 "canvas": true,
 "outgoing-link": true,
 "tag-pane": true,
 "footnotes": false,
 "properties": false,
 "page-preview": true,
 "daily-notes": true,
 "templates": true,
 "note-composer": true,
 "command-palette": true,
 "slash-command": false,
 "editor-status": true,
 "bookmarks": true,
 "markdown-importer": false,
 "zk-prefixer": false,
 "random-note": false,
 "outline": true,
 "word-count": true,
 "slides": false,
 "audio-recorder": false,
 "workspaces": false,
 "file-recovery": true,
 "publish": false,
 "sync": true,
 "bases": true,
 "webviewer": false
}`,
   	"graph.json": `{
 "collapse-filter": true,
 "search": "",
 "showTags": false,
 "showAttachments": false,
 "hideUnresolved": false,
 "showOrphans": true,
 "collapse-color-groups": true,
 "colorGroups": [],
 "collapse-display": true,
 "showArrow": false,
 "textFadeMultiplier": 0,
 "nodeSizeMultiplier": 1,
 "lineSizeMultiplier": 1,
 "collapse-forces": true,
 "centerStrength": 0.518713248970312,
 "repelStrength": 10,
 "linkStrength": 1,
 "linkDistance": 250,
 "scale": 1,
 "close": true
}`,
   	"workspace.json": `{
 "main": {
   "id": "d1e382edbf87edce",
   "type": "split",
   "children": [
     {
       "id": "f5b0ca4f913860db",
       "type": "tabs",
       "children": [
         {
           "id": "016d933bbc1cb39d",
           "type": "leaf",
           "state": {
             "type": "markdown",
             "state": {
               "file": "Welcome.md",
               "mode": "source",
               "source": false
             },
             "icon": "lucide-file",
             "title": "Welcome"
           }
         }
       ],
       "currentTab": 0
     }
   ],
   "direction": "vertical"
 },
 "left": {
   "id": "9de89e28bf65b7ae",
   "type": "split",
   "children": [
     {
       "id": "f1448db9adb9b5b5",
       "type": "tabs",
       "children": [
         {
           "id": "b66419c56a54b47b",
           "type": "leaf",
           "state": {
             "type": "file-explorer",
             "state": {
               "sortOrder": "alphabetical",
               "autoReveal": false
             },
             "icon": "lucide-folder-closed",
             "title": "Files"
           }
         },
         {
           "id": "e7e94e8543f577d4",
           "type": "leaf",
           "state": {
             "type": "search",
             "state": {
               "query": "",
               "matchingCase": false,
               "explainSearch": false,
               "collapseAll": false,
               "extraContext": false,
               "sortOrder": "alphabetical"
             },
             "icon": "lucide-search",
             "title": "Search"
           }
         },
         {
           "id": "cf8128b2ee5a6d15",
           "type": "leaf",
           "state": {
             "type": "bookmarks",
             "state": {},
             "icon": "lucide-bookmark",
             "title": "Bookmarks"
           }
         }
       ]
     }
   ],
   "direction": "horizontal",
   "width": 300
 },
 "right": {
   "id": "021487bd242948c6",
   "type": "split",
   "children": [
     {
       "id": "19769f2b7b7fa15d",
       "type": "tabs",
       "children": [
         {
           "id": "be8d0b8bcf32a675",
           "type": "leaf",
           "state": {
             "type": "backlink",
             "state": {
               "file": "Welcome.md",
               "collapseAll": false,
               "extraContext": false,
               "sortOrder": "alphabetical",
               "showSearch": false,
               "searchQuery": "",
               "backlinkCollapsed": false,
               "unlinkedCollapsed": true
             },
             "icon": "links-coming-in",
             "title": "Backlinks for Welcome"
           }
         },
         {
           "id": "77a759e6bd5cbf8b",
           "type": "leaf",
           "state": {
             "type": "outgoing-link",
             "state": {
               "file": "Welcome.md",
               "linksCollapsed": false,
               "unlinkedCollapsed": true
             },
             "icon": "links-going-out",
             "title": "Outgoing links from Welcome"
           }
         },
         {
           "id": "6a0af20ca5bae924",
           "type": "leaf",
           "state": {
             "type": "tag",
             "state": {
               "sortOrder": "frequency",
               "useHierarchy": true,
               "showSearch": false,
               "searchQuery": ""
             },
             "icon": "lucide-tags",
             "title": "Tags"
           }
         },
         {
           "id": "ed6dcab8808761f7",
           "type": "leaf",
           "state": {
             "type": "outline",
             "state": {
               "file": "Welcome.md",
               "followCursor": false,
               "showSearch": false,
               "searchQuery": ""
             },
             "icon": "lucide-list",
             "title": "Outline of Welcome"
           }
         }
       ]
     }
   ],
   "direction": "horizontal",
   "width": 300,
   "collapsed": true
 },
 "left-ribbon": {
   "hiddenItems": {
     "switcher:Open quick switcher": false,
     "graph:Open graph view": false,
     "canvas:Create new canvas": false,
     "daily-notes:Open today's daily note": false,
     "templates:Insert template": false,
     "command-palette:Open command palette": false,
     "bases:Create new base": false
   }
 },
 "active": "016d933bbc1cb39d",
 "lastOpenFiles": [
   "Welcome.md"
 ]
}`,
   }

   for filename, content := range files {
   	filePath := filepath.Join(obsidianDir, filename)
   	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
   		return fmt.Errorf("failed to write %s: %w", filename, err)
   	}
   }

   return nil
}

func createInitialFiles(vaultPath string) error {
   files := map[string]string{
   	"Welcome.md": `This is your new *vault*.

Make a note of something, [[create a link]], or try [the Importer](https://help.obsidian.md/Plugins/Importer)!

When you're ready, delete this note and make the vault your own.`,
   }

   for filename, content := range files {
   	filePath := filepath.Join(vaultPath, filename)
   	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
   		return fmt.Errorf("failed to write %s: %w", filename, err)
   	}
   }

   return nil
}

func updateGlobalConfig(vaultPath, vaultID string) error {
   configDir, err := UserConfigDirectory()
   if err != nil {
   	return fmt.Errorf("failed to get user config directory: %w", err)
   }

   obsidianConfigPath := filepath.Join(configDir, "obsidian", "obsidian.json")

   // Create obsidian config directory if it doesn't exist
   if err := os.MkdirAll(filepath.Dir(obsidianConfigPath), 0755); err != nil {
   	return fmt.Errorf("failed to create obsidian config directory: %w", err)
   }

   // Read existing config or create new one
   var config ObsidianConfig
   if data, err := os.ReadFile(obsidianConfigPath); err == nil {
   	if err := json.Unmarshal(data, &config); err != nil {
   		return fmt.Errorf("failed to parse existing config: %w", err)
   	}
   } else if !errors.Is(err, os.ErrNotExist) {
   	return fmt.Errorf("failed to read config file: %w", err)
   }

   // Initialize maps if they're nil
   if config.Vaults == nil {
   	config.Vaults = make(map[string]VaultConfig)
   }
   if config.OpenSchemes == nil {
   	config.OpenSchemes = map[string]bool{
   		"vscode":           true,
   		"chrome-extension": true,
   	}
   }

   // Add new vault
   config.Vaults[vaultID] = VaultConfig{
   	Path: vaultPath,
   	Ts:   1757173820641, // Using timestamp from example
   	Open: true,
   }

   // Write updated config
   data, err := json.MarshalIndent(config, "", "  ")
   if err != nil {
   	return fmt.Errorf("failed to marshal config: %w", err)
   }

   if err := os.WriteFile(obsidianConfigPath, data, 0644); err != nil {
   	return fmt.Errorf("failed to write config file: %w", err)
   }

   return nil
}
