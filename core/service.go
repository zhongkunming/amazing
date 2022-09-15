package core

type Service interface {
	Description() string

	Load()
}
