package main

import (
        "bufio"
        "fmt"
        "net"
        "strings"
	"strconv"
	"math"
	"time"
//	"encoding/json"
)

type player struct {
	health int
	x int
	y int
	renderDistance int
}

type block struct {
	class string
	flicker bool
	flickerUp bool
	sinceLastFlicker int
	x int
	y int
	width int
	height int
}

var blockList []block
var playerList []player
var connList []net.Conn

func inRangeOfNumbers( query int, low int, high int) bool {
	// rewrite this to use < and >?????!!!!
	var inRange bool
	for i := low; i <= high; i++ {
		if query == i {
			inRange = true
		}
	}
	return inRange
}

func getDifference(a int, b int) int {
	var c int
	if a > b {
		c = a - b
	} else {
		c = b - a
	}
	return c
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

func getInfoForPlayer(playerId int) string {
	var tempMessage string
	message := "["

	for index, currentObject := range blockList {
		tempMessage = ""
		for x := currentObject.x; x <= currentObject.x + currentObject.width; x++ {
			for y := currentObject.y; y <= currentObject.y + currentObject.height; y++ {
//				fmt.Println("Index:", strconv.Itoa(index), "Width", strconv.Itoa(i), "Height:", strconv.Itoa(j))
				distanceX := getDifference(playerList[playerId].x, x)
				distanceY := getDifference(playerList[playerId].y, y)
				lineDistance := int(math.Sqrt(math.Pow(float64(distanceX), 2) + math.Pow(float64(distanceY), 2)))
				if lineDistance <= playerList[playerId].renderDistance {
					// add here minimal stuff. Should just tell the client the x and y of the visible chunk by adding it to tempMessage
					tempMessage = tempMessage + `"x": ` + strconv.Itoa(x) + `, "y": ` + strconv.Itoa(y)
					if y != currentObject.y + currentObject.height {
						tempMessage = tempMessage + ", "
					}
				}
			}
		}
		if tempMessage != "" {
			// This should add all the information about the entire cube
			message = message + ` { "class": "` + currentObject.class  + `", "x": ` + strconv.Itoa(currentObject.x) + `, y: ` + strconv.Itoa(currentObject.y) + `, "width": ` + strconv.Itoa(currentObject.width) + `, "height": ` + strconv.Itoa(currentObject.height) + `, ` + `"blocks": { ` + tempMessage + `} }`
			if index != len(blockList) - 1 {
				message = message + ", "
			}
		}
		index = index
	//	fmt.Println(playerList[playerId].renderDistance)
	//	fmt.Println(distanceX)
	//	fmt.Println(distanceY)
	//	fmt.Println(lineDistance)
	//	fmt.Println(index)
	}
	message = message + "]"
	return message
}

func handleConnections(connId int) {
	var message string
	playerId := len(playerList)
	playerList = append(playerList, player{health: 20, x: 0, y: 1, renderDistance: 3})
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
		if strings.TrimSpace(string(dArray[0])) == "update" {
			message = getInfoForPlayer(playerId)
                }

		message = message + "\n"
                connList[connId].Write([]byte(message))
        }
}

func gameLoop() {
	for {
		for index, currentObject := range blockList {
//			fmt.Println(index)
			if currentObject.flicker == true {
				if currentObject.sinceLastFlicker >= 5 {
				if currentObject.flickerUp == true {
						blockList[index].y = blockList[index].y - 1
						blockList[index].flickerUp = false
//						fmt.Println("Moved down")
					} else {
						blockList[index].y = blockList[index].y + 1
						blockList[index].flickerUp = true
//						fmt.Println("Moved up")
					}
				} else {
					blockList[index].sinceLastFlicker = blockList[index].sinceLastFlicker + 1
				}
			}
		}
//		time.Sleep(time.Millisecond)
		time.Sleep(time.Second)
//		fmt.Println("passed")
	}
}

func main() {
	blockList = append(blockList, block{ class: "basic", x: 4, y: 4, height: 1, width: 1})
	blockList = append(blockList, block{ class: "flicker", x: 6, y: 4, height: 1, width: 0, flicker: true})
        PORT := ":9876"
        dstream, err := net.Listen("tcp", PORT)
        if err != nil {
                fmt.Println(err)
                return
        }
        defer dstream.Close()
	go gameLoop()
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
