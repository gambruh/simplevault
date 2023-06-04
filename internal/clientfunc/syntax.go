package clientfunc

import "fmt"

func printSetCardSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntacsis: setcard <cardname> <cardnumber> <cardholder name> <cardholder surname> <card valid till date in format 'dd:mm:yyyy'> <cvv code>")
}

func printRegisterSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntacsis: register <login> <password>")
}

func printLoginSyntax() {
	fmt.Println("Wrong input!")
	fmt.Println("Right syntacsis: login <login> <password>")
}
