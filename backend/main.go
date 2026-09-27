package main

import (
	"backend/config"
	"backend/middleware"
	adminroutes "backend/routes/admin"
	userroutes "backend/routes/user"
	"fmt"
)

func main() {
	config.ConnectDB()
	middleware.StartAdminSessionCleanup()
	r := adminroutes.SetupRouter()

	// Daftarkan user routes ke engine yang sama
	userroutes.SetupUserRoutes(r)

	fmt.Println("🚀 Backend server berjalan di http://localhost:8080 🚀")
	r.Run(":8080")

}
