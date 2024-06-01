package main

import (
	"fmt"
	"os"
	"os/user"
	"repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Monkey Programming language.\n", user.Username)
	fmt.Printf("Type in a command\n")
	repl.Start(os.Stdin, os.Stdout)

}
