package InfrastructureService

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	DomainEntity "github.com/fergkz/test-haytek-service-go/src/Domain/Entity"
)

type haytekOrder struct {
}

func NewHaytekOrder() *haytekOrder {
	service := new(haytekOrder)
	return service
}

type HaytekOrderDaOrderRow struct {
	Id        string      `json:"id"`
	AddressId string      `json:"addressId"`
	CarrierId string      `json:"carrierId"`
	Quantity  int         `json:"quantity"`
	CreatedAt interface{} `json:"createdAt"`
}

func (service *haytekOrder) ListAllOrders() ([]DomainEntity.Order, error) {

	var orders []DomainEntity.Order

	req, err := http.NewRequest("GET", "https://stg-api.haytek.com.br/api/v1/test-haytek-api/orders", nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return orders, err
	}

	client := &http.Client{
		Timeout: 360 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return orders, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return orders, fmt.Errorf("bad status: %s", resp.Status)
	}

	daoResponse := make([]HaytekOrderDaOrderRow, 0)

	if err := json.NewDecoder(resp.Body).Decode(&daoResponse); err != nil {
		return orders, err
	}

	for _, daoRow := range daoResponse {

		var createdAt time.Time

		if daoRow.CreatedAt != nil {
			createdAt, err = time.Parse("2006-01-02T15:04:05.000Z", daoRow.CreatedAt.(string))
			if err != nil {
				return orders, err
			}
		}

		order := *DomainEntity.NewOrder(
			daoRow.Id,
			daoRow.AddressId,
			daoRow.CarrierId,
			createdAt,
			daoRow.Quantity,
		)

		orders = append(orders, order)
	}

	return orders, nil
}
