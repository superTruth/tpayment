package applepay

import (
	"crypto/x509"
	"tpayment/pkg/algorithmutils"
)

func CalcCerHash(cerContent []byte) (string, error) {
	cer, err := x509.ParseCertificate(cerContent)
	if err != nil {
		return "", nil
	}

	ret := algorithmutils.Sha256(cer.RawSubjectPublicKeyInfo)

	return ret, nil
}
