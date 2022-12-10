package main

import (
        "bufio"
        "fmt"
        "net"
        "os"
        "strings"
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

        for {
                reader := bufio.NewReader(os.Stdin)
                fmt.Print(">> ")
                text, _ := reader.ReadString('\n')
                fmt.Fprintf(c, text+"\n")

                message, _ := bufio.NewReader(c).ReadString('\n')
                fmt.Print("->: " + message)
                if strings.TrimSpace(string(text)) == "exit" {
			message, _ := bufio.NewReader(c).ReadString('\n')
			fmt.Print("->: " + message)
			if message == `[ { "exit": "successful" } ] []` { //Change this later to parse the json and find the result of the thing
                        	fmt.Println("TCP client exiting...")
				os.Exit(0)
			} else {
				fmt.Println("Error found.")
				os.Exit(1)
			}
                }
        }
}
