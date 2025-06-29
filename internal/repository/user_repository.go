// FILEPATH: /Users/sayanseksenbaev/Programming/PaperExamGrader/internal/repository/user_repository.go

package repository

import (
	"sync"

	"PaperExamGrader/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

var userRepoInstance *UserRepository
var userRepoOnce sync.Once

func NewUserRepository(db *gorm.DB) *UserRepository {
	userRepoOnce.Do(func() {
		userRepoInstance = &UserRepository{db: db}
	})
	return userRepoInstance
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *UserRepository) List() ([]model.User, error) {
	var users []model.User
	err := r.db.Find(&users).Error
	return users, err
}
