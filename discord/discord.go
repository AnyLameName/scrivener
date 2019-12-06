package discord

import (
	"fmt"
	"log"

    "github.com/bwmarrin/discordgo"
    scryfall "github.com/heroku/scrivener/scryfall"
)

func EmbedCard(card scryfall.Card) discordgo.MessageEmbed {
	image := discordgo.MessageEmbedImage {
		URL: card.Images.Normal,
		ProxyURL: card.Images.Normal,
	}

	ret := discordgo.MessageEmbed {
		Title: card.Name,
		Color: getColor(card.ColorIdentity),
		Image: &image,
		URL: card.Link,
	}

	return ret
}

func EmbedChoice(cardList []scryfall.Card) discordgo.MessageEmbed {
	prompt := "Please choose from the list below. (Type the number)"
	for index, card := range cardList {
		prompt = prompt + fmt.Sprintf("\n**[%d]** - %s", index + 1, card.Name)
	}

	ret := discordgo.MessageEmbed {
		Title: "Multiple Matches Found",
		Description: prompt,
	}

	return ret
}

func getColor(identity []string) int {
	log.Printf("Determining color for: '%v'", identity)
	if(len(identity) > 1){
		return 0x997300
	} else if (len(identity) == 0){
		return 0xaaaaaa
	}

	switch(identity[0]){
	case "R":
		return 0xD3202A
	case "G":
		return 0x00733E
	case "U":
		return 0x0E68AB
	case "W":
		return 0xF9FAF4
	case "B":
		return 0x000000
	default:
		return 0x150B00
	}
}