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

type TodoHandler struct {
	db *gorm.DB
}

func NewTodoHandler(db *gorm.DB) *TodoHandler {
	return &TodoHandler{db: db}
}

func (t *TodoHandler) NewTask(c *gin.Context) {
	/*
		s := c.Request.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(s, "Bearer ")

		if err := auth.Protect(tokenString); err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	*/

	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if todo.Title == "sleep" {
		transactionID := c.Request.Header.Get("TransactionID")
		aud, _ := c.Get("aud")
		log.Println(transactionID, aud, "not allow")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "not allow",
		})
	}

	r := t.db.Create(&todo)
	if err := r.Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"ID": todo.Model.ID,
	})
}

func (t *TodoHandler) List(c *gin.Context) {
	id := c.Query("id")
	createdAt := c.Query("created_at")
	updatedAt := c.Query("updated_at")
	title := c.Query("title")

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
		query := t.db
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

func (t *TodoHandler) Remove(c *gin.Context) {
	idParam := c.Query("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	r := t.db.Delete(&Todo{}, id)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (t *TodoHandler) Update(c *gin.Context) {
	idParam := c.Query("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "id not found",
		})
		return
	}

	var todo Todo
	if err := t.db.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Record not found",
		})
		return
	}

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	t.db.Save(&todo)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
