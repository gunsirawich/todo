package todo

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Todo struct {
	Title string `json:"text"`
	gorm.Model
}

func (Todo) TableName() string {
	return "todolist"
}

type storer interface {
	New(*Todo) error
	List(*[]Todo, int, int, string, string, string, string) error
	Save(*Todo) error
	Delete(*Todo, int) error
	GetByID(*Todo, int) error
}

type TodoHandler struct {
	store storer
}

func NewTodoHandler(store storer) *TodoHandler {
	return &TodoHandler{store: store}
}

type Context interface {
	Bind(interface{}) error
	JSON(int, interface{})
	TransactionID() string
	Audience() string
	DefaultQuery(string, string) string
	Query(string) string
}

func (t *TodoHandler) NewTask(c Context) {
	/*
		s := c.Request.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(s, "Bearer ")

		if err := auth.Protect(tokenString); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	*/

	var todo Todo
	//if err := c.ShouldBindJSON(&todo); err != nil {
	if err := c.Bind(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if todo.Title == "sleep" {
		//transactionID := c.Request.Header.Get("TransactionID")
		transactionID := c.TransactionID()
		//aud, _ := c.Get("aud")
		aud := c.Audience()
		log.Println(transactionID, aud, "not allow")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "not allow",
		})
	}

	err := t.store.New(&todo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"ID": todo.Model.ID,
	})
}

/*
func (t *TodoHandler) List(c *gin.Context) {
	id := c.Query("id")
	createdAt := c.Query("created_at")
	updatedAt := c.Query("updated_at")
	title := c.Query("title")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 10
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	query := t.db
	query = query.Offset(offset).Limit(pageSize)

	if id != "" {
		var todo Todo
		if err := t.db.First(&todo, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, todo)
	} else {
		var todos []Todo
		if createdAt != "" {
			query = query.Where("created_at = ?", createdAt)
		}
		if updatedAt != "" {
			query = query.Where("updated_at = ?", updatedAt)
		}
		if title != "" {
			query = query.Where("title = ?", title)
		}
		r := query.Find(&todos)
		if err := r.Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, todos)
	}
}
*/

func (t *TodoHandler) List(c Context) {
	var todos []Todo
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	createdAt := c.Query("created_at")
	updatedAt := c.Query("updated_at")
	title := c.Query("title")
	id := c.Query("id")

	if err := t.store.List(&todos, offset, pageSize, createdAt, updatedAt, title, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func (t *TodoHandler) Remove(c Context) {
	var todo Todo

	idParam := c.Query("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = t.store.Delete(&todo, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (t *TodoHandler) Update(c Context) {
	idParam := c.Query("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id not found",
		})
		return
	}

	var todo Todo
	err = t.store.GetByID(&todo, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Record not found",
		})
		return
	}

	if err := c.Bind(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = t.store.Save(&todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
