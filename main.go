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

	var (
		moviesCheck models.Mov
		movies      []models.Mov
	)

	dbQuery, err := db.Query(`SELECT id, year, country, title, genre, cast, directors FROM movies`)
	dbCheckErr(err)
	defer dbQuery.Close()
	for dbQuery.Next() {
		err := dbQuery.Scan(&moviesCheck.ID, &moviesCheck.Year, &moviesCheck.Country, &moviesCheck.Title, &moviesCheck.Genre, &moviesCheck.Cast, &moviesCheck.Director)
		dbCheckErr(err)
		movies = append(movies, moviesCheck)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)

}

// function where we get a single movie's information
func getMovie(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	var (
		moviesCheck models.Mov
	)

	dbQuery, err := db.Query(`SELECT id, year, country, title, genre, cast, directors FROM movies where id=?`, params["id"])
	dbCheckErr(err)
	defer dbQuery.Close()
	if dbQuery.Next() {
		err := dbQuery.Scan(&moviesCheck.ID, &moviesCheck.Year, &moviesCheck.Country, &moviesCheck.Title, &moviesCheck.Genre, &moviesCheck.Cast, &moviesCheck.Director)
		dbCheckErr(err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(moviesCheck)
		fmt.Fprintf(w, strconv.Itoa(http.StatusOK))
	} else {
		fmt.Fprintf(w, strconv.Itoa(http.StatusNoContent))
	}

}

// Function where we can add a new a movie to DB
func createMovie(w http.ResponseWriter, r *http.Request) {
	var (
		movie models.Mov
	)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		log.Fatalln(err)
	}

	if movie.Title != "" {
		date := time.Now()
		res, err := db.Exec(`INSERT INTO movies (year, country, title, created_at, genre, cast, directors) VALUES (?,?,?,?,?,?,?)`, movie.Year, movie.Country, movie.Title, date, movie.Genre, movie.Cast, movie.Director)
		dbCheckErr(err)

		addedId, err := res.LastInsertId()
		if err != nil {
			log.Fatalln(err)
		}
		id := int(addedId)
		movie.ID = strconv.Itoa(id)

		// this two lines has to comment for testing - check later
		json.NewEncoder(w).Encode(movie)
		fmt.Fprintf(w, "Movie Added to the DB..\n")

		fmt.Fprintf(w, strconv.Itoa(http.StatusCreated))
		w.WriteHeader(http.StatusCreated)
	} else {
		fmt.Fprintf(w, strconv.Itoa(http.StatusNoContent))
		w.WriteHeader(http.StatusNotFound)
	}
}

// Function where we can add a movie to DB via using html form (Merge this function with createMovie func later)
func formAddMovieToDB(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		fmt.Println("ParseForm() err: ", err)
		return
	}
	var (
		movie  models.Mov
		movies []models.Mov
	)

	if r.FormValue("year") != "" {
		movie.Year = r.FormValue("year")
		movie.Country = r.FormValue("country")
		movie.Title = r.FormValue("title")
		date := time.Now()
		movie.Genre = r.FormValue("genre")
		movie.Cast = r.FormValue("cast")
		movie.Director = r.FormValue("directors")
		movies = append(movies, movie)
		_, err := db.Exec(`INSERT INTO movies (year, country, title, created_at, genre, cast, directors) VALUES (?,?,?,?,?,?,?)`, movie.Year, movie.Country, movie.Title, date, movie.Genre, movie.Cast, movie.Director)
		if err != nil {
			dbCheckErr(err)
		}
		http.Redirect(w, r, "/addmovie.html", http.StatusOK)
	} else {
		fmt.Fprintf(w, strconv.Itoa(http.StatusNoContent))
	}

}

// function where we can delete a movie permanently
func deleteMovies(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	res, err := db.Query(`SELECT id from movies where id=?`, params["id"])
	dbCheckErr(err)

	if res.Next() {
		_, err = db.Exec(`DELETE FROM movies WHERE id=?`, params["id"])
		dbCheckErr(err)
		fmt.Fprintf(w, "\nMovie which has an id `"+params["id"]+"` has been deleted\n")
		fmt.Fprintf(w, strconv.Itoa(http.StatusOK))
	} else {
		fmt.Fprintf(w, strconv.Itoa(http.StatusNotFound))
	}

}

// function where we can update a movie's information
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var (
		movie models.Mov
	)

	res, err := db.Query(`SELECT id from movies where id=?`, params["id"])
	dbCheckErr(err)

	if res.Next() {
		err := json.NewDecoder(r.Body).Decode(&movie)
		if err != nil {
			log.Fatalln(err)
		}
		json.NewEncoder(w).Encode(movie)

		date := time.Now()
		_, err = db.Exec(`UPDATE movies SET title = ?,year = ?,country = ?,genre = ?,cast = ?,directors = ?,updated_at = ? WHERE id=?`, movie.Title, movie.Year, movie.Country, movie.Genre, movie.Cast, movie.Director, date, params["id"])
		dbCheckErr(err)
		fmt.Fprintf(w, strconv.Itoa(http.StatusOK))
	} else {
		fmt.Fprintf(w, strconv.Itoa(http.StatusNotFound))
	}

}

// Database error check
func dbCheckErr(err error) {
	if err != nil {
		log.Fatalln(http.StatusServiceUnavailable) // 503
	}
}

var (
	db      *sql.DB
	getUser []models.User
	user    models.User
)

func main() {

	var err error
	db, err = config.GetDB()
	if err != nil {
		log.Fatalln(http.StatusNotFound)
	}

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

	fmt.Printf("Starting http server at port: 8080\nUse localhost:8080 or 127.0.0.1:8080 in your browser\n")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
