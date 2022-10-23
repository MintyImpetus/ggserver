package main

import (
        "bufio"
        "fmt"
        "net"
        "os"
        "strings"
	"github.com/mattn/go-tty"
	"log"
)

func main() {
        arguments := os.Args
        if len(arguments) == 1 {
                fmt.Println("Please provide host:port.")
                return
        }

        CONNECT := arguments[1]
        c, err := net.Dial("tcp", CONNECT)
        if err != nil {
                fmt.Println(err)
                return
        }

	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

        for {
		text := ""
		r, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		if string(r) == "w" {
			text = "move up"
		}
		if string(r) == "a" {
			text = "move left"
		}
		if string(r) == "d" {
			text = "move right"
		}
		if string(r) == "s" {
			text = "move down"
		}
		if string(r) == "e" {
			text = "exit"
			tty.Close()
		}
                fmt.Fprintf(c, string(text)+"\n")

                message, _ := bufio.NewReader(c).ReadString('\n')
                fmt.Print("->: " + message)
                if strings.TrimSpace(string(r)) == "e" {
                        fmt.Println("TCP client exiting...")
                        return
                }
        }
}
