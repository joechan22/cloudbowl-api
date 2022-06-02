package main

import (
	"encoding/json"
	"fmt"
	"log"
	rand2 "math/rand"
	"net/http"
	"os"
	"reflect"

	// "github.com/jinzhu/copier"
)

var (
	throwCMD string = "T"
	forwardCMD string = "F"
	leftCMD string = "L"
	rightCMD string = "R"
	actions []string = []string{"L", "R", "F"}
	consecutive int = 50		// this number needed to change
	hitRange = 3
	totalAttacks = 0
)

var lastS = lastState{}
var gArena PlayerState

func main() {
	port := "8080"
	if v := os.Getenv("PORT"); v != "" {
		port = v
	}
	http.HandleFunc("/", handler)

	log.Printf("The server was started on port :%s", port)
	err := http.ListenAndServe(":"+port, nil)
	log.Fatalf("Error occured: %v", err)
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, "only support HTTP POST")
		return
	}

	var v StateUpdate

	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	// decoder.DisallowUnknownFields()	//all field must be declared under the type struct
	if err := decoder.Decode(&v); err != nil {
		log.Printf("Error: failed to decode JSON in response body: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp1 := decisionTree(v)
	fmt.Fprint(w, resp1)

}


func (l lastState) IsEmpty() bool {
	return reflect.DeepEqual(l,lastState{})
}

func decisionTree(arena StateUpdate) (response string) {
	if lastS.IsEmpty() {
		lastS = lastState{"-", "-", 0}
	}
	selfLink := arena.Links.Self.Href
	gArena = arena.Arena.State[selfLink]
	log.Println("last state ", lastS)
	// if chkRun(arena) && !chkArd(arena) {
	// 	log.Println("being hit, running ")
	// 	next := randMove(arena)
	// 	lastS.lastAction = next
	// 	return next
	// }

	target := canThrow(arena)
	if target != "" {
		if target != lastS.lastTarget {
			lastS = lastState{target, throwCMD, 1}
			return throwCMD
		}
		if lastS.attacks <= consecutive {
			totalAttacks =  lastS.attacks+1
			lastS = lastState{target, throwCMD, totalAttacks}
			return throwCMD
		}
	}

	action := getNearest(arena, 0)
	if action != "" {
		lastS.lastAction = action
		return action
	}

	//for fear that the data is changed by getNearest()
	arena.Arena.State[selfLink] = gArena

	action = randMove(arena)
	lastS.lastAction = action
	return action
}

// func (ds StateUpdate) copy() *StateUpdate {
// 	return &ds
// }

func getNearest(data StateUpdate, depth int) (action string) {
	selfLink := data.Links.Self.Href
	states := data.Arena.State
	myInfo := states[selfLink]

	var copySet StateUpdate
	copySet = data
	// copier.CopyWithOption(&copySet, &data, copier.Option{IgnoreEmpty: true, DeepCopy: true}) // deep copy with copier

	alterWay := ""

	//for moving forward
	myDirection := myInfo.Direction
	myX := myInfo.X
	myY := myInfo.Y

	switch myDirection {
	case "N":
		if myY > 0 {
			myY -= 1
		}
	case "E":
		if myX < data.Arena.Dimensions[1]-1 {
			myX += 1
		}
	case "S":
		if myY < data.Arena.Dimensions[0]-1 {
			myY += 1
		}
	case "W":
		if myX > 0 {
			myX -= 1
		}
	}
	for k, v := range states {
		if k == selfLink || (k == lastS.lastTarget && lastS.attacks >= consecutive) {
			continue
		}
		if canHit(geoInfo{myInfo.Direction, v.X, v.Y, myInfo.X, myInfo.Y, data.Arena.Dimensions[0], data.Arena.Dimensions[1]}) {
			return forwardCMD
		}
	}

	//TODO furtherAnalysis()
	if depth == 0 && alterWay == "" {

		if entry, ok := copySet.Arena.State[selfLink]; ok {
			entry.Direction = myDirection
			entry.X = myX
			entry.Y = myY
		 
			copySet.Arena.State[selfLink] = entry
		}
		alternative := getNearest(copySet, 1)
		if alternative != "" {
			alterWay = forwardCMD
		}

	}

	// for moving right 
	switch myDirection {
	case "N":
		myDirection = "E"
	case "E":
		myDirection = "S"
	case "S":
		myDirection = "W"
	case "W":
		myDirection = "N"
	}
	for k, v := range states {
		if k == selfLink || (k == lastS.lastTarget && lastS.attacks >= consecutive) {
			continue
		}
		if canHit(geoInfo{myInfo.Direction, v.X, v.Y, myInfo.X, myInfo.Y, data.Arena.Dimensions[0], data.Arena.Dimensions[1]}) {
			return rightCMD
		}
	}

	//TODO furtherAnalysis()
	if depth == 0 && alterWay == "" {

		if entry, ok := copySet.Arena.State[selfLink]; ok {
			entry.Direction = myDirection
			entry.X = myX
			entry.Y = myY
		 
			copySet.Arena.State[selfLink] = entry
		}
		alternative := getNearest(copySet, 1)
		if alternative != "" {
			alterWay = rightCMD
		}

	}

	// for moving left
	switch myDirection {
	case "N":
		myDirection = "W"
	case "E":
		myDirection = "N"
	case "S":
		myDirection = "E"
	case "W":
		myDirection = "S"
	}
	for k, v := range states {
		if k == selfLink || (k == lastS.lastTarget && lastS.attacks >= consecutive) {
			continue
		}
		if canHit(geoInfo{myInfo.Direction, v.X, v.Y, myInfo.X, myInfo.Y, data.Arena.Dimensions[0], data.Arena.Dimensions[1]}) {
			return leftCMD
		}
	}

	//TODO furtherAnalysis()
	if depth == 0 && alterWay == "" {

		if entry, ok := copySet.Arena.State[selfLink]; ok {
			entry.Direction = myDirection
			entry.X = myX
			entry.Y = myY
		 
			copySet.Arena.State[selfLink] = entry
		}
		alternative := getNearest(copySet, 1)
		if alternative != "" {
			alterWay = leftCMD
		}

	}


	return alterWay
}

func randMove(data StateUpdate) (action string) {

	selfLink := data.Links.Self.Href
	states := data.Arena.State
	myInfo := states[selfLink]

	rand := rand2.Intn(3)
	nextAction := actions[rand]
	// prevent position looping
	if (lastS.lastAction == leftCMD && nextAction == rightCMD) || (lastS.lastAction == rightCMD && nextAction == leftCMD) {
		nextAction = forwardCMD
	}

	switch myInfo.Direction {
	case "N":
		if nextAction == forwardCMD {
			if myInfo.Y == 0 {
				nextAction = rightCMD
			} else {
				for _, v := range states {
					if v.X == myInfo.X && v.Y == myInfo.Y - 1 {
						nextAction = rightCMD
						break
					}
				}
			}
		}
		if nextAction == rightCMD && myInfo.X == data.Arena.Dimensions[0]-1 {
			nextAction = leftCMD
		} else if nextAction == leftCMD && myInfo.X == 0 {
			nextAction = rightCMD
		}
	case "E":
		if nextAction == forwardCMD {
			if myInfo.X == data.Arena.Dimensions[0]-1 {
				nextAction = rightCMD
			} else {
				for _, v := range states {
					if v.X == myInfo.X+1 && v.Y == myInfo.Y {
						nextAction = rightCMD
						break
					}
				}
			}
		}
		if nextAction == rightCMD && myInfo.Y == data.Arena.Dimensions[1]-1 {
			nextAction = leftCMD
		} else if nextAction == leftCMD && myInfo.Y == 0 {
			nextAction = rightCMD
		}

	case "S":
		if nextAction == forwardCMD {
			if myInfo.Y == data.Arena.Dimensions[1]-1 {
				nextAction = rightCMD
			} else {
				for _, v := range states {
					if v.X == myInfo.X && v.Y == myInfo.Y+1 {
						nextAction = rightCMD
						break
					}
				}
			}
		}
		if nextAction == rightCMD && myInfo.X == 0 {
			nextAction = leftCMD
		} else if nextAction == leftCMD && myInfo.X == data.Arena.Dimensions[0]-1  {
			nextAction = rightCMD
		}
	case "W":
		if nextAction == forwardCMD {
			if myInfo.X == 0 {
				nextAction = rightCMD
			} else {
				for _, v := range states {
					if v.X == myInfo.X-1 && v.Y == myInfo.Y {
						nextAction = rightCMD
						break
					}
				}
			}
		}
		if nextAction == rightCMD && myInfo.Y == 0 {
			nextAction = leftCMD
		} else if nextAction == leftCMD && myInfo.Y == data.Arena.Dimensions[1]-1  {
			nextAction = rightCMD
		}

	}

	return nextAction
}

func canThrow(data StateUpdate) (url string) {
	selfLink := data.Links.Self.Href
	states := data.Arena.State
	myInfo := states[selfLink]
	for k, v := range states {
		if k == selfLink {
			continue
		}
		if canHit(geoInfo{myInfo.Direction, v.X, v.Y, myInfo.X, myInfo.Y, data.Arena.Dimensions[0], data.Arena.Dimensions[1]}) {
			return k
		}
		
	}
	return ""
}

func canHit(geo geoInfo) bool {
	left := 0
	right := geo.boundaryX - 1
	top := 0
	bottom := geo.boundaryY - 1
	if geo.myX - hitRange >= 0 {
		left = geo.myX - hitRange
	}
	if geo.myX + hitRange < geo.boundaryX {
		right = geo.myX + hitRange
	}
	if geo.myY - hitRange >= 0 {
		top = geo.myY - hitRange
	}
	if geo.myY + hitRange < geo.boundaryY {
		bottom = geo.myY + hitRange
	}
	switch geo.direction {
	case "N":
		if geo.targetX == geo.myX && geo.targetY >= top && geo.targetY <= geo.myY {
			return true
		}
	case "S":
		if geo.targetX == geo.myX && geo.targetY >= geo.myY && geo.targetY <= bottom {
			return true
		}
	case "W":
		if geo.targetY == geo.myY && geo.targetX >= left && geo.targetX <= geo.myX {
			return true
		}
	case "E":
		if geo.targetY == geo.myY && geo.targetX >= geo.myX && geo.targetX <= right {
			return true
		}

	}
	return false
}

//TODO need to solve if it is being trapped
func chkRun(data StateUpdate) bool {
	selfLink := data.Links.Self.Href
	states := data.Arena.State
	myInfo := states[selfLink]
	if myInfo.WasHit {
		return true
	}
	return false
}

//TODO check around if there is any player
func chkArd(data StateUpdate) bool {
	selfLink := data.Links.Self.Href
	states := data.Arena.State
	myInfo := states[selfLink]
	myX := myInfo.X
	myY := myInfo.Y
	cnt := 0
	players := 0	// counting the number of players
	for k, v := range states {
		players += 1
		if k == selfLink {
			continue
		}
		if (myX == v.X && myY+1 == v.Y) || (myX == v.X && myY-1 == v.Y) || (myX-1 == v.X && myY == v.Y) || (myX+1 == v.X && myY == v.Y) {
			cnt += 1
		}
		if myX == 0 && myY == 0 && cnt == 2 {
			return true
		} else if myX == data.Arena.Dimensions[0] && myY == data.Arena.Dimensions[1] && cnt == 2 {
			return true
		} else if myX == 0 && myY != 0 && cnt == 3 {
			return true
		} else if myX != 0 && myY == data.Arena.Dimensions[1] && cnt == 3 {
			return true
		} else if cnt == 4 {
			return true
		}

	}
	if players <= 2 {
		return true
	}
	return false
}