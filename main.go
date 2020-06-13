package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"log"
	"net/http"
	sessioninfo "postapp/SessionInfo"
	"strconv"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"postapp/models"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

func db_init(){
	db, err := gorm.Open("mysql", "root:secret@/postapp?charset=utf8&parseTime=True")
	if err != nil {
		panic("failed to connect database\n")
	}
	db.AutoMigrate(&models.User{})
}

func createPost(text string, email string, updatedat string){
	db, err := gorm.Open("mysql", "root:secret@/postapp?charset=utf8&parseTime=True")
	if err != nil {
		panic("failed to connect database\n")
	}
	db.Create(&models.Post{Text: text, UserEmail: email, JPTime: updatedat})
	//db.AutoMigrate(&models.Post{})
}

func createUser(email string, password []byte){
	db, err := gorm.Open("mysql", "root:secret@/postapp?charset=utf8&parseTime=True")
	if err != nil {
		panic("failed to connect database\n")
	}
	db.Create(&models.User{Email: email, Password: password})
}

func getPostALL() []models.Post{
	db, err := gorm.Open("mysql", "root:secret@/postapp?charset=utf8&parseTime=True")
	if err != nil {
		panic("failed to connect database\n")
	}
	var post []models.Post
	db.Find(&post)
	return post
}

func getPostOne(id int) models.Post{
	db, err := gorm.Open("mysql", "root:secret@/postapp?charset=utf8&parseTime=True")
	if err != nil {
		panic("failed to connect database\n")
	}
	var post models.Post
	db.First(&post,id)
	db.Close()
	return post
}

func getUser(email string) models.User{
	db, err := gorm.Open("mysql", "root:secret@/postapp?charset=utf8&parseTime=True")
	if err != nil {
		panic("failed to connect database\n")
	}
	var user models.User
	db.First(&user, "email=?",email)
	db.Close()
	return user
}

func deletePost(id int){
	db, err := gorm.Open("mysql", "root:secret@/postapp?charset=utf8&parseTime=True")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()
	var post models.Post
	db.First(&post,id)
	db.Delete(&post)
}

func updatePost(id int, text string){
	db, err := gorm.Open("mysql", "root:secret@/postapp?charset=utf8&parseTime=True")
	if err != nil {
		panic("failed to connect database\n")
	}
	defer db.Close()
	var post models.Post
	db.First(&post,id)
	post.Text = text
	post.JPTime = time.Now().Format("2006-01-02")
	db.Save(&post)
}

func Login(ctx *gin.Context, Email string  ){
	session := sessions.Default(ctx)
	session.Set("Email",Email)
	session.Save()
}

func Logout(ctx *gin.Context){
	session := sessions.Default(ctx)
	session.Clear()
	log.Println("クリア処理")
	session.Save()
}


func hash(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func verify(hash, s string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(s))
}

func sessionCheck() gin.HandlerFunc{
	return func (ctx *gin.Context){

		session := sessions.Default(ctx)
		LoginInfo.Email = session.Get("Email")

		if LoginInfo.Email == nil {
			log.Println("ログインしていません")
			ctx.Redirect(http.StatusMovedPermanently, "/login")
			ctx.Abort()
		}else {
			ctx.Set("Email",LoginInfo.Email)
			ctx.Next()
		}
		log.Println("ログインチェック終わり")
	}
}

func GetMenu(ctx *gin.Context){
	Email,_ := ctx.Get("Email")
	log.Println("Emailです")
	log.Println(Email)

	//posts := getPostALL()
	db, err := gorm.Open("mysql", "root:secret@/postapp?charset=utf8&parseTime=True")
	if err != nil {
		panic("failed to connect database\n")
	}
	var post []models.Post
	db.Find(&post)
	log.Println(len(post))
	router := gin.Default()
	router.Static("/assets", "./assets")
	ctx.HTML(http.StatusOK, "index.html", gin.H{"posts": post,"Email":Email})
}

func GetSign(ctx *gin.Context){
	session := sessions.Default(ctx)
	LoginInfo.Email = session.Get("Email")
	ctx.Set("Email",LoginInfo.Email)
	Email,_ := ctx.Get("Email")
	log.Println("Emailです")
	log.Println(Email)

	//posts := getPostALL()
	db, err := gorm.Open("mysql", "root:secret@/postapp?charset=utf8&parseTime=True")
	if err != nil {
		panic("failed to connect database\n")
	}
	var post []models.Post
	db.Find(&post)
	log.Println(len(post))
	router := gin.Default()
	router.Static("/assets", "./assets")
	ctx.HTML(http.StatusOK, "index.html", gin.H{"posts": post,"Email":Email})
}

func PostLogout(ctx *gin.Context){
	log.Println("ログアウト処理")
	Logout(ctx)
	ctx.HTML(200, "login.html", gin.H{})
}
/*func getUser(email string, password string) []models.User{
	db, err := gorm.Open("mysql", "root:secret@/postapp?charset=utf8&parseTime=True")
	if err != nil {
		panic("failed to connect database\n")
	}
	return db.Where(&models.User{}{Email:email})
}*/

var LoginInfo sessioninfo.SessionInfo

func main() {
	//router.Router()
	db_init()
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.LoadHTMLGlob("templates/*.html")

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession",store))

	db, err := gorm.Open("mysql", "root:secret@/postapp?charset=utf8&parseTime=True")
	if err != nil {
		panic("failed to connect database\n")
	}

	db.AutoMigrate(&models.User{},&models.Post{})

	menu := router.Group("/menu")
	menu.Use(sessionCheck())
	{
		menu.Static("/assets", "./assets")
		menu.GET("/top",GetMenu)

		menu.POST("/create", func(ctx *gin.Context) {
			text := ctx.PostForm("text")
			email := ctx.PostForm("email")
			updatedat := time.Now().Format("2006-01-02")
			log.Println(email)
			createPost(text,email,updatedat)
			log.Println("ポストを投稿しました")
			ctx.Redirect(302, "/menu/top")
		})

	}

	router.GET("/signup",GetSign)

	router.POST("/register", func(ctx *gin.Context) {
		email := ctx.PostForm("email")
		password := ctx.PostForm("password")
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
		createUser(email, hashed)
		log.Println("ユーザー作成")
		session := sessions.Default(ctx)
		LoginInfo.Email = session.Get("Email")
		ctx.Set("email",LoginInfo.Email)
		ctx.Redirect(302, "/signup")
	})

	/*router.GET("/", func(ctx *gin.Context) {
		posts := getPostALL()
		ctx.HTML(200, "index.html", gin.H{"posts": posts})
	})*/

	router.GET("/register", func(ctx *gin.Context) {
		ctx.HTML(200, "register.html", gin.H{})
	})

	router.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(200, "login.html", gin.H{})
	})

	/*router.POST("/create", func(ctx *gin.Context) {
		text := ctx.PostForm("text")
		email := ctx.PostForm("email")
		createPost(text,email)
		ctx.Redirect(302, "/")
	})*/

	/*router.POST("/register", func(ctx *gin.Context) {
		email := ctx.PostForm("email")
		password := ctx.PostForm("password")
		hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
		createUser(email, hashed)
		ctx.Redirect(302, "/menu/top")
	})*/

	/*router.POST("/login",func(ctx *gin.Context){
		dbPassword := getUser(ctx.PostForm("email")).Password
		//dbID := getUser(ctx.PostForm("email")).ID
		formPassword := ctx.PostForm("password")
		if err := bcrypt.CompareHashAndPassword(dbPassword,[]byte(formPassword)); err != nil{
			log.Println("ログインできませんでした")
			ctx.HTML(http.StatusBadRequest, "login.html", gin.H{"err": err})
			ctx.Abort()
		}else{
		log.Println("ログインできました")
		ctx.Redirect(302, "/")
		}
	})*/

	router.POST("/login",func(ctx *gin.Context) {
		log.Println("ログイン処理")
		Email := ctx.PostForm("email")
		dbPassword := getUser(ctx.PostForm("email")).Password
		//dbID := getUser(ctx.PostForm("email")).ID
		formPassword := ctx.PostForm("password")
		log.Println(formPassword)
		if err := bcrypt.CompareHashAndPassword(dbPassword, []byte(formPassword)); err != nil {
			log.Println("ログインできませんでした")
			ctx.HTML(http.StatusBadRequest, "login.html", gin.H{"err": err})
			ctx.Abort()
		} else {
			log.Println("ログインできました")
			Login(ctx, Email)
			ctx.Redirect(302, "/menu/top")
		}
	})
	router.GET("/delete/:id", func(ctx *gin.Context) {
		n := ctx.Param("id")
		id,err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		deletePost(id)
		ctx.Redirect(302,"/menu/top")
	})

	router.GET("/edit/:id",func(ctx *gin.Context) {

		n := ctx.Param("id")
		id,err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		Editpost := getPostOne(id)
		log.Println(id)
		log.Println(getPostOne(id))
		ctx.HTML(200, "edit.html", gin.H{"Editpost": Editpost})
	})

	router.POST("/edit/:id",func(ctx *gin.Context) {
		n := ctx.Param("id")
		id,err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		text := ctx.PostForm("text")
		updatePost(id,text)
		ctx.Redirect(302, "/menu/top")
	})
	router.GET("/logout",PostLogout)

	router.Run()
}

