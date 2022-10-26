package core

type Loadable interface {
	Judge() bool

	Load()
}
