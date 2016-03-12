package outputs

type Output interface {
	Write(data string)
	Close()
}