package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasks(ctx context.Context, id int) ([]entity.Task, error)
	StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error)
	GetTaskByID(ctx context.Context, id int) (entity.Task, error)
	GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) GetTasks(ctx context.Context, id int) ([]entity.Task, error) {
	var res []entity.Task
	tx := r.db.Raw("SELECT * FROM tasks WHERE user_id = ?", id).Scan(&res)
	if tx.Error != nil {
		return []entity.Task{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return []entity.Task{}, nil
	}
	return res, nil
}

func (r *taskRepository) StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error) {
	tx := r.db.Create(task)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return task.ID, nil
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	var res entity.Task
	tx := r.db.Raw("SELECT * FROM tasks WHERE id = ?", id).Scan(&res)
	if tx.Error != nil {
		return entity.Task{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return entity.Task{}, nil
	}
	return res, nil
}

func (r *taskRepository) GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error) {
	var res []entity.Task
	tx := r.db.Raw("SELECT * FROM tasks WHERE category_id = ?", catId).Scan(&res)
	if tx.Error != nil {
		return []entity.Task{}, tx.Error
	}
	if tx.RowsAffected == 0 {
		return []entity.Task{}, nil
	}
	return res, nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *entity.Task) error {
	tx := r.db.Model(&entity.Task{}).Where("id=?", task.ID).Updates(task)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {
	tx := r.db.Where("id = ?", id).Delete(&entity.Task{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
