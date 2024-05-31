package bootstrap

import (
	"x.com/api/routes"
	"x.com/api/server"
)

func StartupApi() {
	routes := routes.RegisterRoutes()

	serverConfig := server.NewServerConfig("8080")
	httpServer := server.NewServer(serverConfig, routes)
	server.RunServer(httpServer)
}
