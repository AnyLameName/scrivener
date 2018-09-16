package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"

    scryfall "github.com/heroku/scrivener/scryfall"
    slack "github.com/heroku/scrivener/slack"
    "github.com/gin-gonic/gin"
    _ "github.com/heroku/x/hmetrics/onload"
)

type Action struct {
    Name       string `json:"name"`
    Text       string `json:"text"`
    ActionType string `json:"type"`
    Value      string `json:"value"`
}

type Attachment struct {
    Text           string   `json:"text"`
    Fallback       string   `json:"fallback"`
    CallbackID     string   `json:"callback_id"`
    Color          string   `json:"color"`
    AttachmentType string   `json:"attachment_type"`
    Actions        []Action `json:"actions"`
}

type SlackResponse struct {
    Attachments []Attachment `json:"attachments"`
    ResponseType string `json:"response_type"`
    Text string `json:"text"`
}

func cardSearch(c *gin.Context) {
    text := c.PostForm("text")
    log.Printf("Search text: %s", text)
    cardList, err := scryfall.Search(text)
    if(err != nil){
        log.Printf("Error in Scryfall search: %s", err)
        c.String(http.StatusOK, "Sorry, something went wrong.")
        return
    }

    // Still here? Then we found _something_. Let's see what we should do with it.
    if(len(cardList) == 1){
        card := slack.NewCard(cardList[0])
        c.JSON(http.StatusOK, card)
    }else{
        resp := slack.NewCardChoice(text, cardList)
        c.JSON(http.StatusOK, resp)
    }
}

func buttonCallback(c *gin.Context) {
    errorResponse := SlackResponse {
        ResponseType: "ephemeral",
        Text: "Something went wrong processing your response.",
    }

    body, err := c.GetRawData()
    if(err != nil){
        log.Printf("Couldn't read request from callback payload.")
        c.JSON(http.StatusOK, errorResponse)
        return
    }
    log.Printf("Body parsed: %s", body)

    callback := slack.Callback{}
    err = json.Unmarshal(body, &callback)

    if(len(callback.Actions) == 1){
        answer := callback.Actions[0].Value
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
    router.POST("/button/", buttonCallback)

    router.Run(":" + port)
}
