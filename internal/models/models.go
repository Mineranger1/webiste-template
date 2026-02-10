package models

type Category struct {
    ID   int
    Name string
    Slug string
}

type Product struct {
    ID          int
    Name        string
    Description string
    Dosage      string
    Data        string // Additional data like protein/fat/energy (stored as text/markdown for now)
    CategoryID  int
}
