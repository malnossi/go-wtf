# go-wtf

**Web Template Forms** — Pure semantic, styling-free, class-free, and wrapper-free HTML form generation from Go structs.

`go-wtf` reads struct tags on your Go structs and generates clean, accessible HTML5 forms. It integrates seamlessly with Go's `html/template` engine via a custom `FuncMap`, rendering standard inputs, textareas, selects, radios, and checkboxes with a simple `{{render_form .Form}}` template call.

---

## Features

- 🌿 **Pure Semantic HTML**: Zero default CSS classes (no `.wtf-form` or `.wtf-form-input`) and absolutely no predefined stylesheets or styles.
- 📦 **Flat & Wrapper-Free**: Labels and inputs are rendered directly next to each other in a flat HTML structure with **no wrapper `<div>` elements** (including checkboxes and radios).
- 🔐 **Standard HTML5 Self-Closing Syntax**: Standard self-closing formatting (` />`) for all void tags like `<input>` (text, checkbox, radio, submit, reset, hidden).
- 🧩 **`<fieldset>` Grouping**: Embed `Fieldset = true` inside `wtf.FormInfo` to wrap your fields in a `<fieldset>` container, keeping action controls (submit/reset) clean and outside.
- 🗃️ **Label-First Checkbox Order**: Renders labels *before* the checkbox inputs, keeping markup ordering consistent with other field types.
- 🏷️ **Native ARIA Roles**: Parse and render standard `role="..."` accessibility attributes on any form input.
- 🦾 **HTML5 Autocomplete**: Support for native `autocomplete="..."` tags on form fields.
- 🔄 **FormInfo Config & Runtime Overrides**: Configure actions, methods, custom button raw attributes, classes, and ARIA roles statically via tags or dynamically at runtime on `wtf.FormInfo`.

---

## Installation

```bash
go get github.com/malnossi/go-wtf
```

---

## Quick Start

### 1. Define your form as a Go struct

Embed `wtf.FormInfo` to configure form-level properties, submit buttons, reset buttons, and layouts:

```go
type SubscriptionForm struct {
    // 1. Static config: wrap fields in <fieldset>, set custom submit/reset labels and roles
    Info      wtf.FormInfo `form:"_info,fieldset=true,submit_label=Subscribe,submit_role=button,reset_label=Clear,reset_role=button"`
    
    FirstName string       `form:"first_name,placeholder=First name,autocomplete=given-name,required"`
    Email     string       `form:"email,type=email,placeholder=Email,autocomplete=email,required"`
    Remember  bool         `form:"remember,label=Remember me,role=switch"`
}
```

### 2. Register with your template and execute

```go
package main

import (
    "html/template"
    "net/http"

    wtf "github.com/malnossi/go-wtf"
)

var pageHTML = `
<!DOCTYPE html>
<html>
<body>
    <h2>Newsletter Subscription</h2>
    {{render_form .Form}}
</body>
</html>`

func main() {
    renderer := wtf.New()

    tmpl := template.Must(
        template.New("page").Funcs(renderer.FuncMap()).Parse(pageHTML),
    )

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := struct {
            Form SubscriptionForm
        }{
            Form: SubscriptionForm{},
        }
        tmpl.Execute(w, data)
    })

    http.ListenAndServe(":8080", nil)
}
```

### Generated Plain HTML Output:

```html
<form method="POST">
  <fieldset>
  <label for="form-first_name">First name<span>*</span></label>
  <input type="text" id="form-first_name" name="first_name" placeholder="First name" autocomplete="given-name" required />
  <label for="form-email">Email<span>*</span></label>
  <input type="email" id="form-email" name="email" placeholder="Email" autocomplete="email" required />
  <label for="form-remember">Remember me</label>
  <input type="checkbox" id="form-remember" name="remember" role="switch" value="true" />
  </fieldset>
  <input type="submit" value="Subscribe" role="button" />
  <input type="reset" value="Clear" role="button" />
</form>
```

---

## Struct Tag Syntax

Struct tags follow this format:

```
form:"name,key=value,key=value,flag"
```

- **First value**: The HTML field `name` attribute.
- **Subsequent values**: Key-value options or standalone flags.

### Supported Field Options

| Key | Description | Example |
|-----|-------------|---------|
| `type` | HTML input type | `type=email` |
| `label` | Label text (auto-generated from name if omitted) | `label=Email Address` |
| `placeholder` | Placeholder text | `placeholder=you@example.com` |
| `id` | Custom HTML id (defaults to `form-{name}`) | `id=user-email` |
| `class` | Custom CSS class for the input element | `class=my-input` |
| `autocomplete` | HTML autocomplete attribute value | `autocomplete=email` |
| `role` | ARIA accessibility role attribute | `role=switch` |
| `required` | Mark field as required | `required` |
| `disabled` | Mark field as disabled | `disabled` |
| `readonly` | Mark field as read-only | `readonly` |
| `min` | Minimum value (for number inputs) | `min=0` |
| `max` | Maximum value (for number inputs) | `max=100` |
| `step` | Step value (for number inputs) | `step=0.01` |
| `pattern` | Validation regex pattern | `pattern=[0-9]{3}` |
| `options` | Options for select/radio (pipe-separated) | `options=a:Label A\|b:Label B` |
| `rows` | Number of rows (for textarea) | `rows=5` |
| `cols` | Number of columns (for textarea) | `cols=50` |
| `multiple` | Allow multiple selections (for select) | `multiple` |
| `-` | Skip this field entirely | `form:"-"` |

### Supported Input Types

`text`, `password`, `email`, `number`, `tel`, `url`, `date`, `datetime-local`, `time`, `color`, `range`, `search`, `hidden`, `file`, `textarea`, `select`, `radio`, `checkbox`

### Automatic Type Inference

If no `type` is specified, the type is inferred from the Go struct field type:

| Go Type | HTML Input Type |
|---------|-----------------|
| `string` | `text` |
| `int`, `int64`, etc. | `number` |
| `uint`, `uint64`, etc. | `number` |
| `float32`, `float64` | `number` |
| `bool` | `checkbox` |
| `time.Time` | `datetime-local` |

---

## Form-Level Options & Runtime Overrides (`FormInfo`)

Define form-level configurations directly inside the form struct using the `wtf.FormInfo` struct field:

```go
type ContactForm struct {
    // Declared statically via tag:
    Info    wtf.FormInfo `form:"_info,action=/submit,submit_label=Send"`
    
    Name    string `form:"name"`
    Message string `form:"message,type=textarea"`
}
```

Or configure it dynamically at runtime when instantiating the struct (takes priority over tags):

```go
form := ContactForm{
    Info: wtf.FormInfo{
        Action:      "/dynamic-endpoint",
        Fieldset:    true,
        SubmitLabel: "Send Message",
        SubmitAttrs: `hx-post="/contact" up-dismiss`,
        ResetLabel:  "Clear Form",
        ResetRole:   "button",
    },
    Name: "Alice",
}
```

### `FormInfo` Struct Specification

| Field | Tag Option | Description | Example |
|---|---|---|---|
| `Action` | `action` | Form action URL | `action=/submit` |
| `Method` | `method` | HTTP Method (POST/GET) | `method=GET` |
| `Enctype` | `enctype` | Form enctype | `enctype=multipart/form-data` |
| `FormID` | `id` | Form element ID | `id=my-form` |
| `FormClass` | `class` | Custom CSS class for `<form>` | `class=form-dark` |
| `FormAttrs` | `attrs` | Raw HTML attributes on `<form>` | `attrs=hx-target="#div"` |
| `Fieldset`  | `fieldset` | Wrap form fields in a `<fieldset>` element | `fieldset=true` |
| `SubmitLabel` | `submit_label` | Text for submit button | `submit_label=Save` |
| `SubmitClass` | `submit_class` | Additional CSS class for submit | `submit_class=btn-success` |
| `SubmitAttrs` | `submit_attrs` | Raw HTML attributes on submit button | `submit_attrs=up-dismiss hx-post="/save"` |
| `SubmitRole` | `submit_role` | ARIA Accessibility role attribute for submit input | `submit_role=button` |
| `ResetLabel` | `reset_label` | Text for reset button (rendered if set) | `reset_label=Cancel` |
| `ResetClass` | `reset_class` | Additional CSS class for reset | `reset_class=btn-outline` |
| `ResetAttrs` | `reset_attrs` | Raw HTML attributes on reset button | `reset_attrs=hx-get="/cancel"` |
| `ResetRole` | `reset_role` | ARIA Accessibility role attribute for reset input | `reset_role=button` |

---

## API Reference

### `FormRenderer`

```go
// Create a new renderer
renderer := wtf.New(opts ...Option) *FormRenderer

// Get the template FuncMap (registers "render_form")
funcMap := renderer.FuncMap() template.FuncMap

// Render a struct as an HTML form
html := renderer.RenderForm(v interface{}) template.HTML

// Render with a specific action and method (one-off override)
html := renderer.RenderFormWithAction(v interface{}, action, method string) template.HTML

// Render a single field by name
html := renderer.RenderField(v interface{}, fieldName string) template.HTML

// Get parsed field info for custom template rendering
fields := renderer.Fields(v interface{}) []FormFieldInfo

// Get list of field names
names := renderer.FieldNames(v interface{}) []string
```

### `FormFieldInfo`

When using `renderer.Fields()`, you get access to parsed field metadata:

```go
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
    HTML        template.HTML  // Pre-rendered HTML for this field
}
```

---

## Styling & Layout Control

By default, `go-wtf` generates a flat, pure HTML output containing no styling or classes at all, returning standard elements (such as `  <label for="...">` and `  <input type="..." />`) under standard `2-space` prefixes.

This allows you to style your forms entirely using modern utility-first CSS frameworks like **Tailwind CSS**, standard CSS, or layouts by targeting the raw tag names or by applying custom classes directly in struct tags:

```go
type TailwindForm struct {
    Email string `form:"email,type=email,class=w-full px-4 py-2 border rounded-lg focus:ring focus:ring-blue-300"`
}
```

---

## Security

- **Strict HTML Escaping**: All output values, labels, IDs, classes, and ARIA attributes are fully escaped using `html.EscapeString` to prevent XSS vulnerability.
- **CSRF Token Injection**: Easily inject hidden CSRF inputs via renderer `WithCSRF(token)` option.

---

## Verification & Testing

Our codebase runs on **Go 1.26** with **100% test coverage** verified by a suite of **52 unit tests** ensuring correct ARIA role, fieldset wrapping, self-closing tag, and label-first checkbox layouts.

---

## Example

See the [example](./example/main.go) directory for a complete working demo server.

```bash
go run ./example/
# Visit http://localhost:8080
```

---

## License

MIT
