package server

import "backend/config"

// Init server
func Init() error {
	config := config.GetConfig()
	r := NewRouter()

	return r.Run(config.ServerAddress)
}
