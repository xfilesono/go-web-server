package models

type Mov struct {
	ID       string `json:"id"`
	Country  string `json:"country"`
	Year     string `json:"year"`
	Genre    string `json:"genre"`
	Cast     string `json:"cast"`
	Title    string `json:"title"`
	Director string `json:"director"`
}
