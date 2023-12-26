package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/hoisie/mustache"
	"github.com/minpeter/rctf-backend/api"
	"github.com/minpeter/rctf-backend/database"
)

type ClientConfig struct {
	Meta            Meta              `json:"meta"`
	HomeContent     string            `json:"homeContent"`
	Sponsors        []interface{}     `json:"sponsors"`
	GlobalSiteTag   string            `json:"globalSiteTag"`
	CtfName         string            `json:"ctfName"`
	Divisions       map[string]string `json:"divisions"`
	DefaultDivision string            `json:"defaultDivision"`
	Origin          string            `json:"origin"`
	StartTime       int64             `json:"startTime"`
	EndTime         int64             `json:"endTime"`
	EmailEnabled    bool              `json:"emailEnabled"`
	UserMembers     bool              `json:"userMembers"`
	FaviconUrl      string            `json:"faviconUrl"`
	Github          struct {
		ClientId string `json:"clientId"`
	}
}

type Meta struct {
	Description string `json:"description"`
	ImageURL    string `json:"imageUrl"`
}

var clientConfig ClientConfig

func serveIndex(c *gin.Context) {

	jsonData, err := json.Marshal(clientConfig)
	if err != nil {
		log.Fatalf("Error marshalling clientConfig: %v", err)
	}

	rendered := struct {
		JSONConfig string
		Config     ClientConfig
	}{
		JSONConfig: string(jsonData),
		Config:     clientConfig,
	}

	// Use mustache to render the index.html template
	html := mustache.RenderFile("client/dist/index.html", rendered)
	c.Writer.WriteString(html)
}

func loadClientConfig() {
	configFile, err := os.Open("client-config.json")
	if err != nil {
		fmt.Printf("Error opening client-config.json: %v\n", err)
		return
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&clientConfig)

	clientConfig.Github.ClientId = os.Getenv("GITHUB_CLIENT_ID")

	if err != nil {
		fmt.Printf("Error decoding client-config.json: %v\n", err)
		return
	}
}

func main() {
	if err := database.ConnectDatabase(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	loadClientConfig()

	app := api.NewRouter()

	app.GET("/", serveIndex)

	app.Static("/assets", "client/dist/assets")

	app.NoRoute(serveIndex)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	host := ":" + port
	// Removes the “accept incoming network connections?” pop-up on macOS.
	if runtime.GOOS == "darwin" {
		host = "localhost:" + port
	}

	if err := app.Run(host); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
