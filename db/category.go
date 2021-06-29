package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/angelinahung/product-category/utils"
)

// Category definition.
type Category struct {
	ID          int64  `json:"ID"`
	Name        string `json:"name"`
	IsInvisible int16  `json:"is_invisible,omitempty"`
}

// IsBadRequest for required request validation
func (c Category) IsBadRequest() bool {
	if c.ID <= 0 ||
		c.Name == "" {
		return true
	}
	return false
}

// IsRequired to validate data range / allowed value
func (c Category) IsRequired() bool {
	return true
}

// CreateCategory 新增目錄
func CreateCategory(db *sql.DB, tableName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := ioutil.ReadAll(r.Body)

		var category Category
		if err := json.Unmarshal(reqBody, &category); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if category.IsBadRequest() {
			http.Error(w, "fill in required field {id, name} plz", http.StatusBadRequest)
			return
		}

		sql := fmt.Sprintf(`INSERT INTO %s
			(id, name, is_invisible) 
			VALUES (?,?,?) `, tableName)
		result, err := db.Exec(sql,
			category.ID,
			category.Name,
			category.IsInvisible,
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

// QueryCategories 目錄查詢
func QueryCategories(db *sql.DB, tableName string) http.HandlerFunc {
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
		var queryResults []Category
		for rows.Next() {
			var category Category
			if err := rows.Scan(
				&category.ID,
				&category.Name,
				&category.IsInvisible); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			queryResults = append(queryResults, category)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(queryResults)
	}
}

// UpdateCategory 更改目錄
func UpdateCategory(db *sql.DB, tableName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok || id == "" {
			http.Error(w, "id is missing in parameters", http.StatusBadRequest)
			return
		}

		reqBody, _ := ioutil.ReadAll(r.Body)
		var category Category
		if err := json.Unmarshal(reqBody, &category); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sb := utils.NewBuilder("UPDATE ", tableName, " SET ")
		var sets []string
		var params []interface{}
		if category.Name != "" {
			sets = append(sets, "name = ?")
			params = append(params, category.Name)
		}
		if category.IsInvisible > -1 {
			sets = append(sets, "is_invisible = ?")
			params = append(params, category.IsInvisible)
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

// DeleteCategory 刪除目錄
func DeleteCategory(db *sql.DB, tableName string) http.HandlerFunc {
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
