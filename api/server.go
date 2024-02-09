package api

import (
	"fmt"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authGuardedRoutes := router.Group("/").Use(authMidleware(server.tokenMaker))

	authGuardedRoutes.POST("/accounts", server.createAccount)
	authGuardedRoutes.GET("/accounts", server.listAccount)
	authGuardedRoutes.GET("/accounts/:id", server.getAccount)

	authGuardedRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func OK(ctx *gin.Context, obj any) {
	ctx.JSON(http.StatusOK, obj)
}

func BadRequest(ctx *gin.Context, obj any) {
	ctx.JSON(http.StatusBadRequest, obj)
}

func InternalServerError(ctx *gin.Context, obj any) {
	ctx.JSON(http.StatusInternalServerError, obj)
}

func NotFound(ctx *gin.Context, obj any) {
	ctx.JSON(http.StatusNotFound, obj)
}

func Forbidden(ctx *gin.Context, obj any) {
	ctx.JSON(http.StatusForbidden, obj)
}

func Unauthorized(ctx *gin.Context, obj any) {
	ctx.JSON(http.StatusUnauthorized, obj)
}
