package pay_offline

import (
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/models/txn"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func transferHandle(ctx *gin.Context, req *api_define.TxnReq) (*api_define.TxnResp, conf.ResultCode) {
	logger := tlog.GetLogger(ctx)
	var err error

	err = api_define.Validate(ctx, req)
	if err != nil {
		logger.Warn("validate request body error->", err.Error())
		return nil, conf.ParameterError
	}

	// 创建response数据
	resp := preBuildResp(req)

	// 预处理请求数据，解析卡数据
	errorCode := preHandleRequest(ctx, req)
	if errorCode != conf.Success {
		logger.Warn("preHandleRequest fail->", errorCode.String())
		return resp, errorCode
	}

	// 保存交易记录
	err = txn.CreateTransactionAndDetail(req.TxnRecord, req.TxnRecordDetail)
	if err != nil {
		logger.Error(" models.DB().Transaction error->", err.Error())
		return resp, conf.DBError
	}

	logger.Info("mergeRespAfterPreHandle.....")
	// 再次合并数据到返回结果
	mergeRespAfterPreHandle(resp, req)

	return resp, conf.Success
}
