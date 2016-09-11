package main 

import (
    "github.com/ant0ine/go-json-rest/rest"
    log "github.com/Sirupsen/logrus"
    //"strconv"
    "fmt"
    
)

 //var consumer devclienttypes.DeviceProfile

func runTest(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(map[string]string{"Body": "TEST RESPONSE"})
}

func stopServices(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(map[string]string{"Body": "Stopping Service Broadcast"})

	controlChannel <- 0
}

func getServices(w rest.ResponseWriter, r *rest.Request) {

	log.Info("Getting Services")

	if sdk == nil {
		w.WriteJson(map[string]string{"Body": "The vending machine is down"})
		return
	}

	// initiate new card
	// TODO: Pass the ID in the REST request
	
	// IDEA: Could restrict services offered based on card holder detail or give VIP / Sale pricing

	consumerCard := newCard(1)

	// Hard code ip for now.  ??? Still might use discovery here instead of connecting directly ??? Also not sure about that device ID
	// TODO: Pass in local IP on the fly 

	// IDEA: Could have a network of smart vending machines that can all talk to each other and redirect consumers to other machines
	
	err := sdk.InitConsumer("http://", "192.168.10.100", 8080, "", "188aedcf-f294-45ad-4a1b-34495441ce7b", consumerCard)

	serviceDetails, err := sdk.RequestServices()

	if err == nil {
		
		fmt.Println(len(serviceDetails))
		for _, serviceDetail := range serviceDetails {
			log.Info(serviceDetail.ServiceID)
			log.Info(serviceDetail.ServiceDescription)
			//fmt.Printf("%d - %s\n", serviceDetail.ServiceID, serviceDetail.ServiceDescription)
					
		}

	} else {
		log.Error(err)
		w.WriteJson(map[string]error{"Body": err})
	}

	/*
	services, err := sdk.DeviceDiscovery(1000)

	if err == nil {
		serviceList := "Services: \n"
		for _, svc := range services {

			serviceList += svc.Hostname + " " + strconv.Itoa(svc.PortNumber) + " " + svc.UrlPrefix + " Description: " + svc.DeviceDescription + "\n"
		}

		w.WriteJson(map[string]string{"Body": serviceList})

	} else {
		log.Error(err)
		w.WriteJson(map[string]error{"Body": err})

	} */
	


}
