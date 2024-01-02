package web

import (
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	svcmocks "gitee.com/geekbang/basic-go/webook/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
				birthday, err := time.Parse(time.DateOnly, "2000-11-11")
				if err != nil {
					t.Fatal(err)
				}
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().FindById(gomock.Any(), int64(1)).Return(domain.User{
					Nickname: "Fisher",
					Email:    "12334@gmail.com",
					AboutMe:  "this is a test user",
					Birthday: birthday,
				}, nil)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				userClaims := UserClaims{
					RegisteredClaims: jwt.RegisteredClaims{},
					Uid:              1,
					UserAgent:        "test",
				}
				ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
				req, err := http.NewRequest(http.MethodGet, "/users/profile", nil)
				ctx.Request = req
				ctx.Set("user", userClaims)
				req.Header.Set("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			expectCode: 200,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)

			server := gin.Default()
			hdl.RegisterRoutes(server)

			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, req)
			assert.Equal(t, tc.expectCode, recorder.Code)
			assert.Equal(t, tc.expectBody, recorder.Body.String())

		})
	}

}
