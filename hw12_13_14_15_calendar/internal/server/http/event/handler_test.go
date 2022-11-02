package event_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/event/command"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/event/query"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/mock"
	domainevent "github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/server/http/event"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/pkg/api/openapi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_eventHandler_GetEvents(t *testing.T) {
	tests := map[string]struct {
		getApp   func(controller *gomock.Controller) app.Application
		params   openapi.GetEventsParams
		wantCode int
	}{
		"given period is 'day', when ErrStorage returned then should return code 500": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					GetDayEvents(gomock.Any()).
					Return(nil, domainevent.NewErrStorage("storage error"))
				return a
			},
			params:   openapi.GetEventsParams{Period: openapi.Day},
			wantCode: http.StatusInternalServerError,
		},
		"given period is 'day', when no error returned then should return code 200": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					GetDayEvents(gomock.Any()).
					Return(&query.GetDayEventsResponse{}, nil)
				return a
			},
			params:   openapi.GetEventsParams{Period: openapi.Day},
			wantCode: http.StatusOK,
		},
		"given period is 'week', when ErrStorage returned then should return code 500": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					GetWeekEvents(gomock.Any()).
					Return(nil, domainevent.NewErrStorage("storage error"))
				return a
			},
			params:   openapi.GetEventsParams{Period: openapi.Week},
			wantCode: http.StatusInternalServerError,
		},
		"given period is 'week', when no error returned then should return code 200": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					GetWeekEvents(gomock.Any()).
					Return(&query.GetWeekEventsResponse{}, nil)
				return a
			},
			params:   openapi.GetEventsParams{Period: openapi.Week},
			wantCode: http.StatusOK,
		},
		"given period is 'month', when ErrStorage returned then should return code 500": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					GetMonthEvents(gomock.Any()).
					Return(nil, domainevent.NewErrStorage("storage error"))
				return a
			},
			params:   openapi.GetEventsParams{Period: openapi.Month},
			wantCode: http.StatusInternalServerError,
		},
		"given period is 'month', when no error returned then should return code 200": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					GetMonthEvents(gomock.Any()).
					Return(&query.GetMonthEventsResponse{}, nil)
				return a
			},
			params:   openapi.GetEventsParams{Period: openapi.Month},
			wantCode: http.StatusOK,
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			responseRecorder := httptest.NewRecorder()
			h := event.NewEventHandler(tt.getApp(controller), mock.NewMockLogger(controller))

			h.GetEvents(responseRecorder, nil, tt.params)

			require.Equal(t, tt.wantCode, responseRecorder.Result().StatusCode) //nolint:bodyclose
		})
	}
}

func Test_eventHandler_PostEvents(t *testing.T) {
	requestBody := `{
		"title": "Event 1",
		"date": 1679894859,
		"duration": 30,
		"description": "Event 1 description",
  		"notifyBefore": 10
	}`
	tests := map[string]struct {
		getApp      func(controller *gomock.Controller) app.Application
		getLogger   func(controller *gomock.Controller) app.Logger
		requestBody string
		wantCode    int
	}{
		"when request body is invalid then should return code 400": {
			getApp: func(controller *gomock.Controller) app.Application {
				return mock.NewMockApplication(controller)
			},
			getLogger: func(controller *gomock.Controller) app.Logger {
				l := mock.NewMockLogger(controller)
				l.EXPECT().Error(gomock.Any())
				return l
			},
			requestBody: "{",
			wantCode:    http.StatusBadRequest,
		},
		"when request body is valid then should return code 200": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					CreateEvent(gomock.Any()).
					Return(&command.CreateEventResponse{Event: domainevent.Event{}}, nil)
				return a
			},
			getLogger: func(controller *gomock.Controller) app.Logger {
				return mock.NewMockLogger(controller)
			},
			requestBody: requestBody,
			wantCode:    http.StatusOK,
		},
		"when event date is busy then should return code 409": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					CreateEvent(gomock.Any()).
					Return(nil, domainevent.ErrDateBusy)
				return a
			},
			getLogger: func(controller *gomock.Controller) app.Logger {
				return mock.NewMockLogger(controller)
			},
			requestBody: requestBody,
			wantCode:    http.StatusConflict,
		},
		"when event validation was not passed then should return code 400": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					CreateEvent(gomock.Any()).
					Return(nil, command.ErrValidate{})
				return a
			},
			getLogger: func(controller *gomock.Controller) app.Logger {
				return mock.NewMockLogger(controller)
			},
			requestBody: requestBody,
			wantCode:    http.StatusBadRequest,
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			request := httptest.NewRequest(http.MethodPost, "https://example.com", bytes.NewBufferString(tt.requestBody))
			responseRecorder := httptest.NewRecorder()
			h := event.NewEventHandler(tt.getApp(controller), tt.getLogger(controller))

			h.PostEvents(responseRecorder, request)

			require.Equal(t, tt.wantCode, responseRecorder.Result().StatusCode) //nolint:bodyclose
		})
	}
}

func Test_eventHandler_PutEvents(t *testing.T) {
	requestBody := `{
		"id": "1",
		"title": "Event 1",
		"date": 1679894859,
		"duration": 30,
		"description": "Event 1 description",
  		"notifyBefore": 10
	}`
	tests := map[string]struct {
		getApp      func(controller *gomock.Controller) app.Application
		getLogger   func(controller *gomock.Controller) app.Logger
		requestBody string
		wantCode    int
	}{
		"when request body is invalid then should return code 400": {
			getApp: func(controller *gomock.Controller) app.Application {
				return mock.NewMockApplication(controller)
			},
			getLogger: func(controller *gomock.Controller) app.Logger {
				l := mock.NewMockLogger(controller)
				l.EXPECT().Error(gomock.Any())
				return l
			},
			requestBody: "{",
			wantCode:    http.StatusBadRequest,
		},
		"when request body is valid then should return code 200": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					UpdateEvent(gomock.Any()).
					Return(nil)
				return a
			},
			getLogger: func(controller *gomock.Controller) app.Logger {
				return mock.NewMockLogger(controller)
			},
			requestBody: requestBody,
			wantCode:    http.StatusOK,
		},
		"when event is not found then should return code 404": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					UpdateEvent(gomock.Any()).
					Return(domainevent.ErrEventNotFound)
				return a
			},
			getLogger: func(controller *gomock.Controller) app.Logger {
				return mock.NewMockLogger(controller)
			},
			requestBody: requestBody,
			wantCode:    http.StatusNotFound,
		},
		"when ErrStorage returned then should return code 500": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					UpdateEvent(gomock.Any()).
					Return(domainevent.ErrStorage{})
				return a
			},
			getLogger: func(controller *gomock.Controller) app.Logger {
				return mock.NewMockLogger(controller)
			},
			requestBody: requestBody,
			wantCode:    http.StatusInternalServerError,
		},
		"when event date is busy then should return code 409": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					UpdateEvent(gomock.Any()).
					Return(domainevent.ErrDateBusy)
				return a
			},
			getLogger: func(controller *gomock.Controller) app.Logger {
				return mock.NewMockLogger(controller)
			},
			requestBody: requestBody,
			wantCode:    http.StatusConflict,
		},
		"when event validation was not passed then should return code 400": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					UpdateEvent(gomock.Any()).
					Return(command.ErrValidate{})
				return a
			},
			getLogger: func(controller *gomock.Controller) app.Logger {
				return mock.NewMockLogger(controller)
			},
			requestBody: requestBody,
			wantCode:    http.StatusBadRequest,
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			request := httptest.NewRequest(http.MethodPost, "https://example.com", bytes.NewBufferString(tt.requestBody))
			responseRecorder := httptest.NewRecorder()
			h := event.NewEventHandler(tt.getApp(controller), tt.getLogger(controller))

			h.PutEvents(responseRecorder, request)

			require.Equal(t, tt.wantCode, responseRecorder.Result().StatusCode) //nolint:bodyclose
		})
	}
}

func Test_eventHandler_DeleteEventsId(t *testing.T) {
	tests := map[string]struct {
		getApp   func(controller *gomock.Controller) app.Application
		wantCode int
	}{
		"when event is not found then should return code 404": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					DeleteEvent(gomock.Any()).
					Return(domainevent.ErrEventNotFound)
				return a
			},
			wantCode: http.StatusNotFound,
		},
		"when ErrStorage returned then should return code 500": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					DeleteEvent(gomock.Any()).
					Return(domainevent.ErrStorage{})
				return a
			},
			wantCode: http.StatusInternalServerError,
		},
		"when no error returned then should return code 200": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					DeleteEvent(gomock.Any()).
					Return(nil)
				return a
			},
			wantCode: http.StatusOK,
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			request := httptest.NewRequest(http.MethodPost, "https://example.com", nil)
			responseRecorder := httptest.NewRecorder()
			h := event.NewEventHandler(tt.getApp(controller), mock.NewMockLogger(controller))

			h.DeleteEventsID(responseRecorder, request, "")

			require.Equal(t, tt.wantCode, responseRecorder.Result().StatusCode) //nolint:bodyclose
		})
	}
}
