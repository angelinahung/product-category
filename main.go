package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/angelinahung/product-category/db"
)

func main() {
	fmt.Println("-- Welcome to NuEIP API --")

	sqlDB, err := sql.Open("mysql", "root:@tcp(localhost:3306)/product_category")
	if err != nil {
		panic(err.Error())
	}
	// defer the close till after the main function has finished
	// executing
	defer sqlDB.Close()

	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/product", db.CreateProduct(sqlDB, "product")).Methods("POST")
	myRouter.HandleFunc("/product", db.QueryProducts(sqlDB, "product")).Methods("GET")

	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	log.Fatal(http.ListenAndServe(":6666", myRouter))
}
