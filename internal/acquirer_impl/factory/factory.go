package factory

import "tpayment/internal/acquirer_impl/iso8583/standard"

var AcquirerImpls = map[string]interface{}{
	"with_tid_default": &standard.API{},
}
