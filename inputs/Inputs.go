package inputs

type Input interface {
	ReadLine() (string, error)
	Close()
	Flush()
	GetType() string
}