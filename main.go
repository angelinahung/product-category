package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/angelinahung/product-category/db"
	"github.com/angelinahung/product-category/pkg/config"
	"github.com/angelinahung/product-category/pkg/logger"
)

var (
	tableProduct  = "product"
	tableCategory = "category"
)

func main() {
	// 1. Parse configurations
	config.ParseConfigurations()

	// 2. logger initialization.
	logger.Init()

	sqlDB, err := sql.Open("mysql", config.Options.DBSource)
	if err != nil {
		zap.L().Fatal("failed to connect to DB: ", zap.Error(err))
	}
	// defer the close till after the main function has finished
	// executing
	defer sqlDB.Close()

	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	// product CRUD
	myRouter.HandleFunc("/product", db.CreateProduct(sqlDB, tableProduct)).Methods("POST")
	myRouter.HandleFunc("/product", db.QueryProducts(sqlDB, tableProduct)).Methods("GET")
	myRouter.HandleFunc("/product/id/{id}", db.UpdateProduct(sqlDB, tableProduct)).Methods("PATCH")
	myRouter.HandleFunc("/product/id/{id}", db.DeleteProduct(sqlDB, tableProduct)).Methods("DELETE")

	// category CRUD
	myRouter.HandleFunc("/category", db.CreateCategory(sqlDB, tableCategory)).Methods("POST")
	myRouter.HandleFunc("/category", db.QueryCategories(sqlDB, tableCategory)).Methods("GET")
	myRouter.HandleFunc("/category/id/{id}", db.UpdateCategory(sqlDB, tableCategory)).Methods("PATCH")
	myRouter.HandleFunc("/category/id/{id}", db.DeleteCategory(sqlDB, tableCategory)).Methods("DELETE")

	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	addr := fmt.Sprintf("%s:%d", config.Options.Host, config.Options.Port)
	zap.L().Info("serving RESTful API ...", zap.String("address", addr))
	log.Fatal(http.ListenAndServe(addr, myRouter))
}
