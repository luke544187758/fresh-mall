package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"luke544187758/oss-web/settings"
	"luke544187758/oss-web/utils"
	"net/http"
	"net/url"
	"strings"
)

func Token(ctx *gin.Context) {
	res := utils.Get_policy_token()
	ctx.Header("Access-Control-Allow-Methods", "POST")
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.String(http.StatusOK, res)
}

func HandlerRequest(ctx *gin.Context) {
	fmt.Println("\nHandle Post Request ...")

	// Get PublicKey bytes
	bytePublicKey, err := utils.GetPublicKey(ctx)
	if err != nil {
		utils.ResponseFailed(ctx)
		return
	}

	// Get Authorization bytes : decode from Base64String
	byteAuthorization, err := utils.GetAuthorization(ctx)
	if err != nil {
		utils.ResponseFailed(ctx)
		return
	}

	// Get MD5 bytes from Newly Constructed Authrization String.
	byteMD5, bodyStr, err := utils.GetMD5FromNewAuthString(ctx)
	if err != nil {
		utils.ResponseFailed(ctx)
		return
	}

	decodeurl, err := url.QueryUnescape(bodyStr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(decodeurl)
	params := make(map[string]string)
	datas := strings.Split(decodeurl, "&")
	for _, v := range datas {
		sdatas := strings.Split(v, "=")
		fmt.Println(v)
		params[sdatas[0]] = sdatas[1]
	}
	fileName := params["filename"]
	fileUrl := fmt.Sprintf("%s/%s", settings.Conf.OssService.Host, fileName)

	// verifySignature and response to client
	if utils.VerifySignature(bytePublicKey, byteMD5, byteAuthorization) {
		// do something you want accoding to callback_body ...
		ctx.JSON(http.StatusOK, gin.H{
			"url": fileUrl,
		})
	} else {
		utils.ResponseFailed(ctx) // response FAILED : 400
	}
}
