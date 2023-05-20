package models

type Movies struct {
	ID       string    `json:"id"`
	Year     string    `json:"year"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}
