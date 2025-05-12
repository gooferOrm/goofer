package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gooferOrm/goofer/dialect"
	"github.com/gooferOrm/goofer/repository"
	"github.com/gooferOrm/goofer/schema"
)

// Product entity
type Product struct {
	ID          uint      `orm:"primaryKey;autoIncrement"`
	Name        string    `orm:"type:varchar(255);notnull"`
	Description string    `orm:"type:text"`
	Price       float64   `orm:"type:float;notnull"`
	CategoryID  uint      `orm:"index;notnull"`
	InStock     bool      `orm:"type:boolean;default:true"`
	CreatedAt   time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name for the Product entity
func (Product) TableName() string {
	return "products"
}

// Category entity
type Category struct {
	ID          uint      `orm:"primaryKey;autoIncrement"`
	Name        string    `orm:"type:varchar(100);notnull;unique"`
	Description string    `orm:"type:text"`
	CreatedAt   time.Time `orm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

// TableName returns the table name for the Category entity
func (Category) TableName() string {
	return "categories"
}

// ProductSummary is a DTO for summary data
type ProductSummary struct {
	CategoryName  string
	ProductCount  int
	AveragePrice  float64
	MaxPrice      float64
	MinPrice      float64
}

// ProductWithCategory is a DTO for join results
type ProductWithCategory struct {
	ProductID     uint
	ProductName   string
	Price         float64
	CategoryID    uint
	CategoryName  string
}

func main() {
	// Open SQLite database
	db, err := sql.Open("sqlite3", "./custom_queries.db")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create dialect
	sqliteDialect := dialect.NewSQLiteDialect()

	// Register entities
	if err := schema.Registry.RegisterEntity(Product{}); err != nil {
		log.Fatalf("Failed to register Product entity: %v", err)
	}
	if err := schema.Registry.RegisterEntity(Category{}); err != nil {
		log.Fatalf("Failed to register Category entity: %v", err)
	}

	// Create tables
	productMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(Product{}))
	categoryMeta, _ := schema.Registry.GetEntityMetadata(schema.GetEntityType(Category{}))

	_, err = db.Exec(sqliteDialect.CreateTableSQL(categoryMeta))
	if err != nil {
		log.Fatalf("Failed to create categories table: %v", err)
	}

	_, err = db.Exec(sqliteDialect.CreateTableSQL(productMeta))
	if err != nil {
		log.Fatalf("Failed to create products table: %v", err)
	}

	// Create repositories
	productRepo := repository.NewRepository[Product](db, sqliteDialect)
	categoryRepo := repository.NewRepository[Category](db, sqliteDialect)

	// Create sample data
	fmt.Println("Creating sample data...")
	
	// Create categories
	categories := []Category{
		{Name: "Electronics", Description: "Electronic devices and gadgets"},
		{Name: "Books", Description: "Books, e-books, and audiobooks"},
		{Name: "Clothing", Description: "Apparel and accessories"},
		{Name: "Home & Kitchen", Description: "Home appliances and kitchenware"},
	}
	
	for i := range categories {
		if err := categoryRepo.Save(&categories[i]); err != nil {
			log.Fatalf("Failed to save category: %v", err)
		}
		fmt.Printf("Created category: %s (ID: %d)\n", categories[i].Name, categories[i].ID)
	}
	
	// Create products
	products := []Product{
		{Name: "Smartphone", Description: "Latest smartphone", Price: 799.99, CategoryID: categories[0].ID, InStock: true},
		{Name: "Laptop", Description: "High-performance laptop", Price: 1299.99, CategoryID: categories[0].ID, InStock: true},
		{Name: "Headphones", Description: "Noise-cancelling headphones", Price: 199.99, CategoryID: categories[0].ID, InStock: true},
		{Name: "Tablet", Description: "Portable tablet", Price: 499.99, CategoryID: categories[0].ID, InStock: false},
		
		{Name: "Python Programming", Description: "Programming book", Price: 39.99, CategoryID: categories[1].ID, InStock: true},
		{Name: "Science Fiction Novel", Description: "Bestselling sci-fi", Price: 19.99, CategoryID: categories[1].ID, InStock: true},
		{Name: "Cookbook", Description: "Recipe collection", Price: 29.99, CategoryID: categories[1].ID, InStock: true},
		
		{Name: "T-Shirt", Description: "Cotton t-shirt", Price: 24.99, CategoryID: categories[2].ID, InStock: true},
		{Name: "Jeans", Description: "Denim jeans", Price: 49.99, CategoryID: categories[2].ID, InStock: true},
		{Name: "Dress", Description: "Evening dress", Price: 89.99, CategoryID: categories[2].ID, InStock: false},
		
		{Name: "Blender", Description: "Kitchen blender", Price: 79.99, CategoryID: categories[3].ID, InStock: true},
		{Name: "Coffee Maker", Description: "Automatic coffee maker", Price: 129.99, CategoryID: categories[3].ID, InStock: true},
	}
	
	for i := range products {
		if err := productRepo.Save(&products[i]); err != nil {
			log.Fatalf("Failed to save product: %v", err)
		}
	}
	
	fmt.Printf("Created %d products\n", len(products))
	
	fmt.Println("\n=== Basic Query Builder Examples ===")
	
	// Example 1: Find products with price greater than 100
	expensiveProducts, err := productRepo.Find().
		Where("price > ?", 100).
		OrderBy("price DESC").
		All()
	
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	
	fmt.Printf("\nProducts with price > 100 (%d items):\n", len(expensiveProducts))
	for _, p := range expensiveProducts {
		fmt.Printf("- %s: $%.2f\n", p.Name, p.Price)
	}
	
	// Example 2: Find products in a specific category (Electronics)
	electronicsProducts, err := productRepo.Find().
		Where("category_id = ?", categories[0].ID).
		OrderBy("name ASC").
		All()
	
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	
	fmt.Printf("\nElectronics products (%d items):\n", len(electronicsProducts))
	for _, p := range electronicsProducts {
		fmt.Printf("- %s: $%.2f\n", p.Name, p.Price)
	}
	
	// Example 3: Find out of stock products
	outOfStockProducts, err := productRepo.Find().
		Where("in_stock = ?", false).
		All()
	
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	
	fmt.Printf("\nOut of stock products (%d items):\n", len(outOfStockProducts))
	for _, p := range outOfStockProducts {
		fmt.Printf("- %s (Category ID: %d)\n", p.Name, p.CategoryID)
	}
	
	// Example 4: Count products by category
	for _, category := range categories {
		count, err := productRepo.Find().
			Where("category_id = ?", category.ID).
			Count()
		
		if err != nil {
			log.Fatalf("Count query failed: %v", err)
		}
		
		fmt.Printf("Category '%s' has %d products\n", category.Name, count)
	}
	
	fmt.Println("\n=== Advanced/Custom Query Examples ===")
	
	// Example 5: Raw SQL query - product with category join
	var productsWithCategory []ProductWithCategory
	
	rows, err := db.Query(`
		SELECT p.id, p.name, p.price, c.id, c.name 
		FROM products p
		JOIN categories c ON p.category_id = c.id
		ORDER BY p.price DESC
	`)
	
	if err != nil {
		log.Fatalf("Raw query failed: %v", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var pwc ProductWithCategory
		if err := rows.Scan(&pwc.ProductID, &pwc.ProductName, &pwc.Price, &pwc.CategoryID, &pwc.CategoryName); err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		productsWithCategory = append(productsWithCategory, pwc)
	}
	
	fmt.Printf("\nProducts with category details (%d items):\n", len(productsWithCategory))
	for _, p := range productsWithCategory {
		fmt.Printf("- %s ($%.2f) - Category: %s\n", p.ProductName, p.Price, p.CategoryName)
	}
	
	// Example 6: Aggregate data - product summary by category
	var productSummaries []ProductSummary
	
	rows, err = db.Query(`
		SELECT 
			c.name, 
			COUNT(p.id) as product_count, 
			AVG(p.price) as avg_price,
			MAX(p.price) as max_price,
			MIN(p.price) as min_price
		FROM categories c
		LEFT JOIN products p ON c.id = p.category_id
		GROUP BY c.id
		ORDER BY c.name
	`)
	
	if err != nil {
		log.Fatalf("Aggregate query failed: %v", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var ps ProductSummary
		if err := rows.Scan(&ps.CategoryName, &ps.ProductCount, &ps.AveragePrice, &ps.MaxPrice, &ps.MinPrice); err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		productSummaries = append(productSummaries, ps)
	}
	
	fmt.Printf("\nProduct summary by category:\n")
	for _, ps := range productSummaries {
		fmt.Printf("Category: %s\n", ps.CategoryName)
		fmt.Printf("  - Product Count: %d\n", ps.ProductCount)
		fmt.Printf("  - Average Price: $%.2f\n", ps.AveragePrice)
		fmt.Printf("  - Price Range: $%.2f - $%.2f\n", ps.MinPrice, ps.MaxPrice)
	}
	
	// Example 7: Complex query - find products with price above category average
	rows, err = db.Query(`
		SELECT p.id, p.name, p.price, c.name, avg_price
		FROM products p
		JOIN categories c ON p.category_id = c.id
		JOIN (
			SELECT category_id, AVG(price) as avg_price
			FROM products
			GROUP BY category_id
		) avg ON p.category_id = avg.category_id
		WHERE p.price > avg.avg_price
		ORDER BY (p.price - avg.avg_price) DESC
	`)
	
	if err != nil {
		log.Fatalf("Complex query failed: %v", err)
	}
	defer rows.Close()
	
	fmt.Printf("\nProducts with price above category average:\n")
	for rows.Next() {
		var productID uint
		var productName string
		var price float64
		var categoryName string
		var avgPrice float64
		
		if err := rows.Scan(&productID, &productName, &price, &categoryName, &avgPrice); err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		
		fmt.Printf("- %s ($%.2f) - Category: %s (Avg: $%.2f, Diff: +$%.2f)\n", 
			productName, price, categoryName, avgPrice, price-avgPrice)
	}
	
	// Example 8: Transaction with multiple queries
	fmt.Println("\n=== Transaction Example ===")
	
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
	}
	
	// Create a new category
	newCategoryName := "Sale Items"
	_, err = tx.Exec("INSERT INTO categories (name, description) VALUES (?, ?)", 
		newCategoryName, "Products currently on sale")
	
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to insert category: %v", err)
	}
	
	// Get the ID of the inserted category
	var newCategoryID uint
	err = tx.QueryRow("SELECT id FROM categories WHERE name = ?", newCategoryName).Scan(&newCategoryID)
	
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to get new category ID: %v", err)
	}
	
	// Find some products to put on sale (ones > $100)
	rows, err = tx.Query("SELECT id FROM products WHERE price > 100")
	if err != nil {
		tx.Rollback()
		log.Fatalf("Failed to query products: %v", err)
	}
	
	var productIDs []uint
	for rows.Next() {
		var id uint
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			tx.Rollback()
			log.Fatalf("Failed to scan product ID: %v", err)
		}
		productIDs = append(productIDs, id)
	}
	rows.Close()
	
	// Move the products to the sale category
	for _, id := range productIDs {
		_, err = tx.Exec("UPDATE products SET category_id = ? WHERE id = ?", newCategoryID, id)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to update product category: %v", err)
		}
	}
	
	// Commit the transaction
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}
	
	fmt.Printf("Transaction successful: Moved %d products to '%s' category\n", 
		len(productIDs), newCategoryName)
	
	// Verify the changes
	rows, err = db.Query(`
		SELECT p.name, p.price, c.name
		FROM products p
		JOIN categories c ON p.category_id = c.id
		WHERE c.name = ?
	`, newCategoryName)
	
	if err != nil {
		log.Fatalf("Verification query failed: %v", err)
	}
	defer rows.Close()
	
	fmt.Printf("\nProducts in '%s' category:\n", newCategoryName)
	for rows.Next() {
		var name string
		var price float64
		var categoryName string
		
		if err := rows.Scan(&name, &price, &categoryName); err != nil {
			log.Fatalf("Scan failed: %v", err)
		}
		
		fmt.Printf("- %s ($%.2f)\n", name, price)
	}
}