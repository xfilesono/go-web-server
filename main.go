package main

import (
	"encoding/json"
	"fmt"
	"go-web-server/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// gathering information from form.html
func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "Registration Successful..\n")
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	mail := r.FormValue("mail")
	passwd := r.FormValue("passwd")
	HashPassword(passwd)

	fmt.Fprintf(w, "Name = %s\n", name)
	fmt.Fprintf(w, "Phone = %s\n", phone)
	fmt.Fprintf(w, "E-Mail = %s\n", mail)
	fmt.Fprintf(w, "Password = %s\n", "**********")
}

// Encryption user's password by using golang.org/x/crypto
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// function where we get all movies from our database
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

var movies []models.Movies

func main() {

	movies = append(movies, models.Movies{ID: "1", Year: "2005", Title: "Herşey Güzel Olacak", Director: &models.Director{FirstName: "Onur", LastName: "Kruzowa"}})
	movies = append(movies, models.Movies{ID: "2", Year: "1999", Title: "Adam Olacak Çocuk", Director: &models.Director{FirstName: "Aston", LastName: "Martin"}})

	r := mux.NewRouter()

	fileServer := http.FileServer(http.Dir("./static"))
	//http.Handle("/", fileServer)
	//http.HandleFunc("/form", formHandler)

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/form", formHandler)
	r.Handle("/", fileServer)
	r.Handle("/movies.html", fileServer)

	fmt.Printf("Starting server at port: 8080\n")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
