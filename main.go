package main

import (
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

    /*
    action1 := Action {
        Name: "card",
        Text: "Chandra, the Cool One",
        ActionType: "button",
        Value: "cardOne",
    }

    action2 := Action {
        Name: "card",
        Text: "Chandra, the Lame One",
        ActionType: "button",
        Value: "cardTwo",
    }

    attach := Attachment {
        Text: "Please choose one.",
        Fallback: "Please choose one.",
        CallbackID: "choose_card",
        Color: "#acacac",
        AttachmentType: "default",
        Actions: []Action {},
    }

    attach.Actions = append(attach.Actions, action1, action2)

    slack := SlackResponse {
        ResponseType: "ephemeral", //in_channel or ephemeral
        Text: "I found multiple results for 'Chandra'.",
        Attachments: []Attachment {},
    }

    slack.Attachments = append(slack.Attachments, attach)
    */

func cardSearch(c *gin.Context) {
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

    /*
    base := "Were you looking for '%s'?"
    cardName := fuzzy(text)
    resp := fmt.Sprintf(base, cardName)
    */
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
        card := slack.NewCard(cardList[0].Name, cardList[0].Images.Large)
        c.JSON(http.StatusOK, card)
    }else{
        resp := slack.NewCardChoice(text)
        c.JSON(http.StatusOK, resp)
    }
}

func buttonCallback(c *gin.Context) {
    resp := SlackResponse {
        ResponseType: "ephemeral",
        Text: "Thanks for clicking that button!",
    }

    c.JSON(http.StatusOK, resp)
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
