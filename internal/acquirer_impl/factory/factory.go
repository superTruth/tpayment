package factory

import "tpayment/internal/acquirer_impl/iso8583/standard"

var AcquirerImpls = map[string]interface{}{
	"wlb": &standard.API{},
}
