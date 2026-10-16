package models

const HistoryLimit = 20

type HistoryEntry struct {
	ID         int64
	Expression string
	Result     string
}

type State struct {
	Display         string
	Stored          *string
	Operator        *string
	AwaitingOperand bool
	JustCalculated  bool
	Expression      string
	OperandLabel    *string
	Error           *string
	History         []HistoryEntry
	NextID          int64
}

// Action is a transport-independent calculator command.
type Action struct {
	Type      string
	Value     string
	HasValue  bool
	HistoryID int64
}

func InitialState() State {
	return State{
		Display: "0",
		History: []HistoryEntry{},
		NextID:  1,
	}
}

type Calculation struct {
	Operation string
	Left      string
	Right     *string
	Base      *string
}
