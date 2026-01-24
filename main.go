package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings",
	"os"
)

// Category represents a category in the system
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// In-memory storage (sementara, nanti ganti database)
var categories = []Category{
	{ID: 1, Name: "Makanan", Description: "Kategori produk makanan"},
	{ID: 2, Name: "Minuman", Description: "Kategori produk minuman"},
	{ID: 3, Name: "Bumbu", Description: "Kategori produk bumbu dapur"},
}

func main() {
	// GET  /categories/{id}
	// PUT  /categories/{id}
	// DELETE /categories/{id}
	http.HandleFunc("/categories/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getCategoryByID(w, r)
		case http.MethodPut:
			updateCategory(w, r)
		case http.MethodDelete:
			deleteCategory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// GET  /categories
	// POST /categories
	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getAllCategories(w, r)
		case http.MethodPost:
			createCategory(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// /health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

    port := os.Getenv("PORT")
    if port == "" {
      port = "8080"
    }

    fmt.Println("Server running di localhost:", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("gagal running server:", err)
	}
}

func getAllCategories(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, categories)
}

// POST /categories
func createCategory(w http.ResponseWriter, r *http.Request) {
	var payload Category
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	payload.Name = strings.TrimSpace(payload.Name)
	payload.Description = strings.TrimSpace(payload.Description)

	if payload.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	payload.ID = nextCategoryID()
	categories = append(categories, payload)

	writeJSON(w, http.StatusCreated, payload)
}

// GET /categories/{id}
func getCategoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/categories/")
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	for _, c := range categories {
		if c.ID == id {
			writeJSON(w, http.StatusOK, c)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

// PUT /categories/{id}
func updateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/categories/")
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	var payload Category
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	payload.Name = strings.TrimSpace(payload.Name)
	payload.Description = strings.TrimSpace(payload.Description)

	if payload.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	for i := range categories {
		if categories[i].ID == id {
			payload.ID = id
			categories[i] = payload
			writeJSON(w, http.StatusOK, payload)
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

// DELETE /categories/{id}
func deleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r.URL.Path, "/categories/")
	if err != nil {
		http.Error(w, "Invalid Category ID", http.StatusBadRequest)
		return
	}

	for i, c := range categories {
		if c.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			writeJSON(w, http.StatusOK, map[string]string{
				"message": "sukses delete",
			})
			return
		}
	}

	http.Error(w, "Category not found", http.StatusNotFound)
}

// Helpers

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func parseIDFromPath(path string, prefix string) (int, error) {
	idStr := strings.TrimPrefix(path, prefix)
	idStr = strings.Trim(idStr, "/")
	return strconv.Atoi(idStr)
}

func nextCategoryID() int {
	maxID := 0
	for _, c := range categories {
		if c.ID > maxID {
			maxID = c.ID
		}
	}
	return maxID + 1
}
