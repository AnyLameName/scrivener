package main

import (
    "log"
    "os"
    "os/signal"
    "regexp"
    "strconv"
    "strings"
    "syscall"

    choices "github.com/heroku/scrivener/choices"
    discord "github.com/heroku/scrivener/discord"
    scryfall "github.com/heroku/scrivener/scryfall"
    "github.com/bwmarrin/discordgo"
    _ "github.com/heroku/x/hmetrics/onload"
)

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
    if(numCards == 1){
        discord.RespondWithCard(cardList[0], session, msg.ChannelID)
    }else if(numCards > 1){
        embed := discord.EmbedChoice(cardList)
        log.Printf("Choice Message: %v", embed)
        embeddedMsg, err := session.ChannelMessageSendEmbed(msg.ChannelID, &embed)
        if err != nil {
            log.Fatal("Error trying to send embedded message: '%s'", err)
        }
        choices.AddChoiceToDB(msg.Author.ID, cardList, embeddedMsg.ID)
    }
}

func ready(s *discordgo.Session, event *discordgo.Ready){
    s.UpdateStatus(0, "!card")
}

func messageCreate(session *discordgo.Session, msg *discordgo.MessageCreate){
    // Ignore our own messages.
    if(msg.Author.ID == session.State.User.ID){
        return
    }

    // Look for our trigger word to start a search.
    trigger := "!card"
    if(strings.HasPrefix(msg.Content, trigger)){
        log.Println("---")
        // Make sure there is at least a three-character search term and a space after the trigger.
        if(len(msg.Content) < (len(trigger) + 4)){
            session.ChannelMessageSend(msg.ChannelID, "Thanks for using Scrivener. I need to be given a search term of at least three characters.")
            return
        }

        searchText := msg.Content[6:len(msg.Content)]
        log.Printf("Search initiated by '%s': '%s'", msg.Author.Username, searchText)
        discordSearch(session, msg, searchText, false)
    }

    // Look for mid-sentence searches, using [[card]] for the trigger.
    reg := regexp.MustCompile(`(\[\[.*?\]\])|(\{\{.*?\}\})`)
    for _, match := range reg.FindAllString(msg.Content, -1) {
        searchText := match[2:len(match) - 2]
        walkerOnly := match[0] == '{'
        log.Printf("Search initiated by '%s': '%s'", msg.Author.Username, searchText)
        discordSearch(session, msg, searchText, walkerOnly)
    }

    // Look for a line containing only a number, which is how we let people narrow down searches.
    pattern := `^\d+$`
    matched, err := regexp.MatchString(pattern, msg.Content)
    if(err != nil){
        log.Printf("Error attempting regular expression match: '%s'", err)
    } else if(matched) {
        log.Println("---")
        number, _ := strconv.Atoi(msg.Content)
        log.Printf("Found a potential selection from '%s' aka '%s': %d", msg.Author.ID, msg.Author.Username, number)

        // Look up their selection, regardless.
        cardName, msgToDelete := choices.CheckChoice(msg.Author.ID, number)

        // Check for cancel.
        if(number == 0){
            choices.RemoveChoiceFromDB(msg.Author.ID)
            // Delete the original message that prompted the choice.
            session.ChannelMessageDelete(msg.ChannelID, msgToDelete)
            // Delete the message containing the choice, aka the current message.
            err = session.ChannelMessageDelete(msg.ChannelID, msg.ID)
            if err != nil {
                log.Printf("Could not delete message containing selection: '%s'", err)
            }
        }

        // Deal with an actual choice.
        if(cardName != ""){
            log.Printf("Found a choice: '%s'", cardName)
            discordSearch(session, msg, cardName, false)
            choices.RemoveChoiceFromDB(msg.Author.ID)
            if(msgToDelete != ""){
                // Delete the original message that prompted the choice.
                session.ChannelMessageDelete(msg.ChannelID, msgToDelete)
                // Delete the message containing the choice, aka the current message.
                err = session.ChannelMessageDelete(msg.ChannelID, msg.ID)
                if err != nil {
                    log.Printf("Could not delete message containing selection: '%s'", err)
                }
            }
        }
    }
}

func main() {
    // Connect to discord.
    token := os.Getenv("DISCORD_TOKEN")
    if token == "" {
        log.Fatal("$DISCORD_TOKEN must be set or we can't register with discord.")
        return
    }

    discord, err := discordgo.New("Bot " + token)
    if(err != nil){
        log.Printf("Error creating discord session: %s", err)
    }

    // Set up our handlers.
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