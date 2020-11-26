package acquirer_impl

import (
	"tpayment/api/api_define"
	"tpayment/conf"

	"github.com/gin-gonic/gin"
)

type SaleRequest struct {
	TxqReq *api_define.TxnReq
}

type SaleResponse struct {
	TxnResp         *api_define.TxnResp
	AcquirerReconID string
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
