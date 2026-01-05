package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("🔍 Checking sanitize validation tags in DTOs...")

	dtoDir := "common/dto"
	errors := []string{}

	// Walk through all .go files in the dto directory
	err := filepath.Walk(dtoDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-Go files
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Parse the file
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", path, err)
		}

		// Inspect all declarations
		ast.Inspect(node, func(n ast.Node) bool {
			typeSpec, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				return true
			}

			structName := typeSpec.Name.Name

			// Only check DTOs that need sanitization
			if !shouldCheckStruct(structName) {
				return true
			}

			// Check each field in the struct
			for _, field := range structType.Fields.List {
				// Only check string fields
				if !isStringField(field) {
					continue
				}

				// Get field name
				fieldName := ""
				if len(field.Names) > 0 {
					fieldName = field.Names[0].Name
				}

				// Check if the field has a validate tag with "sanitize"
				if field.Tag == nil {
					errors = append(errors, fmt.Sprintf("❌ %s.%s: missing validate tag with 'sanitize'", structName, fieldName))
					continue
				}

				tagValue := field.Tag.Value
				if !strings.Contains(tagValue, "sanitize") {
					errors = append(errors, fmt.Sprintf("❌ %s.%s: validate tag missing 'sanitize'", structName, fieldName))
				}
			}

			return true
		})

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking directory: %v\n", err)
		os.Exit(1)
	}

	// Print results
	if len(errors) > 0 {
		fmt.Println("\n⚠️  Found validation issues:\n")
		for _, e := range errors {
			fmt.Println(e)
		}
		fmt.Printf("\n❌ Failed: %d string fields missing 'sanitize' tag\n", len(errors))
		os.Exit(1)
	}

	fmt.Println("✅ All DTO string fields have 'sanitize' validation tags!")
}

// shouldCheckStruct returns true if the struct should be checked for sanitize tags
func shouldCheckStruct(name string) bool {
	// Skip response DTOs - they contain already-sanitized data from the database
	if strings.HasSuffix(name, "ResponseDto") || strings.HasSuffix(name, "Response") {
		return false
	}

	// Skip specific DTOs that are used only for output (not user input)
	skipList := []string{
		"SportDto",       // Used only for responses
		"CommonStatsDto", // Used only for responses
	}
	for _, skip := range skipList {
		if name == skip {
			return false
		}
	}

	// Check all other DTO structs (input DTOs)
	// This includes: CreateDto, UpdateDto, SendMessageDto, MarkReadDto, etc.
	return strings.HasSuffix(name, "Dto") || name == "Login"
}

// isStringField returns true if the field is a string type
func isStringField(field *ast.Field) bool {
	ident, ok := field.Type.(*ast.Ident)
	return ok && ident.Name == "string"
}
