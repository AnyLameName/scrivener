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

func cardCallback(c *gin.Context) {
    urlBase := "https://api.scryfall.com/cards/named?fuzzy=%s"
    cardName := c.Param("name");
    url := fmt.Sprintf(urlBase, cardName)
    resp, err := http.Get(url)
    if(err != nil) {
        // Sure, sure. Error handling. Of course.
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    var card Card
    err = json.Unmarshal(body, &card)

    log.Printf("Card name: %s", card.Name)

    c.String(http.StatusOK, fmt.Sprintf("Card Name: %s", card.Name));
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
    c.String(http.StatusOK, fmt.Sprintf("Oh hi there, %s.", user))
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
    router.POST("/card/", slackCallback)

    router.Run(":" + port)
}
