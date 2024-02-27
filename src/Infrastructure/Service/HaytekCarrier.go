package InfrastructureService

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	DomainEntity "github.com/fergkz/test-haytek-service-go/src/Domain/Entity"
)

type haytekCarrier struct {
}

func NewHaytekCarrier() *haytekCarrier {
	service := new(haytekCarrier)
	return service
}

type HaytekCarrierDaoCarrierRow struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	CutOfftime string `json:"cutOfftime"`
}

func (service *haytekCarrier) ListAllCarriers() ([]DomainEntity.Carrier, error) {

	var carriers []DomainEntity.Carrier

	req, err := http.NewRequest("GET", "https://stg-api.haytek.com.br/api/v1/test-haytek-api/carriers", nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return carriers, err
	}

	client := &http.Client{
		Timeout: 360 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return carriers, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return carriers, fmt.Errorf("bad status: %s", resp.Status)
	}

	daoResponse := make([]HaytekCarrierDaoCarrierRow, 0)

	if err := json.NewDecoder(resp.Body).Decode(&daoResponse); err != nil {
		return carriers, err
	}

	for _, daoRow := range daoResponse {
		cutSplit := strings.Split(daoRow.CutOfftime, ":")
		hour, err := strconv.Atoi(cutSplit[0])
		if err != nil {
			return carriers, err
		}
		minute, err := strconv.Atoi(cutSplit[1])
		if err != nil {
			return carriers, err
		}

		carrier := *DomainEntity.NewCarrier(
			daoRow.Id,
			daoRow.Name,
			*DomainEntity.NewCarrierCutOffTime(hour, minute),
		)

		carriers = append(carriers, carrier)
	}

	return carriers, nil
}
