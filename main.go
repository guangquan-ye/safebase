package main

import (
    "fmt"
    "net/http"
)

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Bienvenue sur SafeBase! ANAIS JEREMY")
}

func main() {
    http.HandleFunc("/", homePage)
    fmt.Println("Le serveur démarre sur http://localhost:8081")
    http.ListenAndServe(":8081", nil)
}
