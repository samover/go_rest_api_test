package main

import (
    "fmt"
    "log"
    "strconv"
    "net/http"
    "encoding/json"
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

    var err error
    app.DB, err = sql.Open("postgres", connectionString)
    if err != nil {
        log.Fatal(err)
    }

    app.Router = mux.NewRouter()
    app.initializeRoutes()
}

func (app *App) Run(addr string) {
    log.Fatal(http.ListenAndServe(":8000", app.Router))
}

func (app *App) initializeRoutes() {
    app.Router.HandleFunc("/products", app.getProducts).Methods("GET")
    app.Router.HandleFunc("/products", app.createProduct).Methods("POST")
    app.Router.HandleFunc("/product/{id:[0-9]+}", app.getProduct).Methods("GET")
    app.Router.HandleFunc("/product/{id:[0-9]+}", app.updateProduct).Methods("PUT")
    app.Router.HandleFunc("/product/{id:[0-9]+}", app.deleteProduct).Methods("DELETE")
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func (app *App) getProduct(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])

    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid Product Id")
        return
    }

    p := product{ Id: id }

    if err := p.getProduct(app.DB); err != nil {
        switch err {
            case sql.ErrNoRows:
                respondWithError(w, http.StatusNotFound, "Product not found")
            default:
                respondWithError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }

    respondWithJSON(w, http.StatusOK, p)
}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
    count, _ := strconv.Atoi(r.FormValue("count"))
    start, _ := strconv.Atoi(r.FormValue("start"))

    if (count > 10 || count < 1) {
        count = 10
    }
    if start < 0 {
        start = 0
    }

    p := product{}

    products, err := p.getProducts(app.DB, start, count)
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, products)
}

func (app *App) createProduct(w http.ResponseWriter, r *http.Request) {
    var p product
    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&p); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    defer r.Body.Close()

    if err := p.createProduct(app.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusCreated, p)
}

func (app *App) updateProduct(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])

    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid Product Id")
        return
    }

    var p product
    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&p); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    defer r.Body.Close()
    p.Id = id

    if err := p.updateProduct(app.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, p)
}

func (app *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])

    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid Product Id")
        return
    }

    p := product{ Id: id }

    if err := p.deleteProduct(app.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }

    respondWithJSON(w, http.StatusOK, map[string]string{"result":"success"})
}
