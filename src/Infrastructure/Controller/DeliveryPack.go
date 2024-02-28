package InfrastructureController

import (
	"encoding/json"
	"net/http"

	ApplicationContractService "github.com/fergkz/test-haytek-service-go/src/Application/Contract/Service"
	ApplicationUsecase "github.com/fergkz/test-haytek-service-go/src/Application/Usecase"
)

type deliveryPack struct {
	ServiceAddress ApplicationContractService.Address
	ServiceBox     ApplicationContractService.Box
	ServiceCarrier ApplicationContractService.Carrier
	ServiceOrder   ApplicationContractService.Order
}

func NewDeliveryPack(
	ServiceAddress ApplicationContractService.Address,
	ServiceBox ApplicationContractService.Box,
	ServiceCarrier ApplicationContractService.Carrier,
	ServiceOrder ApplicationContractService.Order,
) *deliveryPack {
	controller := new(deliveryPack)
	controller.ServiceAddress = ServiceAddress
	controller.ServiceBox = ServiceBox
	controller.ServiceCarrier = ServiceCarrier
	controller.ServiceOrder = ServiceOrder
	return controller
}

func (controller *deliveryPack) Get(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

	usecase := ApplicationUsecase.NewGroupByDelivery(
		controller.ServiceAddress,
		controller.ServiceBox,
		controller.ServiceCarrier,
		controller.ServiceOrder,
	)

	packages, err := usecase.Run()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type responseDaoTypeBoxOrder struct {
		Id string `json:"id"`
	}

	type responseDaoTypeBox struct {
		BoxType       string                    `json:"boxType"`
		ItemsQuantity int                       `json:"itemsQuantity"`
		Orders        []responseDaoTypeBoxOrder `json:"orders"`
	}

	type responseDaoType struct {
		DeliveryDate string `json:"deliveryDate"`
		CarrierId    struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"carrier"`
		Address struct {
			Street       string `json:"street"`
			Complement   string `json:"complement"`
			Neighborhood string `json:"neighborhood"`
			ZipCode      string `json:"zipCode"`
			City         string `json:"city"`
			State        string `json:"state"`
		} `json:"address"`
		Boxes []responseDaoTypeBox `json:"boxes"`
	}

	var responseDao []responseDaoType

	for _, pack := range packages {
		var packDao responseDaoType

		zipCode := pack.Address.GetZipCode()

		packDao.DeliveryDate = pack.DeliveryDate.Format("2006-01-02")
		packDao.CarrierId.Id = pack.Carrier.GetId()
		packDao.CarrierId.Name = pack.Carrier.GetName()
		packDao.Address.Street = pack.Address.GetStreet()
		packDao.Address.Complement = pack.Address.GetComplement()
		packDao.Address.Neighborhood = pack.Address.GetNeighborhood()
		packDao.Address.ZipCode = zipCode.GetCode()
		packDao.Address.City = pack.Address.GetCity()
		packDao.Address.State = pack.Address.GetState()

		for _, box := range pack.BoxPackages {
			var boxDao responseDaoTypeBox
			boxDao.BoxType = box.Box.GetBoxType()
			boxDao.ItemsQuantity = box.QuantityInBox
			for _, order := range box.Orders {
				var orderDao responseDaoTypeBoxOrder
				orderDao.Id = order.Order.GetId()
				boxDao.Orders = append(boxDao.Orders, orderDao)
			}
			packDao.Boxes = append(packDao.Boxes, boxDao)
		}

		responseDao = append(responseDao, packDao)
	}

	jsonBytes, err := json.Marshal(responseDao)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
