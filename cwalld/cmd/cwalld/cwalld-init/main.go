package main

import (
	"cwalld/internal/senv"
	"fmt"
	"os"
)

func main(){
	DIR := "/home/testgrounds/"

	err := senv.Setup(DIR)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
}
