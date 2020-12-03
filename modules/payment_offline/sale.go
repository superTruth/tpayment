package payment_offline

import (
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func saleHandle(ctx *gin.Context, req *api_define.TxnReq) (*api_define.TxnResp, conf.ResultCode) {
	logger := tlog.GetLogger(ctx)
	var err error

	// 创建response数据
	resp := preBuildResp(req)

	// 预处理请求数据，解析卡数据
	errorCode := preHandleRequest(ctx, req)
	if errorCode != conf.Success {
		logger.Warn("preHandleRequest fail->", errorCode.String())
		return resp, errorCode
	}

	// 保存交易记录
	req.TxnRecord.BaseModel = models.BaseModel{
		Db:  models.DB(),
		Ctx: ctx,
	}
	err = req.TxnRecord.Create(req.TxnRecord)
	if err != nil {
		logger.Warn("create record error->", err.Error())
		return resp, conf.DBError
	}

	logger.Info("mergeRespAfterPreHandle.....")
	// 再次合并数据到返回结果
	mergeRespAfterPreHandle(resp, req)

	return resp, conf.Success
}
