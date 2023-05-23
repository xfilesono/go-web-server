package main

import (
	"encoding/json"
	"fmt"
	"go-web-server/models"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// gathering information from form.html for sign-up
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

	db, err := sql.Open("mysql", "xfilesono:62674819@(127.0.0.1:3306)/allmovies")
	if err = db.Ping(); err != nil {
		dbCheckErr(err)
	}
	defer db.Close()

	var (
		id        string
		year      string
		country   string
		title     string
		genre     string
		cast      string
		directors string
	)
	sorgu := "SELECT id, year, country, title, genre, cast, directors FROM movies WHERE id=?"
	erro, _ := db.Query("SELECT id FROM movies")
	var i = 0

	for erro.Next() {
		i++
		errs := db.QueryRow(sorgu, i).Scan(&id, &year, &country, &title, &genre, &cast, &directors)
		dbCheckErr(errs)
		// Terminal check
		fmt.Println(id, year, country, title, genre, cast, directors)
		movies = append(movies, models.Movies{ID: id, Year: year, Country: country, Title: title, Genre: genre, Cast: cast, Director: &models.Director{Directors: directors}})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)

	emptyMoviesVariable()

}

// This function is for prevent movies's variable's duplications since we are using mysql database
func emptyMoviesVariable() {
	length := len(movies)
	movies = append(movies[:0], movies[length:]...)
	fmt.Println(movies)
}

// function where we can delete a movie permanently
func deleteMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

// function where we get a single movie's information
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

// Function where we can create a new line of a movie
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie models.Movies
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(10000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

// Function where we can add a movie to the mysql database
func formAddMovieToDB(w http.ResponseWriter, r *http.Request) {

	db, err := sql.Open("mysql", "xfilesono:62674819@(127.0.0.1:3306)/allmovies")
	if err = db.Ping(); err != nil {
		dbCheckErr(err)
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	year := r.FormValue("year")
	country := r.FormValue("country")
	title := r.FormValue("title")
	date := time.Now()
	genre := r.FormValue("genre")
	cast := r.FormValue("cast")
	directors := r.FormValue("directors")

	result, err := db.Exec(`INSERT INTO movies (year, country, title, created_at, genre, cast, directors) VALUES (?,?,?,?,?,?,?)`, year, country, title, date, genre, cast, directors)
	dbCheckErr(err)
	addedId, err := result.LastInsertId()
	defer db.Close()

	fmt.Fprintf(w, "Movie Added to the DB..\n")
	fmt.Fprintf(w, "Movie's ID = %d\n", addedId)
	fmt.Fprintf(w, "Movie's Title = %s\n", title)
	fmt.Fprintf(w, "Released = %s\n", year)
	fmt.Fprintf(w, "Country = %s\n", country)
	fmt.Fprintf(w, "Movie's Genre = %s\n", genre)
	fmt.Fprintf(w, "Movie's Cast = %s\n", cast)
	fmt.Fprintf(w, "Director's Name = %s\n", directors)
}

// function where we can update a movie's information
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie models.Movies
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
		}
	}
}

// Database error check
func dbCheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var movies []models.Movies

func main() {

	//movies = append(movies, models.Movies{ID: id, Country: country, Year: year, Genre: genre, Cast: cast, Title: title, Director: &models.Director{Directors: directors}})

	r := mux.NewRouter()

	fileServer := http.FileServer(http.Dir("./static"))

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovies).Methods("DELETE")
	r.HandleFunc("/form", formHandler)
	r.HandleFunc("/addmovie", formAddMovieToDB)
	r.Handle("/form.html", fileServer)
	r.Handle("/", fileServer)
	r.Handle("/movies.html", fileServer)
	r.Handle("/addmovie.html", fileServer)

	fmt.Printf("Starting server at port: 8080\n")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
