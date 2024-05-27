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

func controlla(titolo string) int {
	var tito string
	row := db.QueryRow("SELECT title FROM pagine WHERE titolo = ?", titolo)
	errore_scan := row.Scan(&tito)
	if errore_scan == nil {
        return 1 // pagina non trovata
	}
    return 0 // trovata una pagina
}

// FUNZIONA
func aggiungiPagina(pag Pagina) error {
	errore_controlla := controlla(pag.Title)
	if errore_controlla == 1 {
		result, errore_insert := db.Exec("INSERT INTO pagine (titolo, body) VALUES (?, ?)", pag.Title, pag.Body)
		if errore_insert != nil {
			return errore_insert
		}
		fmt.Println(result)
	} else if errore_controlla == 0 {
		result, errore_update := db.Exec("UPDATE pagine SET body = ? WHERE titolo = ?", pag.Body, pag.Title)
		fmt.Println(result)
		if errore_update != nil {
			return errore_update
		}
	}
	return nil
}

func (p *Pagina) salva(titolo string) {
	err := aggiungiPagina(Pagina{
		Title: titolo,
		Body:  p.Body,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("query fatta")

}

func loadPage(titolo string) (*Pagina, error) {
	var pag Pagina
	row := db.QueryRow("SELECT body FROM pagine WHERE titolo = ?", titolo)
	err := row.Scan(&pag.Body)
	if err != nil {
		return nil, fmt.Errorf("pagina non trovata : %v", err)
	}
	return &Pagina{Title: titolo, Body: pag.Body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, titolo string) {
	p, err := loadPage(titolo)
	if err != nil {
		http.Redirect(w, r, "/nuovo/"+titolo, http.StatusFound)
		return
	}
	renderTemplate(w, "visualizza", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, titolo string) {
	p, err := loadPage(titolo)
	if err != nil {
		p = &Pagina{Title: titolo}
	}
	renderTemplate(w, "modifica", p)
}

func newHandler(w http.ResponseWriter, r *http.Request, titolo string) {
	errore_controlla := controlla(titolo)

		p, err := loadPage(titolo)
		if err != nil {
			p = &Pagina{Title: titolo}
		}
    if errore_controlla == 1 {
		renderTemplate(w, "nuova_pagina", p)
	} else if errore_controlla == 0 {
	    viewHandler(w,r,titolo)
    }
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Pagina{Title: title, Body: []byte(body)}
	p.salva(title)
	http.Redirect(w, r, "/visualizza/"+title, http.StatusFound)
}

var templates = template.Must(template.ParseFiles(
	"modifica.html",
	"visualizza.html",
	"nuova_pagina.html",
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
