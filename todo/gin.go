package todo

import "github.com/gin-gonic/gin"

type GinContext struct {
	*gin.Context
}

func NewGinContext(c *gin.Context) *GinContext {
	return &GinContext{Context: c}
}

func (c *GinContext) Bind(v interface{}) error {
	return c.Context.ShouldBindJSON(v)
}

func (c *GinContext) JSON(statuscode int, v interface{}) {
	c.Context.JSON(statuscode, v)
}

func (c *GinContext) TransactionID() string {
	return c.Request.Header.Get("Transaction")
}

func (c *GinContext) Audience() string {
	if aud, ok := c.Get("aud"); ok {
		if s, ok := aud.(string); ok {
			return s
		}
	}
	return ""
}

func (c *GinContext) DefaultQuery(key, defaultValue string) string {
	return c.Context.DefaultQuery(key, defaultValue)
}

func (c *GinContext) Query(key string) string {
	return c.Context.Query(key)
}

func NewGinHandler(handler func(Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(NewGinContext(c))
	}
}
