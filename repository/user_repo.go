package repository

import (
	"database/sql"
	"task-manager-backend/models"
)

type UserRepository interface {
	Create(user models.User) (models.User, error)
	CreateWithTask(user models.User, task models.Task) error
	FindByEmail(email string) (models.User, error)
}

type userRepo struct {
	db *sql.DB //доступ к базе
}

func NewUserRepo(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user models.User) (models.User, error) {
	err := r.db.QueryRow("INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at", user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt) //id и created_at генерирует база
	return user, err
}

func (r *userRepo) CreateWithTask(user models.User, task models.Task) error {
	tx, err := r.db.Begin() 
	if err != nil {
		return err
	}

	//defer гарантирует rollback если что-то не так
	//выполни этот код в самом конце функци
	defer func() {
		if p := recover(); p != nil { //Это обработка паники.
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	} ()

	var userID int
	err = tx.QueryRow("INSERT INTO users (email, password_hash) VALUES ($1,$2) RETURNING id", user.Email, user.PasswordHash).Scan(&userID)
	if err != nil {return err}
	_, err = tx.Exec("INSERT INTO tasks (title,done,priority,user_id) VALUES ($1,$2,$3,$4)", task.Title, task.Done, task.Priority, userID)
	if err != nil {
		return err
	}
	return tx.Commit() //теперь можно окончательно сохранить изменения
}

func (r *userRepo) FindByEmail(email string) (models.User, error) {
	var user models.User

	err := r.db.QueryRow("SELECT id, email, password_hash, created_at FROM users WHERE email = $1", email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	return user, err
}
