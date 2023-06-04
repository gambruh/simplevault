package clientfunc

import "fmt"

func printSetCardSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax: setcard <cardname> <cardnumber> <cardholder name> <cardholder surname> <card valid till date in format 'dd:mm:yyyy'> <cvv code>")
}

func printRegisterSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax: register <login> <password>")
}

func printLoginSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax: login <login> <password>")
}

func printListCardsSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax: listcards")
}

func printGetCardSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntax: getcard <cardname>")
}
