package wtf

import (
	"html/template"
	"reflect"
	"strings"
	"testing"
	"time"
)

// --- Test Structs ---

type LoginForm struct {
	Username string `form:"username,type=text,placeholder=Enter username,required"`
	Password string `form:"password,type=password,placeholder=Enter password,required"`
	Remember bool   `form:"remember,type=checkbox,label=Remember me"`
}

type ContactForm struct {
	Name    string `form:"name,type=text,label=Full Name,required,placeholder=Your name"`
	Email   string `form:"email,type=email,label=Email Address,required,placeholder=you@example.com"`
	Phone   string `form:"phone,type=tel,label=Phone Number,pattern=[0-9]{3}-[0-9]{3}-[0-9]{4}"`
	Subject string `form:"subject,type=select,options=general:General Inquiry|support:Technical Support|billing:Billing Question"`
	Message string `form:"message,type=textarea,rows=5,placeholder=Your message here...,required"`
}

type SettingsForm struct {
	DisplayName string `form:"display_name,type=text,label=Display Name"`
	Bio         string `form:"bio,type=textarea,rows=4,cols=50,placeholder=Tell us about yourself"`
	Theme       string `form:"theme,type=radio,options=light:Light Mode|dark:Dark Mode|auto:Auto"`
	Newsletter  bool   `form:"newsletter,type=checkbox,label=Subscribe to newsletter"`
	Age         int    `form:"age,type=number,min=13,max=120,label=Your Age"`
	Website     string `form:"website,type=url,placeholder=https://example.com"`
}

type HiddenFieldForm struct {
	ID     string `form:"id,type=hidden"`
	Token  string `form:"token,type=hidden"`
	Name   string `form:"name,type=text"`
}

type SkipFieldForm struct {
	Name     string `form:"name,type=text"`
	Internal string `form:"-"`
	NoTag    string
	Email    string `form:"email,type=email"`
}

type DisabledForm struct {
	ReadonlyField string `form:"readonly_field,type=text,readonly"`
	DisabledField string `form:"disabled_field,type=text,disabled"`
}

type StepForm struct {
	Price  float64 `form:"price,type=number,step=0.01,min=0"`
}

type DateForm struct {
	Birthday time.Time `form:"birthday,type=date,label=Date of Birth"`
	Meeting  time.Time `form:"meeting,type=datetime-local,label=Meeting Time"`
}

type InferredForm struct {
	Name    string  `form:"name"`
	Age     int     `form:"age"`
	Score   float64 `form:"score"`
	Active  bool    `form:"active"`
}

type CustomClassForm struct {
	Name  string `form:"name,type=text,class=my-custom-class"`
	Email string `form:"email,type=email,class=another-class"`
}

type MultiSelectForm struct {
	Colors string `form:"colors,type=select,options=red:Red|green:Green|blue:Blue,multiple"`
}

// --- Tests ---

func TestNew(t *testing.T) {
	r := New()
	if r.method != "POST" {
		t.Errorf("expected default method POST, got %s", r.method)
	}
	if r.submitLabel != "Submit" {
		t.Errorf("expected default submit label Submit, got %s", r.submitLabel)
	}
	if r.classPrefix != "wtf" {
		t.Errorf("expected default class prefix wtf, got %s", r.classPrefix)
	}
}

func TestNewWithOptions(t *testing.T) {
	r := New(
		WithAction("/login"),
		WithMethod("GET"),
		WithCSRF("mytoken123"),
		WithClassPrefix("myapp"),
		WithSubmitLabel("Sign In"),
		WithFormID("login-form"),
		WithFormClass("form-horizontal"),
		WithNoSubmit(true),
		WithEnctype("multipart/form-data"),
	)

	if r.action != "/login" {
		t.Errorf("expected action /login, got %s", r.action)
	}
	if r.method != "GET" {
		t.Errorf("expected method GET, got %s", r.method)
	}
	if r.csrfToken != "mytoken123" {
		t.Errorf("expected csrf token mytoken123, got %s", r.csrfToken)
	}
	if r.classPrefix != "myapp" {
		t.Errorf("expected class prefix myapp, got %s", r.classPrefix)
	}
	if r.submitLabel != "Sign In" {
		t.Errorf("expected submit label Sign In, got %s", r.submitLabel)
	}
	if r.formID != "login-form" {
		t.Errorf("expected form ID login-form, got %s", r.formID)
	}
	if r.formClass != "form-horizontal" {
		t.Errorf("expected form class form-horizontal, got %s", r.formClass)
	}
	if !r.noSubmit {
		t.Error("expected noSubmit to be true")
	}
	if r.enctype != "multipart/form-data" {
		t.Errorf("expected enctype multipart/form-data, got %s", r.enctype)
	}
}

func TestFuncMap(t *testing.T) {
	r := New()
	fm := r.FuncMap()

	if _, ok := fm["render_form"]; !ok {
		t.Error("expected render_form to be registered in FuncMap")
	}
}

func TestRenderFormNil(t *testing.T) {
	r := New()
	result := r.RenderForm(nil)
	if !strings.Contains(string(result), "cannot render nil value") {
		t.Errorf("expected error comment for nil, got %s", result)
	}
}

func TestRenderFormNonStruct(t *testing.T) {
	r := New()
	result := r.RenderForm("not a struct")
	if !strings.Contains(string(result), "expected struct") {
		t.Errorf("expected error comment for non-struct, got %s", result)
	}
}

func TestRenderFormPointerToStruct(t *testing.T) {
	r := New()
	form := &LoginForm{
		Username: "john",
		Password: "secret",
	}
	result := string(r.RenderForm(form))

	if !strings.Contains(result, "<form") {
		t.Error("expected <form> tag")
	}
	if !strings.Contains(result, `name="username"`) {
		t.Error("expected username field")
	}
	if !strings.Contains(result, `value="john"`) {
		t.Error("expected pre-populated username value")
	}
}

func TestRenderLoginForm(t *testing.T) {
	r := New(WithAction("/login"), WithMethod("POST"))
	form := LoginForm{
		Username: "testuser",
	}

	result := string(r.RenderForm(form))

	// Check form attributes
	if !strings.Contains(result, `action="/login"`) {
		t.Error("expected action attribute")
	}
	if !strings.Contains(result, `method="POST"`) {
		t.Error("expected method attribute")
	}

	// Check username field
	if !strings.Contains(result, `type="text"`) {
		t.Error("expected text input type")
	}
	if !strings.Contains(result, `name="username"`) {
		t.Error("expected username name")
	}
	if !strings.Contains(result, `placeholder="Enter username"`) {
		t.Error("expected username placeholder")
	}
	if !strings.Contains(result, `value="testuser"`) {
		t.Error("expected pre-populated username")
	}
	if !strings.Contains(result, ` required`) {
		t.Error("expected required attribute")
	}

	// Check password field
	if !strings.Contains(result, `type="password"`) {
		t.Error("expected password input type")
	}
	if !strings.Contains(result, `name="password"`) {
		t.Error("expected password name")
	}

	// Check checkbox
	if !strings.Contains(result, `type="checkbox"`) {
		t.Error("expected checkbox input")
	}
	if !strings.Contains(result, `name="remember"`) {
		t.Error("expected remember name")
	}

	// Check submit input
	if !strings.Contains(result, `type="submit"`) {
		t.Error("expected submit input")
	}
}

func TestRenderContactForm(t *testing.T) {
	r := New()
	form := ContactForm{
		Name:    "John Doe",
		Subject: "support",
	}

	result := string(r.RenderForm(form))

	// Check labels
	if !strings.Contains(result, "Full Name") {
		t.Error("expected Full Name label")
	}
	if !strings.Contains(result, "Email Address") {
		t.Error("expected Email Address label")
	}

	// Check email type
	if !strings.Contains(result, `type="email"`) {
		t.Error("expected email input type")
	}

	// Check tel type
	if !strings.Contains(result, `type="tel"`) {
		t.Error("expected tel input type")
	}

	// Check pattern
	if !strings.Contains(result, `pattern="[0-9]{3}-[0-9]{3}-[0-9]{4}"`) {
		t.Error("expected pattern attribute")
	}

	// Check select
	if !strings.Contains(result, "<select") {
		t.Error("expected select element")
	}
	if !strings.Contains(result, `<option value="general"`) {
		t.Error("expected general option")
	}
	if !strings.Contains(result, `<option value="support" selected`) {
		t.Error("expected support option to be selected")
	}

	// Check textarea
	if !strings.Contains(result, "<textarea") {
		t.Error("expected textarea element")
	}
	if !strings.Contains(result, `rows="5"`) {
		t.Error("expected rows attribute on textarea")
	}
}

func TestRenderSettingsFormWithRadio(t *testing.T) {
	r := New()
	form := SettingsForm{
		Theme: "dark",
		Age:   25,
	}

	result := string(r.RenderForm(form))

	// Check radio buttons
	if !strings.Contains(result, `type="radio"`) {
		t.Error("expected radio input type")
	}
	if !strings.Contains(result, `value="light"`) {
		t.Error("expected light radio option")
	}
	if !strings.Contains(result, `value="dark"`) {
		t.Error("expected dark radio option")
	}
	if !strings.Contains(result, `value="auto"`) {
		t.Error("expected auto radio option")
	}

	// Check number input with min/max
	if !strings.Contains(result, `min="13"`) {
		t.Error("expected min attribute")
	}
	if !strings.Contains(result, `max="120"`) {
		t.Error("expected max attribute")
	}

	// Check URL type
	if !strings.Contains(result, `type="url"`) {
		t.Error("expected url input type")
	}

	// Check pre-populated age
	if !strings.Contains(result, `value="25"`) {
		t.Error("expected pre-populated age value")
	}
}

func TestRenderHiddenFields(t *testing.T) {
	r := New()
	form := HiddenFieldForm{
		ID:    "123",
		Token: "abc",
		Name:  "Test",
	}

	result := string(r.RenderForm(form))

	// Hidden fields should not have labels or form-group wrappers
	if !strings.Contains(result, `type="hidden"`) {
		t.Error("expected hidden input type")
	}
	if !strings.Contains(result, `value="123"`) {
		t.Error("expected hidden ID value")
	}
	if !strings.Contains(result, `value="abc"`) {
		t.Error("expected hidden token value")
	}
}

func TestSkipFields(t *testing.T) {
	r := New()
	form := SkipFieldForm{
		Name:     "Test",
		Internal: "should be skipped",
		NoTag:    "also skipped",
		Email:    "test@example.com",
	}

	result := string(r.RenderForm(form))

	if !strings.Contains(result, `name="name"`) {
		t.Error("expected name field")
	}
	if !strings.Contains(result, `name="email"`) {
		t.Error("expected email field")
	}
	if strings.Contains(result, "Internal") {
		t.Error("should not render field with '-' tag")
	}
	if strings.Contains(result, "NoTag") {
		t.Error("should not render field without form tag")
	}
}

func TestDisabledAndReadonly(t *testing.T) {
	r := New()
	form := DisabledForm{
		ReadonlyField: "cannot edit",
		DisabledField: "cannot interact",
	}

	result := string(r.RenderForm(form))

	if !strings.Contains(result, " readonly") {
		t.Error("expected readonly attribute")
	}
	if !strings.Contains(result, " disabled") {
		t.Error("expected disabled attribute")
	}
}

func TestStepAttribute(t *testing.T) {
	r := New()
	form := StepForm{Price: 9.99}

	result := string(r.RenderForm(form))

	if !strings.Contains(result, `step="0.01"`) {
		t.Error("expected step attribute")
	}
	if !strings.Contains(result, `min="0"`) {
		t.Error("expected min attribute")
	}
}

func TestDateFields(t *testing.T) {
	r := New()
	birthday, _ := time.Parse("2006-01-02", "1990-05-15")
	meeting, _ := time.Parse("2006-01-02T15:04", "2024-12-25T14:30")

	form := DateForm{
		Birthday: birthday,
		Meeting:  meeting,
	}

	result := string(r.RenderForm(form))

	if !strings.Contains(result, `type="date"`) {
		t.Error("expected date input type")
	}
	if !strings.Contains(result, `value="1990-05-15"`) {
		t.Error("expected formatted date value")
	}
	if !strings.Contains(result, `type="datetime-local"`) {
		t.Error("expected datetime-local input type")
	}
}

func TestInferredTypes(t *testing.T) {
	r := New()
	form := InferredForm{
		Name:   "test",
		Age:    30,
		Score:  98.5,
		Active: true,
	}

	result := string(r.RenderForm(form))

	// String -> text
	if !strings.Contains(result, `type="text"`) {
		t.Error("expected inferred text type for string")
	}
	// Int -> number
	if !strings.Contains(result, `type="number"`) {
		t.Error("expected inferred number type for int/float")
	}
	// Bool -> checkbox
	if !strings.Contains(result, `type="checkbox"`) {
		t.Error("expected inferred checkbox type for bool")
	}
}

func TestCSRFToken(t *testing.T) {
	r := New(WithCSRF("csrf-token-123"))
	form := LoginForm{}

	result := string(r.RenderForm(form))

	if !strings.Contains(result, `name="_csrf"`) {
		t.Error("expected CSRF hidden field")
	}
	if !strings.Contains(result, `value="csrf-token-123"`) {
		t.Error("expected CSRF token value")
	}
}

func TestCustomClassPrefix(t *testing.T) {
	r := New(WithClassPrefix("myapp"))
	if r.classPrefix != "myapp" {
		t.Error("expected class prefix myapp")
	}
}

func TestCustomInputClass(t *testing.T) {
	r := New()
	form := CustomClassForm{Name: "test"}

	result := string(r.RenderForm(form))

	if !strings.Contains(result, "my-custom-class") {
		t.Error("expected custom class on name input")
	}
	if !strings.Contains(result, "another-class") {
		t.Error("expected custom class on email input")
	}
}

func TestMultipleSelectAttribute(t *testing.T) {
	r := New()
	form := MultiSelectForm{}

	result := string(r.RenderForm(form))

	if !strings.Contains(result, " multiple") {
		t.Error("expected multiple attribute on select")
	}
}

func TestNoSubmitButton(t *testing.T) {
	r := New(WithNoSubmit(true))
	form := LoginForm{}

	result := string(r.RenderForm(form))

	if strings.Contains(result, `type="submit"`) {
		t.Error("should not contain submit input when WithNoSubmit is true")
	}
}

func TestFormIDAndClass(t *testing.T) {
	r := New(WithFormID("my-form"), WithFormClass("custom-form-class"))
	form := LoginForm{}

	result := string(r.RenderForm(form))

	if !strings.Contains(result, `id="my-form"`) {
		t.Error("expected form ID")
	}
	if !strings.Contains(result, "custom-form-class") {
		t.Error("expected custom form class")
	}
}

func TestEnctype(t *testing.T) {
	r := New(WithEnctype("multipart/form-data"))
	form := LoginForm{}

	result := string(r.RenderForm(form))

	if !strings.Contains(result, `enctype="multipart/form-data"`) {
		t.Error("expected enctype attribute")
	}
}

func TestRenderFormWithAction(t *testing.T) {
	r := New(WithAction("/default"), WithMethod("POST"))
	form := LoginForm{}

	result := string(r.RenderFormWithAction(form, "/custom-action", "GET"))

	if !strings.Contains(result, `action="/custom-action"`) {
		t.Error("expected overridden action")
	}
	if !strings.Contains(result, `method="GET"`) {
		t.Error("expected overridden method")
	}
}

func TestRenderFormWithActionNil(t *testing.T) {
	r := New()
	result := r.RenderFormWithAction(nil, "/test", "POST")
	if !strings.Contains(string(result), "cannot render nil value") {
		t.Error("expected error for nil value")
	}
}


func TestHTMLEscaping(t *testing.T) {
	r := New()
	form := LoginForm{
		Username: `<script>alert("xss")</script>`,
	}

	result := string(r.RenderForm(form))

	if strings.Contains(result, `<script>`) {
		t.Error("HTML should be escaped in field values")
	}
	if !strings.Contains(result, "&lt;script&gt;") {
		t.Error("expected escaped HTML in value")
	}
}

func TestLabelGeneration(t *testing.T) {
	r := New()
	form := InferredForm{}

	result := string(r.RenderForm(form))

	// "name" -> "Name"
	if !strings.Contains(result, ">Name") {
		t.Error("expected auto-generated label 'Name'")
	}
	// "age" -> "Age"
	if !strings.Contains(result, ">Age") {
		t.Error("expected auto-generated label 'Age'")
	}
	// "score" -> "Score"
	if !strings.Contains(result, ">Score") {
		t.Error("expected auto-generated label 'Score'")
	}
}

func TestDefaultIDs(t *testing.T) {
	r := New()
	form := InferredForm{}

	result := string(r.RenderForm(form))

	if !strings.Contains(result, `id="form-name"`) {
		t.Error("expected auto-generated ID form-name")
	}
	if !strings.Contains(result, `id="form-age"`) {
		t.Error("expected auto-generated ID form-age")
	}
}

func TestRequiredIndicator(t *testing.T) {
	r := New()
	form := LoginForm{}

	result := string(r.RenderForm(form))

	// Required fields should have an asterisk
	if !strings.Contains(result, `*</span>`) {
		t.Error("expected required indicator asterisk")
	}
}

func TestRenderField(t *testing.T) {
	r := New()
	form := LoginForm{Username: "john"}

	result := string(r.RenderField(form, "username"))

	if !strings.Contains(result, `name="username"`) {
		t.Error("expected username field")
	}
	if !strings.Contains(result, `value="john"`) {
		t.Error("expected pre-populated value")
	}
}

func TestRenderFieldNotFound(t *testing.T) {
	r := New()
	form := LoginForm{}

	result := string(r.RenderField(form, "nonexistent"))

	if !strings.Contains(result, "not found") {
		t.Error("expected not found comment")
	}
}

func TestRenderFieldNil(t *testing.T) {
	r := New()
	result := r.RenderField(nil, "test")
	if !strings.Contains(string(result), "cannot render nil value") {
		t.Error("expected error for nil value")
	}
}

func TestFields(t *testing.T) {
	r := New()
	form := LoginForm{Username: "john"}

	fields := r.Fields(form)

	if len(fields) != 3 {
		t.Errorf("expected 3 fields, got %d", len(fields))
	}

	if fields[0].Name != "username" {
		t.Errorf("expected first field name 'username', got '%s'", fields[0].Name)
	}
	if fields[0].Value != "john" {
		t.Errorf("expected first field value 'john', got '%s'", fields[0].Value)
	}
	if fields[0].Type != "text" {
		t.Errorf("expected first field type 'text', got '%s'", fields[0].Type)
	}
	if !fields[0].Required {
		t.Error("expected first field to be required")
	}
	if string(fields[0].HTML) == "" {
		t.Error("expected non-empty HTML for first field")
	}
}

func TestFieldsNil(t *testing.T) {
	r := New()
	fields := r.Fields(nil)
	if fields != nil {
		t.Error("expected nil for nil input")
	}
}

func TestFieldNames(t *testing.T) {
	r := New()
	form := ContactForm{}

	names := r.FieldNames(form)

	expected := []string{"name", "email", "phone", "subject", "message"}
	if len(names) != len(expected) {
		t.Errorf("expected %d field names, got %d", len(expected), len(names))
	}

	for i, name := range names {
		if name != expected[i] {
			t.Errorf("expected field name %s at index %d, got %s", expected[i], i, name)
		}
	}
}

func TestFieldNamesNil(t *testing.T) {
	r := New()
	names := r.FieldNames(nil)
	if names != nil {
		t.Error("expected nil for nil input")
	}
}

func TestSelectWithLabeledOptions(t *testing.T) {
	r := New()
	form := ContactForm{Subject: "billing"}

	result := string(r.RenderForm(form))

	if !strings.Contains(result, `>General Inquiry</option>`) {
		t.Error("expected labeled option 'General Inquiry'")
	}
	if !strings.Contains(result, `>Technical Support</option>`) {
		t.Error("expected labeled option 'Technical Support'")
	}
	if !strings.Contains(result, `<option value="billing" selected>Billing Question</option>`) {
		t.Error("expected billing option to be selected")
	}
}

func TestCheckboxChecked(t *testing.T) {
	r := New()
	form := LoginForm{Remember: true}

	result := string(r.RenderForm(form))

	if !strings.Contains(result, " checked") {
		t.Error("expected checked attribute when bool is true")
	}
}

func TestCheckboxUnchecked(t *testing.T) {
	r := New()
	form := LoginForm{Remember: false}

	result := string(r.RenderForm(form))

	if strings.Contains(result, " checked") {
		t.Error("should not have checked attribute when bool is false")
	}
}

func TestTemplateIntegration(t *testing.T) {
	r := New(WithAction("/submit"))

	tmplStr := `<!DOCTYPE html><html><body>{{render_form .Form}}</body></html>`
	tmpl, err := template.New("test").Funcs(r.FuncMap()).Parse(tmplStr)
	if err != nil {
		t.Fatalf("failed to parse template: %v", err)
	}

	data := struct {
		Form LoginForm
	}{
		Form: LoginForm{
			Username: "testuser",
			Password: "",
			Remember: true,
		},
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, data)
	if err != nil {
		t.Fatalf("failed to execute template: %v", err)
	}

	result := buf.String()

	if !strings.Contains(result, `<form`) {
		t.Error("expected form in template output")
	}
	if !strings.Contains(result, `action="/submit"`) {
		t.Error("expected action in template output")
	}
	if !strings.Contains(result, `value="testuser"`) {
		t.Error("expected pre-populated value in template output")
	}
	if !strings.Contains(result, " checked") {
		t.Error("expected checked checkbox in template output")
	}
}

// --- Benchmark ---

func BenchmarkRenderForm(b *testing.B) {
	r := New()
	form := ContactForm{
		Name:    "John Doe",
		Email:   "john@example.com",
		Phone:   "555-123-4567",
		Subject: "support",
		Message: "Hello, I need help.",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.RenderForm(form)
	}
}

// --- Parser unit tests ---

func TestTitleCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"name", "Name"},
		{"first_name", "First Name"},
		{"email-address", "Email Address"},
		{"firstName", "First Name"},
		{"userID", "User Id"},
		{"URL", "Url"},
		{"a", "A"},
		{"", ""},
	}

	for _, tt := range tests {
		result := titleCase(tt.input)
		if result != tt.expected {
			t.Errorf("titleCase(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestSplitTag(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"name", []string{"name"}},
		{"name,type=text", []string{"name", "type=text"}},
		{"name,type=text,required", []string{"name", "type=text", "required"}},
		{"subject,type=select,options=a|b|c", []string{"subject", "type=select", "options=a|b|c"}},
	}

	for _, tt := range tests {
		result := splitTag(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("splitTag(%q): got %d parts, want %d", tt.input, len(result), len(tt.expected))
			continue
		}
		for i, part := range result {
			if part != tt.expected[i] {
				t.Errorf("splitTag(%q)[%d] = %q, want %q", tt.input, i, part, tt.expected[i])
			}
		}
	}
}

func TestParseOptions(t *testing.T) {
	tests := []struct {
		input    string
		expected []selectOption
	}{
		{"a|b|c", []selectOption{
			{Value: "a", Label: "A"},
			{Value: "b", Label: "B"},
			{Value: "c", Label: "C"},
		}},
		{"admin:Administrator|user:Regular User", []selectOption{
			{Value: "admin", Label: "Administrator"},
			{Value: "user", Label: "Regular User"},
		}},
	}

	for _, tt := range tests {
		result := parseOptions(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("parseOptions(%q): got %d options, want %d", tt.input, len(result), len(tt.expected))
			continue
		}
		for i, opt := range result {
			if opt.Value != tt.expected[i].Value || opt.Label != tt.expected[i].Label {
				t.Errorf("parseOptions(%q)[%d] = {%s,%s}, want {%s,%s}",
					tt.input, i, opt.Value, opt.Label, tt.expected[i].Value, tt.expected[i].Label)
			}
		}
	}
}

func TestParseTagPart(t *testing.T) {
	tests := []struct {
		input    string
		key      string
		value    string
		hasValue bool
	}{
		{"required", "required", "", false},
		{"type=text", "type", "text", true},
		{"placeholder=Enter name", "placeholder", "Enter name", true},
		{"options=a|b|c", "options", "a|b|c", true},
	}

	for _, tt := range tests {
		key, value, hasValue := parseTagPart(tt.input)
		if key != tt.key || value != tt.value || hasValue != tt.hasValue {
			t.Errorf("parseTagPart(%q) = (%q, %q, %v), want (%q, %q, %v)",
				tt.input, key, value, hasValue, tt.key, tt.value, tt.hasValue)
		}
	}
}

func TestInferType(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{"", "text"},
		{0, "number"},
		{int64(0), "number"},
		{uint(0), "number"},
		{float64(0), "number"},
		{false, "checkbox"},
		{time.Time{}, "datetime-local"},
	}

	for _, tt := range tests {
		result := inferType(reflect.TypeOf(tt.input))
		if result != tt.expected {
			t.Errorf("inferType(%T) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}


func TestTextareaWithValue(t *testing.T) {
	r := New()
	form := ContactForm{
		Message: "Hello World",
	}

	result := string(r.RenderForm(form))

	if !strings.Contains(result, ">Hello World</textarea>") {
		t.Error("expected textarea to contain value as inner text")
	}
}

func TestRadioSelectedOption(t *testing.T) {
	r := New()
	form := SettingsForm{
		Theme: "dark",
	}

	result := string(r.RenderForm(form))

	// The dark radio option should be checked
	// We need to verify the "dark" value radio has "checked"
	lines := strings.Split(result, "\n")
	foundDarkChecked := false
	for _, line := range lines {
		if strings.Contains(line, `value="dark"`) && strings.Contains(line, "checked") {
			foundDarkChecked = true
			break
		}
	}
	if !foundDarkChecked {
		t.Error("expected dark radio option to be checked")
	}
}

type StaticFormInfoForm struct {
	Info     FormInfo `form:"_info,action=/submit-info,method=GET,submit_label=Send,submit_class=btn-p,submit_attrs=hx-post='/static' up-dismiss,reset_label=Cancel,reset_class=btn-c,reset_attrs=hx-get='/cancel'"`
	Username string   `form:"username"`
}

type DynamicFormInfoForm struct {
	Info     FormInfo
	Username string `form:"username"`
}

func TestStaticFormInfo(t *testing.T) {
	r := New()
	form := StaticFormInfoForm{}

	result := string(r.RenderForm(form))

	if !strings.Contains(result, `action="/submit-info"`) {
		t.Error("expected custom action from tag FormInfo")
	}
	if !strings.Contains(result, `method="GET"`) {
		t.Error("expected custom method from tag FormInfo")
	}
	if !strings.Contains(result, `class="btn-p"`) {
		t.Error("expected custom submit button class")
	}
	if !strings.Contains(result, `hx-post='/static' up-dismiss`) {
		t.Error("expected custom submit attributes")
	}
	if !strings.Contains(result, `value="Cancel"`) {
		t.Error("expected reset input with label Cancel")
	}
	if !strings.Contains(result, `class="btn-c"`) {
		t.Error("expected custom reset class")
	}
	if !strings.Contains(result, `hx-get='/cancel'`) {
		t.Error("expected custom reset attributes")
	}
}

func TestDynamicFormInfo(t *testing.T) {
	r := New()
	form := DynamicFormInfoForm{
		Info: FormInfo{
			Action:      "/dynamic-action",
			SubmitLabel: "Register Now",
			SubmitAttrs: `up-dismiss hx-post="/dynamic"`,
			ResetLabel:  "Clear",
			ResetAttrs:  `hx-get="/clear"`,
			FormAttrs:   `hx-target="#div"`,
		},
		Username: "bob",
	}

	result := string(r.RenderForm(form))

	if !strings.Contains(result, `action="/dynamic-action"`) {
		t.Error("expected dynamic action")
	}
	if !strings.Contains(result, `hx-target="#div"`) {
		t.Error("expected dynamic raw form attributes")
	}
	if !strings.Contains(result, `value="Register Now"`) {
		t.Error("expected dynamic submit value Register Now")
	}
	if !strings.Contains(result, `up-dismiss hx-post="/dynamic"`) {
		t.Error("expected dynamic submit attributes")
	}
	if !strings.Contains(result, `value="Clear"`) {
		t.Error("expected dynamic reset value Clear")
	}
	if !strings.Contains(result, `hx-get="/clear"`) {
		t.Error("expected dynamic reset attributes")
	}
}

type RoleForm struct {
	Info     FormInfo `form:"_info,submit_role=button,reset_role=button,reset_label=Clear"`
	Username string   `form:"username,role=textbox"`
	Age      int      `form:"age,role=spinbutton"`
	Agree    bool     `form:"agree,role=checkbox"`
	Theme    string   `form:"theme,type=radio,options=light|dark,role=radiogroup"`
}

func TestInputRoleAttribute(t *testing.T) {
	r := New()
	form := RoleForm{}

	result := string(r.RenderForm(form))

	// Verify field-level roles
	if !strings.Contains(result, `role="textbox"`) {
		t.Error("expected role='textbox' on username input")
	}
	if !strings.Contains(result, `role="spinbutton"`) {
		t.Error("expected role='spinbutton' on age input")
	}
	if !strings.Contains(result, `role="checkbox"`) {
		t.Error("expected role='checkbox' on agree input")
	}
	if !strings.Contains(result, `role="radiogroup"`) {
		t.Error("expected role='radiogroup' on theme radio inputs")
	}

	// Verify static FormInfo submit/reset roles
	if !strings.Contains(result, `<input type="submit" value="Submit" role="button" />`) {
		t.Error("expected role='button' on static submit input")
	}
	if !strings.Contains(result, `<input type="reset" value="Clear" role="button" />`) {
		t.Error("expected role='button' on static reset input")
	}
}

func TestDynamicFormInfoRoles(t *testing.T) {
	r := New()
	type DynamicRoleForm struct {
		Info     FormInfo
		Username string `form:"username"`
	}

	form := DynamicRoleForm{
		Info: FormInfo{
			SubmitRole: "submit-btn",
			ResetLabel: "Clear",
			ResetRole:  "reset-btn",
		},
	}

	result := string(r.RenderForm(form))

	if !strings.Contains(result, `role="submit-btn"`) {
		t.Error("expected dynamic submit role")
	}
	if !strings.Contains(result, `role="reset-btn"`) {
		t.Error("expected dynamic reset role")
	}
}

func TestNoDivWrappers(t *testing.T) {
	r := New()
	type WrapperTestForm struct {
		Username string `form:"username"`
		Agree    bool   `form:"agree"`
		Theme    string `form:"theme,type=radio,options=light|dark"`
	}

	form := WrapperTestForm{}
	result := string(r.RenderForm(form))

	if strings.Contains(result, "<div") || strings.Contains(result, "</div>") {
		t.Errorf("expected HTML output to contain no wrapper div elements, got: %s", result)
	}
}

func TestFieldsetAndAutocomplete(t *testing.T) {
	r := New()
	type FieldsetForm struct {
		Info      FormInfo `form:"_info,fieldset=true,submit_label=Subscribe"`
		FirstName string   `form:"first_name,placeholder=First name,autocomplete=given-name"`
		Email     string   `form:"email,type=email,placeholder=Email,autocomplete=email"`
		Remember  bool     `form:"remember,role=switch"`
	}

	form := FieldsetForm{}
	result := string(r.RenderForm(form))

	// Verify fieldset wrapper is present and placed before fields but after form open
	if !strings.Contains(result, "<form") || !strings.Contains(result, "<fieldset>") || !strings.Contains(result, "</fieldset>") {
		t.Error("expected form and fieldset tags")
	}

	// Verify autocomplete attribute is correctly rendered on inputs
	if !strings.Contains(result, `name="first_name" autocomplete="given-name"`) {
		t.Error("expected autocomplete='given-name' attribute")
	}
	if !strings.Contains(result, `name="email" autocomplete="email"`) {
		t.Error("expected autocomplete='email' attribute")
	}

	// Verify checkbox renders label first, then input with role="switch"
	labelIdx := strings.Index(result, "Remember")
	inputIdx := strings.Index(result, `type="checkbox"`)
	if labelIdx == -1 || inputIdx == -1 || labelIdx > inputIdx {
		t.Error("expected label to be rendered before checkbox input")
	}

	// Verify self-closing tags with " />"
	if !strings.Contains(result, `type="checkbox"`) || !strings.Contains(result, `role="switch"`) || !strings.Contains(result, `value="true" />`) {
		t.Error("expected checkbox to be self-closing with ' />'")
	}
	if !strings.Contains(result, `type="submit" value="Subscribe" />`) {
		t.Error("expected submit input to be self-closing with ' />'")
	}
}

