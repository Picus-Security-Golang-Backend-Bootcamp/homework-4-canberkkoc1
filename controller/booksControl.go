package controller

import (
	"ck/handlers"
	"ck/helper"
	"ck/migration"
	"ck/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

var myKey = []byte("yourname")

func GenerateJWT() (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	tokenString, err := token.SignedString(myKey)

	if err != nil {
		fmt.Errorf("something went wrong %s", err.Error())
		return "", err
	}

	return tokenString, nil

}

//? http://localhost:8080/token
func GetToken(w http.ResponseWriter, r *http.Request) {
	validToken, err := GenerateJWT()

	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	fmt.Fprint(w, validToken)
}

//? http://localhost:8080//create/book
func AddBook(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	book := models.Books{
		StockNumber: helper.RandomNumber(1, 45),
		PageNumber:  helper.RandomNumber(1, 200),
		Price:       helper.RandomFloat(1, 100),
		Name:        "",
		StockCode:   helper.RandomString(5),
		Isbn:        helper.RandomString(6),
		AuthorName:  models.Author{}.AuthorName,
	}
	json.NewDecoder(r.Body).Decode(&book)

	if book.Name == "" {
		http.Error(w, "book_name or author_name empty", http.StatusBadRequest)

	} else {

		migration.DB.Create(&book)

		json.NewEncoder(w).Encode(book)

		w.WriteHeader(http.StatusCreated)
	}

}

//? http://localhost:8080/Allbooks
func GetAllBooks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	books, err := handlers.GetAllBook()

	if books == nil {
		http.Error(w, "Object empty", http.StatusInternalServerError)
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(books)

}

//? GET http://localhost/books/{name}

func GetBooksByName(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	param := mux.Vars(r)

	name := param["name"]

	if name == "" {
		http.Error(w, "Please enter the book name ", http.StatusBadRequest)
	}

	booksByName, err := handlers.GetBookByName(name)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(booksByName)

}

//? PUT http://localhost/books/buy/{id}/{stock}
func BuyBook(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])
	buyNumber, _ := strconv.Atoi(params["stock"])

	updateStock, err := handlers.UpdateStock(id, buyNumber)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(updateStock)

}

//? DELETE http://localhost/books/delete/{id}
func DeleteBook(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	param := mux.Vars(r)

	id, _ := strconv.Atoi(param["id"])

	deleteBook, err := handlers.DeleteBookById(id)

	if err != nil {
		http.Error(w, "id not found", http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(deleteBook)

}
