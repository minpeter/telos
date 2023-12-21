package challs

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minpeter/rctf-backend/utils"
)

func Routes(challRoutes *gin.RouterGroup) {

	challRoutes.GET("", getChallsHandler)
	challRoutes.GET("/:id/solves", getChallSolvesHandler)
	challRoutes.POST("/:id/submit", submitChallHandler)

}

func getChallsHandler(c *gin.Context) {

	utils.SendResponse(c, "goodChallenges", []gin.H{
		{
			"files":       []string{},
			"description": "This is a good challenge",
			"author":      "minpeter",
			"points":      100,
			"id":          "34344543-3453-345-5344-34534534534534",
			"name":        "Good Challenge",
			"category":    "pwn",
			"solves":      2,
		},
	})
}

func getChallSolvesHandler(c *gin.Context) {

	c.Status(http.StatusNoContent)
}

func submitChallHandler(c *gin.Context) {

	utils.SendResponse(c, "badEnded", gin.H{})
}