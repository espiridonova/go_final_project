package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const (
	dateLayout       = "20060102"
	searchDateLayout = "02.01.2006"
)

type Task struct {
	ID      int64  `json:"id,string"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(t *Task) (int64, error) {
	var id int64

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := db.Exec(query,
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func Tasks(search string, limit int) ([]*Task, error) {
	where := ""
	args := []any{sql.Named("limit", limit)}
	if search != "" {
		searchDate, err := time.Parse(searchDateLayout, search)
		if err != nil {
			where = "WHERE title LIKE :search OR comment LIKE :search"
			args = append(args, sql.Named("search", fmt.Sprintf("%%%s%%", search)))
		} else {
			where = "WHERE date = :searchDate"
			args = append(args, sql.Named("searchDate", searchDate.Format(dateLayout)))
		}
	}

	query := fmt.Sprintf(`SELECT 
			id, 
			date, 
			title, 
			comment, 
			repeat 
		FROM scheduler 
		%s
		ORDER BY date LIMIT :limit`, where)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]*Task, 0)

	for rows.Next() {
		var t Task
		err = rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return tasks, err
}

func GetTask(id string) (*Task, error) {
	t := Task{}
	row := db.QueryRow(
		`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id`,
		sql.Named("id", id))
	err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &t, err
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = :date, title = :title, 
                     comment = :comment, repeat = :repeat WHERE id = :id`
	res, err := db.Exec(query, sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = :date WHERE id = :id`
	_, err := db.Exec(query, sql.Named("id", id),
		sql.Named("date", next))

	return err
}

func DeleteTask(id string) error {
	query := "DELETE FROM scheduler WHERE id = :id"
	_, err := db.Exec(query, sql.Named("id", id))

	return err
}
