package routes

import (
	"github.com/gin-gonic/gin"
	"luke544187758/goods-web/api/brands"
	"luke544187758/goods-web/middlewares"
)

func InitBrandRouter(r *gin.RouterGroup) {
	b := r.Group("/brands")
	{
		b.GET("/list", brands.BrandList)
		b.DELETE("/:id", middlewares.JWTAuthMiddleware(), middlewares.IsAdminAuth(), brands.DeleteBrand)
		b.POST("/create", middlewares.JWTAuthMiddleware(), middlewares.IsAdminAuth(), brands.NewBrand)
		b.PUT("/:id", middlewares.JWTAuthMiddleware(), middlewares.IsAdminAuth(), brands.UpdateBrand)
	}
	cb := r.Group("/categorybrands")
	{
		cb.GET("/list", brands.CategoryBrandList)
		cb.DELETE("/:id", middlewares.JWTAuthMiddleware(), middlewares.IsAdminAuth(), brands.DeleteCategoryBrand)
		cb.POST("/create", middlewares.JWTAuthMiddleware(), middlewares.IsAdminAuth(), brands.NewCategoryBrand)
		cb.PUT("/:id", middlewares.JWTAuthMiddleware(), middlewares.IsAdminAuth(), brands.UpdateCategoryBrand)
		cb.GET("/:id", brands.GetCategoryBrandList)
	}
}
