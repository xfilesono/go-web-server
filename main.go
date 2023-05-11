package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func formHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "Registration Succesful..\n")
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func main() {

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/form", formHandler)

	fmt.Printf("Starting server at port: 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
