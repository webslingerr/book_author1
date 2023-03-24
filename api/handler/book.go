package handler

import (
	"app/api/models"
	"app/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateBook(c *gin.Context) {
	var createBook models.CreateBook

	err := c.ShouldBindJSON(&createBook)
	if err != nil {
		h.handlerResponse(c, "create book", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.storages.Book().Create(&createBook)
	if err != nil {
		h.handlerResponse(c, "storage.create.book", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.storages.Book().GetById(&models.BookPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.getbyid.book", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create book", http.StatusCreated, resp)
}

func (h *Handler) GetByIdBook(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "get by id book", http.StatusBadRequest, "invalid book id")
		return
	}

	resp, err := h.storages.Book().GetById(&models.BookPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.book.get by id", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id book", http.StatusOK, resp)
}

func (h *Handler) GetListBook(c *gin.Context) {
	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list book", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list book", http.StatusBadRequest, "invalid offset")
		return
	}

	resp, err := h.storages.Book().GetList(&models.GetListBookRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})

	if err != nil {
		h.handlerResponse(c, "storage.book.getlist", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list book response", http.StatusOK, resp)
}

func (h *Handler) UpdateBook(c *gin.Context) {
	var updateBook models.UpdateBook

	id := c.Param("id")
	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "update book", http.StatusBadRequest, "invalid book id")
	}

	err := c.ShouldBindJSON(&updateBook)
	if err != nil {
		h.handlerResponse(c, "update book", http.StatusBadRequest, err.Error())
		return
	}

	updateBook.Id = id

	rowsAffected, err := h.storages.Book().Update(&updateBook)
	if err != nil {
		h.handlerResponse(c, "storage.update book", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.book.update", http.StatusBadRequest, "no rows affected")
		return
	}

	resp, err := h.storages.Book().GetById(&models.BookPrimaryKey{Id: updateBook.Id})
	if err != nil {
		h.handlerResponse(c, "storage.getbyid book", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "update book response", http.StatusAccepted, resp)
}

func (h *Handler) DeleteBook(c *gin.Context) {
	id := c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "delete book", http.StatusBadRequest, "invalid book id")
		return
	}

	err := h.storages.Book().Delete(&models.BookPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.book.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "delete book", http.StatusNoContent, nil)
}