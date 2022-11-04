// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event (interfaces: Storage)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	time "time"

	event "github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event"
	gomock "github.com/golang/mock/gomock"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockStorage) Create(arg0 context.Context, arg1 event.Event) (event.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(event.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockStorageMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockStorage)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockStorage) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockStorageMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStorage)(nil).Delete), arg0, arg1)
}

// DeleteEventsOlderThan mocks base method.
func (m *MockStorage) DeleteEventsOlderThan(arg0 context.Context, arg1 time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEventsOlderThan", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEventsOlderThan indicates an expected call of DeleteEventsOlderThan.
func (mr *MockStorageMockRecorder) DeleteEventsOlderThan(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEventsOlderThan", reflect.TypeOf((*MockStorage)(nil).DeleteEventsOlderThan), arg0, arg1)
}

// GetDayEvents mocks base method.
func (m *MockStorage) GetDayEvents(arg0 context.Context, arg1 time.Time) ([]event.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDayEvents", arg0, arg1)
	ret0, _ := ret[0].([]event.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDayEvents indicates an expected call of GetDayEvents.
func (mr *MockStorageMockRecorder) GetDayEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDayEvents", reflect.TypeOf((*MockStorage)(nil).GetDayEvents), arg0, arg1)
}

// GetEventsNotifyBetween mocks base method.
func (m *MockStorage) GetEventsNotifyBetween(arg0 context.Context, arg1, arg2 time.Time) ([]event.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEventsNotifyBetween", arg0, arg1, arg2)
	ret0, _ := ret[0].([]event.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEventsNotifyBetween indicates an expected call of GetEventsNotifyBetween.
func (mr *MockStorageMockRecorder) GetEventsNotifyBetween(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEventsNotifyBetween", reflect.TypeOf((*MockStorage)(nil).GetEventsNotifyBetween), arg0, arg1, arg2)
}

// GetMonthEvents mocks base method.
func (m *MockStorage) GetMonthEvents(arg0 context.Context, arg1 time.Time) ([]event.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMonthEvents", arg0, arg1)
	ret0, _ := ret[0].([]event.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMonthEvents indicates an expected call of GetMonthEvents.
func (mr *MockStorageMockRecorder) GetMonthEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMonthEvents", reflect.TypeOf((*MockStorage)(nil).GetMonthEvents), arg0, arg1)
}

// GetWeekEvents mocks base method.
func (m *MockStorage) GetWeekEvents(arg0 context.Context, arg1 time.Time) ([]event.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWeekEvents", arg0, arg1)
	ret0, _ := ret[0].([]event.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWeekEvents indicates an expected call of GetWeekEvents.
func (mr *MockStorageMockRecorder) GetWeekEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWeekEvents", reflect.TypeOf((*MockStorage)(nil).GetWeekEvents), arg0, arg1)
}

// Update mocks base method.
func (m *MockStorage) Update(arg0 context.Context, arg1 string, arg2 event.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockStorageMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockStorage)(nil).Update), arg0, arg1, arg2)
}
