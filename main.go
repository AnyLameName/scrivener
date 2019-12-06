package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strings"
    "syscall"

    discord "github.com/heroku/scrivener/discord"
    scryfall "github.com/heroku/scrivener/scryfall"
    slack "github.com/heroku/scrivener/slack"
    "github.com/bwmarrin/discordgo"
    "github.com/gin-gonic/gin"
    _ "github.com/heroku/x/hmetrics/onload"
)

func isLinkOnly(text string) (flag bool, newText string){
    LINK := "--link"
    flag = false
    newText = text
    if(strings.Index(text, LINK) != -1){
        log.Printf("Link-only requested")
        flag = true
        newText = strings.TrimSpace(strings.Replace(text, LINK, "", -1))
    }
    return flag, newText
}

func cardSearch(c *gin.Context) {
    text := c.PostForm("text")
    responseURL := c.PostForm("response_url")
    linkOnly, text := isLinkOnly(text)
    log.Printf("Search text: '%s', responding to: '%s'", text, responseURL)

    c.String(http.StatusOK, "")
    go doSearch(text, responseURL, linkOnly, false)
}


func walkerSearch(c *gin.Context) {
    text := c.PostForm("text")
    responseURL := c.PostForm("response_url")
    linkOnly, text := isLinkOnly(text)
    log.Printf("Planeswalker search text: '%s', responding to: '%s'", text, responseURL)

    c.String(http.StatusOK, "")
    go doSearch(text, responseURL, linkOnly, true)
}

func linkSearch(c *gin.Context) {
    text := c.PostForm("text")
    responseURL := c.PostForm("response_url")
    log.Printf("Link search text: '%s', responding to: '%s'", text, responseURL)

    c.String(http.StatusOK, "")
    go doSearch(text, responseURL, true, false)
}

func doSearch(text string, responseURL string, linkOnly bool, walkerOnly bool) {
    cardList := []scryfall.Card{}
    var err error

    if(walkerOnly) {
        cardList, err = scryfall.WalkerSearch(text)
    }else{
        cardList, err = scryfall.Search(text) 
    }

    if(err != nil){
        log.Printf("Error in Scryfall search: %s", err)
    }

    resp := slack.NewResponse("in_channel", fmt.Sprintf("No cards found matching: '%s'.", text))

    // Still here? Then we at least have results to process.
    numCards := len(cardList)
    if(numCards == 1){
        resp = slack.NewCard(cardList[0], linkOnly)
    }else if(numCards > 1){
        resp = slack.NewCardChoice(text, cardList, linkOnly)
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
        // Check for link switch here.
        linkOnly, answer := isLinkOnly(answer)
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
            card := slack.NewCard(cardList[0], linkOnly)
            c.JSON(http.StatusOK, card)
            return
        }else{
            log.Printf("Somehow, we still didn't get it down to one card.")
            resp := slack.NewCardChoice(answer, cardList, linkOnly)
            c.JSON(http.StatusOK, resp)
            return
        }
    }
    
    c.JSON(http.StatusOK, errorResponse)
}

func discordSearch(session *discordgo.Session, msg *discordgo.MessageCreate, text string, walkerOnly bool) {
    cardList := []scryfall.Card{}
    var err error

    if(walkerOnly) {
        cardList, err = scryfall.WalkerSearch(text)
    }else{
        cardList, err = scryfall.Search(text) 
    }

    if(err != nil){
        log.Printf("Error in Scryfall search: %s", err)
    }

    // Still here? Then we at least have results to process.
    numCards := len(cardList)
    embed := discordgo.MessageEmbed{}
    if(numCards == 1){
        embed = discord.EmbedCard(cardList[0])
    }else if(numCards > 1){
        embed = discord.EmbedChoice(cardList)
    }
    session.ChannelMessageSendEmbed(msg.ChannelID, &embed)
}

func ready(s *discordgo.Session, event *discordgo.Ready){
    s.UpdateStatus(0, "!card")
}

func messageCreate(session *discordgo.Session, msg *discordgo.MessageCreate){
    // Ignore our own messages.
    if(msg.Author.ID == session.State.User.ID){
        return
    }

    // Look for our trigger word.
    trigger := "!card"
    if(strings.HasPrefix(msg.Content, trigger)){
        log.Println("---")
        // Make sure there is at least a three-character search term and a space after the trigger.
        if(len(msg.Content) < (len(trigger) + 4)){
            session.ChannelMessageSend(msg.ChannelID, "Thanks for using Scrivener. I need to be given a search term of at least three characters.")
            return
        }

        searchText := msg.Content[6:len(msg.Content)]
        log.Printf("Search initiated: '%s'", searchText)
        discordSearch(session, msg, searchText, false)
    }
}

func main() {
    token := os.Getenv("DISCORD_TOKEN")
    if token == "" {
        log.Fatal("$DISCORD_TOKEN must be set or we can't register with discord.")
        return
    }


    discord, err := discordgo.New("Bot " + token)
    if(err != nil){
        log.Printf("Error creating discord session: %s", err)
    }

    discord.AddHandler(ready)
    discord.AddHandler(messageCreate)

    err = discord.Open()
    if(err != nil){
        log.Print(err)
    }

    log.Print("Scrivener Discord functionality online. Press ctrl-c to exit.")
    // Wait here until CTRL-C or other term signal is received.
    sc := make(chan os.Signal, 1)
    signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
    <-sc

    // Cleanly close down the Discord session.
    discord.Close()
}