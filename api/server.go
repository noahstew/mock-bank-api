package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/noahstew/bank-api/db/sqlc"
)

// serves http requests
type Server struct {
	store db.Store
	router *gin.Engine
}

// creates new http server and sets up routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/transfers", server.createTransfer)


	server.router = router 

	return server
}

// runs http server on an address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}