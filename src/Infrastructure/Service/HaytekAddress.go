package InfrastructureService

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	DomainEntity "github.com/fergkz/test-haytek-service-go/src/Domain/Entity"
)

type haytekAddress struct {
}

func NewHaytekAddress() *haytekAddress {
	service := new(haytekAddress)
	return service
}

type HaytekAddressDaoAddressRow struct {
	Id           string `json:"id"`
	State        string `json:"state"`
	ZipCode      string `json:"zipcode"`
	Street       string `json:"street"`
	Complement   string `json:"complement"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
}

func (service *haytekAddress) ListAllAddresses() ([]DomainEntity.Address, error) {

	var addresses []DomainEntity.Address

	req, err := http.NewRequest("GET", "https://stg-api.haytek.com.br/api/v1/test-haytek-api/adresses", nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return addresses, err
	}

	client := &http.Client{
		Timeout: 360 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return addresses, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return addresses, fmt.Errorf("bad status: %s", resp.Status)
	}

	daoResponse := make([]HaytekAddressDaoAddressRow, 0)

	if err := json.NewDecoder(resp.Body).Decode(&daoResponse); err != nil {
		return addresses, err
	}

	for _, daoRow := range daoResponse {
		address := *DomainEntity.NewAddress(
			daoRow.Id,
			daoRow.State,
			*DomainEntity.NewAddressZipCode(daoRow.ZipCode),
			daoRow.Street,
			daoRow.Complement,
			daoRow.Neighborhood,
			daoRow.City,
		)

		addresses = append(addresses, address)
	}

	return addresses, nil
}
