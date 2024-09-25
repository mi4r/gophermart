package server

type Gophermart struct {
	*Server
}

func NewGophermart(server *Server) *Gophermart {
	return &Gophermart{
		Server: server,
	}
}

func (s *Gophermart) SetRoutes() {
	gUsers := s.router.Group("/api/user")
	gUsers.POST("/register", s.userRegisterHandler)
	gUsers.POST("/login", s.userLoginHandler)
	gUsers.POST("/orders", s.userPostOrdersHandler)
	gUsers.GET("/orders", s.userGetOrdersHandler)
	gUsers.GET("/balance", s.userGetBalance)
}
