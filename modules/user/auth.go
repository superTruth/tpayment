package user

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/http"
	"tpayment/conf"
	"tpayment/models"
	"tpayment/models/account"
	"tpayment/pkg/algorithmutils"
	"tpayment/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context) (*account.UserBean, *account.AppIdBean, error) {
	tokens := ctx.Request.Header[conf.HeaderTagToken]
	return AuthByToken(ctx, tokens[0])

	//if len(tokens) != 0 {
	//	return AuthByToken(ctx, tokens[0])
	//}
	//
	//return AuthByAccessKey(ctx)
}

// token验证
func AuthByToken(ctx *gin.Context, token string) (*account.UserBean, *account.AppIdBean, error) {
	logger := tlog.GetLogger(ctx)

	// 创建 或者 更新  token记录
	tokenBean, err := account.GetTokenBeanByToken(models.DB(), ctx, token)
	if err != nil {
		logger.Warn("GetTokenBeanByToken fail->", err.Error())
		return nil, nil, err
	}

	if tokenBean == nil { // 没有对应的token记录
		return nil, nil, nil
	}

	//
	accountBean, err := account.GetUserById(models.DB(), ctx, tokenBean.UserId)
	if err != nil {
		logger.Warn("GetUserById fail->", err.Error())
		return nil, nil, err
	}

	//appBean, err := account.GetAppIdByID(models.DB(), ctx, tokenBean.AppId)
	//if err != nil {
	//	logger.Warn("GetUserById fail->", err.Error())
	//	return nil, nil, err
	//}

	return accountBean, nil, nil
}

// accessKey验证
func AuthByAccessKey(ctx *gin.Context) (*account.UserBean, *account.AppIdBean, error) {
	logger := tlog.GetLogger(ctx)
	keys := ctx.Request.Header[conf.HeaderTagAccessKey]
	if len(keys) == 0 {
		return nil, nil, nil
	}

	hashes := ctx.Request.Header[conf.HeaderTagAccessHash]
	if len(hashes) != 0 {
		return nil, nil, nil
	}

	// 获取相关access
	accessBean, err := account.GetUserAccessKeyFromKey(models.DB(), ctx, keys[0])
	if err != nil {
		return nil, nil, err
	}

	secretByte, err := hex.DecodeString(accessBean.Secret)
	if err != nil {
		logger.Error("DecodeString fail access key id->", accessBean.ID)
		return nil, nil, err
	}

	// 计算hash
	method := ctx.Request.Method
	bodyStr := dumpBodyContent(ctx)

	calcHash := algorithmutils.Hmac(secretByte, []byte(method+bodyStr))

	if calcHash != hashes[0] {
		logger.Warn("hash validate fail")
		return nil, nil, errors.New("hash validate fail")
	}

	// 获取真正的账号数据
	userBean, err := account.GetUserById(models.DB(), ctx, accessBean.UserId)
	if err != nil {
		logger.Error("account.GetUserById->", err.Error())
		return nil, nil, err
	}

	return userBean, nil, nil
}

// 拷贝body数据
func dumpBodyContent(ctx *gin.Context) string {
	if ctx.Request.Body == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return ""
	}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(ctx.Request.Body); err != nil {
		return ""
	}
	if err := ctx.Request.Body.Close(); err != nil {
		return ""
	}

	bodyData := buf.Bytes()
	ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(bodyData))

	return string(bodyData)
}
