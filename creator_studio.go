package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

func MakeZipCreativeStudio(portfolio Portfolio) ([]byte, error) {
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
		accent = "#b7ff65"
	}

	font := "Inter, Arial, sans-serif"

	switch strings.ToLower(strings.TrimSpace(portfolio.Font)) {
	case "inter":
		font = "Inter, Arial, sans-serif"
	case "arial":
		font = "Arial, sans-serif"
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
		role = "Creative Developer"
	}

	bio := strings.TrimSpace(portfolio.Bio)
	if bio == "" {
		bio = "I create thoughtful digital experiences and turn ideas into reality."
	}

	email := strings.TrimSpace(portfolio.Email)

	type CreativeProject struct {
		Name        string
		Description string
	}

	projects := make([]CreativeProject, 2)

	defaultProjects := []CreativeProject{
		{
			Name:        "Project One",
			Description: "A creative project built with care and attention to detail.",
		},
		{
			Name:        "Project Two",
			Description: "Another idea brought to life through design and technology.",
		},
	}

	for i := 0; i < 2; i++ {
		projects[i] = defaultProjects[i]

		if i < len(portfolio.Projects) {
			if strings.TrimSpace(portfolio.Projects[i].Name) != "" {
				projects[i].Name = strings.TrimSpace(portfolio.Projects[i].Name)
			}

			if strings.TrimSpace(portfolio.Projects[i].Description) != "" {
				projects[i].Description = strings.TrimSpace(
					portfolio.Projects[i].Description,
				)
			}
		}
	}

	const htmlTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="description" content="{{.Bio}}">
    <title>{{.Name}} — Creative Studio</title>
    <link rel="stylesheet" href="css/style.css">
</head>
<body>

    <nav class="navbar">
        <a class="brand" href="#home">{{.Name}}<span>.</span></a>
        <a class="nav-link" href="#contact">Let's talk ↗</a>
    </nav>

    <main>
        <section class="hero" id="home">
            <div class="eyebrow">
                <span class="status-dot"></span>
                AVAILABLE FOR OPPORTUNITIES
            </div>

            <h1>
                Ideas into<br>
                <span class="highlight">experiences.</span>
            </h1>

            <div class="hero-bottom">
                <p class="role">{{.Role}}</p>
                <p class="bio">{{.Bio}}</p>
            </div>

            <a class="scroll-link" href="#work">SCROLL TO EXPLORE ↓</a>
        </section>

        <section class="work-section" id="work">
            <div class="section-heading">
                <div>
                    <p class="eyebrow">SELECTED WORK</p>
                    <h2>Things I've made<span class="highlight">.</span></h2>
                </div>
                <span class="work-count">02 PROJECTS</span>
            </div>

            <div class="project-grid">
                {{range $index, $project := .Projects}}
                <article class="project-card">
                    <div class="project-art art-{{$index}}">
                        <span class="art-number">0{{add $index 1}}</span>
                        <span class="art-symbol">{{if eq $index 0}}✳{{else}}◈{{end}}</span>
                        <span class="art-caption">CREATIVE WORK / 0{{$index | add 1}}</span>
                    </div>

                    <div class="project-info">
                        <div>
                            <h3>{{$project.Name}}</h3>
                            <p>{{$project.Description}}</p>
                        </div>
                        <span class="project-arrow">↗</span>
                    </div>
                </article>
                {{end}}
            </div>
        </section>

        <section class="about-section">
            <p class="eyebrow">A LITTLE ABOUT ME</p>
            <h2>Curiosity first.<br><span class="highlight">Always.</span></h2>
            <p class="about-copy">{{.Bio}}</p>
        </section>

        <section class="contact-section" id="contact">
            <p class="eyebrow">HAVE SOMETHING IN MIND?</p>
            <h2>Let's make<br>something <span class="highlight">great.</span></h2>

            {{if .Email}}
            <a class="contact-button" href="mailto:{{.Email}}">
                Get in touch <span>↗</span>
            </a>
            {{else}}
            <p class="contact-note">Thanks for stopping by.</p>
            {{end}}
        </section>
    </main>

    <footer>
        <a class="brand" href="#home">{{.Name}}<span>.</span></a>
        <p>DESIGNED WITH INTENTION.</p>
        <a href="#home" class="back-top">BACK TO TOP ↑</a>
    </footer>

</body>
</html>`

	tmpl, err := template.New("creative-studio").Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	}).Parse(htmlTemplate)
	if err != nil {
		return nil, fmt.Errorf("parse Creative Studio HTML: %w", err)
	}

	data := struct {
		Name     string
		Role     string
		Bio      string
		Email    string
		Projects []CreativeProject
	}{
		Name:     name,
		Role:     role,
		Bio:      bio,
		Email:    email,
		Projects: projects,
	}

	var htmlBuffer bytes.Buffer

	if err := tmpl.Execute(&htmlBuffer, data); err != nil {
		return nil, fmt.Errorf("generate Creative Studio HTML: %w", err)
	}

	css := `
@import url('https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800;900&display=swap');

:root {
    --bg: #11120f;
    --panel: #191a16;
    --text: #f3f2eb;
    --muted: #a2a398;
    --border: #303129;
    --accent: ACCENT_COLOR;
    --font: FONT_FAMILY;
}

* {
    box-sizing: border-box;
    scroll-behavior: smooth;
}

body {
    margin: 0;
    background: var(--bg);
    color: var(--text);
    font-family: var(--font);
    -webkit-font-smoothing: antialiased;
}

a {
    color: inherit;
    text-decoration: none;
}

.navbar {
    min-height: 82px;
    padding: 0 8%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    border-bottom: 1px solid var(--border);
}

.brand {
    font-size: 19px;
    font-weight: 800;
    letter-spacing: -0.8px;
}

.brand span,
.highlight {
    color: var(--accent);
}

.nav-link {
    color: var(--muted);
    font-size: 13px;
    transition: color .2s ease;
}

.nav-link:hover {
    color: var(--accent);
}

.hero {
    min-height: 650px;
    padding: 120px 12% 70px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    border-bottom: 1px solid var(--border);
}

.eyebrow {
    margin: 0;
    color: var(--muted);
    font-size: 10px;
    font-weight: 800;
    letter-spacing: 2px;
}

.status-dot {
    display: inline-block;
    width: 7px;
    height: 7px;
    margin-right: 7px;
    border-radius: 50%;
    background: var(--accent);
}

.hero h1 {
    margin: 32px 0 48px;
    font-size: clamp(64px, 10vw, 132px);
    line-height: .94;
    letter-spacing: -8px;
    font-weight: 800;
}

.hero-bottom {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 40px;
    max-width: 760px;
}

.role {
    margin: 0;
    font-size: 17px;
    font-weight: 700;
}

.bio {
    margin: 0;
    color: var(--muted);
    font-size: 14px;
    line-height: 1.9;
}

.scroll-link {
    margin-top: 75px;
    color: var(--muted);
    font-size: 10px;
    font-weight: 700;
    letter-spacing: 1.7px;
}

.work-section {
    padding: 110px 8%;
}

.section-heading {
    display: flex;
    justify-content: space-between;
    align-items: end;
    gap: 24px;
    margin-bottom: 42px;
}

.section-heading h2,
.about-section h2,
.contact-section h2 {
    margin: 14px 0 0;
    font-size: clamp(38px, 5vw, 65px);
    letter-spacing: -3px;
    line-height: 1.08;
}

.work-count {
    color: var(--muted);
    font-size: 10px;
    letter-spacing: 1.5px;
}

.project-grid {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 24px;
}

.project-card {
    min-width: 0;
}

.project-art {
    min-height: 290px;
    position: relative;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    background: var(--panel);
    border: 1px solid var(--border);
    transition: border-color .25s ease;
}

.project-card:hover .project-art {
    border-color: var(--accent);
}

.art-0 {
    background:
        radial-gradient(ellipse at 50% 50%, #343b2b 0%, transparent 60%),
        #191a16;
}

.art-1 {
    background:
        radial-gradient(ellipse at 50% 50%, #34312a 0%, transparent 60%),
        #191a16;
}

.art-number {
    position: absolute;
    top: 20px;
    left: 22px;
    color: var(--muted);
    font-size: 11px;
}

.art-symbol {
    color: var(--accent);
    font-size: clamp(90px, 13vw, 160px);
    line-height: 1;
    transition: transform .35s ease;
}

.project-card:hover .art-symbol {
    transform: rotate(18deg) scale(1.08);
}

.art-caption {
    position: absolute;
    bottom: 20px;
    left: 22px;
    color: var(--muted);
    font-size: 9px;
    letter-spacing: 1.5px;
}

.project-info {
    display: flex;
    align-items: start;
    justify-content: space-between;
    gap: 20px;
    padding: 22px 0 30px;
    border-bottom: 1px solid var(--border);
}

.project-info h3 {
    margin: 0 0 10px;
    font-size: 21px;
    letter-spacing: -.6px;
}

.project-info p {
    max-width: 390px;
    margin: 0;
    color: var(--muted);
    font-size: 13px;
    line-height: 1.8;
}

.project-arrow {
    color: var(--accent);
    font-size: 23px;
}

.about-section {
    padding: 100px 12%;
    background: var(--panel);
    border-top: 1px solid var(--border);
    border-bottom: 1px solid var(--border);
}

.about-section h2 {
    margin-bottom: 28px;
    font-size: clamp(45px, 7vw, 85px);
}

.about-copy {
    max-width: 620px;
    color: var(--muted);
    font-size: 15px;
    line-height: 2;
}

.contact-section {
    padding: 120px 12%;
}

.contact-section h2 {
    margin-bottom: 35px;
    font-size: clamp(52px, 8vw, 100px);
}

.contact-button {
    display: inline-flex;
    min-width: 180px;
    align-items: center;
    justify-content: space-between;
    gap: 35px;
    padding: 17px 20px;
    background: var(--accent);
    color: #11120f;
    font-size: 13px;
    font-weight: 800;
    transition: transform .2s ease;
}

.contact-button:hover {
    transform: translateY(-3px);
}

.contact-note {
    color: var(--muted);
    font-size: 13px;
}

footer {
    min-height: 90px;
    padding: 24px 8%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 20px;
    border-top: 1px solid var(--border);
}

footer p,
.back-top {
    color: var(--muted);
    font-size: 9px;
    font-weight: 700;
    letter-spacing: 1.5px;
}

.back-top:hover {
    color: var(--accent);
}

@media (max-width: 700px) {
    .navbar {
        min-height: 70px;
        padding: 0 6%;
    }

    .hero {
        min-height: 560px;
        padding: 90px 7% 55px;
    }

    .hero h1 {
        font-size: clamp(56px, 14vw, 88px);
        letter-spacing: -4px;
    }

    .hero-bottom {
        grid-template-columns: 1fr;
        gap: 16px;
    }

    .scroll-link {
        margin-top: 50px;
    }

    .work-section {
        padding: 75px 6%;
    }

    .section-heading {
        align-items: start;
        flex-direction: column;
    }

    .project-grid {
        grid-template-columns: 1fr;
        gap: 8px;
    }

    .project-art {
        min-height: 250px;
    }

    .about-section,
    .contact-section {
        padding: 75px 7%;
    }

    footer {
        padding: 24px 6%;
        flex-wrap: wrap;
    }

    footer p {
        display: none;
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
