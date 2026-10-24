package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

func MakeZipNeonTech(portfolio Portfolio) ([]byte, error) {
	accent := strings.TrimSpace(portfolio.Accent)

	validHex := func(value string) bool {
		if len(value) != 4 && len(value) != 7 {
			return false
		}
		if !strings.HasPrefix(value, "#") {
			return false
		}

		for _, char := range value[1:] {
			if !strings.ContainsRune("0123456789abcdefABCDEF", char) {
				return false
			}
		}

		return true
	}

	if !validHex(accent) {
		accent = "#a78bfa"
	}

	font := "Arial, Helvetica, sans-serif"

	switch strings.ToLower(strings.TrimSpace(portfolio.Font)) {
	case "inter":
		font = "Inter, Arial, sans-serif"
	case "arial":
		font = "Arial, Helvetica, sans-serif"
	case "georgia":
		font = "Georgia, serif"
	case "monospace", "mono":
		font = "monospace"
	}

	name := strings.TrimSpace(portfolio.Name)
	if name == "" {
		name = "Your Name"
	}

	role := strings.TrimSpace(portfolio.Role)
	if role == "" {
		role = "Software Developer"
	}

	bio := strings.TrimSpace(portfolio.Bio)
	if bio == "" {
		bio = "I build useful software, explore new technologies, and turn ideas into real products."
	}

	email := strings.TrimSpace(portfolio.Email)

	type NeonProject struct {
		Name        string
		Description string
	}

	projects := []NeonProject{
		{
			Name:        "Project One",
			Description: "A project built to solve a problem with thoughtful technology.",
		},
		{
			Name:        "Project Two",
			Description: "An experiment turned into a useful digital experience.",
		},
	}

	for i := 0; i < 2 && i < len(portfolio.Projects); i++ {
		if strings.TrimSpace(portfolio.Projects[i].Name) != "" {
			projects[i].Name = strings.TrimSpace(portfolio.Projects[i].Name)
		}

		if strings.TrimSpace(portfolio.Projects[i].Description) != "" {
			projects[i].Description = strings.TrimSpace(
				portfolio.Projects[i].Description,
			)
		}
	}

	const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="{{.Bio}}">
    <title>{{.Name}} — Digital Portfolio</title>
    <link rel="stylesheet" href="css/style.css">
</head>
<body>
    <main>
        <nav>
            <a class="logo" href="#home">{{.Name}}<span>.DEV</span></a>

            <div class="nav-links">
                <a href="#about">About</a>
                <a href="#skills">Skills</a>
                <a href="#projects">Projects</a>
            </div>
        </nav>

        <header class="hero" id="home">
            <div class="status">
                <span class="status-dot"></span>
                Available for interesting projects
            </div>

            <p class="role-label">{{.Role}}</p>

            <h1>
                Building the<br>
                <span>next layer</span><br>
                of technology.
            </h1>

            <p class="bio">{{.Bio}}</p>

            <div class="hero-actions">
                <a class="button button-primary" href="#projects">
                    Explore my work ↗
                </a>

                <a class="button button-secondary" href="#contact">
                    Get in touch
                </a>
            </div>
        </header>

        <section id="about">
            <p class="section-label">01 / About</p>
            <h2>Curiosity drives everything.</h2>
            <p class="about-text">{{.Bio}}</p>
        </section>

        <section id="skills">
            <p class="section-label">02 / Technical stack</p>
            <h2>Tools of the trade.</h2>

            <div class="skills">
                <span class="skill">Go</span>
                <span class="skill">Python</span>
                <span class="skill">REST APIs</span>
                <span class="skill">SQL</span>
                <span class="skill">HTML &amp; CSS</span>
                <span class="skill">Git</span>
                <span class="skill">AI / ML</span>
            </div>
        </section>

        <section id="projects">
            <p class="section-label">03 / Selected projects</p>
            <h2>Experiments turned into products.</h2>

            <div class="projects">
                {{range $i, $project := .Projects}}
                <article class="project">
                    <div class="project-top">
                        <span>0{{add $i 1}} — PROJECT</span>
                        <span>DEV</span>
                    </div>

                    <h3>{{$project.Name}}</h3>
                    <p>{{$project.Description}}</p>

                    <a href="#contact">Get in touch ↗</a>
                </article>
                {{end}}
            </div>
        </section>

        <section id="contact" class="contact">
            <p class="section-label">04 / Connection established</p>
            <h2>Have something interesting in mind?</h2>

            <p>
                I'm always curious about ambitious ideas
                and interesting technical challenges.
            </p>

            {{if .Email}}
            <a class="button button-primary" href="mailto:{{.Email}}">
                Send a message ↗
            </a>
            {{else}}
            <p class="contact-note">Thanks for stopping by.</p>
            {{end}}
        </section>

        <footer>
            <span>{{.Name}}<span class="footer-accent">.DEV</span></span>
            <p>Designed for the future · © 2026</p>
        </footer>
    </main>
</body>
</html>`

	tmpl, err := template.New("neon-tech").Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}).Parse(htmlTemplate)
	if err != nil {
		return nil, fmt.Errorf("parse Neon Tech HTML: %w", err)
	}

	data := struct {
		Name     string
		Role     string
		Bio      string
		Email    string
		Projects []NeonProject
	}{
		Name:     name,
		Role:     role,
		Bio:      bio,
		Email:    email,
		Projects: projects,
	}

	var htmlBuffer bytes.Buffer

	if err := tmpl.Execute(&htmlBuffer, data); err != nil {
		return nil, fmt.Errorf("generate Neon Tech HTML: %w", err)
	}

	css := `
:root {
    --background: #080a18;
    --surface: #101329;
    --surface-light: #171b35;
    --violet: ACCENT_COLOR;
    --cyan: #67e8f9;
    --white: #f3f4ff;
    --muted: #969bbd;
    --border: #252947;
    --font: FONT_FAMILY;
}

* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

html {
    scroll-behavior: smooth;
}

body {
    background: var(--background);
    color: var(--white);
    font-family: var(--font);
    line-height: 1.7;
}

a {
    color: inherit;
}

main {
    max-width: 1050px;
    margin: 0 auto;
    padding: 28px;
}

/* Navigation */

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
    color: var(--white);
    font-size: 19px;
    font-weight: 800;
    letter-spacing: -1px;
    text-decoration: none;
}

.logo span,
footer span {
    color: var(--violet);
}

.nav-links {
    display: flex;
    flex-wrap: wrap;
    gap: 22px;
}

nav a {
    color: var(--muted);
    text-decoration: none;
    font-size: 13px;
    transition: color .2s;
}

nav a:hover {
    color: var(--cyan);
}

/* Hero */

.hero {
    position: relative;
    padding: 105px 0 90px;
}

.status {
    display: inline-flex;
    align-items: center;
    gap: 9px;
    padding: 7px 13px;
    border: 1px solid var(--border);
    border-radius: 30px;
    color: var(--muted);
    font-size: 12px;
    background: var(--surface);
}

.status-dot {
    width: 7px;
    height: 7px;
    border-radius: 50%;
    background: var(--cyan);
    box-shadow: 0 0 12px var(--cyan);
}

.role-label {
    margin-top: 28px;
    color: var(--cyan);
    font-size: 13px;
    font-weight: 700;
    letter-spacing: 1.5px;
    text-transform: uppercase;
}

h1 {
    margin: 12px 0 20px;
    font-size: clamp(46px, 8vw, 88px);
    line-height: 1.05;
    letter-spacing: -4px;
}

h1 span {
    color: var(--violet);
    text-shadow: 0 0 30px #a78bfa30;
}

.bio {
    max-width: 650px;
    color: var(--muted);
    font-size: 18px;
    overflow-wrap: anywhere;
}

.hero-actions {
    display: flex;
    flex-wrap: wrap;
    gap: 14px;
    margin-top: 32px;
}

.button {
    display: inline-block;
    padding: 12px 19px;
    border-radius: 8px;
    text-decoration: none;
    font-size: 14px;
    font-weight: bold;
    transition: transform .2s, box-shadow .2s;
}

.button:hover {
    transform: translateY(-2px);
}

.button-primary {
    background: var(--violet);
    color: #100b24;
    box-shadow: 0 0 24px #a78bfa25;
}

.button-primary:hover {
    box-shadow: 0 0 30px #a78bfa50;
}

.button-secondary {
    border: 1px solid var(--border);
    color: var(--white);
}

/* Sections */

section {
    padding: 42px 0;
    border-top: 1px solid var(--border);
}

.section-label {
    color: var(--cyan);
    font-size: 12px;
    font-weight: bold;
    letter-spacing: 2px;
    text-transform: uppercase;
}

h2 {
    margin: 12px 0 20px;
    font-size: clamp(27px, 4vw, 38px);
    letter-spacing: -1px;
}

.about-text {
    max-width: 750px;
    color: var(--muted);
    overflow-wrap: anywhere;
}

/* Skills */

.skills {
    display: flex;
    flex-wrap: wrap;
    gap: 12px;
    margin-top: 24px;
}

.skill {
    padding: 9px 15px;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 8px;
    color: #d8d7f5;
    font-size: 13px;
    transition: border-color .2s;
}

.skill:hover {
    border-color: var(--violet);
}

/* Projects */

.projects {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 18px;
    margin-top: 28px;
}

.project {
    min-width: 0;
    padding: 25px;
    background: linear-gradient(
        145deg,
        var(--surface-light),
        var(--surface)
    );
    border: 1px solid var(--border);
    border-radius: 12px;
    transition: border-color .25s, transform .25s;
}

.project:hover {
    border-color: var(--violet);
    transform: translateY(-4px);
}

.project-top {
    display: flex;
    justify-content: space-between;
    gap: 12px;
    color: var(--cyan);
    font-size: 12px;
}

.project h3 {
    margin: 22px 0 10px;
    font-size: 22px;
    overflow-wrap: anywhere;
}

.project p {
    color: var(--muted);
    font-size: 14px;
    overflow-wrap: anywhere;
}

.project a {
    display: inline-block;
    margin-top: 20px;
    color: var(--violet);
    font-size: 13px;
    font-weight: bold;
    text-decoration: none;
}

.project a:hover {
    color: var(--cyan);
}

/* Contact */

.contact {
    padding: 32px;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 14px;
    margin-top: 18px;
}

.contact p {
    color: var(--muted);
}

.contact .button {
    margin-top: 22px;
}

.contact-note {
    margin-top: 18px;
}

/* Footer */

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

footer .footer-accent {
    color: var(--violet);
}

/* Mobile */

@media (max-width: 650px) {
    main {
        padding: 20px;
    }

    .nav-links {
        gap: 12px;
    }

    nav a {
        font-size: 12px;
    }

    .hero {
        padding: 75px 0 60px;
    }

    h1 {
        letter-spacing: -2px;
    }

    .bio {
        font-size: 16px;
    }

    .projects {
        grid-template-columns: 1fr;
    }

    .contact {
        padding: 24px;
    }
}
`

	css = strings.ReplaceAll(css, "ACCENT_COLOR", accent)
	css = strings.ReplaceAll(css, "FONT_FAMILY", font)

	var zipBuffer bytes.Buffer
	zipWriter := zip.NewWriter(&zipBuffer)

	htmlFile, err := zipWriter.Create("index.html")
	if err != nil {
		return nil, fmt.Errorf("create index.html in ZIP: %w", err)
	}

	if _, err := htmlFile.Write(htmlBuffer.Bytes()); err != nil {
		return nil, fmt.Errorf("write index.html to ZIP: %w", err)
	}

	cssFile, err := zipWriter.Create("css/style.css")
	if err != nil {
		return nil, fmt.Errorf("create CSS file in ZIP: %w", err)
	}

	if _, err := cssFile.Write([]byte(css)); err != nil {
		return nil, fmt.Errorf("write CSS to ZIP: %w", err)
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("finish ZIP archive: %w", err)
	}

	return zipBuffer.Bytes(), nil
}
