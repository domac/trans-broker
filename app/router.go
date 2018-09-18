package app

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

const (
	API_POST_DATAPUSH = "/data/push"
	API_GET_VERSION   = "/version"
)

func pushHandler(c *gin.Context) {
	//读取提交的报文数据
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		getApplicationContext().Log().Error(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"ret_code": -1,
			"ret_data": "bad request",
		})
	}

	//数据异步处理
	getApplicationContext().Put(NewPacket(b))

	c.JSON(http.StatusOK, gin.H{
		"ret_code": 0,
		"ret_msg":  "push success",
		"ret_data": "",
	})
}

func versionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": GetVersion(),
	})
}

func loadRouter(r *gin.Engine) {
	r.GET(API_GET_VERSION, versionHandler)
	r.POST(API_POST_DATAPUSH, pushHandler)
}
