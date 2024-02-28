package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
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

	config := new(Config)
	config.Load("config.yml")

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

	t.Run("Group By Delivery - API", func(t *testing.T) {
		t.Log("Testing Group By Delivery API...")

		testAPICall(t, "GET", "http://localhost:"+config.Server.Port+"/v1/delivery-pack", "", nil, 200, func(payload string) bool {
			return len(payload) > 0
		})

		t.Log("\x1b[32mGroup By Delivery API Test Passed!\x1b[0m")
	})

}

func testAPICall(t *testing.T, method, url, payload string, headers map[string]string, expectedStatus int, expectedPayload interface{}) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		t.Fatalf("\x1b[31mErro ao criar requisição: %v\x1b[0m", err)
	}

	// Adiciona headers à requisição
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Faz a requisição
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("\x1b[31mErro ao realizar requisição: %v\x1b[0m", err)
	}
	defer resp.Body.Close()

	// Verifica o código de status
	if expectedStatus != 0 {
		if resp.StatusCode != expectedStatus {
			t.Fatalf("\x1b[31mCódigo de status inesperado. Esperado: %d, Obtido: %d\x1b[0m", expectedStatus, resp.StatusCode)
		}
	}

	// Converte o corpo da resposta em string
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("\x1b[31mErro ao converter corpo da resposta para string: %v\x1b[0m", err)
	}

	// Remove espaços do payload para comparação
	bodyStr := string(body)

	removeSpaces := func(s string) string {
		return fmt.Sprintf("%s", bytes.Replace([]byte(s), []byte("\n"), []byte(""), -1))
	}

	bodyStr = removeSpaces(bodyStr)
	// log.Println("RESPONSE", resp.StatusCode, bodyStr)

	switch expectedPayload := expectedPayload.(type) {
	case string:
		bodyStr = removeSpaces(bodyStr)
		expectedPayloadStr := removeSpaces(expectedPayload)

		// Compara os payloads
		if !reflect.DeepEqual(expectedPayloadStr, bodyStr) {
			t.Fatalf("\x1b[31mPayload da resposta inesperado. Esperado: %s, Obtido: %s\x1b[0m", expectedPayloadStr, bodyStr)
		}
	case func(string) bool:
		if !expectedPayload(bodyStr) {
			t.Fatalf("\x1b[31mFunção de comparação do payload retornou falso\x1b[0m")
		}
	case func(string) error:
		if err := expectedPayload(bodyStr); err != nil {
			t.Fatalf("\x1b[31m%v\x1b[0m", err)
		}
	}
}
