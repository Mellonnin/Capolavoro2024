package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	//"os"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"regexp"
)

var db *sql.DB

type Pagina struct {
	Title string
	Body  []byte
}

// FUNZIONA
func (pag *Pagina) salva() {
	var tito string
	row := db.QueryRow("SELECT titolo FROM pagine WHERE titolo = ?", pag.Title)
	errore_scan := row.Scan(&tito)
	if errore_scan != nil {
		db.Exec("INSERT INTO pagine (titolo, body) VALUES (?, ?)", pag.Title, pag.Body)
	} else {
		db.Exec("UPDATE pagine SET body = ? WHERE titolo = ?", pag.Body, pag.Title)
	}
}

func loadPage(title string) (*Pagina, error) {
	var pag Pagina
	row := db.QueryRow("SELECT body FROM pagine WHERE titolo = ?", title)
	err := row.Scan(&pag.Body)
	if err != nil {
		return nil, fmt.Errorf("pagina non trovata : %v", err)
	}
	return &Pagina{Title: title, Body: pag.Body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/nuovo/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "visualizza", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/nuovo/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "modifica", p)
}

func newHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err == nil {
		http.Redirect(w, r, "/modifica/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "nuovo", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Pagina{Title: title, Body: []byte(body)}
	p.salva()
	http.Redirect(w, r, "/visualizza/"+title, http.StatusFound)
}

var templates = template.Must(template.ParseFiles(
	"modifica.html",
	"visualizza.html",
	"nuovo.html",
))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Pagina) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(modifica|salva|visualizza|nuovo)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func connessioneDB() {
	cfg := mysql.Config{
		User:                 "mellonnin",
		Passwd:               "arch",
		Net:                  "tcp",
		Addr:                 "localhost:3306",
		DBName:               "wiki",
		AllowNativePasswords: true,
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("connesso!")
}

func main() {
	connessioneDB()
	http.HandleFunc("/visualizza/", makeHandler(viewHandler))
	http.HandleFunc("/modifica/", makeHandler(editHandler))
	http.HandleFunc("/nuovo/", makeHandler(newHandler))
	http.HandleFunc("/salva/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8070", nil))
}
