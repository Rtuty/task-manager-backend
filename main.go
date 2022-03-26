package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type User struct {
	Id                   int
	Name                 string
	Age                  uint16
	Money                int16
	Avg_grades, Happines float64
	Hobbies [] string
}

func (u *User) getAllInfo() string {
	return fmt.Sprintf("ID user: %d. User name is: %s. He is %d And he has money equal: %d", u.Id, u.Name, u.Age, u.Money)
}

func (u *User) setNewName(newName string) {
	u.Name = newName
}

func home_page(page http.ResponseWriter, reque *http.Request) {
	kirill := User{Id: 1, Name: "Kudryavcev Kirill", Age: 18, Money: 100, Avg_grades: 4.5, Happines: 10.0, Hobbies: [] string{"Football", "and skate"}},
	//kirill.setNewName("Kudryavcev Kirill")
	//fmt.Fprintf(page, "<b> Main text </b>")
	tmpl, _ := template.ParseFiles("templates/home_page.html")
	tmpl.Execute(page, kirill)
}

func contacts_page(page http.ResponseWriter, reque *http.Request) {
	fmt.Fprintf(page, "Работает ли кириллица?)")
}

func hadleRequest() { //обработка запроса
	http.HandleFunc("/", home_page)
	http.HandleFunc("/contacts/", contacts_page)
	http.ListenAndServe(":8080", nil)
}

func main() {

	hadleRequest()
}
