package wtf

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// formField represents a parsed form field with all its attributes
// extracted from struct tags and reflection.
type formField struct {
	// Core attributes
	Name  string // HTML name attribute
	Label string // Label text
	Type  string // HTML input type (text, password, email, etc.)
	ID    string // HTML id attribute
	Value string // Current value (pre-populated from struct)

	// Validation
	Required bool
	Disabled bool
	ReadOnly bool
	Pattern  string
	Min      string
	Max      string
	Step     string

	// Appearance
	Placeholder  string
	Class        string
	Autocomplete string

	// Select/Radio options
	Options []selectOption

	// Textarea
	Rows string
	Cols string

	// Multi-value (for slices)
	Multiple bool

	// Accessibility
	Role string

	// Field metadata
	FieldType reflect.Type
}

// selectOption represents an option in a select or radio group.
type selectOption struct {
	Value    string
	Label    string
	Selected bool
}

// cachedField represents the static struct tag metadata of a field,
// cached to avoid using reflection repeatedly.
type cachedField struct {
	Index        int
	Name         string
	Label        string
	Type         string
	ID           string
	Required     bool
	Disabled     bool
	ReadOnly     bool
	Pattern      string
	Min          string
	Max          string
	Step         string
	Placeholder  string
	Class        string
	Autocomplete string
	Options      []selectOption
	Rows         string
	Cols         string
	Multiple     bool
	Role         string
	FieldType    reflect.Type
}

// structMetadata holds both the static FormInfo settings and the cached field metadatas for a struct type.
type structMetadata struct {
	FormInfo FormInfo
	Fields   []cachedField
}

// parseStructMetadata parses the static form tags and metadata for a struct type.
func parseStructMetadata(typ reflect.Type) (structMetadata, error) {
	if typ.Kind() != reflect.Struct {
		return structMetadata{}, fmt.Errorf("wtf: expected struct, got %s", typ.Kind())
	}

	var meta structMetadata

	for i := 0; i < typ.NumField(); i++ {
		structField := typ.Field(i)

		// Skip unexported fields
		if !structField.IsExported() {
			continue
		}

		// Handle FormInfo specially
		if structField.Type == reflect.TypeOf(FormInfo{}) {
			tag := structField.Tag.Get("form")
			if tag != "" && tag != "-" {
				meta.FormInfo = parseFormInfoTag(tag)
			}
			continue
		}

		tag := structField.Tag.Get("form")

		// Skip fields with no form tag
		if tag == "" {
			continue
		}

		// Skip fields explicitly marked with "-"
		if tag == "-" {
			continue
		}

		field := parseTagMetadata(tag, structField, i)
		meta.Fields = append(meta.Fields, field)
	}

	return meta, nil
}

// parseFormInfoTag parses struct tags on the wtf.FormInfo field.
func parseFormInfoTag(tag string) FormInfo {
	parts := splitTag(tag)
	var info FormInfo

	for _, part := range parts[1:] {
		key, value, _ := parseTagPart(part)
		switch key {
		case "action":
			info.Action = value
		case "method":
			info.Method = value
		case "enctype":
			info.Enctype = value
		case "id":
			info.FormID = value
		case "class":
			info.FormClass = value
		case "attrs":
			info.FormAttrs = value
		case "submit_label":
			info.SubmitLabel = value
		case "submit_class":
			info.SubmitClass = value
		case "submit_attrs":
			info.SubmitAttrs = value
		case "submit_role":
			info.SubmitRole = value
		case "reset_label":
			info.ResetLabel = value
		case "reset_class":
			info.ResetClass = value
		case "reset_attrs":
			info.ResetAttrs = value
		case "reset_role":
			info.ResetRole = value
		case "fieldset":
			info.Fieldset = value == "true"
		}
	}
	return info
}

// mergeFormInfoValues merges dynamic FormInfo runtime values with static struct tag ones.
func mergeFormInfoValues(static FormInfo, dynamic FormInfo) FormInfo {
	merged := static
	if dynamic.Action != "" {
		merged.Action = dynamic.Action
	}
	if dynamic.Method != "" {
		merged.Method = dynamic.Method
	}
	if dynamic.Enctype != "" {
		merged.Enctype = dynamic.Enctype
	}
	if dynamic.FormID != "" {
		merged.FormID = dynamic.FormID
	}
	if dynamic.FormClass != "" {
		merged.FormClass = dynamic.FormClass
	}
	if dynamic.FormAttrs != "" {
		merged.FormAttrs = dynamic.FormAttrs
	}
	if dynamic.Fieldset {
		merged.Fieldset = true
	}
	if dynamic.SubmitLabel != "" {
		merged.SubmitLabel = dynamic.SubmitLabel
	}
	if dynamic.SubmitClass != "" {
		merged.SubmitClass = dynamic.SubmitClass
	}
	if dynamic.SubmitAttrs != "" {
		merged.SubmitAttrs = dynamic.SubmitAttrs
	}
	if dynamic.SubmitRole != "" {
		merged.SubmitRole = dynamic.SubmitRole
	}
	if dynamic.ResetLabel != "" {
		merged.ResetLabel = dynamic.ResetLabel
	}
	if dynamic.ResetClass != "" {
		merged.ResetClass = dynamic.ResetClass
	}
	if dynamic.ResetAttrs != "" {
		merged.ResetAttrs = dynamic.ResetAttrs
	}
	if dynamic.ResetRole != "" {
		merged.ResetRole = dynamic.ResetRole
	}
	return merged
}

// parseTagMetadata parses a single struct tag string into static cachedField metadata.
func parseTagMetadata(tag string, structField reflect.StructField, index int) cachedField {
	parts := splitTag(tag)

	field := cachedField{
		Index:     index,
		Name:      parts[0],
		FieldType: structField.Type,
	}

	// Parse key=value pairs and flags
	for _, part := range parts[1:] {
		key, value, hasValue := parseTagPart(part)

		switch key {
		case "type":
			field.Type = value
		case "label":
			field.Label = value
		case "placeholder":
			field.Placeholder = value
		case "id":
			field.ID = value
		case "class":
			field.Class = value
		case "pattern":
			field.Pattern = value
		case "min":
			field.Min = value
		case "max":
			field.Max = value
		case "step":
			field.Step = value
		case "rows":
			field.Rows = value
		case "cols":
			field.Cols = value
		case "options":
			field.Options = parseOptions(value)
		case "required":
			if hasValue {
				field.Required = value != "false"
			} else {
				field.Required = true
			}
		case "disabled":
			if hasValue {
				field.Disabled = value != "false"
			} else {
				field.Disabled = true
			}
		case "readonly":
			if hasValue {
				field.ReadOnly = value != "false"
			} else {
				field.ReadOnly = true
			}
		case "multiple":
			if hasValue {
				field.Multiple = value != "false"
			} else {
				field.Multiple = true
			}
		case "role":
			field.Role = value
		case "autocomplete":
			field.Autocomplete = value
		}
	}

	// Apply defaults
	if field.Type == "" {
		field.Type = inferType(structField.Type)
	}

	if field.Label == "" {
		field.Label = titleCase(field.Name)
	}

	if field.ID == "" {
		field.ID = "form-" + field.Name
	}

	return field
}


// splitTag splits a tag string by commas, but respects values that contain
// pipe characters and doesn't split on commas within option lists.
func splitTag(tag string) []string {
	var parts []string
	var current strings.Builder
	depth := 0

	for _, r := range tag {
		switch r {
		case '(':
			depth++
			current.WriteRune(r)
		case ')':
			depth--
			current.WriteRune(r)
		case ',':
			if depth == 0 {
				parts = append(parts, strings.TrimSpace(current.String()))
				current.Reset()
			} else {
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, strings.TrimSpace(current.String()))
	}

	return parts
}

// parseTagPart parses a single tag part like "key=value" or "flag".
func parseTagPart(part string) (key, value string, hasValue bool) {
	idx := strings.Index(part, "=")
	if idx == -1 {
		return part, "", false
	}
	return part[:idx], part[idx+1:], true
}

// parseOptions parses a pipe-separated list of options.
// Format: "opt1|opt2|opt3" or "value1:Label 1|value2:Label 2"
func parseOptions(opts string) []selectOption {
	var options []selectOption
	for _, opt := range strings.Split(opts, "|") {
		opt = strings.TrimSpace(opt)
		if opt == "" {
			continue
		}

		// Check for value:label format
		if idx := strings.Index(opt, ":"); idx != -1 {
			options = append(options, selectOption{
				Value: opt[:idx],
				Label: opt[idx+1:],
			})
		} else {
			options = append(options, selectOption{
				Value: opt,
				Label: titleCase(opt),
			})
		}
	}
	return options
}

// inferType guesses the HTML input type based on the Go type.
func inferType(t reflect.Type) string {
	// Dereference pointer types
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Check for time.Time
	if t == reflect.TypeOf(time.Time{}) {
		return "datetime-local"
	}

	switch t.Kind() {
	case reflect.Bool:
		return "checkbox"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "number"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "number"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.String:
		return "text"
	case reflect.Slice:
		return "select"
	default:
		return "text"
	}
}

// extractValue gets the current value of a field as a string for HTML rendering.
func extractValue(v reflect.Value, inputType string) string {
	// Dereference pointers
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	// Handle time.Time specially
	if v.Type() == reflect.TypeOf(time.Time{}) {
		t := v.Interface().(time.Time)
		if t.IsZero() {
			return ""
		}
		switch inputType {
		case "date":
			return t.Format("2006-01-02")
		case "time":
			return t.Format("15:04")
		case "datetime-local":
			return t.Format("2006-01-02T15:04")
		default:
			return t.Format("2006-01-02T15:04")
		}
	}

	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return "true"
		}
		return ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n := v.Int()
		if n == 0 {
			return ""
		}
		return fmt.Sprintf("%d", n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n := v.Uint()
		if n == 0 {
			return ""
		}
		return fmt.Sprintf("%d", n)
	case reflect.Float32, reflect.Float64:
		f := v.Float()
		if f == 0 {
			return ""
		}
		return fmt.Sprintf("%g", f)
	case reflect.String:
		return v.String()
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}

// titleCase converts a snake_case or camelCase string to Title Case.
func titleCase(s string) string {
	// Replace underscores and hyphens with spaces
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")

	// Insert spaces before uppercase letters in camelCase
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			prev := rune(s[i-1])
			if prev >= 'a' && prev <= 'z' {
				result.WriteRune(' ')
			}
		}
		result.WriteRune(r)
	}

	// Title case each word
	words := strings.Fields(result.String())
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, " ")
}
