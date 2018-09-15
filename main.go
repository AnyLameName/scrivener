package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    _ "github.com/heroku/x/hmetrics/onload"
)

func cardCallback(c *gin.Context) {
    cardName := c.Param("name");
    response := fmt.Sprintf("Requested: '%s'", cardName)
    c.String(http.StatusOK, response);
}

func main() {
    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }

    router := gin.New()
    router.Use(gin.Logger())

    router.GET("/", func(c *gin.Context) {
        c.String(http.StatusOK, "You should try making an actual request.")
    })

    router.GET("/card/:name", cardCallback);

    router.Run(":" + port)
}
