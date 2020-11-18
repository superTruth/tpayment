package qrcode

type EMVQRDecodeContent struct {
	CardNum string
	ICCData string
}

func DecodeEmvQR(qrCode string) (*EMVQRDecodeContent, error) {

	return nil, nil
}
