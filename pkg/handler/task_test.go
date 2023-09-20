package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"task_manager"
	"task_manager/pkg/service"
	mockService "task_manager/pkg/service/mocks"
	"testing"
	"time"
)

func TestHandler_createTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockBehavior func(r *mockService.MockTaskManagerTask, userId int, task task_manager.CreateTaskInput)

	testTable := []struct {
		name                 string
		userCtx              string
		inputBody            string
		userId               int
		task                 task_manager.CreateTaskInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			userCtx:   "userId",
			inputBody: `{"name": "Task 1"}`,
			userId:    1,
			task:      task_manager.CreateTaskInput{Name: "Task 1"},
			mockBehavior: func(r *mockService.MockTaskManagerTask, userId int, task task_manager.CreateTaskInput) {
				r.EXPECT().Create(userId, task).Return(3, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"id\":3}",
		},
		{
			name:                 "No required field",
			userCtx:              "userId",
			inputBody:            `{"nam": "Task 1"}`,
			userId:               1,
			task:                 task_manager.CreateTaskInput{},
			mockBehavior:         func(r *mockService.MockTaskManagerTask, userId int, task task_manager.CreateTaskInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: "{\"message\":\"validation error requires name field\"}",
		},
		{
			name:      "Fail in service layer",
			userCtx:   "userId",
			inputBody: `{"name": "Task 1"}`,
			userId:    1,
			task:      task_manager.CreateTaskInput{Name: "Task 1"},
			mockBehavior: func(r *mockService.MockTaskManagerTask, userId int, task task_manager.CreateTaskInput) {
				r.EXPECT().Create(userId, task).Return(0, errors.New("failed"))
			},
			expectedStatusCode:   400,
			expectedResponseBody: "{\"message\":\"failed\"}",
		},
		{
			name:                 "Miss user id in context",
			userCtx:              "context",
			inputBody:            `{"name": "Task 1"}`,
			userId:               1,
			mockBehavior:         func(r *mockService.MockTaskManagerTask, userId int, task task_manager.CreateTaskInput) {},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"message\":\"user id not found\"}",
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockTaskManagerTask(c)
			test.mockBehavior(repo, test.userId, test.task)

			services := &service.Service{TaskManagerTask: repo}
			handler := Handler{services}

			r := gin.New()
			r.POST("/tasks", func(c *gin.Context) {
				c.Set(test.userCtx, test.userId)
			}, handler.createTask)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/tasks",
				bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)

		})
	}
}

func TestHandler_getAllTasks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockBehavior func(r *mockService.MockTaskManagerTask, userId int)

	testTable := []struct {
		name                 string
		userId               int
		userCtx              string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "Ok",
			userCtx: "userId",
			userId:  1,
			mockBehavior: func(r *mockService.MockTaskManagerTask, userId int) {
				r.EXPECT().GetAll(userId).Return([]task_manager.Task{
					{
						Id:        1,
						CreatedAt: time.Time{},
						UpdatedAt: time.Time{},
						UserId:    1,
						Name:      "Task1",
						StatusEnd: "Start",
						EndTask:   nil,
					},
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"data\":[{\"id\":1,\"created_at\":\"0001-01-01T00:00:00Z\",\"updated_at\":\"0001-01-01T00:00:00Z\",\"user_id\":1,\"name\":\"Task1\",\"status_end\":\"Start\",\"end_task_at\":null}]}",
		},
		{
			name:    "Fail in service layer",
			userCtx: "userId",
			userId:  1,
			mockBehavior: func(r *mockService.MockTaskManagerTask, userId int) {
				r.EXPECT().GetAll(userId).Return(nil, errors.New("failed"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"message\":\"failed\"}",
		},
		{
			name:                 "Miss user id in context",
			userCtx:              "context",
			userId:               1,
			mockBehavior:         func(r *mockService.MockTaskManagerTask, userId int) {},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"message\":\"user id not found\"}",
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockTaskManagerTask(c)
			test.mockBehavior(repo, test.userId)

			services := &service.Service{TaskManagerTask: repo}
			handler := Handler{services}

			r := gin.New()
			r.GET("/tasks", func(c *gin.Context) {
				c.Set(test.userCtx, test.userId)
			}, handler.getAllTasks)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/tasks", nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)

		})
	}
}

func TestHandler_getTaskById(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockBehavior func(r *mockService.MockTaskManagerTask, userId int, taskId int)

	testTable := []struct {
		name                 string
		userId               int
		taskId               int
		param                string
		userCtx              string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "Ok",
			userCtx: "userId",
			userId:  1,
			taskId:  2,
			param:   "2",
			mockBehavior: func(r *mockService.MockTaskManagerTask, userId, taskId int) {
				r.EXPECT().GetById(userId, taskId).Return(task_manager.Task{
					Id:        2,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
					UserId:    1,
					Name:      "Task1",
					StatusEnd: "Start",
					EndTask:   nil,
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"id\":2,\"created_at\":\"0001-01-01T00:00:00Z\",\"updated_at\":\"0001-01-01T00:00:00Z\",\"user_id\":1,\"name\":\"Task1\",\"status_end\":\"Start\",\"end_task_at\":null}",
		},
		{
			name:    "Fail in service layer",
			userCtx: "userId",
			userId:  1,
			taskId:  2,
			param:   "2",
			mockBehavior: func(r *mockService.MockTaskManagerTask, userId, taskId int) {
				r.EXPECT().GetById(userId, taskId).Return(task_manager.Task{}, errors.New("failed"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"message\":\"failed\"}",
		},
		{
			name:                 "Miss user id in context",
			userCtx:              "context",
			userId:               1,
			taskId:               2,
			param:                "2",
			mockBehavior:         func(r *mockService.MockTaskManagerTask, userId, taskId int) {},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"message\":\"user id not found\"}",
		},
		{
			name:                 "Missing param",
			userCtx:              "userId",
			userId:               1,
			taskId:               2,
			param:                "miss",
			mockBehavior:         func(r *mockService.MockTaskManagerTask, userId, taskId int) {},
			expectedStatusCode:   400,
			expectedResponseBody: "{\"message\":\"invalid id param\"}",
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockTaskManagerTask(c)
			test.mockBehavior(repo, test.userId, test.taskId)

			services := &service.Service{TaskManagerTask: repo}
			handler := Handler{services}

			r := gin.New()
			r.GET("/tasks/:id", func(c *gin.Context) {
				c.Set(test.userCtx, test.userId)
			}, handler.getTaskById)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/tasks/%s", test.param), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)

		})
	}
}

func TestHandler_updateTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockBehavior func(r *mockService.MockTaskManagerTask, userId int, taskId int, input task_manager.UpdateTaskInput)

	testTable := []struct {
		name                 string
		userCtx              string
		inputBody            string
		userId               int
		taskId               int
		param                string
		task                 task_manager.UpdateTaskInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			userCtx:   "userId",
			inputBody: `{"name": "Task 2"}`,
			userId:    1,
			taskId:    1,
			param:     "1",
			task: task_manager.UpdateTaskInput{
				Name: "Task 2",
			},
			mockBehavior: func(r *mockService.MockTaskManagerTask, userId int, taskId int, input task_manager.UpdateTaskInput) {
				r.EXPECT().Update(userId, taskId, input).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"status\":\"ok\"}",
		},
		{
			name:      "Fail in service layer",
			userCtx:   "userId",
			inputBody: `{"name": "Task 2"}`,
			userId:    1,
			taskId:    1,
			param:     "1",
			task:      task_manager.UpdateTaskInput{Name: "Task 2"},
			mockBehavior: func(r *mockService.MockTaskManagerTask, userId int, taskId int, input task_manager.UpdateTaskInput) {
				r.EXPECT().Update(userId, taskId, input).Return(errors.New("failed"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"message\":\"failed\"}",
		},
		{
			name:                 "Miss user id in context",
			userCtx:              "context",
			inputBody:            `{"name": "Task 1"}`,
			userId:               1,
			taskId:               1,
			param:                "1",
			mockBehavior:         func(r *mockService.MockTaskManagerTask, userId int, taskId int, input task_manager.UpdateTaskInput) {},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"message\":\"user id not found\"}",
		},
		{
			name:                 "Missing param",
			userCtx:              "userId",
			inputBody:            `{"name": "Task 1"}`,
			taskId:               1,
			param:                "miss",
			userId:               1,
			mockBehavior:         func(r *mockService.MockTaskManagerTask, userId int, taskId int, input task_manager.UpdateTaskInput) {},
			expectedStatusCode:   400,
			expectedResponseBody: "{\"message\":\"invalid id param\"}",
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockTaskManagerTask(c)
			test.mockBehavior(repo, test.userId, test.taskId, test.task)

			services := &service.Service{TaskManagerTask: repo}
			handler := Handler{services}

			r := gin.New()
			r.PUT("/tasks/:id", func(c *gin.Context) {
				c.Set(test.userCtx, test.userId)
			}, handler.updateTask)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", fmt.Sprintf("/tasks/%s", test.param),
				bytes.NewBufferString(test.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)

		})
	}
}

func TestHandler_deleteTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockBehavior func(r *mockService.MockTaskManagerTask, userId int, taskId int)

	testTable := []struct {
		name                 string
		userId               int
		taskId               int
		param                string
		userCtx              string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "Ok",
			userCtx: "userId",
			userId:  1,
			taskId:  2,
			param:   "2",
			mockBehavior: func(r *mockService.MockTaskManagerTask, userId, taskId int) {
				r.EXPECT().Delete(userId, taskId).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "{\"status\":\"ok\"}",
		},
		{
			name:    "Fail in service layer",
			userCtx: "userId",
			userId:  1,
			taskId:  2,
			param:   "2",
			mockBehavior: func(r *mockService.MockTaskManagerTask, userId, taskId int) {
				r.EXPECT().Delete(userId, taskId).Return(errors.New("failed"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"message\":\"failed\"}",
		},
		{
			name:                 "Miss user id in context",
			userCtx:              "context",
			userId:               1,
			taskId:               2,
			param:                "2",
			mockBehavior:         func(r *mockService.MockTaskManagerTask, userId, taskId int) {},
			expectedStatusCode:   500,
			expectedResponseBody: "{\"message\":\"user id not found\"}",
		},
		{
			name:                 "Missing param",
			userCtx:              "userId",
			userId:               1,
			taskId:               2,
			param:                "miss",
			mockBehavior:         func(r *mockService.MockTaskManagerTask, userId, taskId int) {},
			expectedStatusCode:   400,
			expectedResponseBody: "{\"message\":\"invalid id param\"}",
		},
	}
	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockTaskManagerTask(c)
			test.mockBehavior(repo, test.userId, test.taskId)

			services := &service.Service{TaskManagerTask: repo}
			handler := Handler{services}

			r := gin.New()
			r.DELETE("/tasks/:id", func(c *gin.Context) {
				c.Set(test.userCtx, test.userId)
			}, handler.deleteTask)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/tasks/%s", test.param), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)

		})
	}
}
