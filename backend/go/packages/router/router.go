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

	router.GET("/api/v2/settings", controllers.GetSettingsMongo) // DONE

	userMongo := router.Group("/api/v2/user") // USER DONE
	userMongo.POST("/create", controllers.CreateUserMongo)
	userMongo.POST("/login", controllers.LoginMongo)
	userMongo.Use(AuthMiddlewareMongo())
	userMongo.DELETE("/login", controllers.LogoutMongo)
	userMongo.GET("/data", controllers.GetUserDataMongo)

	buildingMongo := router.Group("/api/v2/building") // buildings DONE
	buildingMongo.GET("/types", controllers.GetBuildingsTypesMongo)
	buildingMongo.GET("/blueprints", controllers.GetBlueprintsMongo)
	buildingMongo.POST("/get", controllers.GetBuildingsMongo) // GET
	buildingMongo.Use(AuthMiddlewareMongo())
	buildingMongo.POST("/construct", controllers.ConstructBuildingMongo)
	buildingMongo.GET("/my", controllers.GetMyBuildingsMongo)
	buildingMongo.POST("/start_work", controllers.StartWorkMongo)
	buildingMongo.POST("/hiring", controllers.SetHiringMongo)
	buildingMongo.DELETE("/destroy", controllers.DestroyBuildingMongo)

	dataMongo := router.Group("/api/v2/data")
	dataMongo.GET("/users_by_prefix", controllers.GetUserNamesByPrefixMongo)

	mapCellMongo := router.Group("/api/v2/map") // MAP DONE
	mapCellMongo.GET("/cell_owners", controllers.GetCellOwnersMongo)
	mapCellMongo.GET("/", controllers.GetMapMongo)
	mapCellMongo.GET("/all_land_lords", controllers.GetAllLandLordsMongo)
	mapCellMongo.Use(AuthMiddlewareMongo())
	mapCellMongo.POST("/buy_land", controllers.BuyLandMongo)
	mapCellMongo.GET("/my", controllers.GetMyLandMongo)

	resourceMongo := router.Group("/api/v2/resource") // RES DONE
	resourceMongo.GET("/types", controllers.GetResourceTypesMongo)
	resourceMongo.Use(AuthMiddlewareMongo())
	resourceMongo.GET("/my", controllers.GetMyResourcesMongo)
	resourceMongo.POST("/move", controllers.ResourceMoveMongo)
	resourceMongo.GET("/my_logistics", controllers.GetMyLogisticsMongo)

	storageMongo := router.Group("/api/v2/storage") // Storage DONE
	storageMongo.Use(AuthMiddlewareMongo())
	storageMongo.GET("/my", controllers.GetMyStoragesMongo)

	storeMongo := router.Group("/api/v2/store") // Store DONE
	storeMongo.GET("/goods/get", controllers.GetStoreGoodsMongo)
	storeMongo.Use(AuthMiddlewareMongo())
	storeMongo.POST("/goods/set", controllers.SetStoreGoodsMongo)

	marketMongo := router.Group("/api/v2/market")
	marketMongo.GET("/order/get", controllers.GetOrdersMongo)
	marketMongo.Use(AuthMiddlewareMongo())
	marketMongo.POST("/order/create", controllers.CreateOrderMongo)
	marketMongo.GET("/order/my", controllers.GetMyOrdersMongo)
	marketMongo.POST("/order/execute", controllers.ExecuteOrderMongo)

	return router
}
