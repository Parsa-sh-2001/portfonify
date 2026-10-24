package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"html"
	"net/mail"
	"regexp"
	"strings"
)

var validAccentColor = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

func makeZipMinimal(portfolio Portfolio) ([]byte, error) {
	css := `* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

:root {
  --background: #080a18;
  --surface: #101329;
  --surface-light: #171b35;
  --violet: ACCENT_PLACEHOLDER;
  --cyan: #67e8f9;
  --white: #f3f4ff;
  --muted: #969bbd;
  --border: #252947;
}

html { scroll-behavior: smooth; }

body {
  background: var(--background);
  color: var(--white);
  font-family: FONT_PLACEHOLDER;
  line-height: 1.7;
}

a { color: inherit; }

main {
  max-width: 1050px;
  margin: 0 auto;
  padding: 28px;
}

nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 18px;
  padding-bottom: 22px;
  border-bottom: 1px solid var(--border);
}

.logo {
  font-size: 19px;
  font-weight: 800;
  letter-spacing: -1px;
}

.logo span { color: var(--violet); }

nav a {
  color: var(--muted);
  text-decoration: none;
  font-size: 13px;
  margin-left: 22px;
}

nav a:hover { color: var(--cyan); }

.hero { padding: 105px 0 90px; }

.eyebrow {
  color: var(--cyan);
  font-size: 12px;
  letter-spacing: 2px;
  text-transform: uppercase;
}

h1 {
  margin: 28px 0 20px;
  font-size: clamp(46px, 8vw, 88px);
  line-height: 1.05;
  letter-spacing: -4px;
  overflow-wrap: anywhere;
}

h1 em {
  color: var(--violet);
  font-weight: normal;
}

.bio {
  max-width: 650px;
  color: var(--muted);
  font-size: 18px;
}

section {
  padding: 42px 0;
  border-top: 1px solid var(--border);
}

.section-number {
  color: var(--cyan);
  font-size: 12px;
  letter-spacing: 2px;
  text-transform: uppercase;
}

h2 {
  margin: 12px 0 20px;
  font-size: clamp(27px, 4vw, 38px);
}

.about-text {
  max-width: 750px;
  color: var(--muted);
  overflow-wrap: anywhere;
}

.projects {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18px;
  margin-top: 28px;
}

.project {
  padding: 25px;
  background: linear-gradient(145deg, var(--surface-light), var(--surface));
  border: 1px solid var(--border);
  border-radius: 12px;
  transition: border-color 0.25s, transform 0.25s;
}

.project:hover {
  border-color: var(--violet);
  transform: translateY(-4px);
}

.project h3 {
  margin: 12px 0 10px;
  font-size: 22px;
  overflow-wrap: anywhere;
}

.project p {
  color: var(--muted);
  font-size: 14px;
  overflow-wrap: anywhere;
}

.contact {
  padding: 32px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 14px;
  margin-top: 18px;
}

.contact p { color: var(--muted); }

.contact a {
  display: inline-block;
  margin-top: 22px;
  color: var(--violet);
  text-decoration: none;
}

footer {
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 45px;
  padding-top: 22px;
  border-top: 1px solid var(--border);
  color: var(--muted);
  font-size: 12px;
}

footer span { color: var(--violet); }

@media (max-width: 650px) {
  main { padding: 20px; }
  nav a { margin-left: 10px; font-size: 12px; }
  .hero { padding: 75px 0 60px; }
  h1 { letter-spacing: -2px; }
  .projects { grid-template-columns: 1fr; }
  .contact { padding: 24px; }
}`

	// Validate the accent color.
	accent := strings.TrimSpace(portfolio.Accent)
	if !validAccentColor.MatchString(accent) {
		accent = "#a78bfa"
	}

	// Only allow fonts supported by the builder.
	font := "Arial, Helvetica, sans-serif"

	switch strings.TrimSpace(portfolio.Font) {
	case "Georgia, serif":
		font = "Georgia, serif"
	case "Arial, sans-serif", "Arial, Helvetica, sans-serif":
		font = "Arial, Helvetica, sans-serif"
	case "Inter, sans-serif":
		font = "Inter, sans-serif"
	case "system-ui, sans-serif":
		font = "system-ui, sans-serif"
	}

	css = strings.ReplaceAll(css, "ACCENT_PLACEHOLDER", accent)
	css = strings.ReplaceAll(css, "FONT_PLACEHOLDER", font)

	// Escape user-provided text before putting it into HTML.
	name := strings.TrimSpace(portfolio.Name)
	role := strings.TrimSpace(portfolio.Role)
	bio := strings.TrimSpace(portfolio.Bio)

	if name == "" {
		name = "My Portfolio"
	}
	if role == "" {
		role = "Portfolio"
	}
	if bio == "" {
		bio = "Welcome to my portfolio."
	}

	name = html.EscapeString(name)
	role = html.EscapeString(role)
	bio = html.EscapeString(bio)

	// Only include a contact link for a valid email address.
	contactHTML := `<p>Add your email address to let people contact you.</p>`

	email := strings.TrimSpace(portfolio.Email)
	if email != "" {
		if parsed, err := mail.ParseAddress(email); err == nil &&
			parsed.Address == email &&
			!strings.ContainsAny(email, "\r\n") {

			contactHTML = fmt.Sprintf(
				`<p>Have an interesting idea or project in mind?</p>
				<a href="mailto:%s">Get in touch ↗</a>`,
				html.EscapeString(email),
			)
		}
	}

	// Read the two project slots, using defaults for missing values.
	projectName := func(index int) string {
		if index >= len(portfolio.Projects) {
			return "Untitled project"
		}

		value := strings.TrimSpace(portfolio.Projects[index].Name)
		if value == "" {
			return "Untitled project"
		}

		return html.EscapeString(value)
	}

	projectDescription := func(index int) string {
		if index >= len(portfolio.Projects) {
			return "Add a description for this project."
		}

		value := strings.TrimSpace(portfolio.Projects[index].Description)
		if value == "" {
			return "Add a description for this project."
		}

		return html.EscapeString(value)
	}

	project1Name := projectName(0)
	project1Description := projectDescription(0)
	project2Name := projectName(1)
	project2Description := projectDescription(1)

	page := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="description" content="%s">
  <title>%s — Portfolio</title>
  <link rel="stylesheet" href="./css/style.css">
</head>
<body>
  <main>
    <nav>
      <div class="logo">%s<span>.</span></div>
      <div>
        <a href="#about">About</a>
        <a href="#projects">Work</a>
        <a href="#contact">Contact</a>
      </div>
    </nav>

    <header class="hero">
      <p class="eyebrow">%s</p>
      <h1>%s<br><em>with purpose.</em></h1>
      <p class="bio">%s</p>
    </header>

    <section id="about">
      <p class="section-number">01 / About</p>
      <h2>A little context.</h2>
      <p class="about-text">%s</p>
    </section>

    <section id="projects">
      <p class="section-number">02 / Selected work</p>
      <h2>A few things I've made.</h2>
      <div class="projects">
        <article class="project">
          <span class="section-number">PROJECT 001</span>
          <h3>%s</h3>
          <p>%s</p>
        </article>
        <article class="project">
          <span class="section-number">PROJECT 002</span>
          <h3>%s</h3>
          <p>%s</p>
        </article>
      </div>
    </section>

    <section id="contact" class="contact">
      <p class="section-number">03 / Contact</p>
      <h2>Let's talk.</h2>
      %s
    </section>

    <footer>
      <span>%s · PORTFOLIO</span>
      <span>MADE WITH PORTFONIFY · © 2026</span>
    </footer>
  </main>
</body>
</html>`,
		bio, name, name, role, name, bio, bio,
		project1Name, project1Description,
		project2Name, project2Description,
		contactHTML, name,
	)

	// Create the ZIP in memory.
	var buffer bytes.Buffer
	archive := zip.NewWriter(&buffer)

	// Write index.html.
	index, err := archive.Create("index.html")
	if err != nil {
		_ = archive.Close()
		return nil, fmt.Errorf("create HTML entry: %w", err)
	}

	if _, err := index.Write([]byte(page)); err != nil {
		_ = archive.Close()
		return nil, fmt.Errorf("write HTML: %w", err)
	}

	// Write the stylesheet.
	style, err := archive.Create("css/style.css")
	if err != nil {
		_ = archive.Close()
		return nil, fmt.Errorf("create CSS entry: %w", err)
	}

	if _, err := style.Write([]byte(css)); err != nil {
		_ = archive.Close()
		return nil, fmt.Errorf("write CSS: %w", err)
	}

	// Finalize the archive before returning its bytes.
	if err := archive.Close(); err != nil {
		return nil, fmt.Errorf("finish ZIP: %w", err)
	}

	return buffer.Bytes(), nil
}
