package processors

type Processor interface {
	Modify(string) string
}

