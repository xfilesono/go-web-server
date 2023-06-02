package models

type User struct {
	ID       string `json:"id"`
	FullName string `json:"fullname"`
	Phone    string `json:"phone"`
	Mail     string `json:"mail"`
	Passwd   string `json:"passwd"`
}
