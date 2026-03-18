package repository

//repository — это слой, который берёт данные из БД и кладёт данные в БД
import (
	"database/sql"
	"task-manager-backend/models"
)

// Интерфейс описывает ЧТО умеет репозиторий.
// "любой репозиторий обязан иметь эти методы"
// handlers/ работает с интерфейсом — не знает про конкретную БД.
// Это позволяет подменять реализацию в тестах (mock).

type TaskRepository interface {
	GetAll(userID int) ([]models.Task, error)     //дай userID → верну список задач
	GetByID(id, userID int) (models.Task, error)  //дай id задачи и юзера → верну одну задачу
	Create(Task models.Task) (models.Task, error) //дай задачу → сохраню → верну её с ID
	Update(task models.Task) error                //обновлю задачу
	Delete(id, userID int) error                  //удалю задачу
}

type taskRepo struct {
	db *sql.DB
}

//создаёт репозиторий и даёт ему доступ к базе
func NewTaskRepo(db *sql.DB) TaskRepository {
	return &taskRepo{db: db}
}

func (r *taskRepo) GetAll(userID int) ([]models.Task, error) {
	rows, err := r.db.Query("SELECT id, title, done, priority, user_id, created_at FROM tasks WHERE user_id=$1 ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //надо после query. Когда функция закончится, строки результата закроются.мчтобы не было утечки ресурсов

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Done, &task.Priority, &task.UserID, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func (r *taskRepo) GetByID(id, userID int) (models.Task, error) {
	var task models.Task

	err := r.db.QueryRow("SELECT id, title, done, priority, user_id, created_at FROM tasks WHERE ID=$1 AND user_id=$2", id, userID).Scan(&task.ID, &task.Title, &task.Done, &task.Priority, &task.UserID, &task.CreatedAt)

	if err != nil {
			return models.Task{}, err
	}

	return task, nil
}

func (r *taskRepo) Create(task models.Task) (models.Task, error) {
	err := r.db.QueryRow("INSERT INTO tasks (title, done, priority, user_id) VALUES ($1,$2,$3,$4) RETURNING id, created_at", task.Title, task.Done, task.Priority, task.UserID).Scan(&task.ID, &task.CreatedAt)
	return task, err
}
