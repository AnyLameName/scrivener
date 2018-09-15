package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    _ "github.com/heroku/x/hmetrics/onload"
)

type Card struct {
    Name string `json:"name"`
}

func fuzzy(text string) string {
    urlBase := "https://api.scryfall.com/cards/named?fuzzy=%s"

    url := fmt.Sprintf(urlBase, text)
    resp, err := http.Get(url)
    if(err != nil) {
        // Sure, sure. Error handling. Of course.
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    var card Card
    err = json.Unmarshal(body, &card)

    log.Printf("Card name: %s", card.Name)

    return card.Name
}

func cardCallback(c *gin.Context) {
    text := c.Param("text");

    result := fuzzy(text)

    c.String(http.StatusOK, fmt.Sprintf("Card Name: %s", result));
}

func slackCallback(c *gin.Context) {
    /*
    token=gIkuvaNzQIHg97ATvDxqgjtO
    &team_id=T0001
    &team_domain=example
    &enterprise_id=E0001
    &enterprise_name=Globular%20Construct%20Inc
    &channel_id=C2147483705
    &channel_name=test
    &user_id=U2147483697
    &user_name=Steve
    &command=/weather
    &text=94070
    &response_url=https://hooks.slack.com/commands/1234/5678
    &trigger_id=13345224609.738474920.8088930838d88f008e0
    */
    user := c.PostForm("user_name")
    text := c.PostForm("text")
    base := "Alright, %s, let's see what I can find for '%s'... did mean '%s'?"

    cardName := fuzzy(text)
    resp := fmt.Sprintf(base, user, text, cardName)

    c.String(http.StatusOK, resp)
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

    router.GET("/card/:text", cardCallback);
    router.POST("/card/", slackCallback)

    router.Run(":" + port)
}
