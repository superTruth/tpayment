package associate

import (
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/merchant"
	"tpayment/modules"
	merchantModule "tpayment/modules/merchant"
	"tpayment/pkg/tlog"
	"tpayment/pkg/utils"

	"github.com/labstack/echo"
)

func AddHandle(ctx echo.Context) error {
	logger := tlog.GetLogger(ctx)

	req := new(merchant.UserMerchantAssociate)

	err := utils.Body2Json(ctx.Request().Body, req)
	if err != nil {
		logger.Warn("Body2Json fail->", err.Error())
		modules.BaseError(ctx, conf.ParameterError)
		return err
	}

	// 判断权限
	err = merchantModule.CheckPermission(ctx, req.MerchantId)
	if err != nil {
		logger.Warn(err.Error())
		modules.BaseError(ctx, conf.NoPermission)
		return err
	}

	err = models.CreateBaseRecord(req)

	if err != nil {
		logger.Info("CreateBaseRecord sql error->", err.Error())
		modules.BaseError(ctx, conf.DBError)
		return err
	}

	modules.BaseSuccess(ctx, nil)

	return nil
}

//func checkPermission(ctx echo.Context, bean *merchant.UserMerchantAssociate) error {
//	// 1. 是否存在这2个账号
//	merchantBean, err := merchant.GetMerchantById(models.DB(), ctx, bean.MerchantId)
//	if err != nil {
//		return err
//	}
//	if merchantBean == nil {
//		return errors.New("merchant id not exist")
//	}
//
//	userBean, err := account.GetUserById(models.DB(), ctx, bean.UserId)
//	if err != nil {
//		return err
//	}
//	if userBean == nil {
//		return errors.New("user id not exist")
//	}
//
//	// 2. 当前用户是否有权限操作这个商户
//	if modules.IsAdmin(ctx) != nil { // 管理员账户可以直接操作
//		return nil
//	}
//
//	agency := modules.IsAgencyAdmin(ctx)
//	if agency != nil && merchantBean.AgencyId == agency.ID { // 相关机构管理员可以操作
//		return nil
//	}
//
//	// 当前账号是否是此merchant下面关联的管理账号
//	currentUserBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)
//	associateBean, err := merchant.GetUserMerchantAssociateByMerchantIdAndUserId(models.DB(), ctx,
//		bean.MerchantId, currentUserBean.ID)
//	if err != nil {
//		return err
//	}
//	if associateBean == nil || associateBean.Role != string(conf.RoleAdmin) {
//		return errors.New("no permission for the merchant")
//	}
//
//	return nil
//}
