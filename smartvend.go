package main

import (
	//"encoding/json"
	"errors"
	//"fmt"
	//"io/ioutil"

	log "github.com/Sirupsen/logrus"
	devclienttypes "github.com/wptechinnovation/worldpay-within-sdk/applications/dev-client/types"
	"github.com/wptechinnovation/worldpay-within-sdk/sdkcore/wpwithin"
	//"github.com/wptechinnovation/worldpay-within-sdk/sdkcore/wpwithin/rpc"
	"github.com/wptechinnovation/worldpay-within-sdk/sdkcore/wpwithin/types"
	"time"
)

var sdk wpwithin.WPWithin
var hceCard types.HCECard

func newCard(cardHolderID int32) *types.HCECard {

	// TODO: Read this from a simple DB based on cardHolderID

	var firstName, lastName, cardNumber, cardType, cvc string
	var expYear, expMonth int32

	if cardHolderID == 1 {

		firstName 	= "Joe"
		lastName 	= "Bloggs"
		expMonth  	= 12
		expYear     = 2017
		cardNumber  = "5555555555554444"
		cardType    = "Card"
		cvc 		= "123"

	}

	card := types.HCECard{

		FirstName:  firstName,
		LastName:   lastName,
		ExpMonth:   expMonth,
		ExpYear:    expYear,
		CardNumber: cardNumber,
		Type:       cardType,
		Cvc:        cvc,
	}

	return &card 
}

func broadcast(timeout int, period int, quit chan int) error {
	if sdk == nil {
		return errors.New(devclienttypes.ErrorDeviceNotInitialised)
	}


	for {
		select { 
			case <- quit:
				log.Info("Stopping Service Broadcast")
				return nil

			default:

				if err := sdk.StartServiceBroadcast(timeout); err != nil {
					return err
				}

				time.Sleep(time.Duration(period))
		
		}
	}
}

func startProducer(controlChannel chan int) error {
	_sdk, err := wpwithin.Initialise("SmartVend", "A proximity aware vending machine")

	if err != nil {
		return err
	}

	sdk = _sdk

	merchantClientKey  := "T_C_03eaa1d3-4642-4079-b030-b543ee04b5af"
	merchantServiceKey := "T_S_f50ecb46-ca82-44a7-9c40-421818af5996"


	_ = sdk.InitProducer(merchantClientKey, merchantServiceKey)

	// Add some services

	// TODO: Read service data from external config

	smartVend, _ := types.NewService()
	smartVend.Name = "SmartVend"
	smartVend.Description = "Smart Vending Machine"
	smartVend.Id = 1

	snickersPrice := types.Price{

		UnitID:          1,
		ID:              1,
		Description:     "Snickers",
		UnitDescription: "Snickers Chocolate bar",
		PricePerUnit: &types.PricePerUnit{
			Amount:       50,
			CurrencyCode: "GBP",
		},
	}

	cokePrice := types.Price{

		UnitID:          1,
		ID:              1,
		Description:     "Coke",
		UnitDescription: "Coca-Cola Soft Drink",
		PricePerUnit: &types.PricePerUnit{
			Amount:       60,
			CurrencyCode: "GBP",
		},
	}
	
	smartVend.AddPrice(snickersPrice)
	smartVend.AddPrice(cokePrice)

	if err := sdk.AddService(smartVend); err != nil {
		return err
	}

	// Start Broadcasting 
	go broadcast(60000, 10, controlChannel)

	return nil
}


