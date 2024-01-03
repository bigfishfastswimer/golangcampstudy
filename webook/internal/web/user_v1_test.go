package web

import (
	"bytes"
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	svcmocks "gitee.com/geekbang/basic-go/webook/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHello(t *testing.T) {
	t.Log("hello, test")
}

func TestUserProfile(t *testing.T) {
	birthday, err := time.Parse(time.DateOnly, "2000-11-11")
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name       string
		mock       func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		reqBuilder func(t *testing.T) *http.Request
		expectCode int
		expectBody string
	}{
		{
			name: "test profile",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().FindById(gomock.Any(), int64(1)).Return(domain.User{
					Nickname: "Fisher",
					Email:    "12334@gmail.com",
					AboutMe:  "this is a test user",
					Birthday: birthday,
				}, nil)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodGet, "/users/profile", bytes.NewReader([]byte("")))
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			expectCode: http.StatusOK,
			expectBody: "{\"nickname\":\"Fisher\",\"email\":\"12334@gmail.com\",\"aboutMe\":\"this is a test user\",\"birthday\":\"2000-11-11\"}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// construct handler
			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)

			// prepare server, register routes
			server := gin.Default()
			server.Use(func(ctx *gin.Context) {
				ctx.Set("user", UserClaims{
					Uid: 1,
				})
			})
			hdl.RegisterRoutes(server)

			// prepare request and http recorder
			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			//start mock server
			server.ServeHTTP(recorder, req)

			// execute test
			assert.Equal(t, tc.expectCode, recorder.Code)
			assert.Equal(t, tc.expectBody, recorder.Body.String())

		})
	}

}
