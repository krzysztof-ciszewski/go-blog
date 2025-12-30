package bootstrap

import (
	dependency_injection "main/internal/Infrastructure/DependencyInjection"
	auth "main/internal/UserInterface/Api/Handler/Auth"
	post "main/internal/UserInterface/Api/Handler/Post"
	middleware "main/internal/UserInterface/Api/Middleware"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func BootstrapGin(container dependency_injection.Container) *gin.Engine {
	r := gin.Default()

	r.Use(otelgin.Middleware(container.Telemetry.GetServiceName()))
	r.Use(container.Telemetry.LogRequest())
	r.Use(container.Telemetry.MeterRequestDuration())
	r.Use(container.Telemetry.MeterRequestsInFlught())
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
			auth.OauthCallback(ctx, container.CommandBus, container.QueryBus, container.Telemetry)
		})
		authGroup.GET("/:provider", func(ctx *gin.Context) {
			auth.OauthInitial(ctx, container.QueryBus, container.Telemetry)
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
			post.CreatePost(ctx, container.CommandBus, container.QueryBus)
		})
		apiGroup.PUT("/posts/:id", func(ctx *gin.Context) {
			post.UpdatePost(ctx, container.CommandBus, container.QueryBus)
		})
		apiGroup.DELETE("/posts/:id", func(ctx *gin.Context) {
			post.DeletePost(ctx, container.CommandBus)
		})
	}

	return r
}
