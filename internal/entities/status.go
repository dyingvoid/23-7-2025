package entities

type Status uint8

const (
	StatusPending Status = iota
	StatusArchived
)

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusArchived:
		return "archived"
	default:
		panic("unhandled default case")
	}
}
