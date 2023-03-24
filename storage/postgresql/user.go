package postgresql

import (
	"app/api/models"
	"app/pkg/helper"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (u *userRepo) Create(req *models.CreateUser) (string, error) {

	var (
		query string
		id    = uuid.New().String()
	)

	query = `
			INSERT INTO "user" (
				id,
				name, 
				balance,
				updated_at
			)
			VALUES($1, $2, $3, NOW())
		`

	_, err := u.db.Exec(query, id, req.Name, req.Balance)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (a *userRepo) GetById(req *models.UserPrimaryKey) (*models.User, error) {
	var (
		query string
		user  models.User
	)

	query = `
		SELECT 
			id,
			name, 
			balance,
			created_at,
			updated_at
		FROM "user"
		WHERE id = $1
	`

	err := a.db.QueryRow(query, req.Id).Scan(
		&user.Id,
		&user.Name,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (b *userRepo) GetList(req *models.GetListUserRequest) (*models.GetListUserResponse, error) {
	resp := &models.GetListUserResponse{}

	var (
		query  string
		filter = " WHERE TRUE "
		offset = " OFFSET 0 "
		limit  = " LIMIT 10"
	)

	query = `
		SELECT 
			id,
			name,
			balance,
			created_at,
			updated_at
		FROM "user"
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

	rows, err := b.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User

		err = rows.Scan(
			&user.Id,
			&user.Name,
			&user.Balance,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		resp.Users = append(resp.Users, &user)
	}

	resp.Count = len(resp.Users)

	return resp, nil
}

func (b *userRepo) Update(req *models.UpdateUser) (int64, error) {
	var (
		query  string
		params map[string]interface{}
	)

	query = `
		UPDATE 
			"user"
		SET 
			name = :name,
			balance = :balance,
			updated_at = now()
		WHERE id = :id
	`

	params = map[string]interface{}{
		"id":      req.Id,
		"name":    req.Name,
		"balance": req.Balance,
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

func (b *userRepo) Delete(req *models.UserPrimaryKey) error {
	var (
		query string
	)

	query = `
		DELETE FROM "user"
		WHERE id = $1
	`

	_, err := b.db.Exec(query, req.Id)
	if err != nil {
		return err
	}

	return nil
}
