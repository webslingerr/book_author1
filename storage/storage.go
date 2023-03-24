package storage

import "app/api/models"

type StorageI interface {
	CloseDB()
	Book() BookRepoI
	Author() AuthorRepoI
	User() UserRepoI
}

type BookRepoI interface {
	Create(*models.CreateBook) (string, error)
	GetById(*models.BookPrimaryKey) (*models.Book, error)
	GetList(*models.GetListBookRequest) (*models.GetListBookResponse, error)
	Update(*models.UpdateBook) (int64, error)
	Delete(*models.BookPrimaryKey) error
}

type AuthorRepoI interface {
	Create(*models.CreateAuthor) (string, error)
	GetById(*models.AuthorPrimaryKey) (*models.Author, error)
	GetList(*models.GetListAuthorRequest) (*models.GetListAuthorResponse, error)
	Update(*models.UpdateAuthor) (int64, error)
	Delete(*models.AuthorPrimaryKey) error
}

type UserRepoI interface {
	Create(*models.CreateUser) (string, error)
	GetById(*models.UserPrimaryKey) (*models.User, error)
	GetList(*models.GetListUserRequest) (*models.GetListUserResponse, error)
	Update(*models.UpdateUser) (int64, error)
	Delete(*models.UserPrimaryKey) error
}