/*

TODO:

Complete the communication protocol and begin working on the client in pygame.
	Make it so the message variable is repeatedly updated to include what a player within render distance does, and every time a player updates, send them the response to their previous action, all the actions of visible players and then all visible block data.



Whenever action is done, run a goroutine function that adds the action to all nearby players list of actions. Have a goroutine running that constantly lowers the tick-till-end number for each of the actions and destroys them when done. 

Whenever the client asks for an update, send them all the actions and how long 'till they end.

At some point in the future, I could attempt to make all major tasks run through goroutines, as it may speed up the program.

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
	visibleActions []action
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

type action struct {
	name string
	duration int
}

var blockList = make(map[string]block)
var entityList = make(map[string]entity)
var playerList = make(map[string]player)
var connList = make(map[string]net.Conn)

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
	fmt.Println(PlayerX, PlayerY)
	if strings.TrimSpace(string(dArray[1])) == "up" {
		fmt.Println("up")
		tmpY = tmpY + 1
	}
	if strings.TrimSpace(string(dArray[1])) == "down" {
		fmt.Println("down")
		tmpY = tmpY - 1
	}
	if strings.TrimSpace(string(dArray[1])) == "left" {
		fmt.Println("left")
		tmpX = tmpX - 1
	}
	if strings.TrimSpace(string(dArray[1])) == "right" {
		fmt.Println("right")
		tmpX = tmpX + 1
	}
	blocked := checkIfBlock(tmpX, tmpY)
	if blocked == true {
		fmt.Println("Blocked")
		tmpX = PlayerX
		tmpY = PlayerY
	}
	fmt.Println(PlayerX, PlayerY)
	return tmpX, tmpY
}

func displayAnimation(name string, duration int) {
	for key, currentPlayer := range playerList {
		if getObjectDistance <= playerList[key] {
			currentPlayer.visibleActions = append(playerList[key].visibleActions, action{ name: name, duration: duration })			
			playerList[key] = currentPlayer
		}
	}
}

func handleActions(connId string, dArray []string) string {

	response := "[ { "
	currentPlayer := playerList[connId]

	if strings.TrimSpace(string(dArray[0])) == "exit" {
        	fmt.Println("Exiting game server!")
        }
	if strings.TrimSpace(string(dArray[0])) == "echo" {
		response = response + " { " + strings.TrimSpace(string(dArray[1])) 
        }
	if strings.TrimSpace(string(dArray[0])) == "health" {
		response = (strconv.Itoa(currentPlayer.health))
	}
	if strings.TrimSpace(string(dArray[0])) == "sethealth" {
		var err error
		currentPlayer.health, err = strconv.Atoi(strings.TrimSpace(string(dArray[1])))
		if err != nil {
			fmt.Println(err)
		}
	}
	if strings.TrimSpace(string(dArray[0])) == "animate" {
		displayAnimation( connId, "exampleAnimation", 40)	
	}
	if strings.TrimSpace(string(dArray[0])) == "move" {
		currentPlayer.x, currentPlayer.y = HandleMove(currentPlayer.x, currentPlayer.y, dArray)
		playerList[connId] = currentPlayer
		fmt.Println("PlayerX:", currentPlayer.x, "PlayerY:", currentPlayer.y)
		fmt.Println("listX:", playerList[connId].x, "listY", playerList[connId].y)

		response = response + `"x": ` + strconv.Itoa(currentPlayer.x) + `, "y": ` + strconv.Itoa(currentPlayer.y)
	}

	response = response + " } ] "

	return response
}

func getActions(playerId string) string {
return ""
}

func updateClient(playerId string) string {
	var tempMessage string
	message := "["
	i := 0
	for key, currentBlock := range blockList {
		tempMessage = ""
		for x := currentBlock.x; x <= currentBlock.x + currentBlock.width; x++ {
			for y := currentBlock.y; y <= currentBlock.y + currentBlock.height; y++ {
				//Make this into a function.
				distanceX := getDifference(playerList[playerId].x, x)
				distanceY := getDifference(playerList[playerId].y, y)
				lineDistance := int(math.Sqrt(math.Pow(float64(distanceX), 2) + math.Pow(float64(distanceY), 2)))
				if lineDistance <= playerList[playerId].renderDistance {
					if tempMessage != "" {
						tempMessage = tempMessage + ", "
					}
					tempMessage = tempMessage + `{"x": ` + strconv.Itoa(x) + `, "y": ` + strconv.Itoa(y) + `}`
				}
			}
		}
		if tempMessage != "" {
			message = message + ` { "id": "` + key + ` ", "blockType": "` + currentBlock.blockType  + `", "x": ` + strconv.Itoa(currentBlock.x) + `, "y": ` + strconv.Itoa(currentBlock.y) + `, "width": ` + strconv.Itoa(currentBlock.width) + `, "height": ` + strconv.Itoa(currentBlock.height) + `, ` + `"blocks": [ ` + tempMessage + ` ] }`
			if i != len(blockList) - 1 {
				message = message + ", "
			}
		}
		i = i + 1
	}
	message = message + "]"
	return message
}


func handleConnections(connId string) {

	var message string

	var response string
	var renderUpdates string

	playerList[connId] = player{health: 20, x: 0, y: 1, renderDistance: 3}
	
//Send opening information to the player.

        for {
		//Get what the player wants to do and then send a response.
		message = ""

		response = ""
		renderUpdates = ""

                data, err := bufio.NewReader(connList[connId]).ReadString('\n')
                if err != nil {
                        fmt.Println(err)
                        return
                }
		dArray := strings.Split(data, " ")
		response = handleActions(connId, dArray)

		actions := getActions(connId)

		renderUpdates = updateClient(connId)
		
		message = response + actions + renderUpdates + "\n"
                connList[connId].Write([]byte(message))
        }
}

func gameLoop() {
	for {
		for key, currentBlock := range blockList {
			if currentBlock.blockType == "flicker" {
				if currentBlock.sinceLastFlicker >= 5 {
				if currentBlock.flickerUp == true {
						currentBlock.y = currentBlock.y - 1
						currentBlock.flickerUp = false
					} else {
						currentBlock.y = currentBlock.y + 1
						currentBlock.flickerUp = true
					}
				} else {
					currentBlock.sinceLastFlicker = currentBlock.sinceLastFlicker + 1
				}
			}
			blockList[key] = currentBlock
		}
		time.Sleep(time.Second)
	}
}

func main() {
	blockList[genUUID()] = block{ blockType: "basic", x: 4, y: 4, height: 0, width: 0}
//	blockList[genUUID()] = block{ blockType: "flicker", x: 6, y: 4, height: 1, width: 0}
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
