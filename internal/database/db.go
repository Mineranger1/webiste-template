package database

import (
	"database/sql"
	"html/template"

	"app/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	createTables := `
    CREATE TABLE IF NOT EXISTS categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        slug TEXT NOT NULL UNIQUE,
        description TEXT,
        image_path TEXT
    );

    CREATE TABLE IF NOT EXISTS products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        description TEXT NOT NULL,
        dosage TEXT,
        data TEXT,
        content TEXT,
        image_path TEXT,
        category_id INTEGER,
        FOREIGN KEY(category_id) REFERENCES categories(id)
    );

    CREATE TABLE IF NOT EXISTS employees (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        position TEXT NOT NULL,
        phone TEXT,
        email TEXT,
        image_path TEXT
    );
    `

	_, err = DB.Exec(createTables)
	if err != nil {
		return err
	}

	// Try to add the column if it doesn't exist (for existing DBs)
	// We ignore the error as it will fail if the column already exists
	_, _ = DB.Exec("ALTER TABLE products ADD COLUMN image_path TEXT")

	return nil
}

func GetCategories() ([]models.Category, error) {
	rows, err := DB.Query("SELECT id, name, slug, description, image_path FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		var description, imagePath sql.NullString
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &description, &imagePath); err != nil {
			return nil, err
		}
		c.Description = description.String
		c.ImagePath = imagePath.String
		categories = append(categories, c)
	}
	return categories, nil
}

func GetProducts() ([]models.Product, error) {
	rows, err := DB.Query("SELECT id, name, description, dosage, data, content, image_path, category_id FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var content, imagePath sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Dosage, &p.Data, &content, &imagePath, &p.CategoryID); err != nil {
			return nil, err
		}
		p.Content = template.HTML(content.String)
		p.ImagePath = imagePath.String
		products = append(products, p)
	}
	return products, nil
}

func GetProductsByCategory(categorySlug string) ([]models.Product, error) {
	rows, err := DB.Query(`
        SELECT p.id, p.name, p.description, p.dosage, p.data, p.content, p.image_path, p.category_id
        FROM products p
        JOIN categories c ON p.category_id = c.id
        WHERE c.slug = ?`, categorySlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var content, imagePath sql.NullString
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Dosage, &p.Data, &content, &imagePath, &p.CategoryID); err != nil {
			return nil, err
		}
		p.Content = template.HTML(content.String)
		p.ImagePath = imagePath.String
		products = append(products, p)
	}
	return products, nil
}

func GetProductByID(id int) (models.Product, error) {
	row := DB.QueryRow("SELECT id, name, description, dosage, data, content, image_path, category_id FROM products WHERE id = ?", id)

	var p models.Product
	var content, imagePath sql.NullString
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Dosage, &p.Data, &content, &imagePath, &p.CategoryID); err != nil {
		return p, err
	}
	p.Content = template.HTML(content.String)
	p.ImagePath = imagePath.String
	return p, nil
}

func GetEmployees() ([]models.Employee, error) {
	rows, err := DB.Query("SELECT id, name, position, phone, email, image_path FROM employees")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var e models.Employee
		var phone, email, imagePath sql.NullString
		if err := rows.Scan(&e.ID, &e.Name, &e.Position, &phone, &email, &imagePath); err != nil {
			return nil, err
		}
		e.Phone = phone.String
		e.Email = email.String
		e.ImagePath = imagePath.String
		employees = append(employees, e)
	}
	return employees, nil
}

func GetCategoryBySlug(slug string) (models.Category, error) {
	row := DB.QueryRow("SELECT id, name, slug, description, image_path FROM categories WHERE slug = ?", slug)
	var c models.Category
	var description, imagePath sql.NullString
	if err := row.Scan(&c.ID, &c.Name, &c.Slug, &description, &imagePath); err != nil {
		return c, err
	}
	c.Description = description.String
	c.ImagePath = imagePath.String
	return c, nil
}

func GetProductByName(name string) (models.Product, error) {
	row := DB.QueryRow("SELECT id, name, description, dosage, data, content, image_path, category_id FROM products WHERE name = ?", name)

	var p models.Product
	var content, imagePath sql.NullString
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Dosage, &p.Data, &content, &imagePath, &p.CategoryID); err != nil {
		return p, err
	}
	p.Content = template.HTML(content.String)
	p.ImagePath = imagePath.String
	return p, nil
}
