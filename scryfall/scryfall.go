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
    Name string `json:"name"`
    Images ImageSet `json:"image_uris"`
}

func FuzzySearch(text string) (Card, error) {
    log.Printf("Fuzzy search requested: '%s'", text)
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

    log.Printf("Card name: %s", card.Name)

    return card, nil
}

