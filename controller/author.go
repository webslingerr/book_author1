package controller

import (
	"app/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (c *Controller) AuthorController(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		c.CreateAuthor(w, r)
	case "GET":
		c.GetListAuthor(w, r)
	case "PUT":
		c.UpdateAuthor(w, r)
	case "DELETE":
		c.DeleteAuthor(w, r)
	}
}

func (c *Controller) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.HandleFuncResponse(w, "Create author", 400, err.Error())
	}

	var createAuthor models.CreateAuthor

	err = json.Unmarshal(body, &createAuthor)
	if err != nil {
		c.HandleFuncResponse(w, "Create author unmarshal", http.StatusBadRequest, err.Error())
		return
	}

	id, err := c.store.Author().Create(&createAuthor)
	if err != nil {
		c.HandleFuncResponse(w, "Storage create author", http.StatusInternalServerError, err.Error())
		return
	}

	author, err := c.store.Author().GetById(&models.AuthorPrimaryKey{Id: id})
	if err != nil {
		c.HandleFuncResponse(w, "Storage get by id author", http.StatusInternalServerError, err.Error())
		return
	}

	body, err = json.Marshal(author)
	if err != nil {
		c.HandleFuncResponse(w, "Storage marshal get by id author", http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func (c *Controller) GetListAuthor(w http.ResponseWriter, r *http.Request) {
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
			c.HandleFuncResponse(w, "Get List author limit", http.StatusBadRequest, err.Error())
			return
		}
	}

	if _, ok := val["offset"]; ok {
		offset, err = strconv.Atoi(val["offset"][0])
		if err != nil {
			c.HandleFuncResponse(w, "Get list author offset", http.StatusBadRequest, err.Error())
			return
		}
	}

	if _, ok := val["Search"]; ok {
		search = val["search"][0]
	}

	authors, err := c.store.Author().GetList(&models.GetListAuthorRequest{
		Limit:  limit,
		Offset: offset,
		Search: search,
	})
	if err != nil {
		c.HandleFuncResponse(w, "Storage get list author", http.StatusInternalServerError, err.Error())
		return
	}

	body, err := json.Marshal(authors)
	if err != nil {
		c.HandleFuncResponse(w, "Storage get list marshal author", http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (c *Controller) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.HandleFuncResponse(w, "Update author", 400, err.Error())
		return
	}

	var updateAuthor models.UpdateAuthor

	err = json.Unmarshal(body, &updateAuthor)
	if err != nil {
		c.HandleFuncResponse(w, "Update author unmarshal", http.StatusBadRequest, err.Error())
		return
	}

	err = c.store.Author().Update(&updateAuthor)
	if err != nil {
		c.HandleFuncResponse(w, "Storage update author", http.StatusInternalServerError, err.Error())
		return
	}

	author, err := c.store.Author().GetById(&models.AuthorPrimaryKey{Id: updateAuthor.Id})
	if err != nil {
		c.HandleFuncResponse(w, "Storage update get by id author", http.StatusInternalServerError, err.Error())
		return
	}

	body, err = json.Marshal(author)
	if err != nil {
		c.HandleFuncResponse(w, "Storage update marshal author", http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func (c *Controller) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.HandleFuncResponse(w, "Delete author", 400, err.Error())
		return
	}

	var deleteAuthor models.AuthorPrimaryKey

	err = json.Unmarshal(body, &deleteAuthor)
	if err != nil {
		c.HandleFuncResponse(w, "Delete unmarshal author", 400, err.Error())
		return
	}

	err = c.store.Author().Delete(&deleteAuthor)
	if err != nil {
		c.HandleFuncResponse(w, "Storage delete author", 400, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
