package jet

import (
	"fmt"
	"sort"
	"strings"
)

type formatMode int

const (
	modeNormal formatMode = iota
	modeFlattened
	modeNormalized
)

type jetWriter struct {
	sb   *strings.Builder
	mode formatMode
}

func format(data interface{}) ([]byte, error) {
	w := &jetWriter{
		sb:   &strings.Builder{},
		mode: modeNormal,
	}
	err := w.writeValue(data, 0)
	if err != nil {
		return nil, err
	}
	return []byte(w.sb.String()), nil
}

func formatFlattened(data interface{}) ([]byte, error) {
	w := &jetWriter{
		sb:   &strings.Builder{},
		mode: modeFlattened,
	}
	err := w.writeValue(data, 0)
	if err != nil {
		return nil, err
	}
	return []byte(w.sb.String()), nil
}

func formatNormalized(data interface{}) ([]byte, error) {
	w := &jetWriter{
		sb:   &strings.Builder{},
		mode: modeNormalized,
	}
	err := w.writeValue(data, 0)
	if err != nil {
		return nil, err
	}
	return []byte(w.sb.String()), nil
}

func (w *jetWriter) writeValue(data interface{}, indentLevel int) error {
	indentStr := strings.Repeat(" ", indentLevel)

	switch v := data.(type) {
	// Handling maps
	case map[string]interface{}:
		keys := make([]string, 0, len(v))
		for key := range v {
			keys = append(keys, key)
		}

		sort.Strings(keys)
		for _, key := range keys {
			value := v[key]
			if subSlice, ok := value.([]interface{}); ok && isTabular(subSlice) {
				w.writeTabularArray(indentStr, key, subSlice, indentLevel)
			} else if subMap, ok := value.(map[string]interface{}); ok {
				// Nested object
				w.sb.WriteString(fmt.Sprintf("%s%s:\n", indentStr, key))
				w.writeValue(subMap, indentLevel+1)
			} else {
				// Simple key-value pair
				w.sb.WriteString(fmt.Sprintf("%s%s: %v\n", indentStr, key, value))
			}
		}
	case []interface{}:
		// Handle non-tabular arrays
		if isTabular(v) {
			// This shouldn't happen in normal flow, but handle it
			w.writeTabularArray(indentStr, "", v, indentLevel)
		} else {
			// Simple array - not implemented in spec, but handle gracefully
			for _, item := range v {
				w.sb.WriteString(fmt.Sprintf("%s- %v\n", indentStr, item))
			}
		}
	default:
		// Scalar value
		w.sb.WriteString(fmt.Sprintf("%s%v\n", indentStr, v))
	}

	return nil
}

func (w *jetWriter) writeTabularArray(indentStr, key string, data []interface{}, indentLevel int) {
	firstRow := data[0].(map[string]interface{})
	schema := make([]string, 0, len(firstRow))
	for k := range firstRow {
		schema = append(schema, k)
	}
	sort.Strings(schema)

	switch w.mode {
	case modeFlattened:
		w.writeTabularArrayFlattened(indentStr, key, data, schema, indentLevel)
	case modeNormalized:
		w.writeTabularArrayNormalized(indentStr, key, data, schema, indentLevel)
	default:
		w.writeTabularArrayNormal(indentStr, key, data, schema, indentLevel)
	}
}

// writeTabularArrayNormal writes the normal format with nested blocks using > sigil
func (w *jetWriter) writeTabularArrayNormal(indentStr, key string, data []interface{}, schema []string, indentLevel int) {
	// Write header
	w.sb.WriteString(fmt.Sprintf("%s%s{%s}:\n", indentStr, key, strings.Join(schema, "|")))

	// Write rows
	rowDataIndent := strings.Repeat(" ", indentLevel+1)
	for _, row := range data {
		rowMap := row.(map[string]interface{})
		w.sb.WriteString(rowDataIndent)

		// Handle nesting
		values := []string{}
		for _, col := range schema {
			val := rowMap[col]
			if _, ok := val.(map[string]interface{}); ok {
				// Ignore
			} else if subSlice, ok := val.([]interface{}); ok && isTabular(subSlice) {
				// Ignore
			} else {
				values = append(values, fmt.Sprintf("%v", val))
			}
		}
		w.sb.WriteString(strings.Join(values, "|"))
		w.sb.WriteString("\n")

		// Handle nesting
		for _, col := range schema {
			val := rowMap[col]
			if subMap, ok := val.(map[string]interface{}); ok {
				w.sb.WriteString(fmt.Sprintf("%s  > %s:\n", rowDataIndent, col))
				w.writeValue(subMap, indentLevel+2)
			} else if subSlice, ok := val.([]interface{}); ok && isTabular(subSlice) {
				w.sb.WriteString(fmt.Sprintf("%s  > ", rowDataIndent))
				w.writeTabularArray("", col, subSlice, indentLevel+2)
			}
		}
	}
}

// writeTabularArrayFlattened writes flattened format with nested scalar objects inline
func (w *jetWriter) writeTabularArrayFlattened(indentStr, key string, data []interface{}, schema []string, indentLevel int) {
	firstRow := data[0].(map[string]interface{})
	// Build flattened schema and collect values
	flatSchema := w.buildFlattenedSchema(schema, firstRow)

	// Write header
	w.sb.WriteString(fmt.Sprintf("%s%s{%s}:\n", indentStr, key, flatSchema))

	// Write rows
	rowDataIndent := strings.Repeat(" ", indentLevel+1)
	for _, row := range data {
		rowMap := row.(map[string]interface{})
		w.sb.WriteString(rowDataIndent)

		values := w.extractFlattenedValues(schema, rowMap)
		w.sb.WriteString(strings.Join(values, "|"))
		w.sb.WriteString("\n")
	}
}

// writeTabularArrayNormalized writes normalized format with nested objects as pipe-delimited blocks
func (w *jetWriter) writeTabularArrayNormalized(indentStr, key string, data []interface{}, schema []string, indentLevel int) {
	firstRow := data[0].(map[string]interface{})
	// Build normalized schema showing nested structure
	normalizedSchema := w.buildNormalizedSchema(schema, firstRow)

	// Write header
	w.sb.WriteString(fmt.Sprintf("%s%s{%s}:\n", indentStr, key, normalizedSchema))

	// Write rows
	rowDataIndent := strings.Repeat(" ", indentLevel+1)
	for _, row := range data {
		rowMap := row.(map[string]interface{})
		w.sb.WriteString(rowDataIndent)

		// Write scalar values
		values := []string{}
		for _, col := range schema {
			val := rowMap[col]
			if _, ok := val.(map[string]interface{}); ok {
				// Skip - will be handled in nested block
			} else if subSlice, ok := val.([]interface{}); ok && isTabular(subSlice) {
				// Skip - will be handled in nested block
			} else {
				values = append(values, fmt.Sprintf("%v", val))
			}
		}
		w.sb.WriteString(strings.Join(values, "|"))
		w.sb.WriteString("\n")

		// Handle nesting with normalized format
		for _, col := range schema {
			val := rowMap[col]
			if subMap, ok := val.(map[string]interface{}); ok {
				if canFlattenObject(subMap) {
					// Write normalized nested object as pipe-delimited values
					w.sb.WriteString(fmt.Sprintf("%s  > %s:\n", rowDataIndent, col))
					w.sb.WriteString(rowDataIndent + "  ")

					subKeys := make([]string, 0, len(subMap))
					for k := range subMap {
						subKeys = append(subKeys, k)
					}
					sort.Strings(subKeys)

					subValues := []string{}
					for _, subKey := range subKeys {
						subValues = append(subValues, fmt.Sprintf("%v", subMap[subKey]))
					}
					w.sb.WriteString(strings.Join(subValues, "|"))
					w.sb.WriteString("\n")
				} else {
					// Cannot normalize - has nested structures, use normal format
					w.sb.WriteString(fmt.Sprintf("%s  > %s:\n", rowDataIndent, col))
					w.writeValue(subMap, indentLevel+2)
				}
			} else if subSlice, ok := val.([]interface{}); ok && isTabular(subSlice) {
				w.sb.WriteString(fmt.Sprintf("%s  > ", rowDataIndent))
				w.writeTabularArray("", col, subSlice, indentLevel+2)
			}
		}
	}
}

// buildFlattenedSchema creates a flattened schema string with nested objects expanded
func (w *jetWriter) buildFlattenedSchema(schema []string, sampleRow map[string]interface{}) string {
	var parts []string

	for _, col := range schema {
		val := sampleRow[col]
		if subMap, ok := val.(map[string]interface{}); ok {
			// Check if this object has only simple scalar values (can be flattened)
			if canFlattenObject(subMap) {
				// Get subkeys
				subKeys := make([]string, 0, len(subMap))
				for k := range subMap {
					subKeys = append(subKeys, k)
				}
				sort.Strings(subKeys)
				parts = append(parts, fmt.Sprintf("%s{%s}", col, strings.Join(subKeys, ",")))
			} else {
				// Cannot flatten - has nested structures, keep as is
				parts = append(parts, col)
			}
		} else if subSlice, ok := val.([]interface{}); ok && isTabular(subSlice) {
			// Nested tabular array - cannot flatten inline, keep as is
			parts = append(parts, col)
		} else {
			parts = append(parts, col)
		}
	}

	return strings.Join(parts, "|")
}

// buildNormalizedSchema creates a normalized schema string with nested objects shown with pipe-delimited structure
func (w *jetWriter) buildNormalizedSchema(schema []string, sampleRow map[string]interface{}) string {
	var parts []string

	for _, col := range schema {
		val := sampleRow[col]
		if subMap, ok := val.(map[string]interface{}); ok {
			// Check if this object has only simple scalar values (can be normalized)
			if canFlattenObject(subMap) {
				// Get subkeys
				subKeys := make([]string, 0, len(subMap))
				for k := range subMap {
					subKeys = append(subKeys, k)
				}
				sort.Strings(subKeys)
				// Use pipe delimiter to indicate values will be pipe-separated in the nested block
				parts = append(parts, fmt.Sprintf("%s{%s}", col, strings.Join(subKeys, "|")))
			} else {
				// Cannot normalize - has nested structures, keep as is
				parts = append(parts, col)
			}
		} else if subSlice, ok := val.([]interface{}); ok && isTabular(subSlice) {
			// Nested tabular array - cannot normalize inline, keep as is
			parts = append(parts, col)
		} else {
			parts = append(parts, col)
		}
	}

	return strings.Join(parts, "|")
}

// extractFlattenedValues extracts values in flattened order
func (w *jetWriter) extractFlattenedValues(schema []string, rowMap map[string]interface{}) []string {
	var values []string

	for _, col := range schema {
		val := rowMap[col]
		if subMap, ok := val.(map[string]interface{}); ok {
			if canFlattenObject(subMap) {
				// Extract nested values in sorted order
				subKeys := make([]string, 0, len(subMap))
				for k := range subMap {
					subKeys = append(subKeys, k)
				}
				sort.Strings(subKeys)
				for _, subKey := range subKeys {
					values = append(values, fmt.Sprintf("%v", subMap[subKey]))
				}
			} else {
				// Cannot flatten - output placeholder or skip
				values = append(values, "[nested]")
			}
		} else if subSlice, ok := val.([]interface{}); ok && isTabular(subSlice) {
			// Nested tabular array - output placeholder
			values = append(values, "[table]")
		} else {
			values = append(values, fmt.Sprintf("%v", val))
		}
	}

	return values
}

// canFlattenObject checks if an object contains only scalar values (no nested objects/arrays)
func canFlattenObject(obj map[string]interface{}) bool {
	for _, v := range obj {
		switch v.(type) {
		case map[string]interface{}, []interface{}:
			return false
		}
	}
	return true
}

func isTabular(slice []interface{}) bool {
	if len(slice) == 0 {
		return false
	}

	firstMap, ok := slice[0].(map[string]interface{})
	if !ok {
		return false
	}

	// Get keys from first map
	firstKeys := make(map[string]bool)
	for k := range firstMap {
		firstKeys[k] = true
	}

	// Verify all items are maps with the same keys
	for i := 1; i < len(slice); i++ {
		currentMap, ok := slice[i].(map[string]interface{})
		if !ok {
			return false
		}

		// Check same number of keys
		if len(currentMap) != len(firstKeys) {
			return false
		}

		// Check all keys match
		for k := range currentMap {
			if !firstKeys[k] {
				return false
			}
		}
	}

	return true
}
