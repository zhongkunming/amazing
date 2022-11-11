package core

type Loadable interface {
	CanLoad() bool

	Load()
}
