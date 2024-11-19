package handlers

import (
	"net/http"
	"time"

	"github.com/AndreHeber/go-sqlite-blog/middleware"
)

func Health(a *middleware.Adapter) error {
	a.ResponseWriter.WriteHeader(http.StatusOK)
	return nil
}

func TimeConsumingHandler(a *middleware.Adapter) error {
	time.Sleep(5 * time.Second)
	return nil
}
