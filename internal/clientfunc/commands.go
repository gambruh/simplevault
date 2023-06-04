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

	err := c.sendCardToDB(cardData)
	if err != nil {
		fmt.Println("error in client sending data to database in SetCardCommand:", err)
	}
	err = c.saveCardInStorage(cardData)
	if err != nil {
		fmt.Println("error in client saving data to storage in SetCardCommand:", err)
	}
}

func (c *Client) GetCardCommand(input []string) {
	if c.AuthCookie == nil && !c.LoggedOffline {
		fmt.Println("please login first")
		return
	}
}

func (c *Client) ListCommands(input []string) {

}
