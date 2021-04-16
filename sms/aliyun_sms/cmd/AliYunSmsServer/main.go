package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"

	"github.com/njmdk/common/logger"
	"github.com/njmdk/common/network/reuseport"
	aliyunsms "github.com/njmdk/common/sms/aliyun_sms"
)

func usage() {
	fmt.Printf("AliYunSmsServer <listen addr>")
	os.Exit(-1)
}

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	log, err := logger.New("AliYunSmsServer", "./log", zap.DebugLevel, true)
	if err != nil {
		panic(err)
	}
	l, err := reuseport.Listen("tcp4", os.Args[1])
	if err != nil {
		log.Panic("listen failed", zap.Error(err), zap.String("listen addr", os.Args[1]))
	}
	e := echo.New()
	e.Use(middleware.Logger())
	e.HideBanner = true
	e.HidePort = true
	e.Listener = l
	handlerSMSFunc := func(c echo.Context) error {
		resp := &aliyunsms.Response{}
		defer func() {
			_ = c.JSON(http.StatusOK, resp)
		}()
		accessKeyId := c.QueryParam("access_key_id")
		accessKeySecret := c.QueryParam("access_key_secret")
		phone := c.QueryParam("phone")
		templateCode := c.QueryParam("template_code")
		code := c.QueryParam("code")
		signName := c.QueryParam("sign_name")
		if accessKeyId == "" {
			resp.Status = "invalid.request.param.access_key_id"
			resp.Message = "无效的请求参数:access_key_id"
			return nil
		}
		if accessKeySecret == "" {
			resp.Status = "invalid.request.param.access_key_secret"
			resp.Message = "无效的请求参数:access_key_secret"
			return nil
		}
		if phone == "" {
			resp.Status = "invalid.request.param.phone"
			resp.Message = "无效的请求参数:phone"
			return nil
		}
		if templateCode == "" {
			resp.Status = "invalid.request.param.template_code"
			resp.Message = "无效的请求参数:template_code"
			return nil
		}
		if code == "" {
			resp.Status = "invalid.request.param.code"
			resp.Message = "无效的请求参数:code"
			return nil
		}
		if signName == "" {
			resp.Status = "invalid.request.param.sign_name"
			resp.Message = "无效的请求参数:sign_name"
			return nil
		}
		client, err := aliyunsms.NewAliYunSms(accessKeyId, accessKeySecret, log)
		if err != nil {
			resp.Status = "create.aliyun.client.error"
			resp.Message = fmt.Sprintf("创建短信client error:%s", err.Error())
			return nil
		}
		out, _ := client.SendSms(phone, signName, templateCode, code)
		resp = &out
		return nil
	}
	e.POST("/", handlerSMSFunc)
	e.GET("/", handlerSMSFunc)
	go func() {
		log.Info("start aliyun sms server", zap.String("listen addr", os.Args[1]))
		if err := e.Start(""); err != nil {
			log.Fatal("server stop error", zap.Error(err))
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutdown Server ...")
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown", zap.Error(err))
	}
}
