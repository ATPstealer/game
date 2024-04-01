package router

import (
	"backend/packages/cfg"
	"backend/packages/controllers"
	"github.com/gin-gonic/gin"
)

func MakeRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	if !cfg.Config.Secure {
		router.Use(CORSMiddleware()) // disable CORS policy, but I need it on prod !!!
	}

	router.GET("/api/v2/settings", controllers.GetSettings)

	user := router.Group("/api/v2/user")
	user.POST("/create", controllers.CreateUser)
	user.POST("/login", controllers.Login)
	user.Use(AuthMiddleware())
	user.DELETE("/login", controllers.Logout)
	user.GET("/data", controllers.GetUserData)

	building := router.Group("/api/v2/building")
	building.GET("/types", controllers.GetBuildingsTypes)
	building.GET("/blueprints", controllers.GetBlueprints)
	building.POST("/get", controllers.GetBuildings)
	building.Use(AuthMiddleware())
	building.POST("/construct", controllers.ConstructBuilding)
	building.GET("/my", controllers.GetMyBuildings)
	building.POST("/start_work", controllers.StartWork)
	building.POST("/hiring", controllers.SetHiring)
	building.DELETE("/destroy", controllers.DestroyBuilding)

	data := router.Group("/api/v2/data")
	data.GET("/users_by_prefix", controllers.GetUserNamesByPrefix)
	data.GET("/evolution/prices", controllers.GetEvolutionPrices)

	mapCell := router.Group("/api/v2/map")
	mapCell.GET("/cell_owners", controllers.GetCellOwners)
	mapCell.GET("/", controllers.GetMap)
	mapCell.GET("/all_land_lords", controllers.GetAllLandLords)
	mapCell.Use(AuthMiddleware())
	mapCell.POST("/buy_land", controllers.BuyLand)
	mapCell.GET("/my", controllers.GetMyLand)

	resource := router.Group("/api/v2/resource")
	resource.GET("/types", controllers.GetResourceTypes)
	resource.Use(AuthMiddleware())
	resource.GET("/my", controllers.GetMyResources)
	resource.POST("/move", controllers.ResourceMove)
	resource.GET("/my_logistics", controllers.GetMyLogistics)

	storage := router.Group("/api/v2/storage")
	storage.Use(AuthMiddleware())
	storage.GET("/my", controllers.GetMyStorages)

	store := router.Group("/api/v2/store")
	store.Use(AuthMiddleware())
	store.POST("/goods/set", controllers.SetStoreGoods)

	market := router.Group("/api/v2/market")
	market.GET("/order/get", controllers.GetOrders)
	market.Use(AuthMiddleware())
	market.POST("/order/create", controllers.CreateOrder)
	market.GET("/order/my", controllers.GetMyOrders)
	market.DELETE("/order/close", controllers.CloseMyOrder)
	market.POST("/order/execute", controllers.ExecuteOrder)

	return router
}
