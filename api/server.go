package api

import (
	"example.com/simple_bank/db/sqlc"
	"example.com/simple_bank/db/util"
	"example.com/simple_bank/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store      *sqlc.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     util.Config
}

func NewServer(config util.Config, store *sqlc.Store) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.SymetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		store:      store,
		tokenMaker: maker,
		config:     config,
	}

	if validator, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validator.RegisterValidation("validateCurrency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {

	router := gin.Default()

	router.POST("/createUser", server.createUser)
	router.POST("/userLogin", server.UserLogin)

	authRoutes := router.Group("/").Use(authMiddleWare(server.tokenMaker))

	authRoutes.POST("/createAccount", server.createAccount)
	authRoutes.GET("/getAccountById/:id", server.getAccountById)
	authRoutes.POST("/getAllAccounts", server.listAccounts)

	authRoutes.POST("/transfer", server.transfer)

	server.router = router
}

func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}
