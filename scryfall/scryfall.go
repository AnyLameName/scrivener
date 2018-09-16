package scryfall

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
)

type ImageSet struct {
    Large string `json:"large"`
}

type Card struct {
    Object string `json:"object"`
    Type string `json:"type"`
    Name string `json:"name"`
    Images ImageSet `json:"image_uris"`
    Faces []Card `json:"card_faces"`
}

type FullSearch struct {
    Object string `json:"object"`
    Size int `json:"total_cards"`
    HasMore bool `json:"has_more"`
    Cards []Card `json:"data"`
}

func fuzzy(text string) (Card, error) {
    log.Printf("Fuzzy search initiated.")
    card := Card{}

    req, err := http.NewRequest("GET", "https://api.scryfall.com/cards/named", nil)
    if err != nil {
        log.Print(err)
        return card, err
    }

    q := url.Values{}
    q.Add("fuzzy", text)

    req.URL.RawQuery = q.Encode()

    log.Printf("Scryfall url: %s", req.URL.String())

    client := &http.Client{}
    resp, err := client.Do(req)
    if(err != nil) {
        return card, err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    err = json.Unmarshal(body, &card)

    return card, nil
}

func full(text string) ([]Card, error) {
    log.Printf("Full search initiated.")
    cardList := []Card {}
    //https://api.scryfall.com/cards/search?order=cmc&q=chandra
    req, err := http.NewRequest("GET", "https://api.scryfall.com/cards/search", nil)
    if err != nil {
        log.Print(err)
        return cardList, err
    }

    q := url.Values{}
    q.Add("order", "rarity")
    q.Add("q", text)

    req.URL.RawQuery = q.Encode()

    log.Printf("Scryfall url: %s", req.URL.String())

    client := &http.Client{}
    resp, err := client.Do(req)
    if(err != nil) {
        return cardList, err
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    fullResults := FullSearch {}
    err = json.Unmarshal(body, &fullResults)

    // Now let's go through this list we just got.
    for i, card := range fullResults.Cards{
        log.Printf("Card %d: %s", i, card.Name)
        if(card.Faces != nil){
            log.Printf("\tThis card has %d faces.", len(card.Faces))
            continue
        }
        log.Printf("\tAdding to list.")
        cardList = append(cardList, card)
    }

    return cardList, nil
}

func Search(text string) ([]Card, error) {
    cardList := []Card{}
    card, err := fuzzy(text)
    if(err != nil){
        return cardList, err
    }

    if(card.Object == "card"){
        cardList = append(cardList, card)
    }

    if(card.Object == "error" && card.Type == "ambiguous"){
        cardList, err = full(text)
    }

    return cardList, nil
}

