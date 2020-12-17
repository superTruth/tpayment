package acquirer_impl

import (
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/models/agency"
	"tpayment/models/payment/acquirer"
	"tpayment/models/payment/merchantaccount"
	"tpayment/models/payment/record"

	"github.com/jinzhu/gorm"

	"github.com/gin-gonic/gin"
)

type SaleRequest struct {
	TxqReq          *api_define.TxnReq        `json:"txq_req"`
	SettlementTotal []*record.SettlementTotal `json:"settlement_total"`
	Keys            []*acquirer.Key           `json:"keys"`
}

type SaleResponse struct {
	TxnResp         *api_define.TxnResp `json:"txn_resp"`
	AcquirerReconID string              `json:"acquirer_recon_id"`
	Keys            []*acquirer.Key     `json:"keys"`
}

type ISale interface {
	Sale(ctx *gin.Context, req *SaleRequest) (*SaleResponse, conf.ResultCode)
}

type IVoid interface {
	Void(ctx *gin.Context, req *SaleRequest) (*SaleResponse, conf.ResultCode)
}

type IRefund interface {
	Refund(ctx *gin.Context, req *SaleRequest) (*SaleResponse, conf.ResultCode)
}

type IPreAuth interface {
	PreAuth(ctx *gin.Context, req *SaleRequest) (*SaleResponse, conf.ResultCode)
}

type IPreAuthComplete interface {
	PreAuthComplete(ctx *gin.Context, req *SaleRequest) (*SaleResponse, conf.ResultCode)
}

type ISettlementInMID interface {
	SettlementInMID(acq *agency.Acquirer, mid *merchantaccount.MerchantAccount,
		dbFunc func(...func(*gorm.DB) *gorm.DB) *gorm.DB) conf.ResultCode
}

type ISettlementInTID interface {
	SettlementInTID(acq *agency.Acquirer, mid *merchantaccount.MerchantAccount,
		tid *acquirer.Terminal, dbFunc func(...func(*gorm.DB) *gorm.DB) *gorm.DB) conf.ResultCode
}
