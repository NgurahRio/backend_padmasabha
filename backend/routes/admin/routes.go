package routes

import (
	"backend/controllers/admin/auth"
	"backend/controllers/admin/category"
	"backend/controllers/admin/destination"
	"backend/controllers/admin/event"
	"backend/controllers/admin/facility"

	packagess "backend/controllers/admin/packages"
	"backend/controllers/admin/review"
	"backend/controllers/admin/sos"
	"backend/controllers/admin/subcategory"
	"backend/controllers/admin/subpackage"
	"backend/controllers/admin/user"
	"backend/middleware"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	frontendURL := strings.TrimRight(os.Getenv("FRONTEND_URL"), "/")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/admin/login", auth.LoginAdmin)
	r.POST("/admin/forgot-password", auth.ForgotPassword)
	r.POST("/admin/verify-reset-otp", auth.VerifyResetOTP)
	r.POST("/admin/reset-password", auth.ResetPassword)
	r.GET("/admin/session", middleware.AdminSessionAuth(), auth.CurrentSession)
	r.POST("/admin/logout", middleware.AdminSessionAuth(), auth.LogoutAdmin)
	r.POST("/admin/logout-all", middleware.AdminSessionAuth(), auth.LogoutAllAdminSessions)

	admin := r.Group("/admin")
	admin.Use(
		middleware.AdminSessionAuth(),
	)

	{
		// user routes
		admin.GET("/users", user.GetAllUsers)
		admin.GET("/users/:id", user.GetUserByID)
		admin.DELETE("/users/:id", user.DeleteUser)

		// SubPackage routes
		admin.POST("/subpackage", subpackage.CreateSubpackage)
		admin.GET("/subpackage", subpackage.GetAllSubpackages)
		admin.GET("/subpackage/:id", subpackage.GetSubpackageByID)
		admin.PUT("/subpackage/:id", subpackage.UpdateSubpackage)
		admin.DELETE("/subpackage/:id", subpackage.DeleteSubpackage)

		// SubCategory routes
		admin.POST("/subcategory", subcategory.CreateSubcategory)
		admin.GET("/subcategory", subcategory.GetAllSubcategories)
		admin.GET("/subcategory/:id", subcategory.GetSubcategoryByID)
		admin.PUT("/subcategory/:id", subcategory.UpdateSubcategory)
		admin.DELETE("/subcategory/:id", subcategory.DeleteSubcategory)

		// Category routes
		admin.POST("/category", category.CreateCategory)
		admin.GET("/category", category.GetAllCategories)
		admin.GET("/category/:id", category.GetCategoryByID)
		admin.PUT("/category/:id", category.UpdateCategory)
		admin.DELETE("/category/:id", category.DeleteCategory)

		// Destination routes
		admin.POST("/destination", destination.CreateDestination)
		admin.GET("/destination", destination.GetAllDestinations)
		admin.GET("/destination/:id", destination.GetDestinationByID)
		admin.PUT("/destination/:id", destination.UpdateDestination)
		admin.DELETE("/destination/:id", destination.DeleteDestination)

		//review routes
		admin.GET("/review", review.GetAllReview)
		admin.GET("/review/:id", review.GetReviewByID)
		admin.DELETE("/review/:id", review.DeleteReview)

		// Event routes
		admin.POST("/event", event.CreateEvent)
		admin.GET("/event", event.GetAllEvents)
		admin.GET("/event/:id", event.GetEventByID)
		admin.PUT("/event/:id", event.UpdateEvent)
		admin.DELETE("/event/:id", event.DeleteEvent)

		// SOS routes
		admin.POST("/sos", sos.CreateSOS)
		admin.GET("/sos", sos.GetAllSOS)
		admin.GET("/sos/:id", sos.GetSOSByID)
		admin.PUT("/sos/:id", sos.UpdateSOS)
		admin.DELETE("/sos/:id", sos.DeleteSOS)

		// Facility routes
		admin.DELETE("/facility/:id", facility.DeleteFacility)
		admin.PUT("/facility/:id", facility.UpdateFacility)
		admin.POST("/facility", facility.CreateFacility)
		admin.GET("/facility", facility.GetAllFacility)
		admin.GET("/facility/:id", facility.GetFacilityByID)

		// Packages routes
		admin.POST("/packages", packagess.CreatePackages)
		admin.GET("/packages", packagess.GetAllPackages)
		admin.GET("/packages/:destinationId", packagess.GetPackageByDestinationID)
		admin.PUT("/packages/:destinationId", packagess.UpdatePackages)
		admin.DELETE("/packages/:destinationId", packagess.DeletePackages)

	}

	return r
}
