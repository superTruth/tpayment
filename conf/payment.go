package conf

const (
	// Request Payment Method
	RequestCreditCard        = "credit_card"
	RequestCreditCardToken   = "credit_card_token"
	RequestConsumerPresentQR = "consumer_present_qr"
	RequestApplePay          = "apple_pay"

	// Real Payment Method
	Visa       = "visa"
	MasterCard = "master_card"
	UnionPay   = "union_pay"
	AE         = "amex"
	JCB        = "jcb"
	WeChatPay  = "wechat_pay"
	Alipay     = "alipay"

	// Payment Entry Type
	Swipe             = "swipe"
	Contact           = "contact"
	ContactLess       = "contact_less"
	ManualInput       = "manual_input"
	Token             = "token"
	ApplePay          = "apple_pay"
	MerchantPresentQR = "merchant_present_qr"
	ConsumerPresentQR = "consumer_present_qr"
	InApp             = "in_app"
	InWeb             = "in_web"

	// Payment Type
	Sale            = "sale"
	Void            = "void"
	Refund          = "refund"
	PreAuth         = "pre_auth"
	PreAuthComplete = "pre_auth_complete"
)
