package storage

import "app/models"

type StorageI interface {
	CloseDB()
	Book() BookRepoI
	Author() AuthorRepoI
}

type BookRepoI interface {
	Create(*models.CreateBook) (string, error)
	GetById(*models.BookPrimaryKey) (*models.Book, error)
	GetList(*models.GetListBookRequest) (*models.GetListBookResponse, error)
	Update(*models.UpdateBook) error
	Delete(*models.BookPrimaryKey) error
}

type AuthorRepoI interface {
	Create(*models.CreateAuthor) (string, error)
	GetById(*models.AuthorPrimaryKey) (*models.Author, error)
	GetList(*models.GetListAuthorRequest) (*models.GetListAuthorResponse, error)
	Update(*models.UpdateAuthor) error
	Delete(*models.AuthorPrimaryKey) error
}
