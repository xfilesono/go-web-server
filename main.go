package main

import (
	"encoding/json"
	"fmt"
	"go-web-server/config"
	"go-web-server/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// gathering information from form.html for sign-up
func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	user.FullName = r.FormValue("name")
	user.Phone = r.FormValue("phone")
	user.Mail = r.FormValue("mail")
	user.Passwd, _ = config.HashPassword(r.FormValue("passwd"))
	date := time.Now()

	if user.FullName != "" {
		res, err := db.Exec(`INSERT INTO users (fullname, phone, created_at, mail, passwd) VALUES (?,?,?,?,?)`, user.FullName, user.Phone, date, user.Mail, user.Passwd)
		dbCheckErr(err)
		addedId, err := res.LastInsertId()
		if err != nil {
			log.Fatalln(err)
		}
		user.ID = strconv.Itoa(int(addedId))

		// add this user's info's to getUser slice for form.html check
		getUser := append(getUser, models.User{ID: user.ID, FullName: user.FullName, Phone: user.Phone, Mail: user.Mail, Passwd: user.Passwd})

		//Terminal Check
		fmt.Println(getUser)
	}

	// redirect to form.html
	// http.Redirect(w, r, "/form.html?"+usersInfo, http.StatusSeeOther)
	http.Redirect(w, r, "/form.html", http.StatusSeeOther)
}

func getSignUpUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

}

// function where we get all movies from our database
func getMovies(w http.ResponseWriter, r *http.Request) {

	var moviesCheck models.Mov
	var i = 0

	dbQuery := "SELECT id, year, country, title, genre, cast, directors FROM movies WHERE id=?"
	rows, err := db.Query("SELECT id FROM movies")
	dbCheckErr(err)
	defer rows.Close()

	for rows.Next() {
		i++
		if err := db.QueryRow(dbQuery, i).Scan(&moviesCheck.ID, &moviesCheck.Year, &moviesCheck.Country, &moviesCheck.Title, &moviesCheck.Genre, &moviesCheck.Cast, &moviesCheck.Director); err != nil {
			dbCheckErr(err)
		}
		movies = append(movies, models.Mov{ID: moviesCheck.ID, Year: moviesCheck.Year, Country: moviesCheck.Country, Title: moviesCheck.Title, Genre: moviesCheck.Genre, Cast: moviesCheck.Cast, Director: moviesCheck.Director})
	}
	// Terminal print
	for i := range movies {
		fmt.Println(movies[i])
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

// function where we get a single movie's information
func getMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	// get spesific Movie by using ID from DB
	var moviesCheck models.Mov
	dbQuery := "SELECT id, year, country, title, genre, cast, directors FROM movies WHERE id=?"

	if err := db.QueryRow(dbQuery, params["id"]).Scan(&moviesCheck.ID, &moviesCheck.Year, &moviesCheck.Country, &moviesCheck.Title, &moviesCheck.Genre, &moviesCheck.Cast, &moviesCheck.Director); err != nil {
		dbCheckErr(err)
	}
	movies = append(movies, models.Mov{ID: moviesCheck.ID, Year: moviesCheck.Year, Country: moviesCheck.Country, Title: moviesCheck.Title, Genre: moviesCheck.Genre, Cast: moviesCheck.Cast, Director: moviesCheck.Director})
	fmt.Println(movies)

	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

// Function where we can add a new a movie to DB
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie models.Mov
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		log.Fatalln(err)
	}
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)

	date := time.Now()
	res, err := db.Exec(`INSERT INTO movies (year, country, title, created_at, genre, cast, directors) VALUES (?,?,?,?,?,?,?)`, movie.Year, movie.Country, movie.Title, date, movie.Genre, movie.Cast, movie.Director)
	dbCheckErr(err)

	addedId, err := res.LastInsertId()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Fprintf(w, "Movie Added to the DB..\n")
	fmt.Fprintf(w, "Movie's ID = %d\n", addedId)
}

// Function where we can add a movie to DB via using html form
func formAddMovieToDB(w http.ResponseWriter, r *http.Request) {

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
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Fprintf(w, "Movie Added to the DB..\n")
	fmt.Fprintf(w, "Movie's ID = %d\n", addedId)
	fmt.Fprintf(w, "Movie's Title = %s\n", title)
	fmt.Fprintf(w, "Released = %s\n", year)
	fmt.Fprintf(w, "Country = %s\n", country)
	fmt.Fprintf(w, "Movie's Genre = %s\n", genre)
	fmt.Fprintf(w, "Movie's Cast = %s\n", cast)
	fmt.Fprintf(w, "Director's Name = %s\n", directors)
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

	_, err := db.Exec(`DELETE FROM movies WHERE id=?`, params["id"])
	dbCheckErr(err)
	_, err = db.Exec(`ALTER TABLE movies DROP id`)
	dbCheckErr(err)
	_, err = db.Exec(`ALTER TABLE movies ADD id int not null auto_increment primary key first`)
	dbCheckErr(err)

	json.NewEncoder(w).Encode(movies)
}

// function where we can update a movie's information
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	getMovie(w, r)

	var movie models.Mov

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
		}
	}

	date := time.Now()
	fmt.Println(movie)
	_, err := db.Exec(`UPDATE movies SET title = ?,year = ?,country = ?,genre = ?,cast = ?,directors = ?,updated_at = ? WHERE id=?`, movie.Title, movie.Year, movie.Country, movie.Genre, movie.Cast, movie.Director, date, params["id"])
	dbCheckErr(err)
	emptyMoviesVariable()
}

// Database error check
func dbCheckErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

var (
	movies  []models.Mov
	db      *sql.DB
	getUser []models.User
	user    models.User
)

func main() {

	config.Connect()
	db = config.GetDB()

	r := mux.NewRouter()

	fileServer := http.FileServer(http.Dir("./static"))

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/users", getSignUpUser).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovies).Methods("DELETE")
	r.HandleFunc("/form", formHandler)
	r.HandleFunc("/addmovie", formAddMovieToDB)
	r.Handle("/custom.js", fileServer)
	r.Handle("/custom2.js", fileServer)
	r.Handle("/form.html", fileServer)
	r.Handle("/", fileServer)
	r.Handle("/movies.html", fileServer)
	r.Handle("/addmovie.html", fileServer)

	fmt.Printf("Starting server at port: 8080\n")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
