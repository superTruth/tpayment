package standard

import (
	"fmt"
	"tpayment/internal/acquirer_impl"
	"tpayment/models"
	"tpayment/models/payment/acquirer"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func updateKey(ctx *gin.Context, req *acquirer_impl.SaleRequest, resp *acquirer_impl.SaleResponse) error {
	if len(resp.Keys) == 0 {
		return nil
	}
	return models.DB().Transaction(func(tx *gorm.DB) error {
		var (
			err error
		)

		keyBean := acquirer.Key{
			BaseModel: models.BaseModel{
				Db: &models.MyDB{DB: tx},
			},
		}
		//
		keyTag := generateKeyTag(req)

		for i, keyNew := range resp.Keys {
			destKey := findKeyFromArray(keyNew.Type, req.Keys)
			if destKey >= 0 { // 找到，判断旧数据是否相同，如果不同，则删除，再插入
				if req.Keys[destKey].Value == resp.Keys[i].Value {
					continue
				}
				err = req.Keys[destKey].Delete()
				if err != nil {
					return err
				}
			}
			// 插入数据
			resp.Keys[i].Tag = keyTag
			err = keyBean.Create(resp.Keys[i])
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func generateKeyTag(req *acquirer_impl.SaleRequest) string {
	return fmt.Sprintf("%d-%s-%s",
		req.TxqReq.PaymentProcessRule.MerchantAccount.Acquirer.ID,
		req.TxqReq.PaymentProcessRule.MerchantAccount.MID,
		req.TxqReq.PaymentProcessRule.MerchantAccount.Terminal.TID)
}

func findKeyFromArray(destKeyType string, keys []*acquirer.Key) int {
	for i, key := range keys {
		if key.Type == destKeyType {
			return i
		}
	}
	return -1
}
