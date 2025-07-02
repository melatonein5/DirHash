package yara

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/melatonein5/DirHash/src/files"
)

// YaraRule represents a YARA rule structure
type YaraRule struct {
	Name        string
	Description string
	Author      string
	Date        string
	Tags        []string
	Strings     []YaraString
	Condition   string
}

// YaraString represents a string definition in YARA
type YaraString struct {
	Name  string
	Value string
	Type  string // "hex", "text", "regex"
}

// YaraRuleSet represents a collection of YARA rules
type YaraRuleSet struct {
	Rules []YaraRule
}

// GenerateYaraRule creates a YARA rule from file hash data
func GenerateYaraRule(files []*files.File, ruleName string) (*YaraRule, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided for YARA rule generation")
	}

	if ruleName == "" {
		ruleName = "generated_rule"
	}

	// Sanitize rule name (YARA rules must be valid identifiers)
	ruleName = sanitizeRuleName(ruleName)

	rule := &YaraRule{
		Name:        ruleName,
		Description: fmt.Sprintf("Generated rule based on %d files", len(files)),
		Author:      "DirHash",
		Date:        time.Now().Format("2006-01-02"),
		Tags:        []string{"generated", "dirhash"},
		Strings:     make([]YaraString, 0),
		Condition:   "",
	}

	// Generate hash-based strings
	hashStrings := generateHashStrings(files)
	rule.Strings = append(rule.Strings, hashStrings...)

	// Generate filename-based strings if applicable
	filenameStrings := generateFilenameStrings(files)
	rule.Strings = append(rule.Strings, filenameStrings...)

	// Generate condition
	rule.Condition = generateCondition(rule.Strings)

	return rule, nil
}

// GenerateYaraRuleFromHashes creates a YARA rule with only hash conditions
func GenerateYaraRuleFromHashes(files []*files.File, ruleName string, hashTypes []string) (*YaraRule, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided for YARA rule generation")
	}

	if ruleName == "" {
		ruleName = "hash_based_rule"
	}

	ruleName = sanitizeRuleName(ruleName)

	rule := &YaraRule{
		Name:        ruleName,
		Description: fmt.Sprintf("Hash-based rule for %d files", len(files)),
		Author:      "DirHash",
		Date:        time.Now().Format("2006-01-02"),
		Tags:        []string{"hash", "generated", "dirhash"},
		Strings:     make([]YaraString, 0),
		Condition:   "",
	}

	// Generate only hash-based strings for specified hash types
	for _, file := range files {
		for _, hashType := range hashTypes {
			if hash, exists := file.Hashes[hashType]; exists && hash != "" {
				stringName := fmt.Sprintf("$%s_%s", hashType, sanitizeStringName(file.FileName))
				rule.Strings = append(rule.Strings, YaraString{
					Name:  stringName,
					Value: hash,
					Type:  "hex",
				})
			}
		}
	}

	if len(rule.Strings) == 0 {
		return nil, fmt.Errorf("no valid hashes found for specified hash types")
	}

	rule.Condition = generateHashCondition(rule.Strings)
	return rule, nil
}

// ToYaraFormat converts the rule to YARA rule format
func (r *YaraRule) ToYaraFormat() string {
	var builder strings.Builder

	// Rule header
	builder.WriteString(fmt.Sprintf("rule %s\n{\n", r.Name))

	// Metadata section
	builder.WriteString("    meta:\n")
	builder.WriteString(fmt.Sprintf("        description = \"%s\"\n", r.Description))
	builder.WriteString(fmt.Sprintf("        author = \"%s\"\n", r.Author))
	builder.WriteString(fmt.Sprintf("        date = \"%s\"\n", r.Date))
	
	if len(r.Tags) > 0 {
		builder.WriteString(fmt.Sprintf("        tags = \"%s\"\n", strings.Join(r.Tags, ", ")))
	}

	// Strings section
	if len(r.Strings) > 0 {
		builder.WriteString("\n    strings:\n")
		for _, str := range r.Strings {
			switch str.Type {
			case "hex":
				builder.WriteString(fmt.Sprintf("        %s = { %s }\n", str.Name, str.Value))
			case "text":
				builder.WriteString(fmt.Sprintf("        %s = \"%s\"\n", str.Name, str.Value))
			case "regex":
				builder.WriteString(fmt.Sprintf("        %s = /%s/\n", str.Name, str.Value))
			default:
				builder.WriteString(fmt.Sprintf("        %s = \"%s\"\n", str.Name, str.Value))
			}
		}
	}

	// Condition section
	builder.WriteString("\n    condition:\n")
	builder.WriteString(fmt.Sprintf("        %s\n", r.Condition))

	builder.WriteString("}\n")
	return builder.String()
}

// generateHashStrings creates YARA strings from file hashes
func generateHashStrings(files []*files.File) []YaraString {
	var strings []YaraString
	
	for _, file := range files {
		baseName := sanitizeStringName(file.FileName)
		
		// Add hash strings for each available hash type
		hashTypes := []string{"md5", "sha1", "sha256", "sha512"}
		for _, hashType := range hashTypes {
			if hash, exists := file.Hashes[hashType]; exists && hash != "" {
				stringName := fmt.Sprintf("$%s_%s", hashType, baseName)
				strings = append(strings, YaraString{
					Name:  stringName,
					Value: formatHashForYara(hash),
					Type:  "hex",
				})
			}
		}
	}
	
	return strings
}

// generateFilenameStrings creates YARA strings from filenames
func generateFilenameStrings(files []*files.File) []YaraString {
	var strings []YaraString
	seenNames := make(map[string]bool)
	
	for _, file := range files {
		fileName := filepath.Base(file.FileName)
		if !seenNames[fileName] {
			seenNames[fileName] = true
			stringName := fmt.Sprintf("$filename_%s", sanitizeStringName(fileName))
			strings = append(strings, YaraString{
				Name:  stringName,
				Value: fileName,
				Type:  "text",
			})
		}
	}
	
	return strings
}

// generateCondition creates a YARA condition from strings
func generateCondition(yaraStrings []YaraString) string {
	if len(yaraStrings) == 0 {
		return "true"
	}

	var hashConditions []string
	var filenameConditions []string

	for _, str := range yaraStrings {
		if str.Type == "hex" {
			hashConditions = append(hashConditions, str.Name)
		} else if strings.Contains(str.Name, "filename_") {
			filenameConditions = append(filenameConditions, str.Name)
		}
	}

	var conditions []string

	if len(hashConditions) > 0 {
		if len(hashConditions) == 1 {
			conditions = append(conditions, hashConditions[0])
		} else {
			conditions = append(conditions, fmt.Sprintf("any of (%s)", strings.Join(hashConditions, ", ")))
		}
	}

	if len(filenameConditions) > 0 {
		if len(filenameConditions) == 1 {
			conditions = append(conditions, filenameConditions[0])
		} else {
			conditions = append(conditions, fmt.Sprintf("any of (%s)", strings.Join(filenameConditions, ", ")))
		}
	}

	if len(conditions) == 0 {
		return "true"
	} else if len(conditions) == 1 {
		return conditions[0]
	} else {
		return strings.Join(conditions, " or ")
	}
}

// generateHashCondition creates a hash-only condition
func generateHashCondition(yaraStrings []YaraString) string {
	if len(yaraStrings) == 0 {
		return "true"
	}

	var conditions []string
	for _, str := range yaraStrings {
		conditions = append(conditions, str.Name)
	}

	if len(conditions) == 1 {
		return conditions[0]
	}
	return fmt.Sprintf("any of (%s)", strings.Join(conditions, ", "))
}

// sanitizeRuleName ensures the rule name is valid for YARA
func sanitizeRuleName(name string) string {
	// Replace invalid characters with underscores
	result := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return '_'
	}, name)

	// Ensure it starts with a letter or underscore
	if len(result) > 0 && (result[0] >= '0' && result[0] <= '9') {
		result = "_" + result
	}

	if result == "" {
		result = "generated_rule"
	}

	return result
}

// sanitizeStringName ensures the string name is valid for YARA
func sanitizeStringName(name string) string {
	// Remove file extension and sanitize
	base := strings.TrimSuffix(name, filepath.Ext(name))
	result := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return '_'
	}, base)

	if len(result) > 0 && (result[0] >= '0' && result[0] <= '9') {
		result = "_" + result
	}

	if result == "" {
		result = "file"
	}

	return result
}

// formatHashForYara formats a hash string for YARA hex format
func formatHashForYara(hash string) string {
	// Convert hash to YARA hex format (space-separated hex bytes)
	var result []string
	hash = strings.ToUpper(hash)
	
	for i := 0; i < len(hash); i += 2 {
		if i+1 < len(hash) {
			result = append(result, hash[i:i+2])
		}
	}
	
	return strings.Join(result, " ")
}

// GetSupportedHashTypes returns the hash types supported for YARA generation
func GetSupportedHashTypes() []string {
	return []string{"md5", "sha1", "sha256", "sha512"}
}