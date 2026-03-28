package handlers

import (
	"apiserver/internal/apiserver/model"
	"apiserver/store"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetHandler godoc
// @Summary      Получить подписку по ID
// @Description  Возвращает информацию о подписке по её идентификатору
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      int  true  "ID подписки"
// @Success      200  {object}  model.Subscription
// @Failure      400  {object}  map[string]interface{} "Неверный формат ID"
// @Failure      500  {object}  map[string]interface{} "Ошибка сервера"
// @Router       /{id} [get]
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
			return
		}
		sub.StartDate = sub.StartDate[3:]

		loger.Info(fmt.Sprintf("Get subscription from DB [%d %s %d %s %s]", sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate))
		c.JSON(200, *sub)
	}
}

// CreateHandler godoc
// @Summary      Создать новую подписку
// @Description  Добавляет новую подписку в систему
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription body      model.Subscription  true  "Данные подписки"
// @Success      200         {object}  model.Subscription
// @Failure      400         {object}  map[string]interface{} "Неверный формат данных"
// @Failure      500         {object}  map[string]interface{} "Ошибка сервера"
// @Router       / [post]
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

// DeleteHandler godoc
// @Summary      Удалить подписку
// @Description  Удаляет подписку по ID
// @Tags         subscriptions
// @Produce      json
// @Param        id   path      int  true  "ID подписки"
// @Success      200  {integer} int  "ID удаленной подписки"
// @Failure      400  {object}  map[string]interface{} "Неверный формат ID"
// @Failure      500  {object}  map[string]interface{} "Ошибка сервера"
// @Router       /{id} [delete]
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

// UpdateHandler godoc
// @Summary      Обновить подписку
// @Description  Обновляет данные существующей подписки
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id            path      int                 true  "ID подписки"
// @Param        subscription  body      model.Subscription  true  "Обновленные данные подписки"
// @Success      200           {object}  model.Subscription
// @Failure      400           {object}  map[string]interface{} "Неверный формат данных или ID"
// @Failure      500           {object}  map[string]interface{} "Ошибка сервера"
// @Router       /{id} [put]
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

// GetList godoc
// @Summary      Получить список всех подписок
// @Description  Возвращает список всех подписок в системе
// @Tags         subscriptions
// @Produce      json
// @Success      200  {array}   model.Subscription
// @Failure      500  {object}  map[string]interface{} "Ошибка сервера"
// @Router       /list [get]
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

// rangeParam представляет параметры для запроса диапазона
type rangeParam struct {
	Service_name string `json:"service_name" example:"Netflix" description:"Название сервиса"`
	User_id      string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba" description:"ID пользователя"`
	Start_date   string `json:"start_date" example:"01-2024" description:"Дата начала периода (MM-YYYY)"`
	End_date     string `json:"end_date" example:"04-2024" description:"Дата окончания периода (MM-YYYY)"`
}

// GetRangeHandler godoc
// @Summary      Получить сумму по диапазону
// @Description  Возвращает общую сумму подписок за указанный период
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        params  body      rangeParam  true  "Параметры поиска"
// @Success      200     {integer} integer     "Общая сумма"
// @Failure      400     {object}  map[string]interface{} "Неверный формат данных"
// @Failure      500     {object}  map[string]interface{} "Ошибка сервера"
// @Router       /range [post]
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
