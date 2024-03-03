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

	router.GET("/api/v2/settings", controllers.GetSettingsMongo)

	userMongo := router.Group("/api/v2/user")
	userMongo.POST("/create", controllers.CreateUserMongo)
	userMongo.POST("/login", controllers.LoginMongo)
	userMongo.Use(AuthMiddlewareMongo())
	userMongo.DELETE("/login", controllers.LogoutMongo)
	userMongo.GET("/data", controllers.GetUserDataMongo)

	buildingMongo := router.Group("/api/v2/building")
	buildingMongo.GET("/types", controllers.GetBuildingsTypesMongo)
	buildingMongo.GET("/blueprints", controllers.GetBlueprintsMongo)
	buildingMongo.POST("/get", controllers.GetBuildingsMongo)
	buildingMongo.Use(AuthMiddlewareMongo())
	buildingMongo.POST("/construct", controllers.ConstructBuildingMongo)
	buildingMongo.GET("/my", controllers.GetMyBuildingsMongo)
	buildingMongo.POST("/start_work", controllers.StartWorkMongo)
	buildingMongo.POST("/hiring", controllers.SetHiringMongo)
	buildingMongo.DELETE("/destroy", controllers.DestroyBuildingMongo)

	dataMongo := router.Group("/api/v2/data")
	dataMongo.GET("/users_by_prefix", controllers.GetUserNamesByPrefixMongo)
	dataMongo.GET("/evolution/prices", controllers.GetEvolutionPricesMongo)

	mapCellMongo := router.Group("/api/v2/map")
	mapCellMongo.GET("/cell_owners", controllers.GetCellOwnersMongo)
	mapCellMongo.GET("/", controllers.GetMapMongo)
	mapCellMongo.GET("/all_land_lords", controllers.GetAllLandLordsMongo)
	mapCellMongo.Use(AuthMiddlewareMongo())
	mapCellMongo.POST("/buy_land", controllers.BuyLandMongo)
	mapCellMongo.GET("/my", controllers.GetMyLandMongo)

	resourceMongo := router.Group("/api/v2/resource")
	resourceMongo.GET("/types", controllers.GetResourceTypesMongo)
	resourceMongo.Use(AuthMiddlewareMongo())
	resourceMongo.GET("/my", controllers.GetMyResourcesMongo)
	resourceMongo.POST("/move", controllers.ResourceMoveMongo)
	resourceMongo.GET("/my_logistics", controllers.GetMyLogisticsMongo)

	storageMongo := router.Group("/api/v2/storage")
	storageMongo.Use(AuthMiddlewareMongo())
	storageMongo.GET("/my", controllers.GetMyStoragesMongo)

	storeMongo := router.Group("/api/v2/store")
	storeMongo.GET("/goods/get", controllers.GetStoreGoodsMongo)
	storeMongo.Use(AuthMiddlewareMongo())
	storeMongo.POST("/goods/set", controllers.SetStoreGoodsMongo)

	marketMongo := router.Group("/api/v2/market")
	marketMongo.GET("/order/get", controllers.GetOrdersMongo)
	marketMongo.Use(AuthMiddlewareMongo())
	marketMongo.POST("/order/create", controllers.CreateOrderMongo)
	marketMongo.GET("/order/my", controllers.GetMyOrdersMongo)
	marketMongo.DELETE("/order/close", controllers.CloseMyOrderMongo)
	marketMongo.POST("/order/execute", controllers.ExecuteOrderMongo)

	return router
}
