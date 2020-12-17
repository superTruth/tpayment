package socket

import (
	"time"
)

// 通信协议接口
type ICommunicateProtocol interface {
	GetRealData(receiveData []byte) ([]byte, error)
}

// 通信接口
type ICommunicate interface {
	Connect(timeOut time.Duration) error
	Disconnect(timeOut time.Duration) error
	SendData(data []byte, timeOut time.Duration) error
	ReadData(buffer []byte, timeOut time.Duration) (int, error)
	ExchangeMsg(sendData []byte, disconnectAfterReceive bool, protocol ICommunicateProtocol, timeOut time.Duration) ([]byte, ErrorCode)
}

// 数据交换过程监听
type IExChangeMsgListener struct {
	OnConnectError   func(err error)
	OnConnectSuccess func()
	OnSendError      func(err error)
	OnSendSuccess    func()
	OnReceiveError   func(err error)
	OnReceiveSuccess func()
	OnCancel         func()
}

//
type InitObject struct {
	URL          string
	SSLModel     string
	TrustUrl     string
	CerCa        string
	CerClient    string
	CerClientKey string
}

type ErrorCode string

const (
	Success     ErrorCode = "success"
	ConnectFail ErrorCode = "connect_fail"
	SendFail    ErrorCode = "send_fail"
	ReadFail    ErrorCode = "read_fail"

	SSLDisable      = "disable"
	SSLEnableSkip   = "enable_skip"
	SSLSingleAuth   = "enable_single_auth"
	SSLMultipleAuth = "enable_multiple_auth"
)
