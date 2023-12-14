package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"task_manager"
	"time"
)

// @Summary Create task
// @Tags tasks
// @Description create task
// @ID create-task
// @Accept  json
// @Produce  json
// @Param input body task_manager.CreateTaskInputModer true "task info"
// @Success 200 {integer} integer 1
// @Failure 400,401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/tasks [post]
func (h *Handler) createTask(c *gin.Context) {
	var err error
	slog.Info("start create task")
	jsonDataBytes, err := ioutil.ReadAll(c.Request.Body)
	slog.Info("data", jsonDataBytes)
	unquoteJsonData := fmt.Sprintf("%s", jsonDataBytes)
	slog.Info("unquote", unquoteJsonData)
	unquoteJsonData = strings.ReplaceAll(unquoteJsonData, "+", " ")
	unquoteJsonData = strings.ReplaceAll(unquoteJsonData, "%3A", ":")
	slog.Info("unquote", unquoteJsonData)
	fields := strings.Split(unquoteJsonData, "&")
	var input task_manager.CreateTaskInputModeration
	var inputTh task_manager.CreateTaskInput

	for _, field := range fields {
		items := strings.Split(field, "=")
		key := items[0]
		value := items[1]
		if key == "telegram_id" {
			inputTh.TelegramId = value
		}
		if key == "text" {
			inputTh.Text = value
		}
		if key == "start_time" {
			inputTh.StartTime, err = time.Parse(time.DateTime, value)
		}
	}

	slog.Info("task", input)
	id, err := h.services.TaskManagerTask.Create(inputTh)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	slog.Info("create task success",
		"task", input,
		"task", id,
	)
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getAllTasksResponse struct {
	Data []task_manager.Task `json:"data"`
}

// @Summary Get All Tasks
// @Tags tasks
// @Description get all tasks
// @ID get-all-tasks
// @Accept  json
// @Produce  json
// @Param id path int true "telegram ID"
// @Success 200 {object} getAllTasksResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/telegram/{id} [get]
func (h *Handler) getTasksByTelegramId(c *gin.Context) {
	slog.Info("start get all tasks")

	telegramId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	tasks, err := h.services.TaskManagerTask.GetAll(telegramId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	slog.Info("get all tasks success")
	c.JSON(http.StatusOK, getAllTasksResponse{
		Data: tasks,
	})
}

// @Summary Get task By Id
// @Tags tasks
// @Description get task by id
// @ID get-task-by-id
// @Accept  json
// @Produce  json
// @Param id path int true "Task ID"
// @Success 200 {object} task_manager.Task
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/tasks/{id} [get]
func (h *Handler) getTaskById(c *gin.Context) {
	slog.Info("start get task by id")

	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	task, err := h.services.TaskManagerTask.GetById(taskId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	slog.Info("get task by id success",
		"task", taskId)
	c.JSON(http.StatusOK, task)
}

// @Summary Delete task
// @Tags tasks
// @Description delete task
// @ID delete-task
// @Accept  json
// @Produce  json
// @Param id path int true "Task ID"
// @Success 200 {string} ok
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/tasks/{id} [delete]
func (h *Handler) deleteTask(c *gin.Context) {
	slog.Info("start delete tasks")

	taskId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.TaskManagerTask.Delete(taskId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	slog.Info("delete task success",
		"task_id", taskId)
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
