package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) InitUserRoute(r *gin.RouterGroup) {
	user := r.Group("/user")
	{
		user.POST("/add", s.addUserHandler)
		user.POST("/del", s.delUserHandler)
		user.GET("/:id", s.getUserHandler)
		user.GET("/", s.getUserByQueryHandler)
	}
}

type addUserParams struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func (s *Server) addUserHandler(c *gin.Context) {
	params := addUserParams{}
	err := c.ShouldBind(&params) // 兼容form-data和x-www-form-urlencoded和json
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("params: %v", params)

	err = s.db.AddUser(params.Username, params.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User added successfully"})
}

type delUserParams struct {
	ID uint64 `form:"id"`
}

func (s *Server) delUserHandler(c *gin.Context) {
	params := delUserParams{}
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = s.db.DelUser(params.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (s *Server) getUserHandler(c *gin.Context) {
	id := c.Param("id") // 获取路径中的:id
	users, err := s.db.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (s *Server) getUserByQueryHandler(c *gin.Context) {
	id := c.Query("id") // 获取查询参数中的id
	users, err := s.db.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
