package clientfunc

import (
	"fmt"

	"github.com/gambruh/gophkeeper/internal/database"
)

// Cards
func (c *Client) SetCardCommand(input []string) {
	if c.AuthCookie == nil && !c.LoggedOffline {
		fmt.Println("please login first")
		return
	}
	if len(input) != 7 {
		printSetCardSyntax()
		return
	}
	var cardData database.Card
	for i, data := range input {
		switch i {
		case 0:
		case 1:
			cardData.Cardname = data
		case 2:
			cardData.Number = data
		case 3:
			cardData.Name = data
		case 4:
			cardData.Surname = data
		case 5:
			cardData.ValidTill = data
		case 6:
			cardData.Code = data
		}
	}

	err := c.saveCardInStorage(cardData)
	if err != nil {
		fmt.Println("error in client saving data to storage in SetCardCommand:", err)
	}
	fmt.Printf("Card %s saved to the storage!\n", cardData.Cardname)
}

func (c *Client) GetCardCommand(input []string) {
	if c.AuthCookie == nil && !c.LoggedOffline {
		fmt.Println("please login first")
		return
	}
	if len(input) != 2 {
		printGetCardSyntax()
		return
	}

	cardname := input[1]

	card, err := c.getCardFromStorage(cardname)
	if err != nil {
		if err == database.ErrDataNotFound {
			fmt.Println("No data in local storage")
			return
		} else {
			fmt.Println("error when trying to get card data from local storage:", err)
			return
		}
	}
	//result of the command, if no errors
	fmt.Printf("%+v\n", card)
}

func (c *Client) ListCardsCommand(input []string) {

	if c.AuthCookie == nil && !c.LoggedOffline {
		fmt.Println("please login first")
		return
	}
	if len(input) != 1 {
		printListCardsSyntax()
		return
	}

	cards, err := c.listCardsFromStorage()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Cards:")
	for _, card := range cards {
		fmt.Println("  ", card)
	}
	fmt.Println("Use getcard <cardname> to acquire card data")
}
