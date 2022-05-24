package main

type StateUpdate struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	Arena struct {
		Dimensions []int                  `json:"dims"`
		State      map[string]PlayerState `json:"state"`
	} `json:"arena"`
}

type PlayerState struct {
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Direction string `json:"direction"`
	WasHit    bool   `json:"wasHit"`
	Score     int    `json:"score"`
}

type lastState struct {
	lastAction	string
	lastTarget	string
	attacks	int
}

type geoInfo struct {
	direction	string
	targetX		int
	targetY		int
	myX		int
	myY		int
	boundaryX	int
	boundaryY	int
}