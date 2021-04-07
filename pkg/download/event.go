package download

type EventKey string

const (
	EventKeyStart    = "start"
	EventKeyPause    = "pause"
	EventKeyContinue = "continue"
	EventKeyProgress = "progress"
	EventKeyError    = "error"
	EventKeyDone     = "done"
	EventKeyFinally  = "finally"
)

type Event struct {
	Key  EventKey  `json:"key"`
	Task *TaskInfo `json:"task"`
	Err  error     `json:"err"`
}
