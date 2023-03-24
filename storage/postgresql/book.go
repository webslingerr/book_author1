package postgresql

import (
	"app/api/models"
	"app/pkg/helper"
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
			count,
			income_price,
			profit_status,
			profit_price,
			author_id,
			updated_at
		)
		VALUES($1, $2, $3, $4, $5, $6, $7, NOW())
	`

	_, err := b.db.Exec(
		query,
		id,
		req.Name,
		req.Count,
		req.IncomePrice,
		req.ProfitStatus,
		req.ProfitPrice,
		req.AuthorId,
	)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (b *bookRepo) GetById(req *models.BookPrimaryKey) (*models.Book, error) {
	var (
		query string
		book  models.Book
	)

	query = `
		SELECT
			b.id,
			b.name, 
			count,
			income_price,
			profit_status,
			profit_price,
			sell_price,
			COALESCE(a.id, ''),
			COALESCE(a.name, ''),
			TO_CHAR(b.created_at, 'YYYY-MM-DD HH24-MI-SS'),
			b.updated_at
		FROM book AS b
		LEFT JOIN author AS a ON a.id = b.author_id
		WHERE b.id = $1
	`

	err := b.db.QueryRow(query, req.Id).Scan(
		&book.Id,
		&book.Name,
		&book.Count,
		&book.IncomePrice,
		&book.ProfitStatus,
		&book.ProfitPrice,
		&book.SellPrice,
		&book.Author.Id,
		&book.Author.Name,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (b *bookRepo) GetList(req *models.GetListBookRequest) (*models.GetListBookResponse, error) {
	resp := &models.GetListBookResponse{}

	var (
		query  string
		filter = " WHERE TRUE "
		offset = " OFFSET 0 "
		limit  = " LIMIT 10"
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
		LEFT JOIN author AS a ON a.id = b.author_id
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
		var book models.Book
		var author models.ReturnAuthor

		err = rows.Scan(
			&book.Id,
			&book.Name,
			&book.Count,
			&book.IncomePrice,
			&book.ProfitStatus,
			&book.ProfitPrice,
			&book.SellPrice,
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

	resp.Count = len(resp.Books)

	return resp, nil
}

func (b *bookRepo) Update(req *models.UpdateBook) (int64, error) {
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE 
			book
		SET 
			name = :name,
			count = :count,
			income_price = :income_price,
			profit_status = :profit_status,
			profit_price = :profit_price,
			author_id = :author_id,
			updated_at = now()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":            req.Id,
		"name":          req.Name,
		"count":         req.Count,
		"income_price":  req.IncomePrice,
		"profit_status": req.ProfitStatus,
		"profit_price":  req.ProfitPrice,
		"author_id":     req.AuthorId,
	}

	query, args := helper.ReplaceQueryParams(query, params)
	result, err := b.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
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
