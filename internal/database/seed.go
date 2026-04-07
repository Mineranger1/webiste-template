package database

import (
	"fmt"
	"log"
)

type SeedProduct struct {
	Name        string
	Description string
	Dosage      string
	Data        string
	Category    string
	Content     string
	ImagePath   string
}

func Seed() error {
	if err := seedProducts(); err != nil {
		return err
	}
	return seedEmployees()
}

func seedEmployees() error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM employees").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	employees := []struct {
		Name      string
		Position  string
		Phone     string
		Email     string
		ImagePath string
	}{
		{
			Name:      "Jane Doe",
			Position:  "CEO // Consultant",
			Phone:     "+1 555 010 011",
			Email:     "jane@example.com",
			ImagePath: "/static/images/placeholder_employee.webp",
		},
		{
			Name:      "John Smith",
			Position:  "Sales Director",
			Phone:     "+1 555 010 012",
			Email:     "john@example.com",
			ImagePath: "/static/images/placeholder_employee.webp",
		},
		{
			Name:      "Alice Johnson",
			Position:  "Customer Support",
			Phone:     "+1 555 010 013",
			Email:     "alice@example.com",
			ImagePath: "/static/images/placeholder_employee.webp",
		},
	}

	for _, e := range employees {
		_, err := DB.Exec("INSERT INTO employees (name, position, phone, email, image_path) VALUES (?, ?, ?, ?, ?)",
			e.Name, e.Position, e.Phone, e.Email, e.ImagePath)
		if err != nil {
			return fmt.Errorf("failed to insert employee %s: %w", e.Name, err)
		}
	}
	return nil
}

func seedProducts() error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil 
	}

	categories := []struct {
		Name        string
		Slug        string
		Description string
		ImagePath   string
	}{
		{
			Name:        "Category One",
			Slug:        "cat-1",
			Description: "First product category.",
			ImagePath:   "/static/images/placeholder.webp",
		},
		{
			Name:        "Category Two",
			Slug:        "cat-2",
			Description: "Second product category.",
			ImagePath:   "/static/images/placeholder.webp",
		},
	}

	categoryIDs := make(map[string]int)

	for _, cat := range categories {
		res, err := DB.Exec("INSERT INTO categories (name, slug, description, image_path) VALUES (?, ?, ?, ?)",
			cat.Name, cat.Slug, cat.Description, cat.ImagePath)
		if err != nil {
			return fmt.Errorf("failed to insert category %s: %w", cat.Name, err)
		}
		id, _ := res.LastInsertId()
		categoryIDs[cat.Slug] = int(id)
	}

	featuredContent := `
<div class="space-y-6">
    <h2 class="text-2xl font-bold text-primary">Featured Product Highlight</h2>
    <p class="text-lg">This is a great product that highlights key features.</p>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-8 my-8">
        <div class="bg-gray-50 p-6 rounded-lg shadow-sm">
            <h3 class="text-xl font-bold text-secondary mb-4">Why choose this?</h3>
            <ul class="space-y-2 list-disc list-inside">
                <li>Feature A</li>
                <li>Feature B</li>
                <li>Feature C</li>
            </ul>
        </div>
        <div class="bg-gray-50 p-6 rounded-lg shadow-sm">
            <h3 class="text-xl font-bold text-secondary mb-4">Specs</h3>
            <ul class="space-y-2">
                <li><strong>Weight:</strong> 1kg</li>
                <li><strong>Material:</strong> Steel</li>
            </ul>
        </div>
    </div>
</div>
`

	products := []SeedProduct{
		{
			Name:        "SUPER FAT", // The code has a bug where it searches for "SUPER FAT" explicitly for featured, we'll keep the name but change the rest, or just keep it so it doesn't break productsHandler
			Category:    "cat-1",
			Description: "Our flagship template product.",
			Content:     featuredContent,
			Data:        "Some data / specific metrics",
			Dosage:      "As needed",
			ImagePath:   "/static/images/placeholder.webp",
		},
		{
			Name:        "Standard Product 1",
			Category:    "cat-1",
			Description: "A standard product description.",
			Data:        "Standard data",
			Dosage:      "2 per day",
		},
		{
			Name:        "Standard Product 2",
			Category:    "cat-2",
			Description: "Another standard product.",
			Data:        "Info here",
			Dosage:      "1 per week",
		},
	}

	for _, p := range products {
		catID, ok := categoryIDs[p.Category]
		if !ok {
			log.Printf("Category not found for product %s: %s", p.Name, p.Category)
			continue
		}

		content := p.Content
		if content == "" {
			content = fmt.Sprintf("<p>%s</p>", p.Description)
		}

		_, err := DB.Exec(`INSERT INTO products (name, description, dosage, data, content, image_path, category_id) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			p.Name, p.Description, p.Dosage, p.Data, content, p.ImagePath, catID)
		if err != nil {
			log.Printf("Failed to insert product %s: %v", p.Name, err)
		}
	}

	return nil
}
