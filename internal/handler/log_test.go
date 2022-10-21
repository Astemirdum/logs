package handler

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	service_mocks "github.com/Astemirdum/logs/internal/handler/mocks"
	"github.com/Astemirdum/logs/models"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestLogHandler_CreateLog(t *testing.T) {

	type input struct {
		body string
		req  *models.CreateLogRequest
	}
	type response struct {
		expectedCode int
		expectedBody string
	}
	type mockBehavior func(r *service_mocks.MockService, inp *models.CreateLogRequest)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		input        input
		response     response
		wantErr      bool
	}{
		{
			name: "ok",
			mockBehavior: func(r *service_mocks.MockService, inp *models.CreateLogRequest) {
				var id int64 = 2
				r.EXPECT().CreateLog(context.Background(), inp.Raw).Return(id, nil)
			},
			input: input{
				body: `{"raw": "test raw"}`,
				req: &models.CreateLogRequest{
					Raw: "test raw",
				},
			},
			response: response{
				expectedCode: http.StatusOK,
				expectedBody: `{"id":2}`,
			},
			wantErr: false,
		},
		{
			name:         "err. empty raw",
			mockBehavior: func(r *service_mocks.MockService, inp *models.CreateLogRequest) {},
			input: input{
				body: `{"raw": ""}`,
				req: &models.CreateLogRequest{
					Raw: "",
				},
			},
			response: response{
				expectedCode: http.StatusBadRequest,
				expectedBody: ``,
			},
			wantErr: true,
		},
		{
			name: "err. internal",
			mockBehavior: func(r *service_mocks.MockService, inp *models.CreateLogRequest) {
				r.EXPECT().CreateLog(context.Background(), inp.Raw).Return(int64(0), errors.New("db internal"))
			},
			input: input{
				body: `{"raw": "dddd"}`,
				req: &models.CreateLogRequest{
					Raw: "dddd",
				},
			},
			response: response{
				expectedCode: http.StatusInternalServerError,
				expectedBody: ``,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			svc := service_mocks.NewMockService(c)
			tt.mockBehavior(svc, tt.input.req)

			log := zap.NewExample().Named("test")
			h := &logHandler{
				svc: svc,
				log: log,
			}

			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}
			r := httptest.NewRequest(
				http.MethodPost, "/api/v1/logs", strings.NewReader(tt.input.body))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			w := httptest.NewRecorder()

			ctx := e.NewContext(r, w)
			err := h.CreateLog(ctx)
			if !tt.wantErr {
				require.NoError(t, err)
				require.Equal(t, tt.response.expectedCode, w.Code)
				require.Equal(t, tt.response.expectedBody, strings.Trim(w.Body.String(), "\n"))
			} else {
				require.Error(t, err)
				er := &echo.HTTPError{}
				if errors.As(err, &er) {
					require.Equal(t, tt.response.expectedCode, er.Code)
					// require.Equal(t, tt.response.expectedBody, er.Message)
				}
			}

		})

	}
}
