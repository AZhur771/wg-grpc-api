// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package app_mocks is a generated GoMock package.
package app_mocks

import (
	context "context"
	reflect "reflect"

	dto "github.com/AZhur771/wg-grpc-api/internal/dto"
	entity "github.com/AZhur771/wg-grpc-api/internal/entity"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	zapcore "go.uber.org/zap/zapcore"
	wgtypes "golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// MockLogger is a mock of Logger interface.
type MockLogger struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerMockRecorder
}

// MockLoggerMockRecorder is the mock recorder for MockLogger.
type MockLoggerMockRecorder struct {
	mock *MockLogger
}

// NewMockLogger creates a new mock instance.
func NewMockLogger(ctrl *gomock.Controller) *MockLogger {
	mock := &MockLogger{ctrl: ctrl}
	mock.recorder = &MockLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogger) EXPECT() *MockLoggerMockRecorder {
	return m.recorder
}

// Debug mocks base method.
func (m *MockLogger) Debug(arg0 string, arg1 ...zapcore.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debug", varargs...)
}

// Debug indicates an expected call of Debug.
func (mr *MockLoggerMockRecorder) Debug(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockLogger)(nil).Debug), varargs...)
}

// Error mocks base method.
func (m *MockLogger) Error(arg0 string, arg1 ...zapcore.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error.
func (mr *MockLoggerMockRecorder) Error(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockLogger)(nil).Error), varargs...)
}

// Info mocks base method.
func (m *MockLogger) Info(arg0 string, arg1 ...zapcore.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockLoggerMockRecorder) Info(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLogger)(nil).Info), varargs...)
}

// Warn mocks base method.
func (m *MockLogger) Warn(arg0 string, arg1 ...zapcore.Field) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warn", varargs...)
}

// Warn indicates an expected call of Warn.
func (mr *MockLoggerMockRecorder) Warn(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warn", reflect.TypeOf((*MockLogger)(nil).Warn), varargs...)
}

// MockPeerService is a mock of PeerService interface.
type MockPeerService struct {
	ctrl     *gomock.Controller
	recorder *MockPeerServiceMockRecorder
}

// MockPeerServiceMockRecorder is the mock recorder for MockPeerService.
type MockPeerServiceMockRecorder struct {
	mock *MockPeerService
}

// NewMockPeerService creates a new mock instance.
func NewMockPeerService(ctrl *gomock.Controller) *MockPeerService {
	mock := &MockPeerService{ctrl: ctrl}
	mock.recorder = &MockPeerServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPeerService) EXPECT() *MockPeerServiceMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockPeerService) Add(ctx context.Context, addPeerDTO dto.AddPeerDTO) (*entity.Peer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, addPeerDTO)
	ret0, _ := ret[0].(*entity.Peer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockPeerServiceMockRecorder) Add(ctx, addPeerDTO interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockPeerService)(nil).Add), ctx, addPeerDTO)
}

// Disable mocks base method.
func (m *MockPeerService) Disable(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Disable", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Disable indicates an expected call of Disable.
func (mr *MockPeerServiceMockRecorder) Disable(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disable", reflect.TypeOf((*MockPeerService)(nil).Disable), ctx, id)
}

// DownloadConfig mocks base method.
func (m *MockPeerService) DownloadConfig(ctx context.Context, id uuid.UUID) (dto.DownloadFileDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadConfig", ctx, id)
	ret0, _ := ret[0].(dto.DownloadFileDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadConfig indicates an expected call of DownloadConfig.
func (mr *MockPeerServiceMockRecorder) DownloadConfig(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadConfig", reflect.TypeOf((*MockPeerService)(nil).DownloadConfig), ctx, id)
}

// DownloadQRCode mocks base method.
func (m *MockPeerService) DownloadQRCode(ctx context.Context, id uuid.UUID) (dto.DownloadFileDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadQRCode", ctx, id)
	ret0, _ := ret[0].(dto.DownloadFileDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadQRCode indicates an expected call of DownloadQRCode.
func (mr *MockPeerServiceMockRecorder) DownloadQRCode(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadQRCode", reflect.TypeOf((*MockPeerService)(nil).DownloadQRCode), ctx, id)
}

// Enable mocks base method.
func (m *MockPeerService) Enable(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Enable", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Enable indicates an expected call of Enable.
func (mr *MockPeerServiceMockRecorder) Enable(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Enable", reflect.TypeOf((*MockPeerService)(nil).Enable), ctx, id)
}

// Get mocks base method.
func (m *MockPeerService) Get(ctx context.Context, id uuid.UUID) (*entity.Peer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*entity.Peer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockPeerServiceMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPeerService)(nil).Get), ctx, id)
}

// GetAll mocks base method.
func (m *MockPeerService) GetAll(ctx context.Context, getPeersDTO dto.GetPeersRequestDTO) (dto.GetPeersResponseDTO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, getPeersDTO)
	ret0, _ := ret[0].(dto.GetPeersResponseDTO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockPeerServiceMockRecorder) GetAll(ctx, getPeersDTO interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockPeerService)(nil).GetAll), ctx, getPeersDTO)
}

// Remove mocks base method.
func (m *MockPeerService) Remove(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockPeerServiceMockRecorder) Remove(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockPeerService)(nil).Remove), ctx, id)
}

// Update mocks base method.
func (m *MockPeerService) Update(ctx context.Context, updatePeerDTO dto.UpdatePeerDTO) (*entity.Peer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, updatePeerDTO)
	ret0, _ := ret[0].(*entity.Peer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockPeerServiceMockRecorder) Update(ctx, updatePeerDTO interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPeerService)(nil).Update), ctx, updatePeerDTO)
}

// MockDeviceService is a mock of DeviceService interface.
type MockDeviceService struct {
	ctrl     *gomock.Controller
	recorder *MockDeviceServiceMockRecorder
}

// MockDeviceServiceMockRecorder is the mock recorder for MockDeviceService.
type MockDeviceServiceMockRecorder struct {
	mock *MockDeviceService
}

// NewMockDeviceService creates a new mock instance.
func NewMockDeviceService(ctrl *gomock.Controller) *MockDeviceService {
	mock := &MockDeviceService{ctrl: ctrl}
	mock.recorder = &MockDeviceServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeviceService) EXPECT() *MockDeviceServiceMockRecorder {
	return m.recorder
}

// AddPeer mocks base method.
func (m *MockDeviceService) AddPeer(peer wgtypes.PeerConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddPeer", peer)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddPeer indicates an expected call of AddPeer.
func (mr *MockDeviceServiceMockRecorder) AddPeer(peer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddPeer", reflect.TypeOf((*MockDeviceService)(nil).AddPeer), peer)
}

// GetDevice mocks base method.
func (m *MockDeviceService) GetDevice() (*entity.Device, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDevice")
	ret0, _ := ret[0].(*entity.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDevice indicates an expected call of GetDevice.
func (mr *MockDeviceServiceMockRecorder) GetDevice() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDevice", reflect.TypeOf((*MockDeviceService)(nil).GetDevice))
}

// GetPeer mocks base method.
func (m *MockDeviceService) GetPeer(publicKey wgtypes.Key) (wgtypes.Peer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPeer", publicKey)
	ret0, _ := ret[0].(wgtypes.Peer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPeer indicates an expected call of GetPeer.
func (mr *MockDeviceServiceMockRecorder) GetPeer(publicKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeer", reflect.TypeOf((*MockDeviceService)(nil).GetPeer), publicKey)
}

// GetPeers mocks base method.
func (m *MockDeviceService) GetPeers() ([]wgtypes.Peer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPeers")
	ret0, _ := ret[0].([]wgtypes.Peer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPeers indicates an expected call of GetPeers.
func (mr *MockDeviceServiceMockRecorder) GetPeers() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPeers", reflect.TypeOf((*MockDeviceService)(nil).GetPeers))
}

// RemovePeer mocks base method.
func (m *MockDeviceService) RemovePeer(peer wgtypes.PeerConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemovePeer", peer)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemovePeer indicates an expected call of RemovePeer.
func (mr *MockDeviceServiceMockRecorder) RemovePeer(peer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemovePeer", reflect.TypeOf((*MockDeviceService)(nil).RemovePeer), peer)
}

// Setup mocks base method.
func (m *MockDeviceService) Setup(ctx context.Context, name, endpoint, address string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Setup", ctx, name, endpoint, address)
	ret0, _ := ret[0].(error)
	return ret0
}

// Setup indicates an expected call of Setup.
func (mr *MockDeviceServiceMockRecorder) Setup(ctx, name, endpoint, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Setup", reflect.TypeOf((*MockDeviceService)(nil).Setup), ctx, name, endpoint, address)
}

// UpdatePeer mocks base method.
func (m *MockDeviceService) UpdatePeer(peer wgtypes.PeerConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePeer", peer)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePeer indicates an expected call of UpdatePeer.
func (mr *MockDeviceServiceMockRecorder) UpdatePeer(peer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePeer", reflect.TypeOf((*MockDeviceService)(nil).UpdatePeer), peer)
}

// MockPeerStorage is a mock of PeerStorage interface.
type MockPeerStorage struct {
	ctrl     *gomock.Controller
	recorder *MockPeerStorageMockRecorder
}

// MockPeerStorageMockRecorder is the mock recorder for MockPeerStorage.
type MockPeerStorageMockRecorder struct {
	mock *MockPeerStorage
}

// NewMockPeerStorage creates a new mock instance.
func NewMockPeerStorage(ctrl *gomock.Controller) *MockPeerStorage {
	mock := &MockPeerStorage{ctrl: ctrl}
	mock.recorder = &MockPeerStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPeerStorage) EXPECT() *MockPeerStorageMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockPeerStorage) Add(ctx context.Context, peer *entity.Peer) (*entity.Peer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", ctx, peer)
	ret0, _ := ret[0].(*entity.Peer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockPeerStorageMockRecorder) Add(ctx, peer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockPeerStorage)(nil).Add), ctx, peer)
}

// CountAll mocks base method.
func (m *MockPeerStorage) CountAll(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountAll", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountAll indicates an expected call of CountAll.
func (mr *MockPeerStorageMockRecorder) CountAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountAll", reflect.TypeOf((*MockPeerStorage)(nil).CountAll), ctx)
}

// Get mocks base method.
func (m *MockPeerStorage) Get(ctx context.Context, id uuid.UUID) (*entity.Peer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*entity.Peer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockPeerStorageMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPeerStorage)(nil).Get), ctx, id)
}

// GetAll mocks base method.
func (m *MockPeerStorage) GetAll(ctx context.Context, skip, limit int) ([]*entity.Peer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, skip, limit)
	ret0, _ := ret[0].([]*entity.Peer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockPeerStorageMockRecorder) GetAll(ctx, skip, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockPeerStorage)(nil).GetAll), ctx, skip, limit)
}

// Remove mocks base method.
func (m *MockPeerStorage) Remove(ctx context.Context, id uuid.UUID) (*entity.Peer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", ctx, id)
	ret0, _ := ret[0].(*entity.Peer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Remove indicates an expected call of Remove.
func (mr *MockPeerStorageMockRecorder) Remove(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockPeerStorage)(nil).Remove), ctx, id)
}

// Update mocks base method.
func (m *MockPeerStorage) Update(ctx context.Context, peer *entity.Peer) (*entity.Peer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, peer)
	ret0, _ := ret[0].(*entity.Peer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockPeerStorageMockRecorder) Update(ctx, peer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockPeerStorage)(nil).Update), ctx, peer)
}
