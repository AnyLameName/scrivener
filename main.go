package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    scryfall "github.com/heroku/scrivener/scryfall"
    slack "github.com/heroku/scrivener/slack"
    "github.com/gin-gonic/gin"
    _ "github.com/heroku/x/hmetrics/onload"
)

func cardSearch(c *gin.Context) {
    text := c.PostForm("text")
    responseURL := c.PostForm("response_url")
    log.Printf("Search text: '%s', responding to: '%s'", text, responseURL)

    ack := slack.NewResponse("in_channel", "Searching...")

    c.JSON(http.StatusOK, ack)
    go doSearch(text, responseURL)
}

func doSearch(text string, responseURL string) {
    cardList, err := scryfall.Search(text)
    if(err != nil){
        log.Printf("Error in Scryfall search: %s", err)
    }

    resp := slack.NewResponse("in_channel", fmt.Sprintf("No cards found matching: '%s'.", text))

    // Still here? Then we at least have results to process.
    numCards := len(cardList)
    if(numCards == 1){
        resp = slack.NewCard(cardList[0])
    }else if(numCards > 1){
        resp = slack.NewCardChoice(text, cardList)
    }

    if(responseURL == "log"){
        logString, err := json.Marshal(resp)
        if(err == nil){
            log.Printf("Reponse we would have sent: %s", logString)
        }else{
            log.Printf("Couldn't marshal response.")
        }
    }else{
        log.Printf("Responding to : '%s", responseURL)
        slack.Respond(resp, responseURL)
    }
}

func slackCallback(c *gin.Context) {
    errorResponse := slack.NewResponse("ephemeral", "Something went wrong with your search.")

    payload := c.PostForm("payload")
    log.Printf("Payload: %v", payload)

    callback := slack.Callback{}
    json.Unmarshal([]byte(payload), &callback)

    if(len(callback.Actions) == 1){
        var answer string
        log.Printf("Slack callback, id: '%s'.", callback.ID)
        if(callback.ID == "cardButton") {
            answer = callback.Actions[0].Value
        } else if(callback.ID == "cardMenu"){
            answer = callback.Actions[0].Selected[0].Value
        } else {
            c.JSON(http.StatusOK, errorResponse)
            return
        }
        log.Printf("They chose '%s'", answer)

        cardList, err := scryfall.Search(answer)
        if(err != nil){
            log.Printf("Error in Scryfall search: %s", err)
            c.JSON(http.StatusOK, errorResponse)
            return
        }

        // Still here? Then we found _something_. Let's see what we should do with it.
        if(len(cardList) == 1){
            log.Printf("And finally we have a card to show!")
            card := slack.NewCard(cardList[0])
            c.JSON(http.StatusOK, card)
            return
        }else{
            log.Printf("Somehow, we still didn't get it down to one card.")
            resp := slack.NewCardChoice(answer, cardList)
            c.JSON(http.StatusOK, resp)
            return
        }
    }
    
    c.JSON(http.StatusOK, errorResponse)
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

    router.POST("/card/", cardSearch)
    router.POST("/button/", slackCallback)

    router.Run(":" + port)
}
