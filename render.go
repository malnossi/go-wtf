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

	if r.formClass != "" {
		fmt.Fprintf(&b, ` class="%s"`, html.EscapeString(r.formClass))
	}

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
		fmt.Fprintf(&b, `  <input type="submit" value="%s"`, html.EscapeString(r.submitLabel))
		if r.submitClass != "" {
			fmt.Fprintf(&b, ` class="%s"`, html.EscapeString(r.submitClass))
		}
		if r.submitRole != "" {
			fmt.Fprintf(&b, ` role="%s"`, html.EscapeString(r.submitRole))
		}
		if r.submitAttrs != "" {
			fmt.Fprintf(&b, " %s", r.submitAttrs)
		}
		b.WriteString(">\n")
	}

	// Reset button
	if r.resetLabel != "" {
		fmt.Fprintf(&b, `  <input type="reset" value="%s"`, html.EscapeString(r.resetLabel))
		if r.resetClass != "" {
			fmt.Fprintf(&b, ` class="%s"`, html.EscapeString(r.resetClass))
		}
		if r.resetRole != "" {
			fmt.Fprintf(&b, ` role="%s"`, html.EscapeString(r.resetRole))
		}
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
	fmt.Fprintf(&b, "  <div>\n")

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
	fmt.Fprintf(&b, `    <label for="%s">%s`,
		html.EscapeString(f.ID),
		html.EscapeString(f.Label))
	if f.Required {
		b.WriteString("<span>*</span>")
	}
	b.WriteString("</label>\n")
	return b.String()
}

// renderInput renders a standard <input> element.
func renderInput(f formField, r *FormRenderer) string {
	var b strings.Builder

	fmt.Fprintf(&b, `    <input type="%s" id="%s" name="%s"`,
		html.EscapeString(f.Type),
		html.EscapeString(f.ID),
		html.EscapeString(f.Name))

	if f.Class != "" {
		fmt.Fprintf(&b, ` class="%s"`, html.EscapeString(f.Class))
	}

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

	fmt.Fprintf(&b, `    <textarea id="%s" name="%s"`,
		html.EscapeString(f.ID),
		html.EscapeString(f.Name))

	if f.Class != "" {
		fmt.Fprintf(&b, ` class="%s"`, html.EscapeString(f.Class))
	}

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

	fmt.Fprintf(&b, `    <select id="%s" name="%s"`,
		html.EscapeString(f.ID),
		html.EscapeString(f.Name))

	if f.Class != "" {
		fmt.Fprintf(&b, ` class="%s"`, html.EscapeString(f.Class))
	}

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

	fmt.Fprintf(&b, "  <div>\n")
	fmt.Fprintf(&b, "    <div>\n")

	fmt.Fprintf(&b, `      <input type="checkbox" id="%s" name="%s"`,
		html.EscapeString(f.ID),
		html.EscapeString(f.Name))

	if f.Class != "" {
		fmt.Fprintf(&b, ` class="%s"`, html.EscapeString(f.Class))
	}

	if f.Role != "" {
		fmt.Fprintf(&b, ` role="%s"`, html.EscapeString(f.Role))
	}

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

	fmt.Fprintf(&b, `      <label for="%s">%s</label>`+"\n",
		html.EscapeString(f.ID),
		html.EscapeString(f.Label))

	b.WriteString("    </div>\n")
	b.WriteString("  </div>\n")
	return b.String()
}

// renderRadio renders a group of radio buttons.
func renderRadio(f formField, r *FormRenderer) string {
	var b strings.Builder

	fmt.Fprintf(&b, "  <div>\n")

	// Label for the group
	b.WriteString(renderLabel(f, r))

	fmt.Fprintf(&b, "    <div>\n")

	for i, opt := range f.Options {
		optID := fmt.Sprintf("%s-%d", f.ID, i)

		fmt.Fprintf(&b, "      <div>\n")

		fmt.Fprintf(&b, `        <input type="radio" id="%s" name="%s" value="%s"`,
			html.EscapeString(optID),
			html.EscapeString(f.Name),
			html.EscapeString(opt.Value))

		if f.Class != "" {
			fmt.Fprintf(&b, ` class="%s"`, html.EscapeString(f.Class))
		}

		if f.Role != "" {
			fmt.Fprintf(&b, ` role="%s"`, html.EscapeString(f.Role))
		}

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

		fmt.Fprintf(&b, `        <label for="%s">%s</label>`+"\n",
			html.EscapeString(optID),
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
