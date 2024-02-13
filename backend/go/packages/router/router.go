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

	router.GET("/api/v1/settings", controllers.GetSettings)

	user := router.Group("/api/v1/user")
	user.POST("/create", controllers.CreateUser)
	user.POST("/login", controllers.Login)
	// Secure aria
	user.Use(AuthMiddleware())
	user.DELETE("/login", controllers.Logout)
	user.GET("/data", controllers.GetUserData)

	building := router.Group("/api/v1/building")
	building.GET("/types", controllers.GetBuildingsTypes)
	building.POST("/get", controllers.GetBuildings) // GET
	building.GET("/blueprints", controllers.GetBlueprints)
	building.Use(AuthMiddleware())
	building.POST("/construct", controllers.CreateBuilding)
	building.GET("/my", controllers.GetMyBuildings)
	building.POST("/start_work", controllers.StartWork)
	building.POST("/hiring", controllers.SetHiring)
	building.DELETE("/destroy", controllers.DestroyBuilding)

	mapCell := router.Group("/api/v1/map")
	mapCell.GET("/", controllers.GetMap)
	mapCell.GET("/cell_owners", controllers.GetCellOwners)
	mapCell.GET("/all_land_lords", controllers.GetAllLandLords)
	mapCell.Use(AuthMiddleware())
	mapCell.GET("/my", controllers.GetMyLand)
	mapCell.POST("/buy_land", controllers.BuyLand)

	resource := router.Group("/api/v1/resource")
	resource.GET("/types", controllers.GetResourceTypes)
	resource.Use(AuthMiddleware())
	resource.GET("/my", controllers.GetMyResources)
	resource.POST("/move", controllers.ResourceMove)
	resource.GET("/my_logistics", controllers.GetMyLogistics)

	storage := router.Group("/api/v1/storage")
	storage.Use(AuthMiddleware())
	storage.GET("/my", controllers.GetMyStorages)

	market := router.Group("/api/v1/market")
	market.GET("/order/get", controllers.GetOrders)
	market.Use(AuthMiddleware())
	market.POST("/order/create", controllers.CreateOrder)
	market.GET("/order/my", controllers.GetMyOrders)
	market.DELETE("/order/close", controllers.CloseMyOrder)
	market.POST("/order/execute", controllers.ExecuteOrder)

	store := router.Group("/api/v1/store")
	store.GET("/goods/get", controllers.GetStoreGoods)
	store.Use(AuthMiddleware())
	store.POST("/goods/set", controllers.SetStoreGoods)

	data := router.Group("/api/v1/data")
	data.GET("/users_by_prefix", controllers.GetUserNamesByPrefix)
	data.GET("/evolution/prices", controllers.GetEvolutionPrices)

	// MONGO

	router.GET("/api/v2/settings", controllers.GetSettingsMongo)

	userMongo := router.Group("/api/v2/user")
	userMongo.POST("/create", controllers.CreateUserMongo)
	userMongo.POST("/login", controllers.LoginMongo)
	userMongo.Use(AuthMiddlewareMongo())
	userMongo.DELETE("/login", controllers.LogoutMongo)
	userMongo.GET("/data", controllers.GetUserDataMongo)

	buildingMongo := router.Group("/api/v2/building")
	buildingMongo.GET("/blueprints", controllers.GetBlueprintsMongo)

	dataMongo := router.Group("/api/v2/data")
	dataMongo.GET("/users_by_prefix", controllers.GetUserNamesByPrefixMongo)

	return router
}
