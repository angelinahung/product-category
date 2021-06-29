package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/angelinahung/product-category/utils"
)

// Product definition.
type Product struct {
	ID            int64     `json:"ID"`
	Name          string    `json:"name"`
	Budget        int64     `json:"budget"`
	Price         int64     `json:"price"`
	Description   string    `json:"description"`
	IsSale        int16     `json:"is_sale,omitempty"`
	StartSaleTime time.Time `json:"start_sale_time"`
	EndSaleTime   time.Time `json:"end_sale_time"`
}

var timeFormat = "2006-01-02 15:04:05"

// IsBadRequest for required request validation
func (p Product) IsBadRequest() bool {
	if p.ID <= 0 ||
		p.Name == "" ||
		p.Description == "" ||
		p.StartSaleTime.IsZero() ||
		p.EndSaleTime.IsZero() {
		return true
	}
	return false
}

// IsRequired to validate data range / allowed value
func (p Product) IsRequired() bool {
	if p.Budget > p.Price ||
		p.EndSaleTime.Before(p.StartSaleTime) {
		return false
	}
	return true
}

// CreateProduct 新增產品
func CreateProduct(db *sql.DB, tableName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := ioutil.ReadAll(r.Body)

		var product Product
		if err := json.Unmarshal(reqBody, &product); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if product.IsBadRequest() {
			http.Error(w, "fill in the required field plz", http.StatusBadRequest)
			return
		}

		startSaleTime := product.StartSaleTime.Format(timeFormat)
		endSaleTime := product.EndSaleTime.Format(timeFormat)

		if !product.IsRequired() {
			http.Error(w, "product budget must less than product price! and "+
				"start sale time must be <= end sale time!", http.StatusPreconditionFailed)
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
			http.Error(w, "no rows affected", http.StatusExpectationFailed)
			return
		}
		// json.NewEncoder(w).Encode(product)
		io.WriteString(w, "created")
	}
}

// QueryProducts 查詢產品
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
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			product.StartSaleTime, err = time.Parse(timeFormat, startSaleTime)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			product.EndSaleTime, err = time.Parse(timeFormat, endSaleTime)
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

// UpdateProduct 更改產品
func UpdateProduct(db *sql.DB, tableName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok || id == "" {
			http.Error(w, "id is missing in parameters", http.StatusBadRequest)
			return
		}

		// TODO: query product first to verify the requirement

		reqBody, _ := ioutil.ReadAll(r.Body)
		var product Product
		if err := json.Unmarshal(reqBody, &product); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sb := utils.NewBuilder("UPDATE ", tableName, " SET ")
		var sets []string
		var params []interface{}
		if product.Name != "" {
			sets = append(sets, "name = ?")
			params = append(params, product.Name)
		}
		if product.Description != "" {
			sets = append(sets, "description = ?")
			params = append(params, product.Description)
		}
		if product.Budget > 0 {
			sets = append(sets, "budget = ?")
			params = append(params, product.Budget)
		}
		if product.Price > 0 {
			sets = append(sets, "price = ?")
			params = append(params, product.Price)
		}
		if product.IsSale > -1 {
			sets = append(sets, "is_sale = ?")
			params = append(params, product.IsSale)
		}
		if !product.StartSaleTime.IsZero() {
			sets = append(sets, "start_sale_time = ?")
			params = append(params, product.StartSaleTime.Format(timeFormat))
		}
		if !product.EndSaleTime.IsZero() {
			sets = append(sets, "end_sale_time = ?")
			params = append(params, product.EndSaleTime.Format(timeFormat))
		}
		for i, set := range sets {
			sb.WriteString(set)
			if i == len(sets)-1 {
				break
			}
			sb.WriteString(", ")
		}
		sb.WriteString(" WHERE id = ?")
		params = append(params, id)
		result, err := db.Exec(sb.String(), params...)
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
			http.Error(w, "no rows affected", http.StatusExpectationFailed)
			return
		}
		io.WriteString(w, "updated")
	}
}

// DeleteProduct 刪除產品
func DeleteProduct(db *sql.DB, tableName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok || id == "" {
			http.Error(w, "id is missing in parameters", http.StatusBadRequest)
			return
		}

		sb := utils.NewBuilder("DELETE FROM ", tableName, " WHERE id = ?")
		result, err := db.Exec(sb.String(), id)
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
			http.Error(w, "no rows affected", http.StatusExpectationFailed)
			return
		}
		io.WriteString(w, "deleted")
	}
}
