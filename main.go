package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "os"

    scryfall "github.com/heroku/scrivener/scryfall"
    "github.com/gin-gonic/gin"
    _ "github.com/heroku/x/hmetrics/onload"
)

type Card struct {
    Name string `json:"name"`
}

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

func fuzzy(text string) string {
    req, err := http.NewRequest("GET", "https://api.scryfall.com/cards/named", nil)
    if err != nil {
        log.Print(err)
        os.Exit(1)
    }

    q := url.Values{}
    q.Add("fuzzy", text)

    req.URL.RawQuery = q.Encode()

    log.Printf("Scryfall url: %s", req.URL.String())

    client := &http.Client{}
    resp, err := client.Do(req)
    if(err != nil) {
        // Sure, sure. Error handling. Of course.
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    card := Card{}
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

    /*
    text := c.PostForm("text")
    log.Printf("Text: %s", text)

    base := "Were you looking for '%s'?"
    cardName := fuzzy(text)
    resp := fmt.Sprintf(base, cardName)
    */
    scryfall.FuzzySearch("dog")

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

    c.JSON(http.StatusOK, slack)
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

    router.GET("/card/:text", cardCallback)
    router.POST("/card/", slackCallback)
    router.POST("/button/", buttonCallback)

    router.Run(":" + port)
}
