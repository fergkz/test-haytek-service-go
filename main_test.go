package main

import (
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"

	ApplicationContractService "github.com/fergkz/test-haytek-service-go/src/Application/Contract/Service"
	ApplicationUsecase "github.com/fergkz/test-haytek-service-go/src/Application/Usecase"
	InfrastructureService "github.com/fergkz/test-haytek-service-go/src/Infrastructure/Service"
)

func saveToFile(filename string, data interface{}) {
	jsonData := spew.Sdump(data)
	os.WriteFile("test-results/"+filename, []byte(jsonData), 0644)
}

func TestExternalServices(t *testing.T) {
	t.Run("Address Service", func(t *testing.T) {
		t.Log("Testing Address Service...")
		var serviceAddress ApplicationContractService.Address
		serviceAddress = InfrastructureService.NewHaytekAddress()
		addRows, err := serviceAddress.ListAllAddresses()
		if err != nil {
			t.Fatalf("\x1b[31mError:\x1b[0m %v", err)
		}
		saveToFile("TestExternalServices-Addresses.log", addRows)
		t.Log("\x1b[32mAddress Service Test Passed!\x1b[0m")
	})

	t.Run("Box Service", func(t *testing.T) {
		t.Log("Testing Box Service...")
		var serviceBox ApplicationContractService.Box
		serviceBox = InfrastructureService.NewHaytekBox()
		boxRows, err := serviceBox.ListAllBoxes()
		if err != nil {
			t.Fatalf("\x1b[31mError:\x1b[0m %v", err)
		}
		saveToFile("TestExternalServices-Boxes.log", boxRows)
		t.Log("\x1b[32mBox Service Test Passed!\x1b[0m")
	})

	t.Run("Carrier Service", func(t *testing.T) {
		t.Log("Testing Carrier Service...")
		var serviceCarrier ApplicationContractService.Carrier
		serviceCarrier = InfrastructureService.NewHaytekCarrier()
		carrierRows, err := serviceCarrier.ListAllCarriers()
		if err != nil {
			t.Fatalf("\x1b[31mError:\x1b[0m %v", err)
		}
		saveToFile("TestExternalServices-Carriers.log", carrierRows)
		t.Log("\x1b[32mCarrier Service Test Passed!\x1b[0m")
	})

	t.Run("Order Service", func(t *testing.T) {
		t.Log("Testing Order Service...")
		var serviceOrder ApplicationContractService.Order
		serviceOrder = InfrastructureService.NewHaytekOrder()
		orderRows, err := serviceOrder.ListAllOrders()
		if err != nil {
			t.Fatalf("\x1b[31mError:\x1b[0m %v", err)
		}
		saveToFile("TestExternalServices-Orders.log", orderRows)
		t.Log("\x1b[32mOrder Service Test Passed!\x1b[0m")
	})
}

func TestGroupByDelivery(t *testing.T) {

	t.Run("Group By Delivery", func(t *testing.T) {
		t.Log("Testing Group By Delivery...")

		usecase := ApplicationUsecase.NewGroupByDelivery(
			InfrastructureService.NewHaytekAddress(),
			InfrastructureService.NewHaytekBox(),
			InfrastructureService.NewHaytekCarrier(),
			InfrastructureService.NewHaytekOrder(),
		)

		packs, err := usecase.Run()
		if err != nil {
			t.Fatalf("\x1b[31mError:\x1b[0m %v", err)
		}

		if len(packs) == 0 {
			t.Fatalf("\x1b[31mError:\x1b[0m No packages found")
		}

		t.Log("\x1b[32mGroup By Delivery Test Passed!\x1b[0m")

	})

}
