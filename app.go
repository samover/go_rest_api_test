package main

import (
    "fmt"
    "log"
    "database/sql"
    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

type App struct {
    Router *mux.Router
    DB *sql.DB
}

func (app *App) Initialize(user, password, dbname string) {
    connectionString :=
        fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)
    fmt.Println(connectionString)

    var err error
    app.DB, err = sql.Open("postgres", connectionString)
    if err != nil {
        log.Fatal(err)
    }

    app.Router = mux.NewRouter()
}

func (app *App) Run(addr string) {}
