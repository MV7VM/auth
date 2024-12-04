package server

func (s *Server) createController() {
	commonGroup := s.serv.Group("/api")

	publicGroup := commonGroup.Group("")

	publicGroup.POST("auth/login", s.GetUserToken)
	//publicGroup.POST("auth/refresh")

	authGroup := commonGroup.Group("") //todo middleware
	authGroup.GET("get-time")
	authGroup.GET("admin")
}
