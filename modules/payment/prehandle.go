package payment

import (
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/payment/key"
	"tpayment/pkg/paymentmethod/decodecardnum/applepay"
	"tpayment/pkg/paymentmethod/decodecardnum/creditcard"
	"tpayment/pkg/paymentmethod/decodecardnum/qrcode"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

// 预处理提交的数据，
// 1. 分析出支付方式，  2. 分析出用卡方式
func preHandleRequest(ctx *gin.Context, txn *api_define.TxnReq) conf.ResultCode {
	logger := tlog.GetLogger(ctx)
	var err error

	// 分类支付方式
	switch txn.PaymentMethod {
	case conf.RequestCreditCard: // 常规信用卡
		txn.RealEntryType = txn.CreditCardBean.CardReaderMode
		txn.RealPaymentMethod, err = creditcard.Decode(txn.CreditCardBean.CardNumber)
		if err != nil {
			logger.Warn("creditcard.Decode error->", err.Error())
			return conf.DecodeCardBrandError
		}
	case conf.RequestApplePay: // apple pay
		txn.RealEntryType = conf.ApplePay

		// 查找配置的Key
		pukHash, err := applepay.GetApplePayKeyHash(txn.ApplePayBean.Token)
		if err != nil {
			logger.Warn("get apple pay hash fail->", err.Error())
			return conf.ParameterError
		}

		applePayKey := new(key.ApplePayKey)

		applePayKey, err = applePayKey.GetKeyByHash(models.DB(), ctx, pukHash)
		if err != nil {
			logger.Error("GetKeyByHash fail->", err.Error())
			return conf.DBError
		}
		if applePayKey == nil {
			logger.Info("use apple pay key->", applePayKey.ID)
			// 解码apple pay数据
			applePayBean, err := applepay.DecodeApplePay(txn.ApplePayBean.Token, &applepay.ConfigKey{
				PublicKey:  applePayKey.PublicKey,
				PrivateKey: applePayKey.PrivateKey,
			})

			if err != nil {
				logger.Warn("applepay.DecodeApplePay fail->", err.Error())
				return conf.DecodeError
			}
			txn.CreditCardBean = &api_define.CreditCardBean{
				CardExpMonth:            applePayBean.ApplicationExpirationDate[2:4],
				CardExpYear:             applePayBean.ApplicationExpirationDate[:2],
				CardExpDay:              applePayBean.ApplicationExpirationDate[4:],
				CardNumber:              applePayBean.ApplicationPrimaryAccountNumber,
				CardHolderName:          applePayBean.CardholderName,
				IccRequest:              applePayBean.PaymentData.EmvData,
				PIN:                     applePayBean.PaymentData.EncryptedPINData,
				ECI:                     applePayBean.PaymentData.EciIndicator,
				OnlinePaymentCryptogram: applePayBean.PaymentData.OnlinePaymentCryptogram,
			}
		} else {
			logger.Info("txn can't find apple pay key")
		}
	case conf.RequestConsumerPresentQR: // 商户扫手机二维码
		txn.RealEntryType = conf.ConsumerPresentQR

		// 优先判断微信，支付宝，国内云闪付版本
		txn.RealPaymentMethod, err = qrcode.Decode(txn.ConsumerPresentQR.Content)
		if err == nil {
			break
		}

		// EMV 银行卡二维码
		emvQRContent, err := qrcode.DecodeEmvQR(txn.ConsumerPresentQR.Content)
		if err == nil {
			txn.RealPaymentMethod, err = creditcard.Decode(emvQRContent.CardNum)
			if err != nil {
				logger.Warn("creditcard.Decode error->", err.Error())
				return conf.DecodeCardBrandError
			}

			txn.CreditCardBean = &api_define.CreditCardBean{
				CardNumber: emvQRContent.CardNum,
				CardTrack2: emvQRContent.Track2,
				CardSn:     emvQRContent.CardSn,
				IccRequest: emvQRContent.ICCData,
			}

			txn.RealPaymentMethod, err = creditcard.Decode(txn.CreditCardBean.CardNumber)
			if err != nil {
				logger.Warn("creditcard.Decode error->", err.Error())
				return conf.DecodeCardBrandError
			}
			return conf.Success
		}

		logger.Warn("can't decode the qr code->", txn.ConsumerPresentQR.Content)
		return conf.DecodeQRError
	}

	return conf.Success
}
