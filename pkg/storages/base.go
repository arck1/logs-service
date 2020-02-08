package storages

type BaseStorage interface {
	Flush() error
	Save(message []byte) error
	Close() error

	Start()
}

type LogMessage struct {
	Message string
	Level   string
	Place   string
	Args    string
	Error   interface{}

	DocumentId string
	RequestId   string
}
