# go-wtf

**Web Template Forms** — Automatically render HTML forms from Go struct tags.

`go-wtf` reads `form` struct tags on your Go structs and generates clean, accessible, styled HTML forms. It integrates seamlessly with Go's `html/template` engine via a custom `FuncMap`.

## Installation

```bash
go get github.com/malnossi/go-wtf
```

## Quick Start

### 1. Define your form as a Go struct

```go
type LoginForm struct {
    Username string `form:"username,type=text,placeholder=Enter username,required"`
    Password string `form:"password,type=password,placeholder=Enter password,required"`
    Remember bool   `form:"remember,type=checkbox,label=Remember me"`
}
```

### 2. Create a renderer and register it with your template

```go
package main

import (
    "html/template"
    "net/http"

    wtf "github.com/malnossi/go-wtf"
)

func main() {
    renderer := wtf.New(
        wtf.WithAction("/login"),
        wtf.WithMethod("POST"),
    )

    tmpl := template.Must(
        template.New("page").Funcs(renderer.FuncMap()).Parse(pageHTML),
    )

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := struct {
            Form LoginForm
        }{
            Form: LoginForm{Username: "defaultuser"},
        }
        tmpl.Execute(w, data)
    })

    http.ListenAndServe(":8080", nil)
}
```

### 3. Use `render_form` in your HTML template

```html
<!DOCTYPE html>
<html>
<body>
    <h1>Login</h1>
    {{render_form .Form}}
</body>
</html>
```

That's it! The form is rendered with labels, input fields, validation attributes, and built-in styling.

## Struct Tag Syntax

Tags follow the format:

```
form:"name,key=value,key=value,flag"
```

- **First value** is the field name (used as the HTML `name` attribute)
- Subsequent values are key=value pairs or boolean flags

### Supported Tag Options

| Key | Description | Example |
|-----|-------------|---------|
| `type` | HTML input type | `type=email` |
| `label` | Label text (auto-generated from name if omitted) | `label=Email Address` |
| `placeholder` | Placeholder text | `placeholder=you@example.com` |
| `id` | Custom HTML id (defaults to `form-{name}`) | `id=user-email` |
| `class` | Additional CSS class for the input | `class=my-input` |
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

If no `type` is specified, the type is inferred from the Go type:

| Go Type | HTML Input Type |
|---------|-----------------|
| `string` | `text` |
| `int`, `int64`, etc. | `number` |
| `uint`, `uint64`, etc. | `number` |
| `float32`, `float64` | `number` |
| `bool` | `checkbox` |
| `time.Time` | `datetime-local` |

### Select/Radio Options Format

Options use pipe-separated values with an optional `value:label` format:

```go
// Simple options (value = label, auto-title-cased)
Subject string `form:"subject,type=select,options=general|support|billing"`

// Labeled options (value:label)
Subject string `form:"subject,type=select,options=general:General Inquiry|support:Technical Support"`
```

## Declaring Form Options & Button Attributes (`FormInfo`)

Instead of configuring options strictly on the `FormRenderer` instance, you can define form-level properties, custom button classes, a reset button, and raw HTML attributes (like HTMX `hx-post`, `hx-target`, or Unpoly `up-dismiss`) directly inside the form struct itself using the `wtf.FormInfo` struct field.

### Usage

Define a field of type `wtf.FormInfo` inside your form struct:

```go
type MyForm struct {
    // 1. Configure statically via struct tags
    Info wtf.FormInfo `form:"_info,action=/submit,submit_label=Register,submit_class=btn-primary,submit_attrs=up-dismiss hx-post='/register'"`
    
    Username string `form:"username,required"`
}
```

Or configure it dynamically at runtime when instantiating the struct:

```go
form := MyForm{
    // 2. Configure dynamically at runtime (takes priority over tags)
    Info: wtf.FormInfo{
        Action:      "/dynamic-submit",
        SubmitLabel: "Send Message",
        SubmitAttrs: `hx-post="/contact" up-dismiss hx-target="#result"`,
        ResetLabel:  "Clear Form",
        ResetClass:  "btn-secondary",
        ResetAttrs:  `hx-get="/clear"`,
        FormAttrs:   `hx-target="#div"`,
    },
    Username: "bob",
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
| `SubmitLabel` | `submit_label` | Text for submit button | `submit_label=Save` |
| `SubmitClass` | `submit_class` | Additional CSS class for submit | `submit_class=btn-success` |
| `SubmitAttrs` | `submit_attrs` | Raw HTML attributes on submit button | `submit_attrs=up-dismiss hx-post="/save"` |
| `ResetLabel` | `reset_label` | Text for reset button (rendered if set) | `reset_label=Cancel` |
| `ResetClass` | `reset_class` | Additional CSS class for reset | `reset_class=btn-outline` |
| `ResetAttrs` | `reset_attrs` | Raw HTML attributes on reset button | `reset_attrs=hx-get="/cancel"` |

---

## Renderer Options

Configure the renderer with functional options:

```go
renderer := wtf.New(
    wtf.WithAction("/submit"),          // Form action URL
    wtf.WithMethod("POST"),             // HTTP method (default: POST)
    wtf.WithCSRF("token123"),           // Add hidden CSRF field
    wtf.WithClassPrefix("myapp"),       // CSS class prefix (default: wtf)
    wtf.WithSubmitLabel("Sign In"),     // Submit button text (default: Submit)
    wtf.WithFormID("login-form"),       // Form element ID
    wtf.WithFormClass("custom-form"),   // Additional form CSS class
    wtf.WithNoSubmit(true),             // Don't render submit button
    wtf.WithEnctype("multipart/form-data"), // For file uploads
)
```

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

When using `renderer.Fields()`, you get access to field metadata:

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
    HTML        template.HTML  // Pre-rendered HTML
}
```

## Pre-populated Values

Struct fields are automatically pre-populated in the form:

```go
form := LoginForm{
    Username: "john@example.com",  // Input will have value="john@example.com"
    Remember: true,                // Checkbox will be checked
}
```

## Custom Rendering with Fields

For full control over layout, use `Fields()` to iterate over field info:

```go
renderer := wtf.New()

funcMap := template.FuncMap{
    "form_fields": renderer.Fields,
}

tmpl := template.Must(
    template.New("page").Funcs(funcMap).Parse(`
    <form method="post">
        {{range form_fields .Form}}
            <div class="my-field-wrapper">
                {{.HTML}}
            </div>
        {{end}}
        <button type="submit">Submit</button>
    </form>
    `),
)
```

## Security

- All output values are properly HTML-escaped using `html.EscapeString`
- XSS attacks through struct values are prevented
- CSRF tokens can be added via `WithCSRF()`

## CSS Classes & Styling

`go-wtf` produces clean, semantic, styling-free HTML form elements. This allows you to style your forms entirely using your own custom CSS or frameworks (like Tailwind CSS or Bootstrap) using the generated CSS classes. The default class prefix is `wtf-` (customizable via `WithClassPrefix()`).

CSS classes added to the HTML markup:
- `wtf-form` — The form element
- `wtf-form-group` — Field wrapper div
- `wtf-form-label` — Label elements
- `wtf-form-input` — Input elements
- `wtf-form-select` — Select elements
- `wtf-form-textarea` — Textarea elements
- `wtf-form-checkbox` / `wtf-form-checkbox-group` — Checkbox elements
- `wtf-form-radio` / `wtf-form-radio-group` / `wtf-form-radio-option` — Radio elements
- `wtf-form-required` — Required field indicator (asterisk)
- `wtf-form-submit` — Submit button
- `wtf-form-reset` — Reset button

## Example

See the [example](./example/main.go) directory for a complete working demo with multiple form types.

```bash
go run ./example/
# Visit http://localhost:8080
```

## License

MIT
