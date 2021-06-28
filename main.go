package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/angelinahung/product-category/db"
)

var (
	tableProduct = "product"
	tableCategory = "category"
)

func main() {
	sqlDB, err := sql.Open("mysql", "root:@tcp(localhost:3306)/product_category")
	if err != nil {
		panic(err.Error())
	}
	// defer the close till after the main function has finished
	// executing
	defer sqlDB.Close()

	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	// product CRUD
	myRouter.HandleFunc("/product", db.CreateProduct(sqlDB,  tableProduct)).Methods("POST")
	myRouter.HandleFunc("/product", db.QueryProducts(sqlDB, tableProduct)).Methods("GET")
	myRouter.HandleFunc("/product/id/{id}", db.UpdateProduct(sqlDB, tableProduct)).Methods("PATCH")
	myRouter.HandleFunc("/product/id/{id}", db.DeleteProduct(sqlDB,  tableProduct)).Methods("DELETE")

	// category CRUD
	myRouter.HandleFunc("/category", db.CreateCategory(sqlDB, tableCategory)).Methods("POST")
	myRouter.HandleFunc("/category", db.QueryCategories(sqlDB, tableCategory)).Methods("GET")
	myRouter.HandleFunc("/category/id/{id}", db.UpdateCategory(sqlDB, tableCategory)).Methods("PATCH")
	myRouter.HandleFunc("/category/id/{id}", db.DeleteCategory(sqlDB, tableCategory)).Methods("DELETE")

	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	log.Fatal(http.ListenAndServe(":8000", myRouter))
}
