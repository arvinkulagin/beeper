package main

import (
	"log"
	"fmt"
	"os"
	"bufio"
	"strings"
	"flag"
	"text/tabwriter"
	"github.com/arvinkulagin/beeper/api"
)

func main() {
	addr := flag.String("s", "localhost:8889", "Remote address")
	flag.Parse()
	c, err := api.NewClient(*addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for {
		command, args := getCommand(">> ")
		switch command {
		case "add":
			err := c.Add(args[0])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Printf("Add: %s\n", args[0])
		case "del":
			err := c.Del(args[0])
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Printf("Delete: %s\n", args[0])
		case "pub":
			err := c.Pub(args[0], strings.Join(args[1:], " "))
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Printf("Publish: %s\n", strings.Join(args[1:], " "))
		case "list":
			list, err := c.List()
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			fmt.Printf("Topics: %s\n", list)
		case "check":
			fmt.Println("There's nothing here yet")
		case "help":
			help()
		default:
			help()
		}
	}
}

func getCommand(greeting string) (string, []string) {
	fmt.Print(greeting)
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	input = strings.Trim(input, "\n")
	result := strings.Split(input, " ")
	return result[0], result[1:]
}

func help() {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Println("Commands:")
	fmt.Fprintln(w, "list\t - Prints list of active topics")
	fmt.Fprintln(w, "check\t - Prints topic status")
	fmt.Fprintln(w, "help\t - Prints help")
	fmt.Fprintln(w, "add <id>\t - Adds a topic with id")
	fmt.Fprintln(w, "del <id>\t - Deletes a topic with id")
	fmt.Fprintln(w, "pub <id> <message>\t - Publishes message to topic with id")
	err := w.Flush()
	if err != nil {
		log.Fatal(err)
	}
}