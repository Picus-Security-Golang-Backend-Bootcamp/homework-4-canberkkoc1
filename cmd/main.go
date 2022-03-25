package main

import (
	"ck/migration"
	"fmt"
	"net/http"

	"ck/controller"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/dgrijalva/jwt-go"
)

var myKey = []byte("yourname")

func isTokenValid(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("error")
				}

				return myKey, nil

			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}

		} else {
			http.Error(w, "not authorized", http.StatusUnauthorized)
		}

	})

}

func main() {

	migration.InitialMigration()

	r := mux.NewRouter()

	r.Handle("/create/book", isTokenValid(controller.AddBook)).Methods("POST")
	r.HandleFunc("/Allbooks", controller.GetAllBooks).Methods("GET")
	r.HandleFunc("/books/{name}", controller.GetBooksByName).Methods("GET")
	r.HandleFunc("/books/buy/{id}/{stock}", controller.BuyBook).Methods("PUT")
	r.HandleFunc("/books/delete/{id}", controller.DeleteBook).Methods("DELETE")
	r.HandleFunc("/token", controller.GetToken).Methods("GET")

	handler := cors.Default().Handler(r)

	http.ListenAndServe(":8080", handler)

}
