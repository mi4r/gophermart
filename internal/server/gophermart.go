package server

import echoSwagger "github.com/swaggo/echo-swagger"

type Gophermart struct {
	*Server
}

func NewGophermart(server *Server) *Gophermart {
	return &Gophermart{
		Server: server,
	}
}

func (s *Gophermart) SetRoutes() {
	// swagger
	s.router.GET("/swagger/*", echoSwagger.WrapHandler)
	s.router.GET("/ping", s.pingHandler)

	// Группа users
	gUsers := s.router.Group("/api/user")
	gUsers.POST("/register", s.userRegisterHandler)
	gUsers.POST("/login", s.userLoginHandler)
	gUsers.POST("/orders", s.userPostOrdersHandler)
	gUsers.GET("/orders", s.userGetOrdersHandler)

}
