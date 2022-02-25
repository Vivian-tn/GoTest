package main

import (
	"fmt"
	"strings"
)

func main() {
	dbButton := "1#2#3##4"
	buttonParts := strings.Split(dbButton, ";")
	fmt.Println(len(buttonParts))
	//var chosedMember {}
	//fmt.Println(chosedMember)
}
