package main

import (
	"math/rand"
	"time"
)

type Post struct {
	Id        string
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *App) GetAllPost() ([]Post, error) {
	posts := make([]Post, 0)

	rows, err := a.Db.Query(`SELECT id, title, content, created_at, updated_at FROM posts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		singlePost := Post{}
		err = rows.Scan(&singlePost.Id, &singlePost.Title, &singlePost.Content, &singlePost.CreatedAt, &singlePost.UpdatedAt)
		if err != nil {
			return posts, err
		}
		posts = append(posts, singlePost)
	}

	return posts, nil
}

func (a *App) GetSinglePost(id string) (Post, error) {
	post := Post{}

	err := a.Db.QueryRow(`SELECT id, title, content, created_at, updated_at FROM posts WHERE id = $1`, id).Scan(&post.Id, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return post, err
	}

	return post, nil
}

func (a *App) CreatePost(post Post) (string, error) {
	postId := createId()
	var returnId string

	err := a.Db.QueryRow(`INSERT INTO posts(id, title, content, created_at, updated_at) VALUES($1, $2, $3, $4, $5) RETURNING id`, postId, post.Title, post.Content, time.Now(), time.Now()).Scan(&returnId)
	if err != nil {
		return "", err
	}

	return postId, nil
}

func (a *App) UpdatePost(postId, content string) (int, error) {
	res, err := a.Db.Exec(`UPDATE posts set content=$1, updated_at=$2 WHERE id=$3 RETURNING id`, content, time.Now(), postId)
	if err != nil {
		return 0, err
	}

	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsUpdated), err
}

func (a *App) DeletePost(id string) (int, error) {
	res, err := a.Db.Exec(`DELETE FROM posts WHERE id=$1`, id)
	if err != nil {
		return 0, err
	}

	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsUpdated), err
}

const charset = "abcdefghijklmnopqrstuvwxyz" + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func createId() string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	length := 10
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
