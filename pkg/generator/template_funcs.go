package generator

import (
	"strings"
	"unicode"
)

// TemplateFuncs returns a map of custom template functions
func TemplateFuncs() map[string]interface{} {
	return map[string]interface{}{
		"title":      title,
		"lower":      strings.ToLower,
		"upper":      strings.ToUpper,
		"camelCase":  camelCase,
		"snakeCase":  snakeCase,
	}
}

// title converts a string to title case (first letter uppercase)
func title(s string) string {
	if s == "" {
		return s
	}
	// Convert to PascalCase
	return toPascalCase(s)
}

// toPascalCase converts a string to PascalCase
func toPascalCase(s string) string {
	if s == "" {
		return s
	}

	// Split by common delimiters
	words := splitWords(s)

	var result strings.Builder
	for _, word := range words {
		if word == "" {
			continue
		}
		// Capitalize first letter of each word
		result.WriteString(strings.ToUpper(string(word[0])))
		if len(word) > 1 {
			result.WriteString(strings.ToLower(word[1:]))
		}
	}

	return result.String()
}

// camelCase converts a string to camelCase
func camelCase(s string) string {
	if s == "" {
		return s
	}

	pascal := toPascalCase(s)
	if pascal == "" {
		return pascal
	}

	// Make first letter lowercase
	runes := []rune(pascal)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// snakeCase converts a string to snake_case
func snakeCase(s string) string {
	if s == "" {
		return s
	}

	words := splitWords(s)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}

	return strings.Join(words, "_")
}

// splitWords splits a string into words by various delimiters
func splitWords(s string) []string {
	var words []string
	var currentWord strings.Builder

	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		// Handle delimiters
		if r == '_' || r == '-' || r == ' ' || r == '.' {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
			continue
		}

		// Handle camelCase/PascalCase transitions
		if i > 0 && unicode.IsUpper(r) && unicode.IsLower(runes[i-1]) {
			if currentWord.Len() > 0 {
				words = append(words, currentWord.String())
				currentWord.Reset()
			}
		}

		currentWord.WriteRune(r)
	}

	if currentWord.Len() > 0 {
		words = append(words, currentWord.String())
	}

	return words
}
