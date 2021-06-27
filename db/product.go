package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Product definition.
type Product struct {
	ID            int64     `json:"ID"`
	Name          string    `json:"name"`
	Budget        int64     `json:"budget"`
	Price         int64     `json:"price"`
	Description   string    `json:"description"`
	IsSale        int16     `json:"is_sale"`
	StartSaleTime time.Time `json:"start_sale_time"`
	EndSaleTime   time.Time `json:"end_sale_time"`
}

// CreateProdcut 新增產品
func CreateProduct(db *sql.DB, tableName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := ioutil.ReadAll(r.Body)

		var product Product
		if err := json.Unmarshal(reqBody, &product); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if product.ID <= 0 ||
			product.Name == "" ||
			product.Description == "" ||
			product.StartSaleTime.IsZero() ||
			product.EndSaleTime.IsZero() {
			http.Error(w, "fill in required field plz", http.StatusBadRequest)
			return
		}

		startSaleTime := product.StartSaleTime.Format("2006-01-02 15:04:05")
		endSaleTime := product.EndSaleTime.Format("2006-01-02 15:04:05")

		if product.Budget > product.Price {
			http.Error(w, "product budget must less than product price!", http.StatusPreconditionFailed)
			return
		}

		sql := fmt.Sprintf(`INSERT INTO %s
			(id, name, budget, price, description, is_sale, start_sale_time, end_sale_time) 
			VALUES (?,?,?,?,?,?,?,?) `, tableName)
		result, err := db.Exec(sql,
			product.ID,
			product.Name,
			product.Budget,
			product.Price,
			product.Description,
			product.IsSale,
			startSaleTime,
			endSaleTime,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if rowsAffected != 1 {
			http.Error(w, err.Error(), http.StatusExpectationFailed)
			return
		}
		json.NewEncoder(w).Encode(product)
	}
}

func QueryProducts(db *sql.DB, tableName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sql := fmt.Sprintf("SELECT * FROM %s", tableName)
		var opts []interface{}
		var where string
		if id := r.URL.Query().Get("id"); id != "" {
			where = " WHERE id = ?"
			opts = append(opts, id)
		}
		if name := r.URL.Query().Get("name"); name != "" {
			if where != "" {
				where += " AND name = ?"
			} else {
				where = " WHERE name = ?"
			}
			opts = append(opts, name)
		}

		rows, err := db.Query(sql+where, opts...)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//切記用完都要做 Close
		defer rows.Close()
		var queryResults []Product
		for rows.Next() {
			var product Product
			var startSaleTime, endSaleTime string
			if err := rows.Scan(
				&product.ID,
				&product.Name,
				&product.Budget,
				&product.Price,
				&product.Description,
				&product.IsSale,
				&startSaleTime,
				&endSaleTime); err != nil {
				log.Fatal(err)
			}
			product.StartSaleTime, err = time.Parse("2006-01-02 15:04:05", startSaleTime)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			product.EndSaleTime, err = time.Parse("2006-01-02 15:04:05", endSaleTime)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			queryResults = append(queryResults, product)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(queryResults)
	}
}

func UpdateProduct(db *sql.DB, tableName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func DeleteProduct(db *sql.DB, tableName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
