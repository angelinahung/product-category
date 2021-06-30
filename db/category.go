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

	ParentID int64 `json:"parent_id"` // 上層目錄
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

// // CategoryLevelOne lists the most upper category IDs
// func CategoryLevelOne(db *sql.DB, tableName string) ([]int64, error) {
// 	sql := fmt.Sprintf("SELECT id FROM %s WHERE parent_id=?", tableName)
// 	rows, err := db.Query(sql, 0) // 找到最上層目錄
// 	if err != nil {
// 		return nil, err
// 	}

// 	//切記用完都要做 Close
// 	defer rows.Close()
// 	var ids []int64
// 	for rows.Next() {
// 		var category Category
// 		if err := rows.Scan(
// 			&category.ID,
// 		); err != nil {
// 			return nil, err
// 		}
// 		ids = append(ids, category.ID)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return ids, nil
// }

// 新增目錄之前檢查即將新增的目錄的父目錄是否已超過5層
func isParentIDAvailable(db *sql.DB, tableName string, parentID int64) (bool, error) {
	sql := fmt.Sprintf("SELECT id FROM %s WHERE parent_id=?", tableName)
	rows, err := db.Query(sql, parentID)
	if err != nil {
		return false, err
	}

	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(
			&id,
		); err != nil {
			return false, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return false, err
	}
	if len(ids) > 4 {
		return false, nil
	}
	return true, nil
}

// 檢查是否其他category的parentID = id
func isTheMostLowerCategory(db *sql.DB, tableName string, id int64) (bool, error) {
	sql := fmt.Sprintf("SELECT id FROM %s WHERE parent_id=?", tableName)
	rows, err := db.Query(sql, id)
	if err != nil {
		return false, err
	}

	defer rows.Close()
	for rows.Next() {
		var id int64
		if err := rows.Scan(
			&id,
		); err != nil {
			return false, err
		}
		if id > 0 {
			return false, nil
		}
	}
	if err := rows.Err(); err != nil {
		return false, err
	}
	return true, nil
}

// CreateCategory 新增目錄
// TODO: 檢查新增新目錄的上層目錄是否有產品list, 若有產品清單故不允許新增
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

		ok, err := isParentIDAvailable(db, tableName, category.ParentID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "each category is limited to 5 level only", http.StatusBadRequest)
			return
		}

		ok, err = IsProductsAvailableOnCategory(db, "product", category.ParentID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			http.Error(w, "the corresponding parent category have products on it", http.StatusBadRequest)
			return
		}

		sql := fmt.Sprintf(`INSERT INTO %s
			(id, name, is_invisible, parent_id) 
			VALUES (?,?,?,?) `, tableName)
		result, err := db.Exec(sql,
			category.ID,
			category.Name,
			category.IsInvisible,
			category.ParentID,
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
				&category.IsInvisible,
				&category.ParentID); err != nil {
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
// 必須給上層目錄,否這更新失敗
// TODO: 檢查該目錄修改是否違反 最多5層目錄的限制
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
		if category.ParentID > -1 {
			sets = append(sets, "parent_id = ?")
			params = append(params, category.ParentID)
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
// TODO: 檢查要刪除的目錄是否有下層目錄/產品清單
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
