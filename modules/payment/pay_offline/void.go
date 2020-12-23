package pay_offline

import (
	"tpayment/api/api_define"
	"tpayment/conf"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func voidHandle(ctx *gin.Context, req *api_define.TxnReq) (*api_define.TxnResp, conf.ResultCode) {
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

	// 没有原始交易
	if req.OrgRecord == nil {
		logger.Warn("can't find org record->", req.OrgTxnID)
		return resp, conf.RecordNotFund
	}

	// 判断是否可以void
	if req.OrgRecord.AcquirerSettlementAt != nil { // 被结算过的交易不能void
		logger.Warn("the record was settled->", req.OrgTxnID)
		return resp, conf.Success
	}

	if req.OrgRecord.VoidAt != nil { // 已经被void过
		logger.Warn("the record was voided")
		return resp, conf.Success
	}

	//// 保存交易记录
	//err = models.DB().Transaction(func(tx *gorm.DB) error {
	//	req.TxnRecord.BaseModel.Db = &models.MyDB{DB: tx}
	//	err = req.TxnRecord.Create(req.TxnRecord)
	//	if err != nil {
	//		logger.Warn("create record error->", err.Error())
	//		return err
	//	}
	//
	//	req.TxnRecordDetail.BaseModel.Db = &models.MyDB{DB: tx}
	//	err = req.TxnRecordDetail.Create(req.TxnRecordDetail)
	//	if err != nil {
	//		logger.Warn("create detail record error->", err.Error())
	//		return err
	//	}
	//
	//	// 更新原始记录
	//	req.OrgRecord.BaseModel.Db = &models.MyDB{DB: tx}
	//
	//	t := time.Now()
	//	req.OrgRecord.VoidAt = &t
	//	err = req.OrgRecord.UpdateVoidStatus()
	//	if err != nil {
	//		logger.Error("UpdateVoidStatus fail->", err.Error())
	//		return err
	//	}
	//
	//	return nil
	//})
	//
	//if err != nil {
	//	logger.Error(" models.DB().Transaction error->", err.Error())
	//	return resp, conf.DBError
	//}

	logger.Info("mergeRespAfterPreHandle.....")
	// 再次合并数据到返回结果
	mergeRespAfterPreHandle(resp, req)

	return resp, conf.Success
}
