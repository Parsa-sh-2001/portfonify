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

func MakeZipEditorial(portfolio Portfolio) ([]byte, error) {
	// Prepare and escape user-provided content.
	name := html.EscapeString(strings.TrimSpace(portfolio.Name))
	role := html.EscapeString(strings.TrimSpace(portfolio.Role))
	bio := html.EscapeString(strings.TrimSpace(portfolio.Bio))
	email := strings.TrimSpace(portfolio.Email)

	if name == "" {
		name = "My Portfolio"
	}
	if role == "" {
		role = "Developer · Creator"
	}
	if bio == "" {
		bio = "Welcome to my portfolio."
	}

	// Portfonify generates exactly TWO projects.
	project1Name := "Untitled project"
	project1Description := "Add a description for this project."
	project2Name := "Untitled project"
	project2Description := "Add a description for this project."

	if len(portfolio.Projects) > 0 {
		if value := strings.TrimSpace(portfolio.Projects[0].Name); value != "" {
			project1Name = html.EscapeString(value)
		}
		if value := strings.TrimSpace(portfolio.Projects[0].Description); value != "" {
			project1Description = html.EscapeString(value)
		}
	}

	if len(portfolio.Projects) > 1 {
		if value := strings.TrimSpace(portfolio.Projects[1].Name); value != "" {
			project2Name = html.EscapeString(value)
		}
		if value := strings.TrimSpace(portfolio.Projects[1].Description); value != "" {
			project2Description = html.EscapeString(value)
		}
	}

	// Validate the email before creating the contact link.
	contactHTML := `<p>Add your email address to let people contact you.</p>`

	if email != "" && !strings.ContainsAny(email, "\r\n") {
		if parsed, err := mail.ParseAddress(email); err == nil && parsed.Address == email {
			contactHTML = fmt.Sprintf(
				`<p>Have an interesting idea or project in mind? I'd be glad to hear about it.</p>
				<a href="mailto:%s">Get in touch ↗</a>`,
				html.EscapeString(email),
			)
		}
	}

	// Original Editorial HTML design, adapted for dynamic content.
	page := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta name="theme-color" content="#f8f6f0">
	<title>%s — Portfolio</title>
	<link rel="stylesheet" href="css/style.css">
</head>
<body>
	<main>
		<nav>
			<div class="logo">%s<span>.</span></div>
			<div class="nav-links">
				<a href="#about">About</a>
				<a href="#projects">Work</a>
				<a href="#skills">Expertise</a>
				<a href="#contact">Contact</a>
			</div>
		</nav>

		<header class="hero">
			<p class="eyebrow">%s</p>

			<h1>
				A practice of<br>
				making things<br>
				<em>with purpose.</em>
			</h1>

			<p class="bio">%s</p>

			<div class="hero-meta">
				<span>Selected portfolio</span>
				<span>2026 — Present</span>
			</div>
		</header>

		<section id="about">
			<div class="section-heading">
				<p class="section-number">01 / About</p>
				<h2>A little context.</h2>
			</div>

			<p class="about-text">%s</p>
		</section>

		<section id="projects">
			<div class="section-heading">
				<p class="section-number">02 / Selected work</p>
				<h2>A few things I've made.</h2>
			</div>

			<article class="project">
				<span class="project-number">PROJECT 001</span>
				<div>
					<h3>%s</h3>
					<p>%s</p>
				</div>
				<a href="#contact">View ↗</a>
			</article>

			<article class="project">
				<span class="project-number">PROJECT 002</span>
				<div>
					<h3>%s</h3>
					<p>%s</p>
				</div>
				<a href="#contact">View ↗</a>
			</article>
		</section>

		<section id="skills">
			<div class="section-heading">
				<p class="section-number">03 / Expertise</p>
				<h2>Things I work with.</h2>
			</div>

			<div class="skills">
				<span class="skill">Creative Thinking</span>
				<span class="skill">Problem Solving</span>
				<span class="skill">Thoughtful Design</span>
				<span class="skill">Independent Work</span>
				<span class="skill">Building Useful Things</span>
				<span class="skill">Attention to Detail</span>
			</div>
		</section>

		<section id="contact" class="contact">
			<p class="section-number">04 / Contact</p>
			<h2>Good work starts with a conversation.</h2>
			%s
		</section>

		<footer>
			<span>%s · PORTFOLIO</span>
			<span>INDEPENDENTLY MADE · © 2026</span>
		</footer>
	</main>
</body>
</html>`,
		name,
		name,
		role,
		bio,
		bio,
		project1Name,
		project1Description,
		project2Name,
		project2Description,
		contactHTML,
		name,
	)

	// Original Editorial CSS.
	cssEditorial := `:root {
	--paper: #f8f6f0;
	--ink: #282921;
	--muted: #777768;
	--olive: ACCENT_PLACEHOLDER;
	--line: #dcd9ce;
	--surface: #efede4;
}

* {
	margin: 0;
	padding: 0;
	box-sizing: border-box;
}

html {
	scroll-behavior: smooth;
	scroll-padding-top: 24px;
}

body {
	background: var(--paper);
	color: var(--ink);
	font-family: Georgia, "Times New Roman", serif;
	line-height: 1.85;
}

main {
	max-width: 960px;
	margin: 0 auto;
	padding: 30px;
}

nav {
	display: flex;
	justify-content: space-between;
	align-items: center;
	flex-wrap: wrap;
	gap: 18px;
	padding-bottom: 20px;
	border-bottom: 1px solid var(--line);
}

.logo {
	font-size: 18px;
	font-weight: bold;
	letter-spacing: 1px;
	overflow-wrap: anywhere;
}

.logo span {
	color: var(--olive);
}

.nav-links {
	display: flex;
	flex-wrap: wrap;
	gap: 20px;
}

nav a {
	color: var(--muted);
	text-decoration: none;
	font-family: Arial, sans-serif;
	font-size: 13px;
}

nav a:hover {
	color: var(--olive);
}

.hero {
	padding: 100px 0 85px;
}

.eyebrow,
.section-number {
	color: var(--olive);
	font-family: Arial, sans-serif;
	font-size: 11px;
	font-weight: bold;
	letter-spacing: 3px;
	text-transform: uppercase;
}

h1 {
	max-width: 800px;
	font-size: clamp(48px, 8vw, 86px);
	font-weight: normal;
	line-height: 1.08;
	letter-spacing: -3px;
	margin: 25px 0;
	overflow-wrap: anywhere;
}

h1 em {
	color: var(--olive);
	font-weight: normal;
}

.bio {
	max-width: 620px;
	color: var(--muted);
	font-size: 19px;
	overflow-wrap: anywhere;
}

.hero-meta {
	display: flex;
	justify-content: space-between;
	flex-wrap: wrap;
	gap: 12px;
	margin-top: 50px;
	padding-top: 18px;
	border-top: 1px solid var(--line);
	color: var(--muted);
	font-family: Arial, sans-serif;
	font-size: 11px;
	letter-spacing: 1px;
	text-transform: uppercase;
}

section {
	padding: 45px 0;
	border-top: 1px solid var(--line);
}

.section-heading {
	display: grid;
	grid-template-columns: 150px 1fr;
	gap: 30px;
	margin-bottom: 25px;
}

.section-number {
	letter-spacing: 2px;
}

h2 {
	font-size: 32px;
	font-weight: normal;
	line-height: 1.3;
	overflow-wrap: anywhere;
}

.about-text {
	max-width: 680px;
	margin-left: 180px;
	font-size: 17px;
	overflow-wrap: anywhere;
}

.project {
	display: grid;
	grid-template-columns: 150px 1fr auto;
	align-items: start;
	gap: 30px;
	padding: 25px 0;
	border-top: 1px solid var(--line);
}

.project-number {
	color: var(--muted);
	font-family: Arial, sans-serif;
	font-size: 11px;
	letter-spacing: 1px;
}

.project h3 {
	font-size: 24px;
	font-weight: normal;
	margin-bottom: 5px;
	overflow-wrap: anywhere;
}

.project p {
	max-width: 480px;
	color: var(--muted);
	font-size: 15px;
	overflow-wrap: anywhere;
}

.project a {
	color: var(--olive);
	text-decoration: none;
	font-family: Arial, sans-serif;
	font-size: 13px;
	white-space: nowrap;
}

.project a:hover {
	text-decoration: underline;
}

.skills {
	display: flex;
	flex-wrap: wrap;
	gap: 10px;
	margin-left: 180px;
}

.skill {
	border: 1px solid var(--line);
	border-radius: 30px;
	padding: 7px 16px;
	color: var(--muted);
	font-size: 13px;
}

.contact {
	background: var(--surface);
	padding: 35px;
	margin-top: 20px;
}

.contact h2 {
	max-width: 600px;
	font-size: clamp(30px, 5vw, 46px);
}

.contact p {
	color: var(--muted);
	margin-top: 12px;
}

.contact .section-number {
	margin-top: 0;
}

.contact a {
	display: inline-block;
	margin-top: 20px;
	color: var(--olive);
	text-decoration: none;
	border-bottom: 1px solid var(--olive);
	overflow-wrap: anywhere;
}

footer {
	display: flex;
	justify-content: space-between;
	flex-wrap: wrap;
	gap: 12px;
	margin-top: 35px;
	padding-top: 20px;
	border-top: 1px solid var(--line);
	color: var(--muted);
	font-family: Arial, sans-serif;
	font-size: 11px;
	letter-spacing: 1px;
}

footer span {
	color: var(--olive);
}

@media (max-width: 650px) {
	main {
		padding: 20px;
	}

	.hero {
		padding: 70px 0 55px;
	}

	h1 {
		letter-spacing: -1.5px;
	}

	.section-heading {
		grid-template-columns: 1fr;
		gap: 8px;
	}

	.about-text,
	.skills {
		margin-left: 0;
	}

	.project {
		grid-template-columns: 1fr auto;
		gap: 8px;
	}

	.project-number {
		grid-column: 1 / -1;
	}

	.project h3 {
		font-size: 21px;
	}

	.contact {
		padding: 25px;
	}
}

@media (max-width: 380px) {
	main {
		padding: 16px;
	}

	.nav-links {
		gap: 12px;
	}

	.hero-meta {
		flex-direction: column;
	}

	.project {
		grid-template-columns: minmax(0, 1fr);
	}

	.project-number {
		grid-column: auto;
	}

	.project a {
		justify-self: start;
	}
}`

	// Use the selected accent color, with a safe fallback.
	accent := strings.TrimSpace(portfolio.Accent)
	validColor := regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

	if !validColor.MatchString(accent) {
		accent = "#737b51"
	}

	cssEditorial = strings.ReplaceAll(
		cssEditorial,
		"ACCENT_PLACEHOLDER",
		accent,
	)

	// Build the downloadable ZIP archive in memory.
	var buffer bytes.Buffer
	archive := zip.NewWriter(&buffer)

	index, err := archive.Create("index.html")
	if err != nil {
		_ = archive.Close()
		return nil, err
	}

	if _, err := index.Write([]byte(page)); err != nil {
		_ = archive.Close()
		return nil, err
	}

	style, err := archive.Create("css/style.css")
	if err != nil {
		_ = archive.Close()
		return nil, err
	}

	if _, err := style.Write([]byte(cssEditorial)); err != nil {
		_ = archive.Close()
		return nil, err
	}

	if err := archive.Close(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
