/*
TODO:
Must refactor code as the newly implemented uuid identity system no longer supports looping through slices.

Complete the communication protocol and begin working on the client in pygame.

*/

package main

import (
        "bufio"
        "fmt"
        "net"
        "strings"
	"strconv"
	"math"
	"time"
	"github.com/google/uuid"
//	"encoding/json"
)

type item struct {
	name string
	itemType string
	damage int
	ddq int
	description string
}

type player struct {
	health int
	inventory []item
	x int
	y int
	renderDistance int
}

type block struct {
	blockType string
	flickerUp bool
	sinceLastFlicker int
	x int
	y int
	width int
	height int
}

type entity struct {
	name string
	id string
	hp int
	x int
	y int
}

blockList := make(map[block]string)
entityList := make(map[entity]string)

playerList := make(map[player]string)
blockList := make(map[net.Conn]string)

func inRangeOfNumbers( query int, low int, high int) bool {
	if (query >= low && query <= high) {
		return true
	} else {
		return false
	}
}

func getDifference(a int, b int) int {
	if a > b {
		return a - b
	} else {
		return b - a
	}
}

func genUUID() string {
	id := uuid.New()
	return id.String()
}

func checkIfBlock(x int, y int) bool {
	exists := false
	for _, currentObject := range blockList {
		if ( inRangeOfNumbers(x, currentObject.x, currentObject.x + currentObject.width) && inRangeOfNumbers(y, currentObject.y, currentObject.y + currentObject.height)) {
			exists = true
		}
	}
	return exists

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
	blocked := checkIfBlock(tmpX, tmpY)
	if blocked == true {
		fmt.Println("Blocked")
		tmpX = PlayerX
		tmpY = PlayerY
	}
	return tmpX, tmpY
}

func updateClient(playerId int) string {
	var tempMessage string
	message := "["
	for index, _ := range blockList {
		tempMessage = ""
		for x := blockList[index].x; x <= blockList[index].x + blockList[index].width; x++ {
			for y := blockList[index].y; y <= blockList[index].y + blockList[index].height; y++ {
				distanceX := getDifference(playerList[playerId].x, x)
				distanceY := getDifference(playerList[playerId].y, y)
				lineDistance := int(math.Sqrt(math.Pow(float64(distanceX), 2) + math.Pow(float64(distanceY), 2)))
				if lineDistance <= playerList[playerId].renderDistance {
					if tempMessage != "" {
						tempMessage = tempMessage + ", "
					}
					tempMessage = tempMessage + `{"x": ` + strconv.Itoa(x) + `, "y": ` + strconv.Itoa(y) + `}`
					// It seems like the only way to make this work is to either create a queue for blocks that are visible that need to be sent, !!!or to add the , to the previous part, and not add it if it is the first one! 
					//if y != blockList[index].y + blockList[index].height {
					//}
				}
			}
		}
		if tempMessage != "" {
			// Each block needs to be surrounded by curly brackets, within square, eg [{x 5 y 6}, {x 6 y 7}]
			message = message + ` { "blockType": "` + blockList[index].blockType  + `", "x": ` + strconv.Itoa(blockList[index].x) + `, "y": ` + strconv.Itoa(blockList[index].y) + `, "width": ` + strconv.Itoa(blockList[index].width) + `, "height": ` + strconv.Itoa(blockList[index].height) + `, ` + `"blocks": [ ` + tempMessage + ` ] }`
			if index != len(blockList) - 1 {
				message = message + ", "
			}
		}
	}
	message = message + " ]"
	return message
}

func handleConnections(connId int) {
	var message string
	playerId := genUUID()
	playerList[playerId] = player{health: 20, x: 0, y: 1, renderDistance: 3}
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
			message = updateClient(playerId)
                }

		message = message + "\n"
                connList[connId].Write([]byte(message))
        }
}

func gameLoop() {
	for {
		for index, _ := range blockList {
			if blockList[index].blockType == "flicker" {
				if blockList[index].sinceLastFlicker >= 5 {
				if blockList[index].flickerUp == true {
						blockList[index].y = blockList[index].y - 1
						blockList[index].flickerUp = false
					} else {
						blockList[index].y = blockList[index].y + 1
						blockList[index].flickerUp = true
					}
				} else {
					blockList[index].sinceLastFlicker = blockList[index].sinceLastFlicker + 1
				}
			}
		}
		time.Sleep(time.Second)
	}
}

func main() {
	blockList[genUUID()] = block{ blockType: "basic", x: 4, y: 4, height: 1, width: 1}
	blockList[genUUID()] = block{ blockType: "flicker", x: 6, y: 4, height: 1, width: 0}
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
			connId := genUUID()
			connList[connId] = conn
			
			go handleConnections(connId)
		}
	}
}
