package bot

import (
	"api/anticaptcha"
	"net/http"

	"github.com/gin-gonic/gin"
	"wyvern.pw/api/cloudproxy"
	"wyvern.pw/controllers/setting"
)

var (
	maxWorker     = 300
	currentWorker = 0
)

func GetCurrentWorker() int {
	return currentWorker
}

func SolveChallenge(c *gin.Context) {
	if currentWorker >= maxWorker {
		c.JSON(http.StatusTooManyRequests, SolveResponse{
			Code:    http.StatusTooManyRequests,
			Message: "Wait some times...",
		})
		c.Abort()
		return
	}

	currentWorker++
	var request SolveRequest
	err := c.BindJSON(&request)
	if err != nil {
		currentWorker--
		c.JSON(http.StatusInternalServerError, SolveResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		c.Abort()
		return
	}

	apiKey, err := setting.GetApiKey()
	if err != nil {
		currentWorker--
		c.JSON(http.StatusInternalServerError, SolveResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		c.Abort()
		return
	}
	captchaKey, err := anticaptcha.SolveCaptcha()
	currentWorker--
	c.JSON(http.StatusOK, SolveResponse{
		Code:    http.StatusOK,
		Message: "",
		CaptchaKey: captchaKey,
	})
	c.Abort()
	return
}
