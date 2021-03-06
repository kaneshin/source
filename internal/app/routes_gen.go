// Code generated by gen-routes.
// routes/api_v1.tml
// routes/routes.tml
// routes/web.tml
// DO NOT EDIT

package main

import (
	"github.com/gin-gonic/gin"

	"github.com/gophergala2016/source/core/foundation"
	"github.com/gophergala2016/source/core/net/context/accessor"

	"github.com/gophergala2016/source/internal/controllers/api"
	"github.com/gophergala2016/source/internal/controllers/web"
)

func router() *foundation.Engine {
	router := foundation.Default()

	// GET /v1/user/:id
	// api.APIUserController.GetUser
	router.GET("/v1/user/:id", func(c *gin.Context) {
		controller := &api.APIUserController{}
		ctx := foundation.NewContext(c)
		accessor.SetAction(ctx, "APIUserController.GetUser")
		controller.SetContext(ctx)
		controller.GetUser(c.Param("id"))
	})

	// GET /v1/me
	// api.APIMeController.GetMe
	router.GET("/v1/me", func(c *gin.Context) {
		controller := &api.APIMeController{}
		ctx := foundation.NewContext(c)
		accessor.SetAction(ctx, "APIMeController.GetMe")
		controller.SetContext(ctx)
		controller.GetMe()
	})

	// POST /v1/me
	// api.APIMeController.LoginMe
	router.POST("/v1/me", func(c *gin.Context) {
		controller := &api.APIMeController{}
		ctx := foundation.NewContext(c)
		accessor.SetAction(ctx, "APIMeController.LoginMe")
		controller.SetContext(ctx)
		controller.LoginMe()
	})

	// GET /v1/item/:id
	// api.APIItemController.GetItem
	router.GET("/v1/item/:id", func(c *gin.Context) {
		controller := &api.APIItemController{}
		ctx := foundation.NewContext(c)
		accessor.SetAction(ctx, "APIItemController.GetItem")
		controller.SetContext(ctx)
		controller.GetItem(c.Param("id"))
	})

	// GET /v1/items
	// api.APIItemController.GetItemList
	router.GET("/v1/items", func(c *gin.Context) {
		controller := &api.APIItemController{}
		ctx := foundation.NewContext(c)
		accessor.SetAction(ctx, "APIItemController.GetItemList")
		controller.SetContext(ctx)
		controller.GetItemList()
	})

	// GET /v1/favorites
	// api.APIItemController.GetItemFavoriteList
	router.GET("/v1/favorites", func(c *gin.Context) {
		controller := &api.APIItemController{}
		ctx := foundation.NewContext(c)
		accessor.SetAction(ctx, "APIItemController.GetItemFavoriteList")
		controller.SetContext(ctx)
		controller.GetItemFavoriteList()
	})

	// POST /v1/favorite/:id
	// api.APIItemController.CreateItemFavorite
	router.POST("/v1/favorite/:id", func(c *gin.Context) {
		controller := &api.APIItemController{}
		ctx := foundation.NewContext(c)
		accessor.SetAction(ctx, "APIItemController.CreateItemFavorite")
		controller.SetContext(ctx)
		controller.CreateItemFavorite(c.Param("id"))
	})

	// POST /v1/item
	// api.APIItemController.CreateItem
	router.POST("/v1/item", func(c *gin.Context) {
		controller := &api.APIItemController{}
		ctx := foundation.NewContext(c)
		accessor.SetAction(ctx, "APIItemController.CreateItem")
		controller.SetContext(ctx)
		controller.CreateItem()
	})

	// GET /v1/tags
	// api.APITagController.GetTagList
	router.GET("/v1/tags", func(c *gin.Context) {
		controller := &api.APITagController{}
		ctx := foundation.NewContext(c)
		accessor.SetAction(ctx, "APITagController.GetTagList")
		controller.SetContext(ctx)
		controller.GetTagList()
	})

	// POST /v1/tag
	// api.APITagController.CreateTag
	router.POST("/v1/tag", func(c *gin.Context) {
		controller := &api.APITagController{}
		ctx := foundation.NewContext(c)
		accessor.SetAction(ctx, "APITagController.CreateTag")
		controller.SetContext(ctx)
		controller.CreateTag()
	})

	// GET /
	// web.WebIndexController.Index
	router.GET("/", func(c *gin.Context) {
		controller := &web.WebIndexController{}
		ctx := foundation.NewContext(c)
		accessor.SetAction(ctx, "WebIndexController.Index")
		controller.SetContext(ctx)
		controller.Index()
	})

	return router
}
