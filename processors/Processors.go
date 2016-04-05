package processors

type Processor interface {
	Process(string) string
	Matches(string) bool
    ToString() string
}
