package standard

import (
	"strconv"
	"tpayment/conf"
	"tpayment/internal/acquirer_impl"
	"tpayment/models"
	"tpayment/models/payment/acquirer"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func (wlb *API) Sale(ctx *gin.Context, req *acquirer_impl.SaleRequest) (*acquirer_impl.SaleResponse, conf.ResultCode) {
	logger := tlog.GetLogger(ctx)

	var err error
	// 获取账号信息
	account := &acquirer.Account{
		BaseModel: models.BaseModel{
			Db:  models.DB(),
			Ctx: ctx,
		},
	}
	account, err = account.GetOrCreate(
		strconv.Itoa(int(req.TxqReq.PaymentProcessRule.MerchantAccountID)),
		GetAccountTag(req))
	if err != nil {
		logger.Error("account.GetOrCreate error->", err.Error())
		return nil, conf.DBError
	}

	// 拼接发送数据

	// 流水号增加

	//

	return nil, conf.SUCCESS
}
