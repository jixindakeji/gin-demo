package main

import (
	"fmt"
	"gin-demo/pkg/setting"
	"gin-demo/routers"
	"net/http"
	"time"
)

func main() {
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.GConfig.Server.Port),
		Handler:        router,
		ReadTimeout:    setting.GConfig.Server.ReadTimeout * time.Second,
		WriteTimeout:   setting.GConfig.Server.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
