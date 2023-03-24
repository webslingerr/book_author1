package models

type Book struct {
	Id           string       `json:"id"`
	Name         string       `json:"name"`
	Count        int          `json:"count"`
	IncomePrice  float64      `json:"income_price"`
	ProfitStatus string       `json:"profit_status"`
	ProfitPrice  float64      `json:"profit_price"`
	SellPrice    float64      `json:"sell_price"`
	Author       ReturnAuthor `json:"author"`
	CreatedAt    string       `json:"created_at"`
	UpdatedAt    string       `json:"updated_at"`
}

type BookPrimaryKey struct {
	Id string `json:"id"`
}

type CreateBook struct {
	Name         string  `json:"name"`
	Count        int     `json:"count"`
	IncomePrice  float64 `json:"income_price"`
	ProfitStatus string  `json:"profit_status"`
	ProfitPrice  float64 `json:"profit_price"`
	AuthorId     string  `json:"author_id"`
}

type UpdateBook struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Count        int     `json:"count"`
	IncomePrice  float64 `json:"income_price"`
	ProfitStatus string  `json:"profit_status"`
	ProfitPrice  float64 `json:"profit_price"`
	AuthorId  string  `json:"author_id"`
	UpdatedAt string  `json:"updated_at"`
}

type GetListBookRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}

type GetListBookResponse struct {
	Count int     `json:"count"`
	Books []*Book `json:"books"`
}
