package main

import (
    "github.com/ant0ine/go-json-rest/rest"
    //"log"
    log "github.com/Sirupsen/logrus"
    "net/http"
)

var controlChannel chan int

func main() {
    
    log.Info("Starting SmartVendProducer IoT Device")
    controlChannel = make(chan int)
    startProducer(controlChannel)

    log.Info("Starting SmartVend API")
	startAPI()

}

func startAPI() {

	api := rest.NewApi()
    api.Use(rest.DefaultDevStack...)

    router, err := rest.MakeRouter(
    	rest.Get("/test", runTest),
    	rest.Get("/stop", stopServices),
    	rest.Get("/services", getServices),
    )

    if err != nil {
    	log.Fatal(err)
    }

    api.SetApp(router)
    log.Fatal(http.ListenAndServe(":4242", api.MakeHandler()))
}

