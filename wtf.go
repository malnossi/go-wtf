package wtf

import (
	"fmt"
	"html/template"
	"reflect"
	"sync"
)

// FormInfo allows declaring form-level options, submit button attributes,
// and reset button attributes directly inside the form struct.
type FormInfo struct {
	Action      string
	Method      string
	Enctype     string
	FormID      string
	FormClass   string
	FormAttrs   string // Custom raw attributes on <form> tag, e.g. "hx-target='#result'"
	SubmitLabel string // Text for submit button
	SubmitClass string // CSS class for submit button
	SubmitAttrs string // Custom attributes for submit button, e.g. "up-dismiss hx-post='/login'"
	SubmitRole  string // Role attribute for submit button, e.g. "button"
	ResetLabel  string // Text for reset button (only rendered if non-empty)
	ResetClass  string // CSS class for reset button
	ResetAttrs  string // Custom attributes for reset button, e.g. "up-dismiss hx-get='/reset'"
	ResetRole   string // Role attribute for reset button, e.g. "button"
}

// FormRenderer is the main entry point for rendering HTML forms from Go structs.
// It holds configuration and provides methods for template integration.
type FormRenderer struct {
	// Form attributes
	action  string
	method  string
	enctype string
	formID  string
	formClass string

	// Features
	csrfToken     string
	classPrefix   string
	submitLabel   string
	noSubmit      bool

	// FormInfo options
	formAttrs   string
	submitClass string
	submitAttrs string
	submitRole  string
	resetLabel  string
	resetClass  string
	resetAttrs  string
	resetRole   string

	// Cache for parsed struct tags (struct type -> parsed fields)
	cache sync.Map
}

// New creates a new FormRenderer with default settings.
// Use functional options (WithAction, WithMethod, etc.) to customize behavior.
func New(opts ...Option) *FormRenderer {
	r := &FormRenderer{
		method:        "POST",
		submitLabel:   "Submit",
		classPrefix:   "wtf",
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// FuncMap returns a template.FuncMap that registers the "render_form" function
// for use in html/template templates.
//
// Usage:
//
//	renderer := wtf.New(wtf.WithAction("/submit"))
//	tmpl := template.Must(
//	    template.New("page").Funcs(renderer.FuncMap()).Parse(pageHTML),
//	)
//
// In the template:
//
//	{{render_form .Form}}
func (r *FormRenderer) FuncMap() template.FuncMap {
	return template.FuncMap{
		"render_form": r.RenderForm,
	}
}

// RenderForm renders a struct as an HTML form. The struct's fields are inspected
// for "form" struct tags which control how each field is rendered.
//
// The function returns template.HTML so the output is not double-escaped
// when used in html/template.
//
// Returns an HTML comment with an error message if the input is invalid.
// RenderForm renders a struct as an HTML form. The struct's fields are inspected
// for "form" struct tags which control how each field is rendered.
//
// The function returns template.HTML so the output is not double-escaped
// when used in html/template.
//
// Returns an HTML comment with an error message if the input is invalid.
func (r *FormRenderer) RenderForm(v interface{}) template.HTML {
	if v == nil {
		return template.HTML("<!-- wtf: cannot render nil value -->")
	}

	fields, renderer, err := r.parseWithCache(v)
	if err != nil {
		return template.HTML("<!-- wtf: " + err.Error() + " -->")
	}

	html := renderForm(fields, renderer)
	return template.HTML(html)
}

// RenderFormWithAction renders a struct as an HTML form with a specific action and method,
// overriding the renderer's defaults for this single call.
func (r *FormRenderer) RenderFormWithAction(v interface{}, action, method string) template.HTML {
	// Create a new renderer with overridden action/method (avoids copying sync.Map)
	tmp := &FormRenderer{
		action:        action,
		method:        method,
		enctype:       r.enctype,
		formID:        r.formID,
		formClass:     r.formClass,
		csrfToken:     r.csrfToken,
		classPrefix:   r.classPrefix,
		submitLabel:   r.submitLabel,
		noSubmit:      r.noSubmit,
		formAttrs:     r.formAttrs,
		submitClass:   r.submitClass,
		submitAttrs:   r.submitAttrs,
		submitRole:    r.submitRole,
		resetLabel:    r.resetLabel,
		resetClass:    r.resetClass,
		resetAttrs:    r.resetAttrs,
		resetRole:     r.resetRole,
	}

	if v == nil {
		return template.HTML("<!-- wtf: cannot render nil value -->")
	}

	fields, renderer, err := tmp.parseWithCache(v)
	if err != nil {
		return template.HTML("<!-- wtf: " + err.Error() + " -->")
	}

	html := renderForm(fields, renderer)
	return template.HTML(html)
}

// mergeFormInfo merges static and dynamic FormInfo configurations into a FormRenderer copy.
func (r *FormRenderer) mergeFormInfo(info FormInfo) *FormRenderer {
	merged := &FormRenderer{
		action:        r.action,
		method:        r.method,
		enctype:       r.enctype,
		formID:        r.formID,
		formClass:     r.formClass,
		csrfToken:     r.csrfToken,
		classPrefix:   r.classPrefix,
		submitLabel:   r.submitLabel,
		noSubmit:      r.noSubmit,
		formAttrs:     r.formAttrs,
		submitClass:   r.submitClass,
		submitAttrs:   r.submitAttrs,
		submitRole:    r.submitRole,
		resetLabel:    r.resetLabel,
		resetClass:    r.resetClass,
		resetAttrs:    r.resetAttrs,
		resetRole:     r.resetRole,
	}

	if info.Action != "" {
		merged.action = info.Action
	}
	if info.Method != "" {
		merged.method = info.Method
	}
	if info.Enctype != "" {
		merged.enctype = info.Enctype
	}
	if info.FormID != "" {
		merged.formID = info.FormID
	}
	if info.FormClass != "" {
		merged.formClass = info.FormClass
	}
	if info.FormAttrs != "" {
		merged.formAttrs = info.FormAttrs
	}
	if info.SubmitLabel != "" {
		merged.submitLabel = info.SubmitLabel
	}
	if info.SubmitClass != "" {
		merged.submitClass = info.SubmitClass
	}
	if info.SubmitAttrs != "" {
		merged.submitAttrs = info.SubmitAttrs
	}
	if info.SubmitRole != "" {
		merged.submitRole = info.SubmitRole
	}
	if info.ResetLabel != "" {
		merged.resetLabel = info.ResetLabel
	}
	if info.ResetClass != "" {
		merged.resetClass = info.ResetClass
	}
	if info.ResetAttrs != "" {
		merged.resetAttrs = info.ResetAttrs
	}
	if info.ResetRole != "" {
		merged.resetRole = info.ResetRole
	}

	return merged
}

// parseWithCache parses struct tags with caching for performance.
// On the first call for a given struct type, it parses the tags and caches the field
// descriptors. Subsequent calls for the same type reuse the cached metadata but
// re-extract values from the provided instance.
func (r *FormRenderer) parseWithCache(v interface{}) ([]formField, *FormRenderer, error) {
	val := reflect.ValueOf(v)
	typ := val.Type()

	// Dereference pointers
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, nil, fmt.Errorf("wtf: cannot render nil pointer")
		}
		val = val.Elem()
		typ = val.Type()
	}

	if typ.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("wtf: expected struct, got %s", typ.Kind())
	}

	// Load from cache or parse metadata
	var metadata structMetadata
	if cached, ok := r.cache.Load(typ); ok {
		metadata = cached.(structMetadata)
	} else {
		var err error
		metadata, err = parseStructMetadata(typ)
		if err != nil {
			return nil, nil, err
		}
		r.cache.Store(typ, metadata)
	}

	// Extract dynamic FormInfo value if present
	dynamicInfo := FormInfo{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.Type == reflect.TypeOf(FormInfo{}) {
			dynamicInfo = val.Field(i).Interface().(FormInfo)
			break
		}
	}

	// Merge static and dynamic FormInfo config
	mergedInfo := mergeFormInfoValues(metadata.FormInfo, dynamicInfo)
	renderer := r.mergeFormInfo(mergedInfo)

	// Build form fields with dynamic values
	fields := make([]formField, len(metadata.Fields))
	for i, m := range metadata.Fields {
		fieldValue := val.Field(m.Index)
		valStr := extractValue(fieldValue, m.Type)

		// Copy options and set dynamic Selected state
		opts := make([]selectOption, len(m.Options))
		for j, opt := range m.Options {
			opts[j] = selectOption{
				Value:    opt.Value,
				Label:    opt.Label,
				Selected: opt.Value == valStr,
			}
		}

		fields[i] = formField{
			Name:        m.Name,
			Label:       m.Label,
			Type:        m.Type,
			ID:          m.ID,
			Value:       valStr,
			Required:    m.Required,
			Disabled:    m.Disabled,
			ReadOnly:    m.ReadOnly,
			Pattern:     m.Pattern,
			Min:         m.Min,
			Max:         m.Max,
			Step:        m.Step,
			Placeholder: m.Placeholder,
			Class:       m.Class,
			Options:     opts,
			Rows:        m.Rows,
			Cols:        m.Cols,
			Multiple:    m.Multiple,
			Role:        m.Role,
			FieldType:   m.FieldType,
		}
	}

	return fields, renderer, nil
}

// prefix returns a CSS class name with the configured prefix.
func (r *FormRenderer) prefix(class string) string {
	if r.classPrefix == "" {
		return class
	}
	return r.classPrefix + "-" + class
}

// MustFuncMap is like FuncMap but panics if there is an error.
// This is a convenience for template.Must usage patterns.
func (r *FormRenderer) MustFuncMap() template.FuncMap {
	return r.FuncMap()
}

// RenderField renders a single struct field as HTML. This is useful when you want
// more control over form layout in your templates.
func (r *FormRenderer) RenderField(v interface{}, fieldName string) template.HTML {
	if v == nil {
		return template.HTML("<!-- wtf: cannot render nil value -->")
	}

	fields, renderer, err := r.parseWithCache(v)
	if err != nil {
		return template.HTML("<!-- wtf: " + err.Error() + " -->")
	}

	for _, f := range fields {
		if f.Name == fieldName {
			return template.HTML(renderField(f, renderer))
		}
	}

	return template.HTML("<!-- wtf: field '" + fieldName + "' not found -->")
}

// Fields returns the parsed form fields for a struct. This is useful for
// custom rendering in templates using range.
func (r *FormRenderer) Fields(v interface{}) []FormFieldInfo {
	if v == nil {
		return nil
	}

	fields, renderer, err := r.parseWithCache(v)
	if err != nil {
		return nil
	}

	result := make([]FormFieldInfo, len(fields))
	for i, f := range fields {
		result[i] = FormFieldInfo{
			Name:        f.Name,
			Label:       f.Label,
			Type:        f.Type,
			ID:          f.ID,
			Value:       f.Value,
			Required:    f.Required,
			Disabled:    f.Disabled,
			ReadOnly:    f.ReadOnly,
			Placeholder: f.Placeholder,
			Role:        f.Role,
			HTML:        template.HTML(renderField(f, renderer)),
		}
	}

	return result
}

// FormFieldInfo is a public-facing struct that exposes parsed field information
// for use in templates with custom rendering.
type FormFieldInfo struct {
	Name        string
	Label       string
	Type        string
	ID          string
	Value       string
	Required    bool
	Disabled    bool
	ReadOnly    bool
	Placeholder string
	Role        string
	HTML        template.HTML // Pre-rendered HTML for this field
}

// FieldNames returns a list of form field names for a struct.
func (r *FormRenderer) FieldNames(v interface{}) []string {
	if v == nil {
		return nil
	}

	val := reflect.ValueOf(v)
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	typ := val.Type()
	var names []string

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		
		// Skip FormInfo field
		if field.Type == reflect.TypeOf(FormInfo{}) {
			continue
		}

		tag := field.Tag.Get("form")
		if tag == "" || tag == "-" || !field.IsExported() {
			continue
		}
		parts := splitTag(tag)
		if len(parts) > 0 {
			names = append(names, parts[0])
		}
	}

	return names
}
