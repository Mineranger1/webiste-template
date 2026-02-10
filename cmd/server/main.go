package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

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
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/contact", contactHandler)
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

	// base.html defines the layout and executes "content" block (which pageFile defines)
	// We execute "base.html" (filename as template name)
	if err := tmpl.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("Error rendering page %s: %v", pageFile, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func renderPartial(w http.ResponseWriter, parentFile string, partialName string, data interface{}) {
	// For partials defined inside a page file (e.g. "product_list" inside "products.html")
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

type ProductsPageData struct {
	Categories      []models.Category
	Products        []models.Product
	CurrentCategory string
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	categorySlug := r.URL.Query().Get("category")

	var products []models.Product
	var err error

	if categorySlug != "" && categorySlug != "all" {
		products, err = database.GetProductsByCategory(categorySlug)
	} else {
		products, err = database.GetProducts()
	}

	if err != nil {
		log.Printf("Error fetching products: %v", err)
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	// If HTMX request, render only the product list partial
	if r.Header.Get("HX-Request") == "true" {
		renderPartial(w, "products.html", "product_list", products)
		return
	}

	// Otherwise render full page
	categories, err := database.GetCategories()
	if err != nil {
		log.Printf("Error fetching categories: %v", err)
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	data := ProductsPageData{
		Categories:      categories,
		Products:        products,
		CurrentCategory: categorySlug,
	}
	renderPage(w, "products.html", data)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "about.html", nil)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "contact.html", nil)
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
</urlset>`))
}
