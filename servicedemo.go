package main

import (
	"github.com/lypwonderful/servicedemo/pkgs/core/dubbogo"
	"github.com/lypwonderful/servicedemo/pkgs/usercenter"
)

func main() {
	initService()

	dubbogo.DubboServiceRun()
}

func initService() {
	servo := dubbogo.InitDubbo()
	err := servo.Handle(&usercenter.UserProvider{})
	if err != nil {
		panic(err)
		return
	}
}
