package types

type List[T any] struct {
	Items []T
}

type ModuleCommand struct {
	Name       string
	Type       string
	Controller bool
}
