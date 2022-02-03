package rest

import (
	"context"
	"dcamachoj/time-tracker-rest/common"
	"dcamachoj/time-tracker-rest/dbx"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type ckRest string

var employeeEntity = &dbx.EmployeeEntity{}
var timeDataEntity = &dbx.TimeDataEntity{}

const ckGin = ckRest("gin")

func withGin(ctx context.Context, gCtx *gin.Context) context.Context {
	return context.WithValue(ctx, ckGin, gCtx)
}

func getGin(ctx context.Context) *gin.Context {
	var gCtx, ok = ctx.Value(ckGin).(*gin.Context)
	if !ok {
		return nil
	}
	return gCtx
}

func ToRest(hnd common.Resthandler) gin.HandlerFunc {
	return gin.HandlerFunc(func(gCtx *gin.Context) {
		var ctx = context.Background()
		ctx = withGin(ctx, gCtx)
		var res = safeHandler(ctx, hnd)
		gCtx.JSON(res.Status, res)
	})
}
func safeHandler(ctx context.Context, hnd common.Resthandler) (res *common.Response) {
	defer func() {
		if rErr := common.RecoverError(recover(), "handle Rest"); rErr != nil {
			res = common.WrapResponse(ctx, rErr)
		}
		if res == nil {
			res = common.ResponseOK(ctx)
		}
	}()
	return hnd(ctx)
}

func ExecuteServer() {
	var r = gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	var api = r.Group("/api")
	var v0 = api.Group("/v0")

	cfgEmployee(v0.Group("/employee"))

	var addr = common.GetEnv("PORT", ":8081")
	if !strings.HasPrefix(addr, ":") {
		addr = ":" + addr
	}
	var srv = &http.Server{
		Addr:    addr,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Printf("Listening on %s\r\n", addr)
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
