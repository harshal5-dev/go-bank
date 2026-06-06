package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/go-bank/internal/db/sqlc"
	"github.com/go-bank/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string             `json:"username"`
	FullName          string             `json:"full_name"`
	Email             string             `json:"email"`
	PasswordChangedAt pgtype.Timestamptz `json:"password_changed_at"`
	CreatedAt         pgtype.Timestamptz `json:"created_at"`
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		FullName:     req.FullName,
		HashPassword: hashedPassword,
		Email:        req.Email,
		Username:     req.Username,
	}

	createdUser, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			switch pqErr.Code {
			case "23505":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := newUserResponse(createdUser)

	ctx.JSON(http.StatusCreated, res)
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
	}

	err = utils.CheckPassword(req.Password, user.HashPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
	}

	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
		if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
	}

	rsp := loginUserResponse {
		AccessToken: accessToken,
		User: newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}
