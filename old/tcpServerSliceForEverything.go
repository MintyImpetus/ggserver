package main

import (
        "bufio"
        "fmt"
        "net"
        "strings"
	"strconv"
)

type player struct {
	health int
	x int
	y int
}

type block struct {
	x int
	y int
	width int
	height int
}

var blockList []block
var playerList []player
var connList []net.Conn

func inRangeOfNumbers( query int, low int, high int) bool {
	var inRange bool
	for i := low; i <= high; i++ {
		if query == i {
			inRange = true
		}
	}
	return inRange
}

func handleConnections(connId int) {
	var message string
	playerId := len(playerList)
	playerList = append(playerList, player{health: 20, x: 0, y: 1})
        for {
		message = ""
                data, err := bufio.NewReader(connList[connId]).ReadString('\n')
                if err != nil {
                        fmt.Println(err)
                        return
                }
		dArray := strings.Split(data, " ")
		if strings.TrimSpace(string(dArray[0])) == "exit" {
                        fmt.Println("Exiting game server!")
                        return
                }
		if strings.TrimSpace(string(dArray[0])) == "echo" {
			message = strings.TrimSpace(string(dArray[1]))
                }
		if strings.TrimSpace(string(dArray[0])) == "health" {
			message = (strconv.Itoa(playerList[playerId].health))
                }
		if strings.TrimSpace(string(dArray[0])) == "sethealth" {
			var err error
			playerList[playerId].health, err = strconv.Atoi(strings.TrimSpace(string(dArray[1])))
			if err != nil {
				fmt.Println(err)
			}
		}
		if strings.TrimSpace(string(dArray[0])) == "move" {
		playerList[playerId].x, playerList[playerId].y = HandleMove(playerList[playerId].x, playerList[playerId].y, dArray)
		fmt.Println("PlayerX:", playerList[playerId].x, "PlayerY:", playerList[playerId].y)
		}

		message = message + "\n"
                connList[connId].Write([]byte(message))
        }
}

func HandleMove(PlayerX int, PlayerY int, dArray []string) (int, int) {
tmpX := PlayerX
tmpY := PlayerY
if strings.TrimSpace(string(dArray[1])) == "up" {
	tmpY = tmpY + 1
}
if strings.TrimSpace(string(dArray[1])) == "down" {
	tmpY = tmpY - 1
}
if strings.TrimSpace(string(dArray[1])) == "left" {
	tmpX = tmpX - 1
}
if strings.TrimSpace(string(dArray[1])) == "right" {
	tmpX = tmpX + 1
}
canMove := 1
for index, currentObject := range blockList {
//if ( tmpX == currentObject.x && tmpY == currentObject.y) {
if ( inRangeOfNumbers(tmpX, currentObject.x, currentObject.x + currentObject.width) && inRangeOfNumbers(tmpY, currentObject.y, currentObject.y + currentObject.height)) {
	fmt.Println("Can't move")
canMove = 0
}
fmt.Println("Index:", index, "currentObject:", currentObject)
}
if canMove != 1 {
	tmpX = PlayerX
	tmpY = PlayerY
}
return tmpX, tmpY
}

func main() {
	blockList = append(blockList, block{x: 4, y: 4, height: 2, width: 3})
        PORT := ":9876"
        dstream, err := net.Listen("tcp", PORT)
        if err != nil {
                fmt.Println(err)
                return
        }
        defer dstream.Close()
	for {
		conn, err := dstream.Accept()
		if err != nil {
			fmt.Println(err)
                continue
		} else {
			connId := len(connList)
			connList = append(connList, conn)
			
			go handleConnections(connId)
		}
	}
}
