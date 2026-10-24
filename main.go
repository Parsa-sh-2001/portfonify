// Aval dafter be nam'e izade dana || Sane' va parvardegar Hayy tavana //

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

type Portfolio struct {
	Name     string    `json:"name"`
	Bio      string    `json:"bio"`
	Theme    string    `json:"theme"`
	Accent   string    `json:"accent"`
	Font     string    `json:"font"`
	Role     string    `json:"role"`
	Image    string    `json:"image"`
	Email    string    `json:"email"`
	Projects []Project `json:"projects"`
}

func GenerateSite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var portfolio Portfolio

	err := json.NewDecoder(r.Body).Decode(&portfolio)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = json.MarshalIndent(portfolio, "", "    ")
	if err != nil {
		http.Error(w, "Couldnt Unmarshal...", http.StatusBadRequest)
		return
	}

	fmt.Println("Name: ", portfolio.Name)
	fmt.Println("Theme: ", portfolio.Theme)
	fmt.Println("Accent: ", portfolio.Accent)
	fmt.Println("Font: ", portfolio.Font)
	fmt.Println("Projects: ", len(portfolio.Projects))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var ZipByte []byte
	fmt.Println(portfolio.Theme)
	switch strings.ToLower(portfolio.Theme) {
	case "minimal":
		ZipByte, err = makeZipMinimal(portfolio)
	case "editorial":
		ZipByte, err = MakeZipEditorial(portfolio)
	case "creative":
		ZipByte, err = MakeZipCreativeStudio(portfolio)
	case "neon":
		ZipByte, err = MakeZipNeonTech(portfolio)
	default:
		http.Error(w, "Invalid theme!", http.StatusBadRequest)
		return
	}

	if err != nil {
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="portfonify-portfolio.zip"`)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(ZipByte)))

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(ZipByte)
	if err != nil {
		fmt.Println("Couldnt send zip... : ", err)
		return
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "./Assets/home_page.html")
	} else {
		return
	}
}

func CreatorDashboard(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./Assets/creator_deploy.html")
}

func main() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/create", CreatorDashboard)
	http.HandleFunc("/api/create", GenerateSite)
	http.HandleFunc("/api/review", ReviewZip)
	http.HandleFunc("/preview/creativestudio", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./Assets/creative_studio.html")
	})
	http.HandleFunc("/preview/editorial", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./Assets/editorial.html")
	})
	http.HandleFunc("/preview/minimal", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./Assets/minimal.html")
	})
	http.HandleFunc("/preview/neon_tech", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./Assets/neon_tech.html")
	})
	http.HandleFunc("/reviewer", func(w http.ResponseWriter, r *http.Request) {
		path := "./Assets/reviewer.html"
		http.ServeFile(w, r, path)
	})
	port := os.Getenv("PORT")
	if port == "" {
	    port = "8080"
	}
	
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
