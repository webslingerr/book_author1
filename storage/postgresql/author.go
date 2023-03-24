package postgresql

import (
	"app/models"
	"fmt"

	"database/sql"

	"github.com/google/uuid"
)

type authorRepo struct {
	db *sql.DB
}

func NewAuthorRepo(db *sql.DB) *authorRepo {
	return &authorRepo{
		db: db,
	}
}

func (a *authorRepo) Create(req *models.CreateAuthor) (string, error) {
	var (
		query string
		id    = uuid.New().String()
	)

	query = `
		INSERT INTO author(id, name, updated_at)
		VALUES($1, $2, NOW())
	`

	_, err := a.db.Exec(query, id, req.Name)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (a *authorRepo) GetById(req *models.AuthorPrimaryKey) (*models.Author, error) {
	var (
		query string
		author  models.Author
	)

	query = `
		SELECT
			id,
			name, 
			created_at,
			updated_at
		FROM author
		WHERE id = $1
	`

	err := a.db.QueryRow(query, req.Id).Scan(
		&author.Id,
		&author.Name,
		&author.CreatedAt,
		&author.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &author, nil
}

func (a *authorRepo) GetList(req *models.GetListAuthorRequest) (*models.GetListAuthorResponse, error) {
	resp := &models.GetListAuthorResponse{}

	var (
		query string
		filter = " WHERE TRUE "
		offset = " OFFSET 0 "
		limit = " LIMIT 10"

		author models.Author
	)

	query = `
		SELECT 
			COUNT(*) OVER(),
			id,
			name,
			created_at,
			updated_at
		FROM author
	`

	if len(req.Search) > 0 {
		filter += " AND name ILIKE '%' || '" + req.Search + "' || '%' "
	}

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	query += filter + offset + limit

	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(
			&resp.Count,
			&author.Id,
			&author.Name,
			&author.CreatedAt,
			&author.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Authors = append(resp.Authors, &author)
	}

	return resp, nil
}

func (a *authorRepo) Update(req *models.UpdateAuthor) error {
	var (
		query string
	)

	query = `
		UPDATE 
			author
		SET name = $1, updated_at = NOW()
		WHERE id = $2
	`

	_, err := a.db.Exec(query, req.Name, req.Id)
	if err != nil {
		return err
	}

	return nil
}

func (a *authorRepo) Delete(req *models.AuthorPrimaryKey) error {
	var (
		query string
	)

	query = `
		DELETE FROM author
		WHERE id = $1
	`

	_, err := a.db.Exec(query, req.Id)
	if err != nil {
		return err
	}

	return nil
}