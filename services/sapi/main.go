package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type Product struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	db, err := sql.Open("postgres", "host=postgresdb.default.svc.cluster.local user=pqgotest password=password dbname=pqgotest port=5432 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			rows, err := db.Query("SELECT id, name FROM products")
			if err != nil {
				http.Error(w, http.StatusText(500), 500)
				log.Println(err)
				return
			}
			defer rows.Close()

			products := []Product{}
			for rows.Next() {
				var p Product
				err := rows.Scan(&p.ID, &p.Name)
				if err != nil {
					http.Error(w, http.StatusText(500), 500)
					log.Println(err)
					return
				}
				products = append(products, p)
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(products)

		case "POST":
			var p Product
			err := json.NewDecoder(r.Body).Decode(&p)
			if err != nil {
				http.Error(w, http.StatusText(400), 400)
				log.Println(err)
				return
			}
			db.Exec("CREATE TABLE products")
			_, err = db.Exec("INSERT INTO products (name) VALUES ($1)", p.Name)
			if err != nil {
				http.Error(w, http.StatusText(500), 500)
				log.Println(err)
				return
			}

			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, http.StatusText(405), 405)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
 
