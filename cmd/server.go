package main

import (
	"context"
	post_command "main/internal/Application/Command/Post"
	user_command "main/internal/Application/Command/User"
	query "main/internal/Application/Query"
	dependency_injection "main/internal/Infrastructure/DependencyInjection"
	middleware "main/internal/UserInterface/Api/Middleware"
	request "main/internal/UserInterface/Api/Request"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
)

func main() {
	container := dependency_injection.GetContainer()
	commandBus := container.CommandBus

	r := gin.Default()

	store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	store.MaxAge(int(12 * time.Hour / time.Second))
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false
	store.Options.SameSite = http.SameSiteLaxMode

	gothic.Store = store

	auth := r.Group("/auth")
	api := r.Group("/api/v1", middleware.RequireAuth())

	{
		githubPrivder := github.New(
			os.Getenv("GITHUB_CLIENT_ID"),
			os.Getenv("GITHUB_CLIENT_SECRET"),
			"http://localhost:8080/auth/github/callback",
			"user:email",
		)

		goth.UseProviders(githubPrivder)

		auth.GET("/:provider/callback", func(ctx *gin.Context) {
			query := ctx.Request.URL.Query()
			query.Add("provider", ctx.Param("provider"))
			ctx.Request.URL.RawQuery = query.Encode()

			gothUser, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error authenticating user",
					"error":   err.Error(),
				})
				return
			}

			session, err := gothic.Store.New(ctx.Request, os.Getenv("SESSION_NAME"))
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error stroing user session",
					"error":   err.Error(),
				})
				return
			}

			id, err := uuid.NewRandom()
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error generating user ID",
					"error":   err.Error(),
				})
				return
			}

			container.CommandBus.Send(context.Background(), user_command.NewCreateUserCommand(
				id,
				gothUser.Email,
				"",
				gothUser.Provider,
				gothUser.Name,
				gothUser.FirstName,
				gothUser.LastName,
				gothUser.UserID,
				gothUser.AvatarURL,
				gothUser.AccessToken,
				gothUser.AccessTokenSecret,
				gothUser.RefreshToken,
				gothUser.ExpiresAt,
				gothUser.IDToken,
			))

			session.Values["user"] = gothUser

			if err = session.Save(ctx.Request, ctx.Writer); err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error saving user session",
					"error":   err.Error(),
				})
				return
			}

			ctx.Redirect(http.StatusTemporaryRedirect, os.Getenv("CLIENT_URL"))
		})

		auth.GET("/:provider", func(ctx *gin.Context) {
			query := ctx.Request.URL.Query()
			query.Add("provider", ctx.Param("provider"))
			ctx.Request.URL.RawQuery = query.Encode()

			gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
		})

		auth.GET("/logout/:provider", func(ctx *gin.Context) {
			gothic.Logout(ctx.Writer, ctx.Request)
			ctx.Writer.Header().Set("Location", "/")
			ctx.Writer.WriteHeader(http.StatusTemporaryRedirect)
		})
	}

	{
		api.GET("/posts", func(ctx *gin.Context) {
			text := ctx.Query("text")
			author := ctx.Query("author")

			var result any
			var err error

			if text != "" {
				q := query.NewFindAllByTextQuery(text)
				result, err = container.QueryBus.Execute(context.Background(), q)
			} else if author != "" {
				q := query.NewFindAllByAuthorQuery(author)
				result, err = container.QueryBus.Execute(context.Background(), q)
			} else {
				q := query.NewFindAllQuery()
				result, err = container.QueryBus.Execute(context.Background(), q)
			}

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			ctx.JSON(http.StatusOK, result)
		})

		api.GET("/posts/:id", func(ctx *gin.Context) {
			id := ctx.Param("id")

			parsedUUID, err := uuid.Parse(id)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
				return
			}

			q := query.NewGetPostQuery(parsedUUID)
			post, err := container.QueryBus.Execute(context.Background(), q)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

				return
			}

			ctx.JSON(http.StatusOK, post)
		})

		api.GET("/posts/slug/:slug", func(ctx *gin.Context) {
			slug := ctx.Param("slug")

			q := query.NewFindBySlugQuery(slug)
			post, err := container.QueryBus.Execute(context.Background(), q)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			ctx.JSON(http.StatusOK, post)
		})

		api.POST("/posts", func(ctx *gin.Context) {
			var request request.CreatePostRequest

			if err := ctx.ShouldBindJSON(&request); err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

				return
			}

			command := post_command.NewCreatePostCommand(
				uuid.MustParse(request.Id),
				request.Slug,
				request.Title,
				request.Content,
				request.Author,
			)

			commandBus.Send(context.Background(), command)

			ctx.JSON(http.StatusAccepted, gin.H{"message": "Post created"})
		})

		api.DELETE("/posts/:id", func(ctx *gin.Context) {
			id := ctx.Param("id")

			command := post_command.NewDeletePostCommand(uuid.MustParse(id))
			commandBus.Send(context.Background(), command)

			ctx.JSON(http.StatusAccepted, gin.H{"message": "Post deleted"})
		})
	}
	r.Run()
}
