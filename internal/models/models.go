package models

import "html/template"

type Category struct {
    ID          int
    Name        string
    Slug        string
    Description string
    ImagePath   string
}

type Product struct {
    ID          int
    Name        string
    Description string
    Dosage      string
    Data        string        // Additional data like protein/fat/energy (stored as text/markdown for now)
    Content     template.HTML // Full HTML content for the product page
    ImagePath   string
    CategoryID  int
}
