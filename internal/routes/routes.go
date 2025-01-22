package routes

import (
	"fibo_go_server/internal/controllers"
	"fibo_go_server/internal/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())
	db, err := utils.InitDB()
	if err != nil {
		panic("Could not connect to the database")
	}

	r.POST("/signUp", controllers.SignUpUser(db))
	r.POST("/login", controllers.LoginUser(db))

	r.POST("/createPost", controllers.CreatePost(db))
	r.GET("/getPostsList", controllers.GetPostsList(db))
	r.GET("/getPost", controllers.GetPostDetails(db))

	r.GET("/categories", controllers.GetCategories(db))

	r.POST("/comment", controllers.CreateComment(db))

	r.POST("/like", controllers.AddLike(db))

	r.POST("/salary", controllers.CalculateSalary(db))
	return r
}
