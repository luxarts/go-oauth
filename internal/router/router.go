package router

import (
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"go-oauth/internal/controller"
	"go-oauth/internal/defines"
	"go-oauth/internal/repository"
	"go-oauth/internal/service"
)

func New() *gin.Engine {
	r := gin.Default()

	mapRoutes(r)

	return r
}

func mapRoutes(r *gin.Engine) {
	rc := resty.New()

	stateRepo := repository.NewStateRepository()
	repo := repository.NewOAuthRepository(rc)

	svc := service.NewOAuthService(repo, stateRepo)

	ctrl := controller.NewOAuthController(svc)

	r.GET(defines.EndpointGetLogin, ctrl.Login)
	r.GET(defines.EndpointCallback, ctrl.Callback)
}
