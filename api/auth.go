package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	hasher "github.com/ulerdogan/pickaxe/utils/hasher"
)

func (r *ginServer) LoginAdmin(ctx *gin.Context) {
	var req AuthRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.Username != "admin" {
		ctx.JSON(http.StatusUnauthorized, "invalid user")
		return
	}

	pwd, err := r.store.GetHashedIndexerPwd(context.Background())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = hasher.CheckPassword(req.Password, pwd)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := r.token.CreateToken(
		req.Username,
		r.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := AuthResponse{
		AccessToken: accessToken,
		ExpiresIn:   time.Now().Add(r.config.AccessTokenDuration).Unix()}

	ctx.JSON(http.StatusOK, rsp)
}
