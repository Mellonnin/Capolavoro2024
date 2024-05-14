package main

import (
	"os"
	"github.com/gin-gonic/gin"
	"net/http"
    "fmt"
    "encoding/json"
    "io/ioutil"
)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}


func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func postAlbums(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
	    return
	}
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
    save(albums[len(albums)-1])
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func save(p album) {
	file, err := os.Create("people.json")
	if err != nil {
		fmt.Println("Errore durante la creazione del file:", err)
		return 
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(albums)
	if err != nil {
		fmt.Println("Errore durante la scrittura JSON:", err)
		return 
	}
	fmt.Println("Dati JSON scritti su file 'people.json'")
    return 
}

func read()  {
	// Leggere i dati JSON da un file
	fileData, err := ioutil.ReadFile("people.json")
	if err != nil {
		fmt.Println("Errore durante la lettura del file:", err)
		return
	}

	var peopleFromFile []Person
	err = json.Unmarshal(fileData, &peopleFromFile)
	if err != nil {
		fmt.Println("Errore durante il parsing JSON:", err)
		return
	}

	fmt.Println("Dati JSON letti dal file:")
	for _, p := range peopleFromFile {
		fmt.Printf("Nome: %s, Età: %d, Email: %s\n", p.Name, p.Age, p.Email)
	}
}
func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)
    router.GET("/albums/:id", getAlbumByID)
	router.Run("localhost:8090")
}
