package InfrastructureService

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	DomainEntity "github.com/fergkz/test-haytek-service-go/src/Domain/Entity"
)

type haytekBox struct {
}

func NewHaytekBox() *haytekBox {
	service := new(haytekBox)
	return service
}

type HaytekBoxDaoBoxRow struct {
	Type        string `json:"type"`
	MaxQuantity string `json:"maxQuantity"`
}

func (service *haytekBox) ListAllBoxes() ([]DomainEntity.Box, error) {

	var boxes []DomainEntity.Box

	req, err := http.NewRequest("GET", "https://stg-api.haytek.com.br/api/v1/test-haytek-api/boxes", nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return boxes, err
	}

	client := &http.Client{
		Timeout: 360 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return boxes, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return boxes, fmt.Errorf("bad status: %s", resp.Status)
	}

	daoResponse := make([]HaytekBoxDaoBoxRow, 0)

	if err := json.NewDecoder(resp.Body).Decode(&daoResponse); err != nil {
		return boxes, err
	}

	for _, daoRow := range daoResponse {
		maxQuantity, err := strconv.Atoi(daoRow.MaxQuantity)
		if err != nil {
			return boxes, err
		}

		address := *DomainEntity.NewBox(
			daoRow.Type,
			maxQuantity,
		)

		boxes = append(boxes, address)
	}

	return boxes, nil
}
