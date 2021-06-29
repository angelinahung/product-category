# Product Category

Product-Category Management API using Golang

## Prerequisite

- [Go 1.16](https://golang.org/doc/install)
- MySQL admin (I used phpMyAdmin to access and set DB)

## Package Dependencies

1. go get github.com/go-sql-driver/mysql
2. go get github.com/gorilla/mux
3. github.com/go-openapi/swag
4. github.com/jessevdk/go-flags
5. go.uber.org/zap

## Compilation

1. serve RESTful API at:

    By default:

    ```browser
    http://localhost:8000
    ```

    if you prefer to serve at specify address:

    ```bash
    go run . --host=${{YOUR_HOST}} --port=${{YOUR_PORT}}
    ```

    For development mode environment for better logging experience:

    ```bash
    go run . --dev_mode
    ```

    To specify DB Source Name,

    ```bash
    go run . --db_source=${{YOUR_DB_SOURCE_NAME}}
    ```

2. RESTful API Doc:

    PLEASE SEE [swagger.yaml](swagger.yaml) on [swagger editor](https://editor.swagger.io/)