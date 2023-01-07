package parser

type Parser[T interface{}] interface {
	Parse(body []byte) T
}
