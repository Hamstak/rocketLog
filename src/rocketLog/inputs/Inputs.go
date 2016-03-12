package inputs

type Input interface {
	ReadLine() string
	Close()
}