package todo

import "gorm.io/gorm"

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (s *GormStore) New(todo *Todo) error {
	return s.db.Create(todo).Error
}

func (s *GormStore) List(todos *[]Todo, offset, limit int, createdAt, updatedAt, title string, id string) error {

	query := s.db.Offset(offset).Limit(limit)

	if createdAt != "" {
		query = query.Where("created_at = ?", createdAt)
	}
	if updatedAt != "" {
		query = query.Where("updated_at = ?", updatedAt)
	}
	if title != "" {
		query = query.Where("title = ?", title)
	}
	if id != "" {
		query = query.Where("id = ?", id)
	}

	return query.Find(todos).Error
}

func (s *GormStore) GetByID(todo *Todo, id int) error {
	return s.db.First(todo, id).Error
}

func (s *GormStore) Save(todo *Todo) error {
	return s.db.Save(todo).Error
}

func (s *GormStore) Delete(todo *Todo, id int) error {
	return s.db.Delete(todo, id).Error
}
