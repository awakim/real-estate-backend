package api

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func getIP(ctx *gin.Context) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := ctx.GetHeader("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := ctx.GetHeader("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(ctx.Request.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("no valid ip found")
}

func (server *Server) loginRateLimiter(ctx *gin.Context) {

	IP, err := getIP(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, errorResponse(err))
		return
	}

	b, err := server.Cache.IsRateLimited(ctx, IP)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if b {
		ctx.AbortWithStatusJSON(
			http.StatusTooManyRequests,
			errorResponse(errors.New("too many requests. Try again in 15 minutes")))
		return
	}

	ctx.Next()
}
