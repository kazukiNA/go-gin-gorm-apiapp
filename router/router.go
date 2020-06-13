package router

import (
	"github.com/gin-gonic/gin"
	"postapp/models"
)

func Router() {
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*.html")

	posts := GetPost()

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{"posts": posts})
	})

	router.GET("/register", func(ctx *gin.Context) {
		ctx.HTML(200, "register.html", gin.H{})
	})

	router.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(200, "login.html", gin.H{})
	})

	router.POST("/create",func(ctx *gin.Context){
		text := ctx.PostForm("text")

		CreatePost(text)
		ctx.Redirect(302,"/")
	})

	router.POST("/register",func(ctx *gin.Context){
		email := ctx.PostForm("email")
		password := ctx.PostForm("password")

		CreateUser(email,password)
		ctx.Redirect(302,"/")
	})

	router.GET("/delete/:id", func(ctx *gin.Context) {
		post := models.Post{}
		id := ctx.Param("id")

		DeletePost(post,id)
		ctx.Redirect(302,"/")
	})
}