package repository

import (
	"database/sql"
	"fmt"
	entity "main/internal/Domain/Entity"
	repository "main/internal/Domain/Repository"
	"time"

	"github.com/google/uuid"
)

type postRepository struct {
	db *sql.DB
}

func (p postRepository) Save(post entity.Post) error {

	fmt.Printf("Saving post: %+v\n", post)
	_, err := p.db.Exec(`
		INSERT INTO posts (
			id,
			created_at,
			updated_at,
			slug,
			title,
			content,
			author)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		`,
		post.Id().String(),
		post.CreatedAt().Format(time.RFC3339),
		post.UpdatedAt().Format(time.RFC3339),
		post.Slug(),
		post.Title(),
		post.Content(),
		post.Author(),
	)

	return err
}

func (p postRepository) FindByID(id uuid.UUID) (entity.Post, error) {
	row := p.db.QueryRow(`
		SELECT id, created_at, updated_at, slug, title, content, author FROM posts WHERE id = ?
	`, id)

	var postId uuid.UUID
	var createdAt time.Time
	var updatedAt time.Time
	var slug string
	var title string
	var content string
	var author string

	err := row.Scan(&postId, &createdAt, &updatedAt, &slug, &title, &content, &author)

	if err != nil {
		return entity.Post{}, err
	}

	return entity.NewPost(postId, createdAt, updatedAt, slug, title, content, author), nil
}

func (p postRepository) FindBySlug(postSlug string) (entity.Post, error) {
	row := p.db.QueryRow(`
		SELECT id, created_at, updated_at, slug, title, content, author FROM posts WHERE slug = ?
	`, postSlug)

	var postId uuid.UUID
	var createdAt time.Time
	var updatedAt time.Time
	var slug string
	var title string
	var content string
	var author string

	err := row.Scan(&postId, &createdAt, &updatedAt, &slug, &title, &content, &author)

	if err != nil {
		return entity.Post{}, err
	}

	return entity.NewPost(postId, createdAt, updatedAt, slug, title, content, author), nil
}

func (p postRepository) FindAll() ([]entity.Post, error) {
	rows, err := p.db.Query(`
		SELECT id, created_at, updated_at, slug, title, content, author FROM posts
	`)

	if err != nil {
		return nil, err
	}

	var posts []entity.Post
	for rows.Next() {
		var postId uuid.UUID
		var createdAt time.Time
		var updatedAt time.Time
		var slug string
		var title string
		var content string
		var author string

		err := rows.Scan(&postId, &createdAt, &updatedAt, &slug, &title, &content, &author)
		if err != nil {
			return nil, err
		}

		posts = append(posts, entity.NewPost(postId, createdAt, updatedAt, slug, title, content, author))
	}

	return posts, nil
}

func (p postRepository) FindAllByAuthor(author string) ([]entity.Post, error) {
	rows, err := p.db.Query(`
		SELECT id, created_at, updated_at, slug, title, content, author FROM posts WHERE author = ?
	`, author)

	if err != nil {
		return nil, err
	}

	var posts []entity.Post
	for rows.Next() {
		var postId uuid.UUID
		var createdAt time.Time
		var updatedAt time.Time
		var slug string
		var title string
		var content string
		var author string

		err := rows.Scan(&postId, &createdAt, &updatedAt, &slug, &title, &content, &author)
		if err != nil {
			return nil, err
		}

		posts = append(posts, entity.NewPost(postId, createdAt, updatedAt, slug, title, content, author))
	}

	return posts, nil
}

func (p postRepository) FindAllByText(text string) ([]entity.Post, error) {
	rows, err := p.db.Query(`
		SELECT id, created_at, updated_at, slug, title, content, author FROM posts WHERE title LIKE ? OR content LIKE ?
	`, "%"+text+"%", "%"+text+"%")

	if err != nil {
		return nil, err
	}

	var posts []entity.Post
	for rows.Next() {
		var postId uuid.UUID
		var createdAt time.Time
		var updatedAt time.Time
		var slug string
		var title string
		var content string
		var author string

		err := rows.Scan(&postId, &createdAt, &updatedAt, &slug, &title, &content, &author)
		if err != nil {
			return nil, err
		}

		posts = append(posts, entity.NewPost(postId, createdAt, updatedAt, slug, title, content, author))
	}

	return posts, nil
}

func (p postRepository) Delete(id uuid.UUID) error {
	_, err := p.db.Exec(`
		DELETE FROM posts WHERE id = ?
	`, id)

	return err
}

func NewPostRepository(db *sql.DB) repository.PostRepository {
	return &postRepository{db: db}
}
