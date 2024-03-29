// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/notification/queue (interfaces: Producer)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	notification "github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/notification"
	gomock "github.com/golang/mock/gomock"
)

// MockProducer is a mock of Producer interface.
type MockProducer struct {
	ctrl     *gomock.Controller
	recorder *MockProducerMockRecorder
}

// MockProducerMockRecorder is the mock recorder for MockProducer.
type MockProducerMockRecorder struct {
	mock *MockProducer
}

// NewMockProducer creates a new mock instance.
func NewMockProducer(ctrl *gomock.Controller) *MockProducer {
	mock := &MockProducer{ctrl: ctrl}
	mock.recorder = &MockProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProducer) EXPECT() *MockProducerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockProducer) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockProducerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockProducer)(nil).Close))
}

// Put mocks base method.
func (m *MockProducer) Put(arg0 notification.Notification) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockProducerMockRecorder) Put(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockProducer)(nil).Put), arg0)
}
