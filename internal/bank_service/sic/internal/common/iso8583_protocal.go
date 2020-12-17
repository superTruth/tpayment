package common

import (
	"errors"
	"fmt"
	"tpayment/pkg/utils/convert_utils"
	"tpayment/pkg/utils/mix_utils"
)

const Iso8583Offset = 13

type Cn8583Protocol struct {
	recvBuf []byte
	index   int
}

func (p *Cn8583Protocol) init() {
	p.index = 0
	p.recvBuf = make([]byte, 1024)
}

func (p *Cn8583Protocol) GetRealData(receiveData []byte) ([]byte, error) {
	if p.recvBuf == nil {
		p.init()
	}
	fmt.Println("GetRealData->", convert_utils.Bytes2HexString(receiveData))

	mix_utils.BytesArrayCopy(receiveData, 0, p.recvBuf, p.index, len(receiveData))

	p.index += len(receiveData)

	isFinish, err := p.isRecvOver()

	if err != nil {
		return nil, err
	}

	if !isFinish {
		return nil, nil
	}

	return mix_utils.BytesArrayCopyArrange(p.recvBuf, Iso8583Offset, p.index), nil
}

func (p *Cn8583Protocol) isRecvOver() (bool, error) {
	if p.index < Iso8583Offset {
		return false, nil
	}
	needLen := int(convert_utils.BytesHex2Long(p.recvBuf, 0, 2))
	if (p.index - 2) < needLen {
		return false, nil
	}
	if (p.index - 2) > needLen {
		return false, errors.New("format error")
	}
	return true, nil
}
