package merchant

import (
	"errors"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/models/merchant"
	"tpayment/modules"

	"github.com/labstack/echo"
)

func CheckPermission(ctx echo.Context, merchantId uint) error {
	// 1. 是否存在这2个账号
	merchantBean, err := merchant.GetMerchantById(models.DB(), ctx, merchantId)
	if err != nil {
		return err
	}
	if merchantBean == nil {
		return errors.New("merchant id not exist")
	}

	// 2. 当前用户是否有权限操作这个商户
	if modules.IsAdmin(ctx) != nil { // 管理员账户可以直接操作
		return nil
	}

	agency := modules.IsAgencyAdmin(ctx)
	if agency != nil && merchantBean.AgencyId == agency.ID { // 相关机构管理员可以操作
		return nil
	}

	currentUserBean := ctx.Get(conf.ContextTagUser).(*account.UserBean)
	// 当前账号是否是此merchant下面关联的管理账号
	associateBean, err := merchant.GetUserMerchantAssociateByMerchantIdAndUserId(models.DB(), ctx,
		merchantId, currentUserBean.ID)
	if err != nil {
		return err
	}
	if associateBean == nil || associateBean.Role != string(conf.RoleAdmin) {
		return errors.New("no permission for the merchant")
	}

	return nil
}
