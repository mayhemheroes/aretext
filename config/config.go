package config

import (
	"errors"
	"log"
)

const DefaultSyntaxLanguage = "undefined"
const DefaultTabSize = 4
const DefaultTabExpand = false
const DefaultAutoIndent = false

// Config is a configuration for the editor.
type Config struct {
	// Language used for syntax highlighting.
	SyntaxLanguage string

	// Size of a tab character in columns.
	TabSize int

	// If enabled, the tab key inserts spaces.
	TabExpand bool

	// If enabled, indent a new line to match indentation of the previous line.
	AutoIndent bool

	// User-defined commands to include in the menu.
	MenuCommands []MenuCommandConfig

	// Directories to exclude from file search.
	HideDirectories []string
}

// MenuCommandConfig is a configuration for a user-defined menu item.
type MenuCommandConfig struct {
	// Name is the displayed name of the menu.
	Name string

	// ShellCmd is the shell command to execute when the menu item is selected.
	// The command's output will be piped to a pager (usually `less`),
	// so it should be non-interactive.
	ShellCmd string
}

// ConfigFromUntypedMap constructs a configuration from an untyped map.
// The map is usually loaded from a JSON document.
func ConfigFromUntypedMap(m map[string]interface{}) Config {
	return Config{
		SyntaxLanguage:  stringOrDefault(m, "syntaxLanguage", DefaultSyntaxLanguage),
		TabSize:         intOrDefault(m, "tabSize", DefaultTabSize),
		TabExpand:       boolOrDefault(m, "tabExpand", DefaultTabExpand),
		AutoIndent:      boolOrDefault(m, "autoIndent", DefaultAutoIndent),
		MenuCommands:    menuCommandsFromSlice(sliceOrNil(m, "menuCommands")),
		HideDirectories: stringSliceOrNil(m, "hideDirectories"),
	}
}

// DefaultConfig is a configuration with all keys set to default values.
func DefaultConfig() Config {
	return ConfigFromUntypedMap(nil)
}

// Validate checks that the values in the configuration are valid.
func (c Config) Validate() error {
	if c.TabSize < 1 {
		return errors.New("TabSize must be greater than zero")
	}

	return nil
}

func stringOrDefault(m map[string]interface{}, key string, defaultVal string) string {
	v, ok := m[key]
	if !ok {
		return defaultVal
	}

	s, ok := v.(string)
	if !ok {
		log.Printf("Could not decode string for config key '%s'\n", key)
		return defaultVal
	}

	return s
}

func intOrDefault(m map[string]interface{}, key string, defaultVal int) int {
	v, ok := m[key]
	if !ok {
		return defaultVal
	}

	switch v.(type) {
	case int:
		return v.(int)
	case float64:
		return int(v.(float64))
	default:
		log.Printf("Could not decode int for config key '%s'\n", key)
		return defaultVal
	}
}

func boolOrDefault(m map[string]interface{}, key string, defaultVal bool) bool {
	v, ok := m[key]
	if !ok {
		return defaultVal
	}

	b, ok := v.(bool)
	if !ok {
		log.Printf("Could not decode bool for config key '%s'\n", key)
		return defaultVal
	}

	return b
}

func sliceOrNil(m map[string]interface{}, key string) []interface{} {
	v, ok := m[key]
	if !ok {
		return nil
	}

	s, ok := v.([]interface{})
	if !ok {
		log.Printf("Could not decode slice for config key '%s'\n", key)
		return nil
	}

	return s
}

func stringSliceOrNil(m map[string]interface{}, key string) []string {
	slice := sliceOrNil(m, key)
	if slice == nil {
		return nil
	}

	stringSlice := make([]string, 0, len(slice))
	for i := 0; i < len(slice); i++ {
		s, ok := (slice[i]).(string)
		if !ok {
			log.Printf("Could not decode string in slice for config key '%s'\n", key)
			continue
		}
		stringSlice = append(stringSlice, s)
	}
	return stringSlice
}

func menuCommandsFromSlice(s []interface{}) []MenuCommandConfig {
	result := make([]MenuCommandConfig, 0, len(s))
	for _, m := range s {
		menuMap, ok := m.(map[string]interface{})
		if !ok {
			log.Printf("Could not decode menu command map from %v\n", m)
			continue
		}

		result = append(result, MenuCommandConfig{
			Name:     stringOrDefault(menuMap, "name", "[nil]"),
			ShellCmd: stringOrDefault(menuMap, "shellCmd", ""),
		})
	}
	return result
}
