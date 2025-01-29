package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album
	fmt.Println("newAlbum :", newAlbum)

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}
	fmt.Println("albums :", albums)
	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

var collection *mongo.Collection
var ctx = context.TODO()

func main() {
	// mongodb://serveradm:mypassword@localhost:27017/learngolang
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/learngolang")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected OK")
	var testdataFind = client.Database("learngolang").Collection("test")
	var podcast bson.M
	if err = testdataFind.FindOne(ctx, bson.M{}).Decode(&podcast); err != nil {
		log.Fatal(err)
	}
	fmt.Println("podcast: ", podcast)

	err = client.Disconnect(ctx)

	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}
