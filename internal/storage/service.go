package storage

type Storage interface {
}

type storage struct {
}

func NewStorage() Storage {
	return &storage{}
}
