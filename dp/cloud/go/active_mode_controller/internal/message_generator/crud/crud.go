package crud

type ActionType uint8

const (
	Delete ActionType = iota
)

type Action struct {
	Type         ActionType
	SerialNumber string
}
