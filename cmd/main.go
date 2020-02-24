package main

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"github.com/jpillora/backoff"

	"github.com/theonlyrob/vercer/webserver/api"
	//_ "github.com/theonlyrob/vercer/webserver/services/user/client.go"

	globalhandler "github.com/theonlyrob/vercer/webserver/cmd/handler"
	globalindex "github.com/theonlyrob/vercer/webserver/cmd/index"
	globalserver "github.com/theonlyrob/vercer/webserver/cmd/server"
	globalstore "github.com/theonlyrob/vercer/webserver/cmd/store"
)

func main() {
	// Run global server.
	_ = globalserver.Singleton().Run(globalhandler.Singleton())
}
