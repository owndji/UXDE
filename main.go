package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"

	// "math/rand"
	"net/http"

	// "github.com/bxcodec/faker/v4"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db *sql.DB
)

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
	ImageURL    string
}

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./sqlite/dataUX.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Crée les tables si elles n'existent pas
	// files := []string{"product.sql"}
	// for _, file := range files {
	// 	sqlFile, err := ioutil.ReadFile("sqlite/" + file)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	_, err = db.Exec(string(sqlFile))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Printf("File %s executed successfully", file)
	// }

	http.HandleFunc("/home", Home)

	// products := []struct {
	// 	Name        string
	// 	Description string
	// 	Price       float64
	// 	ImageURL    string
	// }{
	// 	{"Laptop", "A high-performance laptop.", 999.99, "https://imgs.search.brave.com/OWvZ1XkBn3LUXcXtK_rHgOQbzLLXFhB0VYPAMqZ_BJE/rs:fit:500:0:0:0/g:ce/aHR0cHM6Ly9tZWRp/YS5pc3RvY2twaG90/by5jb20vaWQvMTY0/NjM3MDc5L3Bob3Rv/L2xhcHRvcC5qcGc_/cz02MTJ4NjEyJnc9/MCZrPTIwJmM9S1U1/NEJxQ2wwcmZNNjJW/QlJ2WXZtUG9qWTRN/V29rSXhUYTBRTnJ4/SFF6az0"},
	// 	{"Smartphone", "A modern smartphone with a great camera.", 699.99, "https://example.com/images/smartphone.jpg"},
	// 	{"Headphones", "Noise-cancelling over-ear headphones.", 199.99, "https://example.com/images/headphones.jpg"},
	// 	{"Keyboard", "A mechanical keyboard with RGB lighting.", 89.99, "https://example.com/images/keyboard.jpg"},
	// 	{"Monitor", "A 27-inch 4K UHD monitor.", 299.99, "https://example.com/images/monitor.jpg"},
	// }

	// stmt, err := db.Prepare(`INSERT INTO products (name, description, price, image_url) VALUES (?, ?, ?, ?)`)
	// if err != nil {
	// 	log.Fatalf("Error preparing statement: %v", err)
	// }
	// defer stmt.Close()

	// for _, p := range products {
	// 	_, err := stmt.Exec(p.Name, p.Description, p.Price, p.ImageURL)
	// 	if err != nil {
	// 		log.Printf("Error inserting product %s: %v", p.Name, err)
	// 	} else {
	// 		fmt.Printf("Inserted product: %s\n", p.Name)
	// 	}
	// }

	// fmt.Println("Data insertion completed successfully!")
	// for i := 1; i <= 600; i++ {
	// 	product := Product{
	// 		Name:        faker.Word(),
	// 		Description: faker.Sentence(),
	// 		Price:       rand.Float64()*900 + 100,
	// 		ImageURL:    faker.URL(),
	// 	}

	// 	_, err := stmt.Exec(product.Name, product.Description, product.Price, product.ImageURL)
	// 	if err != nil {
	// 		log.Printf("Error inserting product %d: %v", i, err)
	// 	} else {
	// 		fmt.Printf("Inserted product %d: %s\n", i, product.Name)
	// 	}
	// }

	fmt.Println("Server started at http://localhost:8081/home") // URL principale
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func getProduct() ([]Product, error) {
	rows, err := db.Query("SELECT id, name, description, price, image_url FROM products")
	if err != nil {
		return nil, fmt.Errorf("error querying products: %v", err)
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var prod Product
		if err := rows.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.ImageURL); err != nil {
			return nil, fmt.Errorf("error scanning products: %v", err)
		}
		products = append(products, prod)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over products: %v", err)
	}

	return products, nil
}

func Home(w http.ResponseWriter, r *http.Request) {
	products, err := getProduct()
	if err != nil {
		log.Printf("Error retrieving products: %v", err)
		http.Error(w, "Error retrieving products", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Products": products,
	}

	renderTemplate(w, "tmpl/home.html", data)
}

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Printf("Error parsing template %s: %v", tmpl, err)
		http.Error(w, "Error loading page", http.StatusInternalServerError)
		return
	}
	if err := t.Execute(w, data); err != nil {
		log.Printf("Error executing template %s: %v", tmpl, err)
		http.Error(w, "Error rendering page", http.StatusInternalServerError)
	}
}
