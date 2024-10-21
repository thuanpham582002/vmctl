package resource

type Visitor interface {
	Visit(VisitorFunc) error
}

type VisitorFunc func(*Info, error) error
