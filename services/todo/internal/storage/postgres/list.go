package postgres

type ListStorage struct {
	*Storage
}

func NewListStorage(s *Storage) *ListStorage {
	return &ListStorage{s}
}
