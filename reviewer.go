package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func ReviewZip(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uploaded, _, err := r.FormFile("zip")
	if err != nil {
		http.Error(w, "Could not receive ZIP file", http.StatusBadRequest)
		return
	}
	defer uploaded.Close()

	data, err := io.ReadAll(uploaded)
	if err != nil {
		http.Error(w, "Could not read ZIP file", http.StatusBadRequest)
		return
	}

	reader, err := zip.NewReader(
		bytes.NewReader(data),
		int64(len(data)),
	)
	if err != nil {
		http.Error(w, "Invalid ZIP file", http.StatusBadRequest)
		return
	}

	for _, file := range reader.File {
		fmt.Println(file.Name)
	}

	fmt.Fprintln(w, "ZIP read successfully!")
}
