package game

import "github.com/google/uuid"

type TaskID string

func NewTaskID() TaskID {
	return TaskID(uuid.NewString())
}

type TaskType int

const (
	TaskTypeCollectItem TaskType = iota
	TaskTypeStayInZone
	TaskTypeSabotageOther
	// we’ll expand this as we design tasks
)

type Task struct {
	ID      TaskID
	Type    TaskType
	OwnerID PlayerID // who got this task
	// Later: time limit, parameters (zone id, target player, etc.)
}
