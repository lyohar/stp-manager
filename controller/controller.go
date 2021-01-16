package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"stpCommon/model"
	"stpManager/repo"
)

type about struct {
	Name    string `json:"name" binding:"required"`
	Version string `json:"version" binding:"required"`
}

type Controller struct {
	repo *repo.PostgresRepo
}

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

func NewController(r *repo.PostgresRepo) *Controller {
	return &Controller{r}
}

func (c *Controller) InitRouter() *gin.Engine {
	router := gin.New()
	router.GET("/about", c.about)
	router.GET("/get-export", c.getExport)
	router.POST("/set-status", c.setExportStatus)
	return router
}

func (c *Controller) about(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, about{Name: "stp-manager", Version: "0.1"})
}


func (c *Controller) getExport(ctx *gin.Context) {

	export, err := c.repo.GetExport()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, export)
}

func (c *Controller) setExportStatus(ctx *gin.Context)  {
	var status model.ExportStatus
	if err := ctx.BindJSON(&status); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if err := c.repo.SetExportStatus(&status); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"table_schema": status.TableSchema,
		"table_name": status.TableName,
		"order_column_value": status.OrderColumnToValue,
	})
}
