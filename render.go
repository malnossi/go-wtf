package wtf

import (
	"fmt"
	"html"
	"strings"
)

// renderForm builds the complete HTML string for a form from parsed fields.
func renderForm(fields []formField, r *FormRenderer) string {
	var b strings.Builder
	// Open <form> tag
	b.WriteString("<form")
	if r.formID != "" {
		fmt.Fprintf(&b, ` id="%s"`, html.EscapeString(r.formID))
	}

	// Build class list
	formClass := r.prefix("form")
	if r.formClass != "" {
		formClass += " " + r.formClass
	}
	fmt.Fprintf(&b, ` class="%s"`, html.EscapeString(formClass))

	if r.action != "" {
		fmt.Fprintf(&b, ` action="%s"`, html.EscapeString(r.action))
	}
	fmt.Fprintf(&b, ` method="%s"`, html.EscapeString(r.method))

	if r.enctype != "" {
		fmt.Fprintf(&b, ` enctype="%s"`, html.EscapeString(r.enctype))
	}

	// Render custom raw form attributes (e.g. hx-target, etc.)
	if r.formAttrs != "" {
		fmt.Fprintf(&b, " %s", r.formAttrs)
	}

	b.WriteString(">\n")

	// CSRF token hidden field
	if r.csrfToken != "" {
		fmt.Fprintf(&b, `  <input type="hidden" name="_csrf" value="%s">`+"\n", html.EscapeString(r.csrfToken))
	}

	// Render each field
	for _, field := range fields {
		b.WriteString(renderField(field, r))
	}

	// Submit button
	if !r.noSubmit {
		submitClass := r.prefix("form-submit")
		if r.submitClass != "" {
			submitClass += " " + r.submitClass
		}
		fmt.Fprintf(&b, `  <input type="submit" class="%s" value="%s"`, 
			html.EscapeString(submitClass), 
			html.EscapeString(r.submitLabel))
		if r.submitAttrs != "" {
			fmt.Fprintf(&b, " %s", r.submitAttrs)
		}
		b.WriteString(">\n")
	}

	// Reset button
	if r.resetLabel != "" {
		resetClass := r.prefix("form-reset")
		if r.resetClass != "" {
			resetClass += " " + r.resetClass
		}
		fmt.Fprintf(&b, `  <input type="reset" class="%s" value="%s"`, 
			html.EscapeString(resetClass), 
			html.EscapeString(r.resetLabel))
		if r.resetAttrs != "" {
			fmt.Fprintf(&b, " %s", r.resetAttrs)
		}
		b.WriteString(">\n")
	}

	b.WriteString("</form>")

	return b.String()
}

// renderField renders a single form field wrapped in a form group div.
func renderField(f formField, r *FormRenderer) string {
	var b strings.Builder

	switch f.Type {
	case "hidden":
		b.WriteString(renderHidden(f, r))
		return b.String()
	case "checkbox":
		b.WriteString(renderCheckbox(f, r))
		return b.String()
	}

	// Wrap in form group
	fmt.Fprintf(&b, `  <div class="%s">`+"\n", html.EscapeString(r.prefix("form-group")))

	// Label
	b.WriteString(renderLabel(f, r))

	// Input element
	switch f.Type {
	case "textarea":
		b.WriteString(renderTextarea(f, r))
	case "select":
		b.WriteString(renderSelect(f, r))
	case "radio":
		b.WriteString(renderRadio(f, r))
	default:
		b.WriteString(renderInput(f, r))
	}

	b.WriteString("  </div>\n")
	return b.String()
}

// renderLabel renders a <label> element for a field.
func renderLabel(f formField, r *FormRenderer) string {
	var b strings.Builder
	fmt.Fprintf(&b, `    <label for="%s" class="%s">%s`,
		html.EscapeString(f.ID),
		html.EscapeString(r.prefix("form-label")),
		html.EscapeString(f.Label))
	if f.Required {
		fmt.Fprintf(&b, `<span class="%s">*</span>`, html.EscapeString(r.prefix("form-required")))
	}
	b.WriteString("</label>\n")
	return b.String()
}

// renderInput renders a standard <input> element.
func renderInput(f formField, r *FormRenderer) string {
	var b strings.Builder
	inputClass := r.prefix("form-input")
	if f.Class != "" {
		inputClass += " " + f.Class
	}

	fmt.Fprintf(&b, `    <input type="%s" id="%s" name="%s" class="%s"`,
		html.EscapeString(f.Type),
		html.EscapeString(f.ID),
		html.EscapeString(f.Name),
		html.EscapeString(inputClass))

	if f.Role != "" {
		fmt.Fprintf(&b, ` role="%s"`, html.EscapeString(f.Role))
	}

	if f.Value != "" {
		fmt.Fprintf(&b, ` value="%s"`, html.EscapeString(f.Value))
	}
	if f.Placeholder != "" {
		fmt.Fprintf(&b, ` placeholder="%s"`, html.EscapeString(f.Placeholder))
	}
	if f.Required {
		b.WriteString(` required`)
	}
	if f.Disabled {
		b.WriteString(` disabled`)
	}
	if f.ReadOnly {
		b.WriteString(` readonly`)
	}
	if f.Pattern != "" {
		fmt.Fprintf(&b, ` pattern="%s"`, html.EscapeString(f.Pattern))
	}
	if f.Min != "" {
		fmt.Fprintf(&b, ` min="%s"`, html.EscapeString(f.Min))
	}
	if f.Max != "" {
		fmt.Fprintf(&b, ` max="%s"`, html.EscapeString(f.Max))
	}
	if f.Step != "" {
		fmt.Fprintf(&b, ` step="%s"`, html.EscapeString(f.Step))
	}

	b.WriteString(">\n")
	return b.String()
}

// renderTextarea renders a <textarea> element.
func renderTextarea(f formField, r *FormRenderer) string {
	var b strings.Builder
	textareaClass := r.prefix("form-textarea")
	if f.Class != "" {
		textareaClass += " " + f.Class
	}

	fmt.Fprintf(&b, `    <textarea id="%s" name="%s" class="%s"`,
		html.EscapeString(f.ID),
		html.EscapeString(f.Name),
		html.EscapeString(textareaClass))

	if f.Role != "" {
		fmt.Fprintf(&b, ` role="%s"`, html.EscapeString(f.Role))
	}

	if f.Rows != "" {
		fmt.Fprintf(&b, ` rows="%s"`, html.EscapeString(f.Rows))
	}
	if f.Cols != "" {
		fmt.Fprintf(&b, ` cols="%s"`, html.EscapeString(f.Cols))
	}
	if f.Placeholder != "" {
		fmt.Fprintf(&b, ` placeholder="%s"`, html.EscapeString(f.Placeholder))
	}
	if f.Required {
		b.WriteString(` required`)
	}
	if f.Disabled {
		b.WriteString(` disabled`)
	}
	if f.ReadOnly {
		b.WriteString(` readonly`)
	}

	fmt.Fprintf(&b, ">%s</textarea>\n", html.EscapeString(f.Value))
	return b.String()
}

// renderSelect renders a <select> element with <option> children.
func renderSelect(f formField, r *FormRenderer) string {
	var b strings.Builder
	selectClass := r.prefix("form-select")
	if f.Class != "" {
		selectClass += " " + f.Class
	}

	fmt.Fprintf(&b, `    <select id="%s" name="%s" class="%s"`,
		html.EscapeString(f.ID),
		html.EscapeString(f.Name),
		html.EscapeString(selectClass))

	if f.Role != "" {
		fmt.Fprintf(&b, ` role="%s"`, html.EscapeString(f.Role))
	}

	if f.Required {
		b.WriteString(` required`)
	}
	if f.Disabled {
		b.WriteString(` disabled`)
	}
	if f.Multiple {
		b.WriteString(` multiple`)
	}

	b.WriteString(">\n")

	// Add a default empty option if not required to have a value
	if !f.Required {
		fmt.Fprintf(&b, `      <option value="">-- Select %s --</option>`+"\n", html.EscapeString(f.Label))
	}

	for _, opt := range f.Options {
		if opt.Selected {
			fmt.Fprintf(&b, `      <option value="%s" selected>%s</option>`+"\n",
				html.EscapeString(opt.Value),
				html.EscapeString(opt.Label))
		} else {
			fmt.Fprintf(&b, `      <option value="%s">%s</option>`+"\n",
				html.EscapeString(opt.Value),
				html.EscapeString(opt.Label))
		}
	}

	b.WriteString("    </select>\n")
	return b.String()
}

// renderCheckbox renders a checkbox input with an inline label.
func renderCheckbox(f formField, r *FormRenderer) string {
	var b strings.Builder

	fmt.Fprintf(&b, `  <div class="%s">`+"\n", html.EscapeString(r.prefix("form-group")))
	fmt.Fprintf(&b, `    <div class="%s">`+"\n", html.EscapeString(r.prefix("form-checkbox-group")))

	inputClass := r.prefix("form-checkbox")
	if f.Class != "" {
		inputClass += " " + f.Class
	}

	fmt.Fprintf(&b, `      <input type="checkbox" id="%s" name="%s" class="%s"`,
		html.EscapeString(f.ID),
		html.EscapeString(f.Name),
		html.EscapeString(inputClass))

	if f.Value == "true" {
		b.WriteString(` checked`)
	}
	if f.Required {
		b.WriteString(` required`)
	}
	if f.Disabled {
		b.WriteString(` disabled`)
	}

	b.WriteString(" value=\"true\">\n")

	fmt.Fprintf(&b, `      <label for="%s" class="%s">%s</label>`+"\n",
		html.EscapeString(f.ID),
		html.EscapeString(r.prefix("form-checkbox-label")),
		html.EscapeString(f.Label))

	b.WriteString("    </div>\n")
	b.WriteString("  </div>\n")
	return b.String()
}

// renderRadio renders a group of radio buttons.
func renderRadio(f formField, r *FormRenderer) string {
	var b strings.Builder

	fmt.Fprintf(&b, `  <div class="%s">`+"\n", html.EscapeString(r.prefix("form-group")))

	// Label for the group
	b.WriteString(renderLabel(f, r))

	fmt.Fprintf(&b, `    <div class="%s">`+"\n", html.EscapeString(r.prefix("form-radio-group")))

	for i, opt := range f.Options {
		optID := fmt.Sprintf("%s-%d", f.ID, i)

		fmt.Fprintf(&b, `      <div class="%s">`+"\n", html.EscapeString(r.prefix("form-radio-option")))

		radioClass := r.prefix("form-radio")
		fmt.Fprintf(&b, `        <input type="radio" id="%s" name="%s" value="%s" class="%s"`,
			html.EscapeString(optID),
			html.EscapeString(f.Name),
			html.EscapeString(opt.Value),
			html.EscapeString(radioClass))

		if opt.Selected {
			b.WriteString(` checked`)
		}
		if f.Required {
			b.WriteString(` required`)
		}
		if f.Disabled {
			b.WriteString(` disabled`)
		}
		b.WriteString(">\n")

		fmt.Fprintf(&b, `        <label for="%s" class="%s">%s</label>`+"\n",
			html.EscapeString(optID),
			html.EscapeString(r.prefix("form-radio-label")),
			html.EscapeString(opt.Label))

		b.WriteString("      </div>\n")
	}

	b.WriteString("    </div>\n")
	b.WriteString("  </div>\n")
	return b.String()
}

// renderHidden renders a hidden input field (no label, no wrapper).
func renderHidden(f formField, r *FormRenderer) string {
	return fmt.Sprintf(`  <input type="hidden" id="%s" name="%s" value="%s">`+"\n",
		html.EscapeString(f.ID),
		html.EscapeString(f.Name),
		html.EscapeString(f.Value))
}
