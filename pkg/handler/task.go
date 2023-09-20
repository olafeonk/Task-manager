package handler

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"task_manager"
)

// @Summary Create task
// @Security BasicAuth
// @Tags tasks
// @Description create task
// @ID create-task
// @Accept  json
// @Produce  json
// @Param input body task_manager.CreateTaskInput true "task info"
// @Success 200 {integer} integer 1
// @Failure 400,401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/tasks [post]
func (h *Handler) createTask(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	slog.Info("start create task",
		"user", userId)

	var input task_manager.CreateTaskInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, "validation error requires name field")
		return
	}

	id, err := h.services.TaskManagerTask.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	slog.Info("create task success",
		"user", userId,
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
// @Security BasicAuth
// @Tags tasks
// @Description get all tasks
// @ID get-all-tasks
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllTasksResponse
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/tasks [get]
func (h *Handler) getAllTasks(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	slog.Info("start get all tasks",
		"user", userId)

	tasks, err := h.services.TaskManagerTask.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	slog.Info("get all tasks success",
		"user", userId)
	c.JSON(http.StatusOK, getAllTasksResponse{
		Data: tasks,
	})
}

// @Summary Get task By Id
// @Security BasicAuth
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
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	slog.Info("start get task by id",
		"user", userId)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	task, err := h.services.TaskManagerTask.GetById(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	slog.Info("get task by id success",
		"user", userId)
	c.JSON(http.StatusOK, task)
}

// @Summary Update task
// @Security BasicAuth
// @Tags tasks
// @Description update task
// @ID update-task
// @Accept  json
// @Produce  json
// @Param id path int true "Task ID"
// @Param input body task_manager.UpdateTaskInput true "task info"
// @Success 200 {string} 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/tasks/{id} [put]
func (h *Handler) updateTask(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	slog.Info("start update task",
		"user", userId)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input task_manager.UpdateTaskInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.TaskManagerTask.Update(userId, id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	slog.Info("update task success",
		"user", userId)
	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// @Summary Delete task
// @Security BasicAuth
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
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	slog.Info("start delete tasks",
		"user", userId)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.TaskManagerTask.Delete(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	slog.Info("delete task success",
		"user", userId)
	c.JSON(http.StatusOK, statusResponse{
		Status: "ok",
	})
}
