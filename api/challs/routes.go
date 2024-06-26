package challs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	"github.com/gin-gonic/gin"
	"github.com/minpeter/telos/auth"
	"github.com/minpeter/telos/database"
	"github.com/minpeter/telos/utils"
)

func Routes(challRoutes *gin.RouterGroup) {

	challRoutes.GET("", auth.AuthTokenMiddleware(), getChallsHandler)
	challRoutes.GET("/:id/solves", getChallSolvesHandler)
	challRoutes.POST("/:id/submit", auth.AuthTokenMiddleware(), submitChallHandler)

	// dynamic router
	challRoutes.POST("/:id/start", auth.AuthTokenMiddleware(), createChallHandler)
	challRoutes.POST("/:id/stop", auth.AuthTokenMiddleware(), deleteChallHandler)
}

func parseEnv(env string) []string {
	envs := strings.Split(env, ",")
	return envs
}

func createChallHandler(c *gin.Context) {

	cli, err := client.NewClientWithOpts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "docker client error - 1",
		})
		return
	}

	challengeID := c.Param("id")

	if has, err := database.IsDynamic(challengeID); err != nil || !has {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "not dynamic challenge",
		})
		return
	}

	host := strings.Split(c.Request.Host, ":")

	if len(host) == 1 {
		if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" || strings.Contains(c.Request.Referer(), "https") {
			// HTTPS인 경우 443번 포트로 설정
			host = append(host, "443")
		} else {
			// HTTP인 경우 80번 포트로 설정
			host = append(host, "80")
		}
	}

	// get hostname from url

	if challengeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id is empty",
		})
		return
	}

	challengeData, err := database.GetChallengeById(challengeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad challenge id",
		})
		return
	}

	ctx := context.Background()

	imageName := challengeData.Dynamic.Image

	hashId := utils.GenerateId()

	utils.PullImage(imageName)

	config := &container.Config{
		Image: imageName,
		Labels: map[string]string{
			"traefik.enable": "true",
			"traefik.tcp.routers." + hashId + ".rule": "HostSNI(`" + hashId + "." + host[0] + "`)",
			"traefik.tcp.routers." + hashId + ".tls":  "true",
			"dynamic":                                 "true",
		},
		Env: parseEnv(challengeData.Dynamic.Env),
	}

	if host[1] == "443" {
		config.Labels["traefik.tcp.routers."+hashId+".tls"] = "true"
	}

	if challengeData.Dynamic.Type == "http" {
		config.Labels = map[string]string{
			"traefik.enable": "true",
			"traefik.http.routers." + hashId + ".rule": "Host(`" + hashId + "." + host[0] + "`)",
			"dynamic": "true",
		}
		if host[1] == "443" {
			config.Labels["traefik.http.routers."+hashId+".tls"] = "true"
		}

	}

	hostConfig := &container.HostConfig{
		NetworkMode: "traefik",
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, nil, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "docker client error - 2",
		})
		return
	}

	sandboxID := resp.ID

	// Start the container
	if err := cli.ContainerStart(ctx, sandboxID, types.ContainerStartOptions{}); err != nil {
		fmt.Println("Failed to start container:", err) // 에러 메시지 출력
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "docker client error - 3: failed to start container",
		})
		return
	}

	fmt.Println("create sandbox: " + sandboxID[0:12])

	utils.OnlineSandboxIds = append(utils.OnlineSandboxIds, sandboxID[0:12])

	utils.Tq.Enqueue(sandboxID[0:12])

	connection := hashId + "." + host[0]

	if host[1] != "443" {
		connection += ":" + host[1]
		connection = "http://" + connection
	} else {
		connection = "https://" + connection
	}

	if challengeData.Dynamic.Type == "http" {
		utils.SendResponse(c, "goodStartInstance", gin.H{
			"connection": connection,
			"id":         sandboxID[0:12],
			"type":       "http",
		})
		return

	} else {
		utils.SendResponse(c, "goodStartInstance", gin.H{
			"connection": []struct {
				Type    string `json:"type"`
				Command string `json:"command"`
			}{
				{
					Type:    "ncat",
					Command: "ncat --ssl " + hashId + "." + host[0] + " " + host[1],
				},
				{
					Type:    "openssl",
					Command: "openssl s_client -connect " + hashId + "." + host[0] + ":" + host[1],
				},
				{
					Type:    "socat",
					Command: "socat openssl:" + hashId + "." + host[0] + ":" + host[1] + ",verify=0 -",
				},
				{
					Type:    "gnutls",
					Command: "gnutls-cli --insecure " + hashId + "." + host[0] + ":" + host[1],
				},
				{
					Type:    "pwn",
					Command: "remote('" + hashId + "." + host[0] + "', " + host[1] + ", ssl=True)",
				},
			},
			"id":   sandboxID[0:12],
			"type": "tcp",
		})
	}

}

func deleteChallHandler(c *gin.Context) {

	sandboxId := c.Param("id")

	message := utils.RemoveSandbox(sandboxId)

	fmt.Println(message)

	utils.SendResponse(c, "goodStopInstance", gin.H{})
}

func getChallsHandler(c *gin.Context) {

	challs, err := database.GetCleanedChallenges()
	if err != nil {
		utils.SendResponse(c, "internalError", gin.H{})
		return
	}

	if challs == nil {
		challs = []database.CleanedChallenge{}
	}

	utils.SendResponse(c, "goodChallenges", challs)
}

func getChallSolvesHandler(c *gin.Context) {

	c.Status(http.StatusNoContent)
}

func submitChallHandler(c *gin.Context) {

	ChallengeId := c.Param("id")

	var req struct {
		Flag string `json:"flag" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendResponse(c, "badRequest", gin.H{})
		return
	}

	challenge, err := database.GetChallengeById(ChallengeId)
	if err != nil {
		utils.SendResponse(c, "badChallenge", gin.H{})
		return
	}

	fmt.Printf("user submitted flag: %s | correct flag: %s\n", req.Flag, challenge.Flag)

	if req.Flag == challenge.Flag {

		solver := database.Solve{
			Challengeid: ChallengeId,
			Userid:      c.MustGet("userid").(string),
		}

		err := database.NewSolve(solver)
		if err != nil {
			utils.SendResponse(c, "internalError", gin.H{})
			return
		}

		utils.SendResponse(c, "goodFlag", gin.H{})
		return
	}

	utils.SendResponse(c, "badFlag", gin.H{})
}
