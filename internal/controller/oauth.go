package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/luxarts/go-oauth/internal/service"
	"net/http"
)

type OAuthController interface {
	Login(ctx *gin.Context)
	Callback(ctx *gin.Context)
}

type oauthController struct {
	svc service.OAuthService
}

func NewOAuthController(svc service.OAuthService) OAuthController {
	return &oauthController{
		svc: svc,
	}
}

func (ctrl *oauthController) Login(ctx *gin.Context) {
	url, err := ctrl.svc.GetLoginURL()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Header("Location", *url)
	ctx.Status(http.StatusFound)
}
func (ctrl *oauthController) Callback(ctx *gin.Context) {
	state, ok := ctx.GetQuery("state")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing query param 'state'"})
		return
	}
	code, ok := ctx.GetQuery("code")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing query param 'code'"})
		return
	}

	token, err := ctrl.svc.Callback(state, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, *token)
}
