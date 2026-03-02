package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"biomix/internal/database"
	"biomix/internal/models"
)

func main() {
	// Initialize Database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "biomix.db"
	}

	if err := database.InitDB(dbPath); err != nil {
		log.Fatal(err)
	}

	// Seed Database
	if err := database.Seed(); err != nil {
		log.Fatal(err)
	}

	// Static Files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/product/", productHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/brand", brandHandler)
	http.HandleFunc("/sitemap.xml", sitemapHandler)
	http.HandleFunc("/robots.txt", robotsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func renderPage(w http.ResponseWriter, pageFile string, data interface{}) {
	// Parse base layout + specific page
	tmpl, err := template.ParseFiles("web/templates/base.html", "web/templates/"+pageFile)
	if err != nil {
		log.Printf("Error parsing templates for %s: %v", pageFile, err)
		http.Error(w, "Template Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Error rendering page %s: %v", pageFile, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func renderPartial(w http.ResponseWriter, parentFile string, partialName string, data interface{}) {
	tmpl, err := template.ParseFiles("web/templates/" + parentFile)
	if err != nil {
		log.Printf("Error parsing template %s for partial %s: %v", parentFile, partialName, err)
		http.Error(w, "Template Error", http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, partialName, data); err != nil {
		log.Printf("Error rendering partial %s: %v", partialName, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	renderPage(w, "home.html", nil)
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	categorySlug := r.URL.Query().Get("category")

	if categorySlug != "" && categorySlug != "all" {
		// Single Category View
		category, err := database.GetCategoryBySlug(categorySlug)
		if err != nil {
			log.Printf("Category not found: %s", categorySlug)
			http.NotFound(w, r)
			return
		}

		products, err := database.GetProductsByCategory(categorySlug)
		if err != nil {
			log.Printf("Error fetching products: %v", err)
			http.Error(w, "Database Error", http.StatusInternalServerError)
			return
		}

		data := struct {
			Category models.Category
			Products []models.Product
		}{
			Category: category,
			Products: products,
		}
		renderPage(w, "category.html", data)
		return
	}

	// All Categories View (Main Products Page)
	categories, err := database.GetCategories()
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	// Fetch featured product
	featuredProduct, err := database.GetProductByName("SUPER FAT")
	var featured *models.Product
	if err == nil {
		featured = &featuredProduct
	} else {
		log.Printf("Featured product not found: %v", err)
	}

	data := struct {
		Categories      []models.Category
		FeaturedProduct *models.Product
	}{
		Categories:      categories,
		FeaturedProduct: featured,
	}

	renderPage(w, "products.html", data)
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path /product/{id}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.NotFound(w, r)
		return
	}
	idStr := pathParts[2]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	product, err := database.GetProductByID(id)
	if err != nil {
		log.Printf("Error fetching product %d: %v", id, err)
		http.NotFound(w, r)
		return
	}

	renderPage(w, "product_detail.html", product)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	employees, err := database.GetEmployees()
	if err != nil {
		log.Printf("Error fetching employees: %v", err)
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}
	renderPage(w, "about.html", employees)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "contact.html", nil)
}

func brandHandler(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "brand.html", nil)
}

func robotsHandler(w http.ResponseWriter, r *http.Request) {
	robots := `User-agent: *
Allow: /
Sitemap: https://biomixpoland.pl/sitemap.xml`
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(robots))
}

func sitemapHandler(w http.ResponseWriter, r *http.Request) {
	baseUrl := "https://biomixpoland.pl"

	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
   <url><loc>` + baseUrl + `/</loc></url>
   <url><loc>` + baseUrl + `/products</loc></url>
   <url><loc>` + baseUrl + `/about</loc></url>
   <url><loc>` + baseUrl + `/contact</loc></url>
   <url><loc>` + baseUrl + `/brand</loc></url>
</urlset>`))
}
