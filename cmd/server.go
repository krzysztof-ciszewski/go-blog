package main

import (
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

func main() {
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

		authGroup.GET("/:provider/callback", auth.OauthCallback)
		authGroup.GET("/:provider", auth.OauthInitial)
		authGroup.GET("/logout/:provider", auth.OauthLogout)
	}

	{
		apiGroup.GET("/posts", post.ListPosts)
		apiGroup.GET("/posts/:id", post.GetPostById)
		apiGroup.GET("/posts/slug/:slug", post.GetPostBySlug)
		apiGroup.POST("/posts", post.CreatePost)
		apiGroup.DELETE("/posts/:id", post.DeletePost)
	}
	r.Run()
}
