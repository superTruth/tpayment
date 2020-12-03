package standard

import (
	"strconv"
	"strings"
	"tpayment/models/agency"
	"tpayment/models/payment/acquirer"
	"tpayment/models/payment/merchantaccount"
)

func GetAccountTag(acq *agency.Acquirer, mid *merchantaccount.MerchantAccount, tid *acquirer.Terminal) string {
	sb := strings.Builder{}

	sb.WriteString(strconv.Itoa(int(acq.ID)))
	sb.WriteString("-")
	sb.WriteString(strconv.Itoa(int(mid.ID)))
	if tid != nil {

	}

	return sb.String()
}
