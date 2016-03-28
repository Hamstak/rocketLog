package inputs

type Input interface {
	ReadLine() (string, error)
	Close()
	GetType() string
	GetInputName() string
}
