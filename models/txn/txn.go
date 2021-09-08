package txn

import (
	"tpayment/models"
	"tpayment/models/payment/acquirer"
	"tpayment/models/payment/record"
	"tpayment/pkg/tlog"

	"gorm.io/gorm"
)

func CreateTransactionAndDetail(txnRecord *record.TxnRecord, detail *record.TxnRecordDetail) error {
	logger := tlog.GetGoroutineLogger()
	return models.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(txnRecord).Create(txnRecord).Error
		if err != nil {
			logger.Error("create detail record error->", err.Error())
			return err
		}

		err = tx.Model(detail).Create(detail).Error
		if err != nil {
			logger.Error("create detail record error->", err.Error())
			return err
		}
		return nil
	})
}

func UpdateSaleResult(txnRecord *record.TxnRecord, detail *record.TxnRecordDetail) error {
	return models.DB.Transaction(func(tx *gorm.DB) error {
		err := txnRecord.UpdateTxnResult(tx)
		if err != nil {
			return err
		}

		if err = detail.Update(tx); err != nil {
			return err
		}

		return nil
	})
}

func CreateAndUpdateKey(createKeys []*acquirer.Key, deleteKeys []*acquirer.Key) error {
	return models.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		for i := 0; i < len(createKeys); i++ {
			err = tx.Create(createKeys[i]).Error
			if err != nil {
				return err
			}
		}

		for i := 0; i < len(deleteKeys); i++ {
			err = tx.Delete(deleteKeys[i]).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
}
