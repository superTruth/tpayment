package socket

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

const sessionPukFilePath = "/Users/truth/project/tpayment/pkg/paymentmethod/decodecardnum/applepay/cer/tls_rsa.crt"
const sessionPrivateFilePath = "/Users/truth/project/tpayment/pkg/paymentmethod/decodecardnum/applepay/cer/www.fang.com.key"

// 双向认证测试
func TestPaymentSocketConnect(t *testing.T) {
	pukFile, _ := os.Open(sessionPukFilePath)
	pukBytes, _ := ioutil.ReadAll(pukFile)

	priKeyFile, _ := os.Open(sessionPrivateFilePath)
	priKeyBytes, _ := ioutil.ReadAll(priKeyFile)
	_, err := tls.X509KeyPair(pukBytes, priKeyBytes)
	fmt.Println("LoadX509KeyPair->", err == nil)

	s := &PaymentSocket{
		InitObject: &InitObject{
			URL:          "cn-apple-pay-gateway.apple.com:443",
			SSLModel:     SSLMultipleAuth,
			TrustUrl:     "",
			CerCa:        "",
			CerClient:    string(pukBytes),
			CerClientKey: string(priKeyBytes),
		},
	}

	err = s.Connect(time.Second * 5)

	if err != nil {
		t.Error("connect fail->", err.Error())
		return
	}
}
