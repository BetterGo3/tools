package a

type Person struct {
	Name string
}

type Stringer interface {
	String() string
}
