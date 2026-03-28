package handlers

import (
	"apiserver/internal/apiserver/model"
	"apiserver/store"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetHandler(loger *slog.Logger, store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		loger.Debug(fmt.Sprintf("Parsed id from param: %d", id))
		var sub *model.Subscription
		if sub, err = store.User().Get(id); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		}
		sub.StartDate = sub.StartDate[3:]

		loger.Info(fmt.Sprintf("Get subscription from DB [%d %s %d %s %s]", sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate))
		c.JSON(200, *sub)
	}
}

func CreateHandler(loger *slog.Logger, store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {

		var newSub model.Subscription

		if err := c.ShouldBindJSON(&newSub); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		loger.Info(fmt.Sprintf("Parsed from JSON [%s %d %s %s]", newSub.ServiceName, newSub.Price, newSub.UserID, newSub.StartDate))

		newSub.StartDate = newSub.StartDate[:3] + "01-" + newSub.StartDate[3:]

		if err := store.User().Create(&newSub); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		newSub.StartDate = newSub.StartDate[:3] + newSub.StartDate[6:]
		loger.Info(fmt.Sprintf("Sent Response: 200 [%d %s %d %s %s]", newSub.ID, newSub.ServiceName, newSub.Price, newSub.UserID, newSub.StartDate))
		c.JSON(200, newSub)

	}
}

func DeleteHandler(loger *slog.Logger, store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		loger.Debug(fmt.Sprintf("Parsed id from param: %d", id))

		if err := store.User().Delete(id); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		loger.Info(fmt.Sprintf("Sent Response: 200 [%d]", id))

		c.JSON(200, id)
	}
}

func UpdateHandler(loger *slog.Logger, store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		var newSub model.Subscription

		if err := c.ShouldBindJSON(&newSub); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		newSub.ID = id
		loger.Info(fmt.Sprintf("Parsed from JSON [%s %d %s %s]", newSub.ServiceName, newSub.Price, newSub.UserID, newSub.StartDate))

		newSub.StartDate = newSub.StartDate[:3] + "01-" + newSub.StartDate[3:]

		loger.Info(fmt.Sprintf("Parsed id from param: %d", id))
		if err := store.User().Update(&newSub); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		newSub.StartDate = newSub.StartDate[:3] + newSub.StartDate[6:]
		c.JSON(200, newSub)
	}

}

func GetList(loger *slog.Logger, store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var list []model.Subscription
		list, err := store.User().List()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		loger.Info("Get list from DB")
		for idx := range list {
			list[idx].StartDate = list[idx].StartDate[3:]
		}
		c.JSON(200, list)
	}
}

type rangeParam struct {
	Service_name string `json:"service_name"`
	User_id      string `json:"user_id"`
	Start_date   string `json:"start_date"`
	End_date     string `json:"end_date"`
}

func GetRangeHandler(loger *slog.Logger, store *store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		var param rangeParam
		if err := c.ShouldBindJSON(&param); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		loger.Info(fmt.Sprintf("Parsed params from JSON[%s %s %s %s]", param.Service_name, param.User_id, param.Start_date, param.End_date))

		param.Start_date = param.Start_date[:3] + "01-" + param.Start_date[3:]
		param.End_date = param.End_date[:3] + "01-" + param.End_date[3:]

		price, err := store.User().SummaryByIdAndPeriod(param.User_id, param.Service_name, param.Start_date, param.End_date)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		loger.Info(fmt.Sprintf("Get price: %d", price))

		c.JSON(200, price)
	}
}
