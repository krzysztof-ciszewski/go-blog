package bootstrap

import (
	dependency_injection "main/internal/Infrastructure/DependencyInjection"
	auth "main/internal/UserInterface/Api/Handler/Auth"
	post "main/internal/UserInterface/Api/Handler/Post"
	middleware "main/internal/UserInterface/Api/Middleware"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

func BootstrapGin(container dependency_injection.Container) *gin.Engine {
	r := gin.Default()

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

	store.MaxAge(int(12 * time.Hour / time.Second))
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false
	store.Options.SameSite = http.SameSiteLaxMode

	gothic.Store = store

	authGroup := r.Group("/auth")
	apiGroup := r.Group("/api/v1", middleware.RequireAuth())

	{
		githubPrivder := github.New(
			os.Getenv("GITHUB_CLIENT_ID"),
			os.Getenv("GITHUB_CLIENT_SECRET"),
			os.Getenv("API_URL")+"/auth/github/callback",
			"user:email",
		)

		goth.UseProviders(githubPrivder)

		authGroup.GET("/:provider/callback", func(ctx *gin.Context) {
			auth.OauthCallback(ctx, container.CommandBus, container.QueryBus)
		})
		authGroup.GET("/:provider", func(ctx *gin.Context) {
			auth.OauthInitial(ctx, container.QueryBus)
		})
		authGroup.GET("/logout/:provider", auth.OauthLogout)
	}

	{
		apiGroup.GET("/posts", func(ctx *gin.Context) {
			post.ListPosts(ctx, container.QueryBus)
		})
		apiGroup.GET("/posts/:id", func(ctx *gin.Context) {
			post.GetPostById(ctx, container.QueryBus)
		})
		apiGroup.POST("/posts", func(ctx *gin.Context) {
			post.CreatePost(ctx, container.CommandBus)
		})
		apiGroup.DELETE("/posts/:id", func(ctx *gin.Context) {
			post.DeletePost(ctx, container.CommandBus)
		})
	}

	return r
}
