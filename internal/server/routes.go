package server

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	v1 := r.Group("/api/v1")
	v1.GET("/", s.HelloWorldHandler)

	userV1 := v1.Group("/user")
	userV1.POST("/add", s.addUserHandler)
	userV1.POST("/del", s.delUserHandler)
	userV1.GET("/:id", s.getUserHandler)
	userV1.GET("/", s.getUserByQueryHandler)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
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
