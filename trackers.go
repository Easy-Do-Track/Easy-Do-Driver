package main

const (
	KeyHead       = "head"
	KeyChest      = "chest"
	KeyElbowLeft  = "elbow_left"
	KeyElbowRight = "elbow_right"
	KeyWrestLeft  = "wrest_left"
	KeyWrestRight = "wrest_right"
	KeyKneeLeft   = "knee_left"
	KeyKneeRight  = "knee_right"
	KeyAnkleLeft  = "ankle_left"
	KeyAnkleRight = "ankle_right"
)

type PosVector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}
