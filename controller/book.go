package controller

import (
	"app/models"
	
	"strconv"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (c *Controller) BookController(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		c.CreateBook(w, r)
	case "GET":
		c.GetListBook(w, r)
	case "PUT":
		c.UpdateBook(w, r)
	case "DELETE":
		c.DeleteBook(w, r)
	}
}

func (c *Controller) CreateBook(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.HandleFuncResponse(w, "Create book", 400, err.Error())
	}

	var createBook models.CreateBook

	err = json.Unmarshal(body, &createBook)
	if err != nil {
		c.HandleFuncResponse(w, "Create book unmarshal", http.StatusBadRequest, err.Error())
		return
	}

	id, err := c.store.Book().Create(&createBook)
	if err != nil {
		c.HandleFuncResponse(w, "Storage create book", http.StatusInternalServerError, err.Error())
		return
	}

	book, err := c.store.Book().GetById(&models.BookPrimaryKey{Id: id})
	if err != nil {
		c.HandleFuncResponse(w, "Storage get by id book", http.StatusInternalServerError, err.Error())
		return
	}

	body, err = json.Marshal(book)
	if err != nil {
		c.HandleFuncResponse(w, "Storage marshal get by id book", http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func (c *Controller) GetListBook(w http.ResponseWriter, r *http.Request) {
	var (
		val    = r.URL.Query()
		limit  int
		offset int
		search string
		err    error
	)

	if _, ok := val["limit"]; ok {
		limit, err = strconv.Atoi(val["limit"][0])
		if err != nil {
			c.HandleFuncResponse(w, "Get List book limit", http.StatusBadRequest, err.Error())
			return
		}
	}

	if _, ok := val["offset"]; ok {
		offset, err = strconv.Atoi(val["offset"][0])
		if err != nil {
			c.HandleFuncResponse(w, "Get list book offset", http.StatusBadRequest, err.Error())
			return
		}
	}

	if _, ok := val["Search"]; ok {
		search = val["search"][0]
	}

	books, err := c.store.Book().GetList(&models.GetListBookRequest{
		Limit:  limit,
		Offset: offset,
		Search: search,
	})
	if err != nil {
		c.HandleFuncResponse(w, "Storage get list book", http.StatusInternalServerError, err.Error())
		return
	}

	body, err := json.Marshal(books)
	if err != nil {
		c.HandleFuncResponse(w, "Storage get list marshal book", http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (c *Controller) UpdateBook(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.HandleFuncResponse(w, "Update book", 400, err.Error())
		return
	}

	var updateBook models.UpdateBook

	err = json.Unmarshal(body, &updateBook)
	if err != nil {
		c.HandleFuncResponse(w, "Update book unmarshal", http.StatusBadRequest, err.Error())
		return
	}

	err = c.store.Book().Update(&updateBook)
	if err != nil {
		c.HandleFuncResponse(w, "Storage update book", http.StatusInternalServerError, err.Error())
		return
	}

	book, err := c.store.Book().GetById(&models.BookPrimaryKey{Id: updateBook.Id})
	if err != nil {
		c.HandleFuncResponse(w, "Storage update get by id book", http.StatusInternalServerError, err.Error())
		return
	}

	body, err = json.Marshal(book)
	if err != nil {
		c.HandleFuncResponse(w, "Storage update marshal book", http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func (c *Controller) DeleteBook(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.HandleFuncResponse(w, "Delete book", 400, err.Error())
		return
	}

	var deleteBook models.BookPrimaryKey

	err = json.Unmarshal(body, &deleteBook)
	if err != nil {
		c.HandleFuncResponse(w, "Delete unmarshal book", 400, err.Error())
		return
	}

	err = c.store.Book().Delete(&deleteBook)
	if err != nil {
		c.HandleFuncResponse(w, "Storage delete book", 400, err.Error()) 
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}