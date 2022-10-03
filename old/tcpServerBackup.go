package main

import (
        "bufio"
        "fmt"
        "net"
        "strings"
)

func main() {
        PORT := ":9876"
        dstream, err := net.Listen("tcp", PORT)
        if err != nil {
                fmt.Println(err)
                return
        }
        defer dstream.Close()

        conn, err := dstream.Accept()
        if err != nil {
                fmt.Println(err)
                return
        }

        for {
		var message string
                data, err := bufio.NewReader(conn).ReadString('\n')
                if err != nil {
                        fmt.Println(err)
                        return
                }
		dArray := strings.Split(data, " ")
                	if strings.TrimSpace(string(dArray[0])) == "exit" {
                        fmt.Println("Exiting TCP server!")
                        return
                }
		if strings.TrimSpace(string(dArray[0])) == "echo" {
			message = strings.TrimSpace(string(dArray[1]))
                }

		message = message + "\n"
                conn.Write([]byte(message))
        }
}
