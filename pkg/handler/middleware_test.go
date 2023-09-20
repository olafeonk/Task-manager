package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http/httptest"
	"task_manager/pkg/service"
	mockService "task_manager/pkg/service/mocks"
	"testing"
)

func TestHandler_userIdentity(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockBehavior func(r *mockService.MockAuthorization, user, password string)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		user                 string
		password             string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
		authHeader           string
		expectedHeader       string
	}{
		{
			name:        "Ok",
			headerName:  "Authorization",
			headerValue: "Basic dXNlcjpwYXNzd29yZA==",
			user:        "user",
			password:    "password",
			mockBehavior: func(r *mockService.MockAuthorization, user, password string) {
				r.EXPECT().GetUserId(user, password).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: "1",
			authHeader:           "WWW-Authenticate",
			expectedHeader:       "",
		},
		{
			name:                 "Invalid Header Name",
			headerName:           "",
			headerValue:          "Basic dXNlcjpwYXNzd29yZA==",
			user:                 "user",
			password:             "password",
			mockBehavior:         func(r *mockService.MockAuthorization, user, password string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"empty auth header"}`,
			authHeader:           "WWW-Authenticate",
			expectedHeader:       "Basic",
		},
		{
			name:                 "Invalid Header Value",
			headerName:           "Authorization",
			headerValue:          "Basii dXNlcjpwYXNzd29yZA==",
			user:                 "user",
			password:             "password",
			mockBehavior:         func(r *mockService.MockAuthorization, user, password string) {},
			expectedStatusCode:   401,
			expectedResponseBody: `{"message":"invalid auth header"}`,
			authHeader:           "WWW-Authenticate",
			expectedHeader:       "Basic",
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {

			c := gomock.NewController(t)
			defer c.Finish()

			repo := mockService.NewMockAuthorization(c)
			test.mockBehavior(repo, test.user, test.password)

			services := &service.Service{Authorization: repo}
			handler := Handler{services}

			// Init Endpoint
			r := gin.New()
			r.GET("/identity", handler.userIdentity, func(c *gin.Context) {
				id, _ := c.Get(userCtx)
				c.String(200, "%d", id)
			})

			// Init Test Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/identity", nil)
			req.Header.Set(test.headerName, test.headerValue)

			r.ServeHTTP(w, req)

			res := w.Result()

			// Asserts
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
			assert.Equal(t, res.Header.Get(test.authHeader), test.expectedHeader)
		})
	}
}

func TestGetUserId(t *testing.T) {
	var getContext = func(id int) *gin.Context {
		ctx := &gin.Context{}
		ctx.Set(userCtx, id)
		return ctx
	}

	testTable := []struct {
		name       string
		ctx        *gin.Context
		id         int
		shouldFail bool
	}{
		{
			name: "Ok",
			ctx:  getContext(1),
			id:   1,
		},
		{
			ctx:        &gin.Context{},
			name:       "Empty",
			shouldFail: true,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			id, err := getUserId(test.ctx)
			if test.shouldFail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, id, test.id)
		})
	}
}
