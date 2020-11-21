package acquirer_impl

import (
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/models/agency"
	"tpayment/models/merchant"
)

type SaleRequest struct {
	TxqReq       *api_define.TxnReq
	AgencyInfo   *agency.Agency
	MerchantInfo *merchant.Merchant
}

type SaleResponse struct {
	TxnResp *api_define.TxnResp
}

type ISale interface {
	Sale(req *SaleRequest) (*SaleResponse, conf.ResultCode)
}

type IVoid interface {
	Void()
}
