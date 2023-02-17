package main
import (
	"net/http" 
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"encoding/json"
	"strings"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	// es, err := elasticsearch.NewDefaultClient()
	cfg := elasticsearch7.Config{
		Addresses: []string{
			"http://localhost:9200/",
		},
		Username: "elastic",
		Password: "vbhr3RY9t7bdWdeCjoEV",
	}
	es, err := elasticsearch7.NewClient(cfg)
  if err != nil {
    log.Fatalf("Error creating the client: %s", err)
  }

  // Print client and server version numbers.
  log.Printf("Client: %s", elasticsearch7.Version)
  log.Println(strings.Repeat("~", 37))


	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	router.GET("/elastic", ESHandler(es))
	http.HandleFunc("/", homeHandler)
	router.Run("localhost:8080")
	
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}
func ESHandler(esCleint *elasticsearch7.Client) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var (
			r  map[string]interface{}
		)
		// 1. Get cluster info
		//
		res, err := esCleint.Info()
		if err != nil {
			log.Fatalf("Error getting response: %s", err)
		}
		defer res.Body.Close()
		// Check response status
		if res.IsError() {
			log.Fatalf("Error: %s", res.String())
		}
		// Deserialize the response into a map.
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}
		log.Printf("Server: %s", r["version"].(map[string]interface{})["number"])
	}
	return gin.HandlerFunc(fn)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
			return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
			if a.ID == id {
					c.IndentedJSON(http.StatusOK, a)
					return
			}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
			errorHandler(w, r, http.StatusNotFound)
			return
	}
	fmt.Fprint(w, "welcome home")
}
func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
			fmt.Fprint(w, "Page not found 404")
	}
}