package server

type AccrualSystem struct {
	*Server
}

func NewAccrualSystem(server *Server) *AccrualSystem {
	return &AccrualSystem{
		Server: server,
	}
}

func (s *AccrualSystem) SetRoutes() {
	// TODO
	// gApi := s.router.Group("/api")
	// gApi.GET("/orders/:number", s.ordersGetHandler)
	// gApi.POST("/orders", s.ordersPostHandler)
	// gApi.POST("/goods", s.goodsPostHandler)
}
