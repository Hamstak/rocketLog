package processors

type Processor interface {
	Process(string) string
}

