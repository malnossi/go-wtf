// Package wtf (Web Template Forms) provides automatic HTML form generation
// from Go structs using struct tags. It integrates with html/template via
// a FuncMap, enabling {{render_form .Form}} syntax in templates.
package wtf

// Option is a functional option for configuring the FormRenderer.
type Option func(*FormRenderer)

// WithAction sets the form's action attribute (the URL the form submits to).
func WithAction(action string) Option {
	return func(r *FormRenderer) {
		r.action = action
	}
}

// WithMethod sets the form's HTTP method (GET or POST). Defaults to "POST".
func WithMethod(method string) Option {
	return func(r *FormRenderer) {
		r.method = method
	}
}

// WithCSRF adds a hidden CSRF token field to the rendered form.
func WithCSRF(token string) Option {
	return func(r *FormRenderer) {
		r.csrfToken = token
	}
}

// WithClassPrefix sets a CSS class prefix for all generated HTML elements.
// For example, WithClassPrefix("my") would generate classes like "my-form-group",
// "my-form-input", etc.
func WithClassPrefix(prefix string) Option {
	return func(r *FormRenderer) {
		r.classPrefix = prefix
	}
}

// WithSubmitLabel sets the text for the form's submit button.
// Defaults to "Submit".
func WithSubmitLabel(label string) Option {
	return func(r *FormRenderer) {
		r.submitLabel = label
	}
}


// WithFormID sets the id attribute of the <form> element.
func WithFormID(id string) Option {
	return func(r *FormRenderer) {
		r.formID = id
	}
}

// WithFormClass sets additional CSS classes on the <form> element.
func WithFormClass(class string) Option {
	return func(r *FormRenderer) {
		r.formClass = class
	}
}

// WithNoSubmit disables the automatic submit button rendering.
func WithNoSubmit(noSubmit bool) Option {
	return func(r *FormRenderer) {
		r.noSubmit = noSubmit
	}
}

// WithEnctype sets the form's enctype attribute (e.g., "multipart/form-data" for file uploads).
func WithEnctype(enctype string) Option {
	return func(r *FormRenderer) {
		r.enctype = enctype
	}
}
