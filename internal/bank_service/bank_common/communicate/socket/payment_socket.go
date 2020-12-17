package socket

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"time"
	"tpayment/pkg/utils/mix_utils"
)

type PaymentSocket struct {
	InitObject    *InitObject
	socket        net.Conn
	receiveBuffer []byte
	IsConnected   bool
}

// 初始化
func Generate(initObject *InitObject) *PaymentSocket {
	return &PaymentSocket{
		InitObject:    initObject,
		receiveBuffer: make([]byte, 1024),
		IsConnected:   false,
	}
}

// 连接
func (p *PaymentSocket) Connect(timeOut time.Duration) error {
	var err error

	if p.InitObject == nil || p.InitObject.SSLModel == SSLDisable {
		p.socket, err = net.DialTimeout("tcp", p.InitObject.URL, timeOut)
	} else {
		var conf *tls.Config

		switch p.InitObject.SSLModel {
		case SSLEnableSkip:
			conf = &tls.Config{
				InsecureSkipVerify: true,
			}
		case SSLSingleAuth:
			certPool := x509.NewCertPool()
			loadCertFileRet := certPool.AppendCertsFromPEM([]byte(p.InitObject.CerCa))
			if !loadCertFileRet {
				return errors.New("cert file format err")
			}

			conf = &tls.Config{
				InsecureSkipVerify: false,
				RootCAs:            certPool,
				ServerName:         p.InitObject.TrustUrl,
			}
		case SSLMultipleAuth:
			conf = &tls.Config{}

			if p.InitObject.CerCa == "" {
				conf.InsecureSkipVerify = true
			} else {
				certPool := x509.NewCertPool()
				loadCertFileRet := certPool.AppendCertsFromPEM([]byte(p.InitObject.CerCa))
				if !loadCertFileRet {
					return errors.New("cert file format err")
				}

				conf.InsecureSkipVerify = false
				conf.RootCAs = certPool
			}

			if p.InitObject.TrustUrl != "" {
				conf.ServerName = p.InitObject.TrustUrl
			}

			cer, err := tls.X509KeyPair([]byte(p.InitObject.CerClient), []byte(p.InitObject.CerClientKey))
			if err != nil {
				return errors.New("tls key error->" + err.Error())
			}
			conf.Certificates = []tls.Certificate{
				cer,
			}
		default:
			return errors.New("didn't init socket")
		}

		p.socket, err = tls.DialWithDialer(&net.Dialer{Timeout: timeOut}, "tcp", p.InitObject.URL, conf)
	}

	if err == nil {
		p.IsConnected = true // connect sucess
	} else {
		err = errors.New("connect fail")
	}

	//tls
	return err
}

// 断开连接
func (p *PaymentSocket) Disconnect(timeOut time.Duration) (errRet error) {
	if p.IsConnected {
		errRet = p.socket.Close()
	}

	p.IsConnected = false // connect success

	if errRet != nil {
		errRet = errors.New("disconnect fail")
	}

	return errRet
}

// 发送数据
func (p *PaymentSocket) SendData(data []byte, timeOut time.Duration) (errRet error) {
	if !p.IsConnected {
		return errors.New("Not connect")
	}

	p.socket.SetWriteDeadline(time.Now().Add(timeOut))
	_, errRet = p.socket.Write(data)

	if errRet != nil {
		p.Disconnect(time.Second * 1)
		errRet = errors.New("send fail")
	}

	return
}

// 读取数据
func (p *PaymentSocket) ReadData(buffer []byte, timeOut time.Duration) (readLen int, errRet error) {
	if !p.IsConnected {
		return 0, errors.New("Not connect")
	}

	p.socket.SetReadDeadline(time.Now().Add(timeOut))

	readLen, errRet = p.socket.Read(buffer)

	if errRet != nil {
		p.Disconnect(time.Second * 1)
		errRet = errors.New("read fail")
	}

	return
}

// 数据交换
func (p *PaymentSocket) ExchangeMsg(sendData []byte, disconnectAfterReceive bool, protocol ICommunicateProtocol, timeOut time.Duration) ([]byte, ErrorCode) {
	var err error

	// 建立连接
	if !p.IsConnected {
		for i := 0; i < 3; i++ { // 3次连接机会
			err = p.Connect(timeOut)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 200) // 连接失败时，稍微等待200毫秒进行重连
		}
		if err != nil {
			return nil, ConnectFail
		}
	}

	// 方法结束时，断开连接
	if disconnectAfterReceive {
		defer p.Disconnect(timeOut)
	}

	// 发送数据
	err = p.SendData(sendData, timeOut)
	if err != nil {
		return nil, SendFail
	}

	// 读取数据
	var ret []byte
	readLen := 0
	for true {
		// 读取数据
		readLen, err = p.ReadData(p.receiveBuffer, timeOut)
		if err != nil {
			return nil, ReadFail
		}

		if readLen <= 0 {
			return nil, ReadFail
		}

		// 根据协议截取头尾
		ret, err = protocol.GetRealData(mix_utils.BytesArrayCopyArrange(p.receiveBuffer, 0, readLen))
		if err != nil {
			return nil, ReadFail
		}

		if ret != nil {
			break
		}
	}
	return ret, Success
}
