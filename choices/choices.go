package choices

import (
    "encoding/json"
    "log"
    "os"

    scryfall "github.com/heroku/scrivener/scryfall"
    "github.com/gomodule/redigo/redis"
)

type CardChoice struct {
    Number int
    Name string
}

func AddChoiceToDB(user string, cardList []scryfall.Card, msgID string) {
    log.Printf("Adding a choice to the db for user '%s'", user)

    // Connect to redis.
    db, err := redis.DialURL(os.Getenv("REDIS_URL"))
    if err != nil {
        log.Printf("Could not connect to redis: '%s'", err)
        return
    }
    defer db.Close()

    // Make our JSON string to store.
    choices := []CardChoice{}
    for i, card := range cardList {
        choice := CardChoice {
            Name: card.Name,
            Number: i + 1,
        }
        choices = append(choices, choice)
    }
    jsonString, _ := json.Marshal(choices)

    _, err = db.Do("HMSET", user, "choices", jsonString, "msgID", msgID)
    if err != nil {
        log.Printf("Could not add choice to db: '%s'", err)
    }
}

func RemoveChoiceFromDB(user string) {
    log.Printf("Adding a choice to the db for user '%s'", user)

    // Connect to redis.
    db, err := redis.DialURL(os.Getenv("REDIS_URL"))
    if err != nil {
        log.Printf("Could not connect to redis: '%s'", err)
        return
    }
    defer db.Close()

    _, err = db.Do("HDEL", user, "choices")
    if err != nil {
        log.Printf("Could not add choice to db: '%s'", err)
    }
}

func CheckChoice(user string, number int) (cardName string, msgID string) {
    cardName = ""
    msgID = ""

    // Connect to redis.
    db, err := redis.DialURL(os.Getenv("REDIS_URL"))
    if err != nil {
        log.Printf("Could not connect to redis: '%s'", err)
        return cardName, msgID
    }
    defer db.Close()

    // Try to find a result.
    choicesJson, err := redis.Bytes(db.Do("HGET", user, "choices"))
    if err != nil {
        log.Printf("Could not fetch choices from redis: '%s'", err)
        return cardName, msgID
    }

    cardList := []CardChoice {}
    err = json.Unmarshal(choicesJson, &cardList)

    index := number - 1
    if(index >= 0 && number <= len(cardList)){
        cardName = cardList[index].Name
    }

    // Fetch the ID of message asking for the choice to be made.
    msgID, err = redis.String(db.Do("HGET", user, "msgID"))
    if err != nil {
        log.Printf("Could not fetch message ID from redis: '%s'", err)
        return cardName, msgID
    }

    return cardName, msgID
}