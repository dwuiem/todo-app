package postgres

type TaskStorage struct {
	*Storage
}

func NewTaskStorage(s *Storage) *TaskStorage {
	return &TaskStorage{s}
}
