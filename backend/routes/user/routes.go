package routes

import (
	"backend/controllers/user/activity"
	"backend/controllers/user/auth"
	"backend/controllers/user/category"
	"backend/controllers/user/event"
	"backend/controllers/user/packages"
	"backend/controllers/user/restaurant"
	"backend/controllers/user/subcategory"
	"backend/controllers/user/subpackage"
	"backend/controllers/user/villa"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(r *gin.Engine) {

	// Auth
	r.POST("/user/register", auth.RegisterUser)
	r.POST("/user/login", auth.LoginUser)

	r.GET("/villas", villa.GetAllVillasUser)
	r.GET("/villas/:id", villa.GetVillaByIDUser)

	r.GET("/activities", activity.GetAllActivitysUser)
	r.GET("/activities/:id", activity.GetActivityByIDUser)

	r.GET("/restaurants", restaurant.GetAllRestaurantsUser)
	r.GET("/restaurants/:id", restaurant.GetRestaurantByIDUser)

	r.GET("/events", event.GetAllEventsUser)
	r.GET("/events/:id", event.GetEventByIDUser)

	r.GET("/categories", category.GetAllCategoriesUser)
	r.GET("/categories/:id", category.GetCategoryByIDUser)

	r.GET("/subcategories", subcategory.GetAllSubcategoriesUser)
	r.GET("/subcategories/:id", subcategory.GetSubcategoryByIDUser)

	r.GET("/packages", packages.GetAllPackagesUser)
	r.GET("/packages/:villaId", packages.GetPackageByVillaIDUser)

	r.GET("/subpackage", subpackage.GetAllSubpackagesUser)
	r.GET("/subpackage/:id", subpackage.GetSubpackageByIDUser)

	user := r.Group("/user",
		middleware.JWTAuthMiddleware(),
		middleware.RoleOnly(1),
	)

	user.GET("/profile", auth.GetProfile)
	user.PUT("/profile", auth.UpdateProfile)
	user.DELETE("/profile", auth.DeleteProfile)

	user.POST("/logout", auth.LogoutUser)
}
