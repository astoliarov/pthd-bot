// Code generated by MockGen. DO NOT EDIT.
// Source: pthd-bot/pkg/interfaces (interfaces: IBotKillLogDAO)

// Package mocks is a generated GoMock package.
package mocks

import (
	entities "pthd-bot/pkg/entities"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIBotKillLogDAO is a mock of IBotKillLogDAO interface.
type MockIBotKillLogDAO struct {
	ctrl     *gomock.Controller
	recorder *MockIBotKillLogDAOMockRecorder
}

// MockIBotKillLogDAOMockRecorder is the mock recorder for MockIBotKillLogDAO.
type MockIBotKillLogDAOMockRecorder struct {
	mock *MockIBotKillLogDAO
}

// NewMockIBotKillLogDAO creates a new mock instance.
func NewMockIBotKillLogDAO(ctrl *gomock.Controller) *MockIBotKillLogDAO {
	mock := &MockIBotKillLogDAO{ctrl: ctrl}
	mock.recorder = &MockIBotKillLogDAOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIBotKillLogDAO) EXPECT() *MockIBotKillLogDAOMockRecorder {
	return m.recorder
}

// GetTopVictims mocks base method.
func (m *MockIBotKillLogDAO) GetTopVictims(arg0 string) ([]*entities.TopVictimLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTopVictims", arg0)
	ret0, _ := ret[0].([]*entities.TopVictimLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTopVictims indicates an expected call of GetTopVictims.
func (mr *MockIBotKillLogDAOMockRecorder) GetTopVictims(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTopVictims", reflect.TypeOf((*MockIBotKillLogDAO)(nil).GetTopVictims), arg0)
}

// Save mocks base method.
func (m *MockIBotKillLogDAO) Save(arg0 *entities.BotKill) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockIBotKillLogDAOMockRecorder) Save(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIBotKillLogDAO)(nil).Save), arg0)
}
