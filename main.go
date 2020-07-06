package main

import (
	"./models"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
	"strings"
)

var db *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = ""
	password = ""
	dbname   = "wallester"
)

type Env struct {
	db models.Datastore
}

const (
	minAge = 18
	maxAge = 60
	limit = 2
)

type PageData struct {
    PageTitle 	string
	Search 		string
	PrevPage 	int
	NextPage 	int
	Customers 	[]*models.Customer
	OrderBy 	string
	Page		int
	Order		string
}

func inTimeSpan(start, end, check time.Time) bool {
    return check.After(start) && check.Before(end)
}

func isBirthdayValid(birthDay string) bool {
	now := time.Now()
	from := now.AddDate(-maxAge, 0, 0)
	to := now.AddDate(-minAge, 0, 0)
	date, _ := time.Parse("2006-01-02", birthDay)
	isValid := inTimeSpan(from, to, date)
	return isValid
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := models.InitDB(psqlInfo)
	if err != nil {
		panic(err)
	}

	env := &Env{db}

	http.HandleFunc("/", env.CustomersList)
	http.HandleFunc("/add", env.AddCustomer)
	http.HandleFunc("/edit", env.EditCustomer)
	http.HandleFunc("/search", env.SearchCustomer)
	http.HandleFunc("/delete", env.DeleteCustomer)

	http.ListenAndServe(":8080", nil)
}

func (env *Env) CustomersList(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	customers, err := env.db.GetCustomers()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	files := []string{
		"templates/customersList.gohtml",
		"templates/base.gohtml",
		"templates/nav.gohtml",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = tmpl.Execute(w, customers)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (env *Env) AddCustomer(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/add" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case "GET":
		files := []string{
			"templates/addCustomer.gohtml",
			"templates/base.gohtml",
			"templates/nav.gohtml",
		}

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	case "POST":
		if err := r.ParseForm(); err != nil {
			log.Println(w, "ParseForm() err: %v", err)
			return
		}
		customer := models.Customer{}
		customer.FirstName = r.FormValue("firstName")
		customer.LastName = r.FormValue("lastName")
		days := r.FormValue("days")
		months := r.FormValue("months")
		years := r.FormValue("years")
		birthDate := fmt.Sprintf("%s-%s-%s", years, months, days)

		if isBirthdayValid(birthDate) != true {
			msg := fmt.Sprintf("The user should be older than %[1]d and younger than %[2]d!", minAge, maxAge)
			log.Println(msg)
			http.Error(w, msg, 500)
			return
		}
		customer.BirthDate = birthDate
		customer.Gender = r.FormValue("gender")
		customer.Email = r.FormValue("email")
		customer.Address = r.FormValue("address")

		err := env.db.AddCustomer(customer)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Such user is already exists.", 500)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (env *Env) EditCustomer(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/edit" {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))

	switch r.Method {
	case "GET":		
		if err != nil {
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}

		customer, err := env.db.GetCustomerById(id)
		if err != nil {
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}

		files := []string{
			"templates/editCustomer.gohtml",
			"templates/base.gohtml",
			"templates/nav.gohtml",
		}

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}

		err = tmpl.Execute(w, customer)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	case "POST":
		if err := r.ParseForm(); err != nil {
			log.Println(w, "ParseForm() err: %v", err)
			return
		}
		customer := models.Customer{}
		customer.ID, _ = strconv.Atoi(r.FormValue("id"))
		customer.FirstName = r.FormValue("firstName")
		customer.LastName = r.FormValue("lastName")
		birthDate := strings.Split(r.FormValue("birthDate"), ".")
		day := birthDate[0]
		month := birthDate[1]
		year := birthDate[2]
		customer.BirthDate = fmt.Sprintf("%s-%s-%s", year, month, day)
		if isBirthdayValid(customer.BirthDate) != true {
			msg := fmt.Sprintf("The user should be older than %[1]d and younger than %[2]d!", minAge, maxAge)
			log.Println(msg)
			http.Error(w, msg, 500)
			return
		}
		customer.Gender = r.FormValue("gender")
		customer.Email = r.FormValue("email")
		customer.Address = r.FormValue("address")

		err := env.db.UpdateCustomer(customer)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Such user is already exists.", 500)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)		
	}

}

func (env *Env) SearchCustomer(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/search" {
		http.NotFound(w, r)
		return
	}

	var prevPage, nextPage int

	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}

	nextPage = page + 1
	prevPage = page - 1

	skips := limit * (page - 1)

	search := r.FormValue("q")

	orderBy := "id"	
	orderByParam := r.FormValue("orderBy")
	if orderByParam != ""  {
		orderBy = orderByParam
	}

	order := "ASC"
	orderParam := r.FormValue("order")
	if orderParam != ""  {
		if orderParam == "ASC" {
			order = "DESC"
		} else {
			order = "ASC"
		}
	}

	customers, err := env.db.GetCustomersByName(search, orderBy, order, limit, skips)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	data := PageData{
		PageTitle: "Search results",
		Customers: customers,
		OrderBy: orderBy,
		PrevPage: prevPage,
		NextPage: nextPage,
		Search: search,
		Page: page,
		Order: order,
	}

	files := []string{
		"templates/searchCustomers.gohtml",
		"templates/base.gohtml",
		"templates/nav.gohtml",
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (env *Env) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/delete" {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	err = env.db.DeleteCustomerById(id)
	if err != nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
