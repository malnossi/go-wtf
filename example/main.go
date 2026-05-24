// Package main demonstrates the go-wtf form rendering module.
// Run this example and visit http://localhost:8080 in your browser.
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	wtf "github.com/malnossi/go-wtf"
)

// UserRegistrationForm demonstrates a typical registration form with various field types.
type UserRegistrationForm struct {
	FirstName   string    `form:"first_name,type=text,label=First Name,placeholder=Enter your first name,required"`
	LastName    string    `form:"last_name,type=text,label=Last Name,placeholder=Enter your last name,required"`
	Email       string    `form:"email,type=email,label=Email Address,placeholder=you@example.com,required"`
	Password    string    `form:"password,type=password,label=Password,placeholder=Minimum 8 characters,required,min=8"`
	DateOfBirth time.Time `form:"dob,type=date,label=Date of Birth"`
	Phone       string    `form:"phone,type=tel,label=Phone Number,placeholder=555-123-4567,pattern=[0-9]{3}-[0-9]{3}-[0-9]{4}"`
	Website     string    `form:"website,type=url,label=Personal Website,placeholder=https://example.com"`
	Bio         string    `form:"bio,type=textarea,label=About You,placeholder=Tell us a bit about yourself...,rows=4"`
	Role        string    `form:"role,type=select,label=Account Type,options=personal:Personal|business:Business|developer:Developer,required"`
	Theme       string    `form:"theme,type=radio,label=Preferred Theme,options=light:Light Mode|dark:Dark Mode|system:System Default"`
	Newsletter  bool      `form:"newsletter,type=checkbox,label=Subscribe to our newsletter"`
	Terms       bool      `form:"terms,type=checkbox,label=I agree to the Terms of Service,required"`
}

// LoginForm demonstrates a simple login form.
type LoginForm struct {
	Info     wtf.FormInfo `form:"_info,submit_label=Sign In,submit_class=btn-signin,submit_attrs=up-dismiss hx-post='/submit'"`
	Username string       `form:"username,type=text,placeholder=Username or email,required"`
	Password string       `form:"password,type=password,placeholder=Password,required"`
	Remember bool         `form:"remember,type=checkbox,label=Remember me"`
}

// ContactForm demonstrates a contact form with textarea and select.
type ContactForm struct {
	Info    wtf.FormInfo
	Name    string `form:"name,type=text,label=Your Name,required,placeholder=Full name"`
	Email   string `form:"email,type=email,label=Email,required,placeholder=your@email.com"`
	Subject string `form:"subject,type=select,label=Subject,options=general:General Inquiry|support:Technical Support|feedback:Feedback|bug:Bug Report,required"`
	Urgency string `form:"urgency,type=radio,label=Urgency,options=low:Low|medium:Medium|high:High|critical:Critical"`
	Message string `form:"message,type=textarea,label=Message,required,placeholder=Describe your issue or question...,rows=6"`
}

const pageTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>go-wtf Form Rendering Demo</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700&display=swap" rel="stylesheet">
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
            background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #0f172a 100%);
            color: #e2e8f0;
            min-height: 100vh;
            padding: 2rem;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
        }
        h1 {
            font-size: 2.5rem;
            font-weight: 700;
            background: linear-gradient(135deg, #60a5fa, #a78bfa, #f472b6);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
            text-align: center;
            margin-bottom: 0.5rem;
        }
        .subtitle {
            text-align: center;
            color: #94a3b8;
            margin-bottom: 2.5rem;
            font-size: 1.1rem;
        }
        .form-section {
            background: rgba(30, 41, 59, 0.8);
            backdrop-filter: blur(12px);
            border: 1px solid rgba(148, 163, 184, 0.1);
            border-radius: 16px;
            padding: 2rem;
            margin-bottom: 2rem;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.3);
        }
        .form-section h2 {
            font-size: 1.5rem;
            font-weight: 600;
            color: #f1f5f9;
            margin-bottom: 0.25rem;
        }
        .form-section .description {
            color: #94a3b8;
            font-size: 0.9rem;
            margin-bottom: 1.5rem;
            padding-bottom: 1rem;
            border-bottom: 1px solid rgba(148, 163, 184, 0.1);
        }
        .tabs {
            display: flex;
            gap: 0.5rem;
            margin-bottom: 2rem;
            flex-wrap: wrap;
        }
        .tab {
            padding: 0.625rem 1.25rem;
            background: rgba(51, 65, 85, 0.5);
            border: 1px solid rgba(148, 163, 184, 0.15);
            border-radius: 8px;
            color: #94a3b8;
            cursor: pointer;
            font-size: 0.9rem;
            font-weight: 500;
            transition: all 0.2s;
            font-family: inherit;
        }
        .tab:hover {
            background: rgba(71, 85, 105, 0.5);
            color: #e2e8f0;
        }
        .tab.active {
            background: linear-gradient(135deg, #3b82f6, #6366f1);
            border-color: transparent;
            color: #fff;
        }
        .form-panel {
            display: none;
        }
        .form-panel.active {
            display: block;
        }
        /* Override wtf default styles for dark theme */
        .wtf-form {
            max-width: 100% !important;
        }
        .wtf-form-label {
            color: #e2e8f0 !important;
        }
        .wtf-form-input,
        .wtf-form-select,
        .wtf-form-textarea {
            background-color: rgba(15, 23, 42, 0.6) !important;
            border-color: rgba(148, 163, 184, 0.2) !important;
            color: #e2e8f0 !important;
        }
        .wtf-form-input:focus,
        .wtf-form-select:focus,
        .wtf-form-textarea:focus {
            border-color: #3b82f6 !important;
            box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2) !important;
        }
        .wtf-form-input::placeholder,
        .wtf-form-textarea::placeholder {
            color: #64748b !important;
        }
        .wtf-form-checkbox-label,
        .wtf-form-radio-label {
            color: #cbd5e1 !important;
        }
        .wtf-form-submit {
            background: linear-gradient(135deg, #3b82f6, #6366f1) !important;
            border: none !important;
            padding: 0.625rem 2rem !important;
            font-weight: 600 !important;
            border-radius: 8px !important;
            transition: all 0.2s !important;
        }
        .wtf-form-submit:hover {
            transform: translateY(-1px);
            box-shadow: 0 4px 12px rgba(59, 130, 246, 0.4) !important;
        }
        .code-section {
            margin-top: 1rem;
            padding-top: 1rem;
            border-top: 1px solid rgba(148, 163, 184, 0.1);
        }
        .code-toggle {
            background: transparent;
            border: 1px solid rgba(148, 163, 184, 0.2);
            color: #94a3b8;
            padding: 0.5rem 1rem;
            border-radius: 6px;
            cursor: pointer;
            font-size: 0.85rem;
            font-family: inherit;
            transition: all 0.2s;
        }
        .code-toggle:hover {
            border-color: #3b82f6;
            color: #e2e8f0;
        }
        pre {
            background: rgba(15, 23, 42, 0.8);
            border: 1px solid rgba(148, 163, 184, 0.1);
            border-radius: 8px;
            padding: 1rem;
            margin-top: 0.75rem;
            overflow-x: auto;
            font-size: 0.85rem;
            line-height: 1.6;
            display: none;
        }
        pre.visible {
            display: block;
        }
        code {
            font-family: 'JetBrains Mono', 'Fira Code', 'Cascadia Code', monospace;
            color: #a5b4fc;
        }
        .kw { color: #c084fc; }
        .str { color: #86efac; }
        .tag { color: #fbbf24; }
        .cmt { color: #64748b; }
        footer {
            text-align: center;
            color: #64748b;
            font-size: 0.85rem;
            margin-top: 2rem;
            padding-top: 1.5rem;
            border-top: 1px solid rgba(148, 163, 184, 0.1);
        }
        footer a {
            color: #60a5fa;
            text-decoration: none;
        }
        footer a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>go-wtf</h1>
        <p class="subtitle">Automatic HTML form rendering from Go struct tags</p>

        <div class="tabs">
            <button class="tab active" onclick="showTab('registration')">Registration</button>
            <button class="tab" onclick="showTab('login')">Login</button>
            <button class="tab" onclick="showTab('contact')">Contact</button>
        </div>

        <div id="registration" class="form-panel active">
            <div class="form-section">
                <h2>User Registration</h2>
                <p class="description">A comprehensive registration form demonstrating text, email, password, date, tel, url, textarea, select, radio, and checkbox inputs.</p>
                {{render_form .Registration}}
                <div class="code-section">
                    <button class="code-toggle" onclick="toggleCode('reg-code')">View Go Struct Code</button>
                    <pre id="reg-code"><code><span class="kw">type</span> UserRegistrationForm <span class="kw">struct</span> {
    FirstName   <span class="kw">string</span>    <span class="tag">` + "`" + `form:"first_name,type=text,label=First Name,placeholder=Enter your first name,required"` + "`" + `</span>
    LastName    <span class="kw">string</span>    <span class="tag">` + "`" + `form:"last_name,type=text,label=Last Name,placeholder=Enter your last name,required"` + "`" + `</span>
    Email       <span class="kw">string</span>    <span class="tag">` + "`" + `form:"email,type=email,label=Email Address,placeholder=you@example.com,required"` + "`" + `</span>
    Password    <span class="kw">string</span>    <span class="tag">` + "`" + `form:"password,type=password,label=Password,placeholder=Minimum 8 characters,required,min=8"` + "`" + `</span>
    DateOfBirth time.Time <span class="tag">` + "`" + `form:"dob,type=date,label=Date of Birth"` + "`" + `</span>
    Phone       <span class="kw">string</span>    <span class="tag">` + "`" + `form:"phone,type=tel,label=Phone Number,placeholder=555-123-4567"` + "`" + `</span>
    Website     <span class="kw">string</span>    <span class="tag">` + "`" + `form:"website,type=url,label=Personal Website,placeholder=https://example.com"` + "`" + `</span>
    Bio         <span class="kw">string</span>    <span class="tag">` + "`" + `form:"bio,type=textarea,label=About You,placeholder=Tell us a bit about yourself...,rows=4"` + "`" + `</span>
    Role        <span class="kw">string</span>    <span class="tag">` + "`" + `form:"role,type=select,label=Account Type,options=personal:Personal|business:Business|developer:Developer,required"` + "`" + `</span>
    Theme       <span class="kw">string</span>    <span class="tag">` + "`" + `form:"theme,type=radio,label=Preferred Theme,options=light:Light Mode|dark:Dark Mode|system:System Default"` + "`" + `</span>
    Newsletter  <span class="kw">bool</span>      <span class="tag">` + "`" + `form:"newsletter,type=checkbox,label=Subscribe to our newsletter"` + "`" + `</span>
    Terms       <span class="kw">bool</span>      <span class="tag">` + "`" + `form:"terms,type=checkbox,label=I agree to the Terms of Service,required"` + "`" + `</span>
}</code></pre>
                </div>
            </div>
        </div>

        <div id="login" class="form-panel">
            <div class="form-section">
                <h2>Login</h2>
                <p class="description">A minimal login form with just username, password, and a remember-me checkbox.</p>
                {{render_form .Login}}
                <div class="code-section">
                    <button class="code-toggle" onclick="toggleCode('login-code')">View Go Struct Code</button>
                    <pre id="login-code"><code><span class="kw">type</span> LoginForm <span class="kw">struct</span> {
    Info     wtf.FormInfo <span class="tag">` + "`" + `form:"_info,submit_label=Sign In,submit_class=btn-signin,submit_attrs=up-dismiss hx-post='/submit'"` + "`" + `</span>
    Username <span class="kw">string</span>       <span class="tag">` + "`" + `form:"username,type=text,placeholder=Username or email,required"` + "`" + `</span>
    Password <span class="kw">string</span>       <span class="tag">` + "`" + `form:"password,type=password,placeholder=Password,required"` + "`" + `</span>
    Remember <span class="kw">bool</span>         <span class="tag">` + "`" + `form:"remember,type=checkbox,label=Remember me"` + "`" + `</span>
}</code></pre>
                </div>
            </div>
        </div>

        <div id="contact" class="form-panel">
            <div class="form-section">
                <h2>Contact Us</h2>
                <p class="description">A contact form demonstrating dynamic, runtime configuration via FormInfo, rendering a custom submit button with HTMX attributes, and a styled Reset button.</p>
                {{render_form .Contact}}
                <div class="code-section">
                    <button class="code-toggle" onclick="toggleCode('contact-code')">View Go Struct Code</button>
                    <pre id="contact-code"><code><span class="kw">type</span> ContactForm <span class="kw">struct</span> {
    Info    wtf.FormInfo
    Name    <span class="kw">string</span> <span class="tag">` + "`" + `form:"name,type=text,label=Your Name,required,placeholder=Full name"` + "`" + `</span>
    Email   <span class="kw">string</span> <span class="tag">` + "`" + `form:"email,type=email,label=Email,required,placeholder=your@email.com"` + "`" + `</span>
    Subject <span class="kw">string</span> <span class="tag">` + "`" + `form:"subject,type=select,label=Subject,options=general:General Inquiry|support:Technical Support|feedback:Feedback|bug:Bug Report,required"` + "`" + `</span>
    Urgency <span class="kw">string</span> <span class="tag">` + "`" + `form:"urgency,type=radio,label=Urgency,options=low:Low|medium:Medium|high:High|critical:Critical"` + "`" + `</span>
    Message <span class="kw">string</span> <span class="tag">` + "`" + `form:"message,type=textarea,label=Message,required,placeholder=Describe your issue...,rows=6"` + "`" + `</span>
}

<span class="cmt">// Set dynamically at runtime:</span>
form := ContactForm{
    Info: wtf.FormInfo{
        Action:      <span class="str">"/submit"</span>,
        SubmitLabel: <span class="str">"Send Message"</span>,
        SubmitAttrs: <span class="str">` + "`" + `hx-post="/submit" up-dismiss hx-target="#div"` + "`" + `</span>,
        ResetLabel:  <span class="str">"Clear Form"</span>,
        ResetAttrs:  <span class="str">` + "`" + `hx-get="/reset"` + "`" + `</span>,
    },
}</code></pre>
                </div>
            </div>
        </div>

        <footer>
            <p>Built with <a href="https://github.com/malnossi/go-wtf">go-wtf</a> — Transform Go structs into HTML forms with a single template call.</p>
        </footer>
    </div>

    <script>
        function showTab(name) {
            document.querySelectorAll('.form-panel').forEach(p => p.classList.remove('active'));
            document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
            document.getElementById(name).classList.add('active');
            event.target.classList.add('active');
        }
        function toggleCode(id) {
            document.getElementById(id).classList.toggle('visible');
        }
    </script>
</body>
</html>`

func main() {
	// Create the form renderer with configuration
	renderer := wtf.New(
		wtf.WithAction("/submit"),
		wtf.WithMethod("POST"),
		wtf.WithSubmitLabel("Submit"),
	)

	// Parse the template with the renderer's FuncMap
	tmpl := template.Must(
		template.New("page").Funcs(renderer.FuncMap()).Parse(pageTemplate),
	)

	// Serve the form page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Registration UserRegistrationForm
			Login        LoginForm
			Contact      ContactForm
		}{
			Registration: UserRegistrationForm{
				Theme: "system",
			},
			Login:   LoginForm{},
			Contact: ContactForm{
				Info: wtf.FormInfo{
					Action:      "/submit",
					SubmitLabel: "Send Message",
					SubmitAttrs: `hx-post="/submit" up-dismiss hx-target="#div"`,
					ResetLabel:  "Clear Form",
					ResetAttrs:  `hx-get="/reset"`,
				},
			},
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Handle form submission
	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Form submitted successfully!\n\nReceived values:\n")
		for key, values := range r.PostForm {
			fmt.Fprintf(w, "  %s = %s\n", key, values)
		}
	})

	addr := ":8080"
	fmt.Printf("🚀 go-wtf example server running at http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
