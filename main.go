package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1{
		fmt.Println("no parametrs!")
		return 
	}	
	url := args[0]
	fmt.Println(url);
}
