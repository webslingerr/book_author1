package postgresql

import (
	"app/models"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type bookRepo struct {
	db *sql.DB
}

func NewBookRepo(db *sql.DB) *bookRepo {
	return &bookRepo{
		db: db,
	}
}

func (b *bookRepo) Create(req *models.CreateBook) (string, error) {
	var (
		query string
		id    = uuid.New().String()
	)

	query = `
		INSERT INTO book(
			id, 
			name,
			price, 
			author_id,
			updated_at
		)
		VALUES($1, $2, $3, $4, NOW())
	`

	_, err := b.db.Exec(query, id, req.Name, req.Price, req.AuthorId)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (b *bookRepo) GetById(req *models.BookPrimaryKey) (*models.Book, error) {
	var (
		query string
		book  models.Book
		author models.ReturnAuther
	)

	query = `
		SELECT
			b.id,
			b.name, 
			price,
			a.id,
			a.name,
			b.created_at,
			b.updated_at
		FROM book AS b
		JOIN author AS a ON a.id = b.author_id
		WHERE b.id = $1
	`

	err := b.db.QueryRow(query, req.Id).Scan(
		&book.Id,
		&book.Name,
		&book.Price,
		&author.Id,
		&author.Name,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	book.Author = author

	return &book, nil
}

func (b *bookRepo) GetList(req *models.GetListBookRequest) (*models.GetListBookResponse, error) {
	resp := &models.GetListBookResponse{}

	var (
		query  string
		filter = " WHERE TRUE "
		offset = " OFFSET 0 "
		limit  = " LIMIT 10"

		book models.Book
		author models.ReturnAuther
	)

	query = `
		SELECT 
			COUNT(*) OVER(),
			b.id,
			b.name, 
			price,
			a.id,
			a.name,
			b.created_at,
			b.updated_at
		FROM book AS b
		JOIN author AS a ON a.id = b.author_id
	`

	if len(req.Search) > 0 {
		filter += " AND b.name ILIKE '%' || '" + req.Search + "' || '%' "
	}

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	query += filter + offset + limit

	rows, err := b.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&resp.Count,
			&book.Id,
			&book.Name,
			&book.Price,
			&author.Id,
			&author.Name,
			&book.CreatedAt,
			&book.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		book.Author = author
		resp.Books = append(resp.Books, &book)
	}

	return resp, nil
}

func (b *bookRepo) Update(req *models.UpdateBook) error {
	var (
		query string
	)

	query = `
		UPDATE 
			book 
		SET name = $1, price = $2, author_id = $3
		WHERE id = $4
	`

	_, err := b.db.Exec(query, req.Name, req.Price, req.AuthorId, req.Id)
	if err != nil {
		return err
	}

	return nil
}

func (b *bookRepo) Delete(req *models.BookPrimaryKey) error {
	var (
		query string
	)

	query = `
		DELETE FROM book
		WHERE id = $1
	`

	_, err := b.db.Exec(query, req.Id)
	if err != nil {
		return err
	}

	return nil
}
