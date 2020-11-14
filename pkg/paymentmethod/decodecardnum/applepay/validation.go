package applepay

import (
	"bytes"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/zhulingbiezhi/pkcs7"
	"sync"
	"time"
)

func ValidateRsa(p7 *pkcs7.PKCS7, orgBean *applePayOrgBean) error {
	wrappedKeyBytes, err := base64.StdEncoding.DecodeString(orgBean.Header.WrappedKey)
	if err != nil {
		fmt.Println("DecodeString2 fail->", err.Error())
		return err
	}

	dataBytes, err := base64.StdEncoding.DecodeString(orgBean.Data)
	if err != nil {
		return err
	}

	transactionIdBytes, err := hex.DecodeString(orgBean.Header.TransactionId)
	if err != nil {
		return err
	}

	applicationDataBytes, err := hex.DecodeString(orgBean.Header.ApplicationData)
	if err != nil {
		return err
	}

	sb := bytes.Buffer{}
	sb.Write(wrappedKeyBytes)
	sb.Write(dataBytes)
	sb.Write(transactionIdBytes)
	sb.Write(applicationDataBytes)

	p7.Content = sb.Bytes()

	return p7.Verify(x509.SHA256WithRSA)
}

func ValidateEcc(p7 *pkcs7.PKCS7, orgBean *applePayOrgBean) error {
	fmt.Println("ValidateEcc")
	ephemeralPubKeyBytes := []byte("-----BEGIN PUBLIC KEY-----\n" + orgBean.Header.EphemeralPublicKey + "\n-----END PUBLIC KEY-----")
	pukByte, _ := pem.Decode(ephemeralPubKeyBytes)
	if pukByte == nil {
		return errors.New("decode EphemeralPublicKey fail")
	}

	dataBytes, err := base64.StdEncoding.DecodeString(orgBean.Data)
	if err != nil {
		return err
	}

	transactionIdBytes, err := hex.DecodeString(orgBean.Header.TransactionId)
	if err != nil {
		return err
	}

	applicationDataBytes, err := hex.DecodeString(orgBean.Header.ApplicationData)
	if err != nil {
		return err
	}

	sb := bytes.Buffer{}
	sb.Write(ephemeralPubKeyBytes)
	sb.Write(dataBytes)
	sb.Write(transactionIdBytes)
	sb.Write(applicationDataBytes)

	p7.Content = sb.Bytes()

	return p7.Verify(x509.ECDSAWithSHA256)
}

var (
	oidData                   = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 1}
	oidSignedData             = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 2}
	oidEnvelopedData          = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 3}
	oidSignedAndEnvelopedData = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 4}
	oidDigestedData           = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 5}
	oidEncryptedData          = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 7, 6}
	oidAttributeContentType   = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 3}
	oidAttributeMessageDigest = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 4}
	oidAttributeSigningTime   = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 9, 5}

	loadApplePayCAProcess sync.Once
	applePayCA            *x509.Certificate
)

const (
	maxValidateTime = 5 * time.Minute
	applePayCa      = `MIICQzCCAcmgAwIBAgIILcX8iNLFS5UwCgYIKoZIzj0EAwMwZzEbMBkGA1UEAwwSQXBwbGUgUm9vdCBDQSAtIEczMSYwJAYDVQQLDB1BcHBsZSBDZXJ0aWZpY2F0aW9uIEF1dGhvcml0eTETMBEGA1UECgwKQXBwbGUgSW5jLjELMAkGA1UEBhMCVVMwHhcNMTQwNDMwMTgxOTA2WhcNMzkwNDMwMTgxOTA2WjBnMRswGQYDVQQDDBJBcHBsZSBSb290IENBIC0gRzMxJjAkBgNVBAsMHUFwcGxlIENlcnRpZmljYXRpb24gQXV0aG9yaXR5MRMwEQYDVQQKDApBcHBsZSBJbmMuMQswCQYDVQQGEwJVUzB2MBAGByqGSM49AgEGBSuBBAAiA2IABJjpLz1AcqTtkyJygRMc3RCV8cWjTnHcFBbZDuWmBSp3ZHtfTjjTuxxEtX/1H7YyYl3J6YRbTzBPEVoA/VhYDKX1DyxNB0cTddqXl5dvMVztK517IDvYuVTZXpmkOlEKMaNCMEAwHQYDVR0OBBYEFLuw3qFYM4iapIqZ3r6966/ayySrMA8GA1UdEwEB/wQFMAMBAf8wDgYDVR0PAQH/BAQDAgEGMAoGCCqGSM49BAMDA2gAMGUCMQCD6cHEFl4aXTQY2e3v9GwOAEZLuN+yRhHFD/3meoyhpmvOwgPUnPWTxnS4at+qIxUCMG1mihDK1A3UT82NQz60imOlM27jbdoXt2QfyFMm+YhidDkLF1vLUagM6BgD56KyKA==`
)

func validateSignature(orgBean *applePayOrgBean, validateSignatureFunc func(*pkcs7.PKCS7, *applePayOrgBean) error) error {
	signData, err := base64.StdEncoding.DecodeString(orgBean.Signature)
	if err != nil {
		return errors.New("get signature data fail")
	}

	loadApplePayCAProcess.Do(func() {
		var err error
		caBytes, _ := base64.StdEncoding.DecodeString(applePayCa)
		applePayCA, err = x509.ParseCertificate(caBytes)
		if err != nil {
			fmt.Println("apple pay ca parseCertificate fail->", err.Error())
			panic("apple pay ca parseCertificate fail")
			return
		}
		if !applePayCA.IsCA {
			panic("apple pay ca not root ca")
			return
		}
	})

	p7, err := pkcs7.Parse(signData)
	if err != nil {
		fmt.Println("pkcs7.Parse->", err.Error())
		return err
	}
	// a. find leaf inter cer
	var (
		leafCer, inter *x509.Certificate
	)
	for i, cer := range p7.Certificates {
		for _, info := range cer.Extensions {
			idStr := info.Id.String()
			if idStr == LeafOID {
				leafCer = p7.Certificates[i]
				break
			}
			if idStr == CAOID {
				inter = p7.Certificates[i]
				break
			}
		}
	}
	if leafCer == nil || inter == nil {
		return errors.New("certification not correct")
	}

	//// b. check intermediate cer inter is signed from apple ca
	//if inter.CheckSignatureFrom(applePayCA) != nil {
	//	return errors.New("intermediate is not signed by apple pay CA G3")
	//}
	//
	//// c. check leaf cer is signed from intermediate cer
	//if leafCer.CheckSignatureFrom(inter) != nil {
	//	return errors.New("leaf is not signed by intermediate cer")
	//}

	// d. Validate the token’s signature
	err = validateSignatureFunc(p7, orgBean)
	if err != nil {
		fmt.Println("validateSignatureFunc->", err.Error())
		return err
	}

	// e. validate sign time
	var signTime *time.Time
	for _, signer := range p7.Signers {
		for _, attr := range signer.AuthenticatedAttributes {
			if attr.Type.Equal(oidAttributeSigningTime) {
				signTime = new(time.Time)
				_, err := asn1.Unmarshal(attr.Value.Bytes, signTime)
				if err != nil {
					signTime = nil
					return errors.New("unmarshal time fail->" + err.Error())
				}
			}
		}
	}
	if signTime == nil {
		return errors.New("can't find sign time")
	}
	//if time.Now().Sub(*signTime) > maxValidateTime {
	//	return errors.New("the token is expired")
	//}

	return nil
}
