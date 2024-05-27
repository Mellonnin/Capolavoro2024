package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

var db *sql.DB

type Pagina struct {
	Titolo string
	Body   string
}

func controlla(titolo string) int {
	var tito string
	row := db.QueryRow("SELECT titolo FROM pagine WHERE titolo = ?", titolo)
	errore_scan := row.Scan(&tito)
	if errore_scan == sql.ErrNoRows {
		return 1 // pagina non trovata
	}
	return 0 // trovata una pagina
}

//func aggiungiPagina(c *gin.Context) error {
//    titolo:= c.Param("titolo")
//	errore_controlla := controlla(pag.Titolo)
//	if errore_controlla == 1 {
//		result, errore_insert := db.Exec("INSERT INTO pagine (titolo, body) VALUES (?, ?)", pag.Titolo, pag.Body)
//		if errore_insert != nil {
//			return errore_insert
//		}
//		fmt.Println(result)
//	} else if errore_controlla == 0 {
//		result, errore_update := db.Exec("UPDATE pagine SET body = ? WHERE titolo = ?", pag.Body, pag.Titolo)
//		fmt.Println(result)
//		if errore_update != nil {
//			return errore_update
//		}
//	}
//	return nil
//}

//func (p *Pagina) salva (titolo string){
//	err := aggiungiPagina(Pagina{
//		Titolo: titolo,
//		Body:  p.Body,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("query fatta")
//}

//func loadPage(titolo string) (*Pagina, error) {
//	var pag Pagina
//	row := db.QueryRow("SELECT body FROM pagine WHERE titolo = ?", titolo)
//	err := row.Scan(&pag.Body)
//	if err != nil {
//		return nil, fmt.Errorf("pagina non trovata : %v", err)
//	}
//	return &Pagina{Titolo: titolo, Body: pag.Body}, nil
//}

//func newHandler(titolo string) {
//	errore_controlla := controlla(titolo)
//
//		p, err := loadPage(titolo)
//		if err != nil {
//			p = &Pagina{Titolo: titolo}
//		}
//    if errore_controlla == 1 {
//        postPagine(titolo)
//	} else if errore_controlla == 0 {
//        getPagine(titolo)
//    }
//}

func getPagine(c *gin.Context) {
	titolo := c.Param("titolo")
	if controlla(titolo) == 0 {
		var pag Pagina
		var row = db.QueryRow("SELECT titolo, body FROM pagine WHERE titolo = ?", titolo)
		row.Scan(&pag.Titolo, &pag.Body)
		c.IndentedJSON(http.StatusOK, pag)
	} else if controlla(titolo) == 1 {
		c.IndentedJSON(http.StatusOK, "la pagina non esiste")
	}
}

func postPagine(c *gin.Context) {
	//completa metti il content della pagina
	// modifica il content della pagina
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
	router := gin.Default()
	router.GET("/pagine/:titolo", getPagine)
	router.POST("/pagine", postPagine)
	router.Run("localhost:8090")
}
