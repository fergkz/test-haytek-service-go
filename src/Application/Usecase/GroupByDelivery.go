package ApplicationUsecase

import (
	"sort"
	"time"

	ApplicationContractService "github.com/fergkz/test-haytek-service-go/src/Application/Contract/Service"
	DomainEntity "github.com/fergkz/test-haytek-service-go/src/Domain/Entity"
)

type groupByDelivery struct {
	ServiceAddress ApplicationContractService.Address
	ServiceBox     ApplicationContractService.Box
	ServiceCarrier ApplicationContractService.Carrier
	ServiceOrder   ApplicationContractService.Order
}

type DeliveryPackageBoxOrderRow struct {
	Order    DomainEntity.Order
	Quantity int
}

type DeliveryPackageBox struct {
	Box           DomainEntity.Box
	QuantityInBox int
	Orders        []DeliveryPackageBoxOrderRow
}

type DeliveryPackage struct {
	DeliveryDate time.Time
	Carrier      DomainEntity.Carrier
	Address      DomainEntity.Address
	BoxPackages  []DeliveryPackageBox
	AllOrders    []DomainEntity.Order
}

func NewGroupByDelivery(
	ServiceAddress ApplicationContractService.Address,
	ServiceBox ApplicationContractService.Box,
	ServiceCarrier ApplicationContractService.Carrier,
	ServiceOrder ApplicationContractService.Order,
) *groupByDelivery {
	service := new(groupByDelivery)
	service.ServiceAddress = ServiceAddress
	service.ServiceBox = ServiceBox
	service.ServiceCarrier = ServiceCarrier
	service.ServiceOrder = ServiceOrder
	return service
}

func (usecase *groupByDelivery) Run() (packages []DeliveryPackage, err error) {

	// Normalização dos Slices para MAPs e ordenação
	addressses, err := usecase.ServiceAddress.ListAllAddresses()
	if err != nil {
		return packages, err
	}
	addressesMap := make(map[string]DomainEntity.Address)
	for _, address := range addressses {
		addressesMap[address.GetId()] = address
	}
	addressses = nil

	boxes, err := usecase.ServiceBox.ListAllBoxes()
	if err != nil {
		return packages, err
	}
	sort.Slice(boxes, func(i, j int) bool {
		return boxes[i].GetMaxQuantity() > boxes[j].GetMaxQuantity()
	})

	carriers, err := usecase.ServiceCarrier.ListAllCarriers()
	if err != nil {
		return packages, err
	}
	carriersMap := make(map[string]DomainEntity.Carrier)
	for _, carrier := range carriers {
		carriersMap[carrier.GetId()] = carrier
	}
	carriers = nil

	orders, err := usecase.ServiceOrder.ListAllOrders()
	if err != nil {
		return packages, err
	}

	// Agrumapento por Transportadora, Endereço e Data
	packGroup := make(map[string]map[string]map[string]DeliveryPackage)
	for _, order := range orders {
		if _, ok := packGroup[order.GetCarrierId()]; !ok {
			packGroup[order.GetCarrierId()] = make(map[string]map[string]DeliveryPackage)
		}

		if _, ok := packGroup[order.GetCarrierId()][order.GetAddressId()]; !ok {
			packGroup[order.GetCarrierId()][order.GetAddressId()] = make(map[string]DeliveryPackage)
		}

		deliveryDate := usecase.getDeliveryDate(order, carriersMap[order.GetCarrierId()])
		deliveryDateStr := deliveryDate.Format("2006-01-02")

		if _, ok := packGroup[order.GetCarrierId()][order.GetAddressId()][deliveryDateStr]; !ok {
			packGroup[order.GetCarrierId()][order.GetAddressId()][deliveryDateStr] = DeliveryPackage{
				DeliveryDate: deliveryDate,
				Carrier:      carriersMap[order.GetCarrierId()],
				Address:      addressesMap[order.GetAddressId()],
				BoxPackages:  make([]DeliveryPackageBox, 0),
				AllOrders:    make([]DomainEntity.Order, 0),
			}
		}

		CurrentDeliveryPackage := packGroup[order.GetCarrierId()][order.GetAddressId()][deliveryDateStr]

		CurrentDeliveryPackage.AllOrders = append(
			CurrentDeliveryPackage.AllOrders,
			order,
		)

		packGroup[order.GetCarrierId()][order.GetAddressId()][deliveryDateStr] = CurrentDeliveryPackage
	}

	// Para cada pacote (dia), faz a separação em caixas
	for _, carrierPack := range packGroup {
		for _, addressPack := range carrierPack {
			for _, deliveryDatePack := range addressPack {
				deliveryDatePack.BoxPackages = usecase.splitOrdersInBoxes(deliveryDatePack.AllOrders, boxes)
				packages = append(packages, deliveryDatePack)
			}
		}
	}

	return packages, nil
}

// Distribui os pedidos em caixas
func (usecase *groupByDelivery) splitOrdersInBoxes(Orders []DomainEntity.Order, BoxesOptions []DomainEntity.Box) []DeliveryPackageBox {
	packBoxes := []DeliveryPackageBox{}

	// Precisamos ornear os pedidos por quantidade (descrescente) para garantir que os maiores pedidos sejam alocados primeiro
	sort.Slice(Orders, func(i, j int) bool {
		return Orders[i].GetQuantity() > Orders[j].GetQuantity()
	})

	// Calcula o total de itens e mapeia os pedidos por quantidade
	// O total de itens será utilizado para saber quantas caixas serão necessárias no total para esse pacote
	// Utilizamos um map para facilitar a busca do pedido pelo ID e manipular a quantidade sem alterar o objeto original
	totalItems := 0
	orderQuantityMap := map[string]int{}
	for _, order := range Orders {
		totalItems += order.GetQuantity()
		orderQuantityMap[order.GetId()] = order.GetQuantity()
	}

	// Ordena as caixas por quantidade (descrescente) para garantir que as maiores caixas sejam calculadas primeiro
	sort.Slice(BoxesOptions, func(i, j int) bool {
		return BoxesOptions[i].GetMaxQuantity() < BoxesOptions[j].GetMaxQuantity()
	})

	// Baseado na quantidade de items total, faz a separação de quantas e quais caixas serão necessárias
	// Faremos em duas etapas para prevenir um loop adicional e ocupar menos memória abaixo
	for totalItems > 0 {
		for ibox, box := range BoxesOptions {

			// If not exists a big of that
			if ibox == len(BoxesOptions)-1 {
				packBoxes = append(packBoxes, DeliveryPackageBox{
					Box:           box,
					QuantityInBox: 0,
					Orders:        make([]DeliveryPackageBoxOrderRow, 0),
				})

				totalItems -= box.GetMaxQuantity()
				break
			}

			if totalItems <= box.GetMaxQuantity() {
				packBoxes = append(packBoxes, DeliveryPackageBox{
					Box:           box,
					QuantityInBox: 0,
					Orders:        make([]DeliveryPackageBoxOrderRow, 0),
				})

				totalItems = 0
				break
			}
		}
	}

	// Varre todos os pedidos do pacote e os aloca nas caixas mais adequadas (maiores primeiro)
	// Também agrupa os pedidos por caixa, para aproveitar o espaço de cada caixa
	currentBox := 0
	for _, order := range Orders {
		orderId := order.GetId()

		for ibox := currentBox; ibox < len(packBoxes); ibox++ {
			currentBox = ibox

			packBoxes[ibox].Orders = append(packBoxes[ibox].Orders, DeliveryPackageBoxOrderRow{
				Order:    order,
				Quantity: orderQuantityMap[orderId],
			})

			currentSpaceInBox := packBoxes[ibox].Box.GetMaxQuantity() - packBoxes[ibox].QuantityInBox

			// Se o espaço da caixa é maior ou igual ao total do pedido
			if currentSpaceInBox >= orderQuantityMap[orderId] {
				packBoxes[ibox].QuantityInBox += orderQuantityMap[orderId] // Aumentamos o espaço utilizado na caixa para o total do pedido
				orderQuantityMap[orderId] = 0                              // Zeramos o pedido
				break                                                      // Direciona para processar o próximo pedido
			}

			// Se o espaço da caixa é menor que o total do pedido
			packBoxes[ibox].QuantityInBox += currentSpaceInBox // Aumentamos o espaço utilizado na caixa até o limite da caixa
			orderQuantityMap[orderId] -= currentSpaceInBox     // Diminuímos os itens que foram alocados na caixa do total do pedido

			// Se o pedido foi totalmente alocado na caixa
			if orderQuantityMap[orderId] == 0 {
				break // Direciona para processar o próximo pedido
			}

		}

	}

	return packBoxes
}

// Calcula a data de entrega baseado na data do pedido e no horário de corte da transportadora
func (usecase *groupByDelivery) getDeliveryDate(Order DomainEntity.Order, Carrier DomainEntity.Carrier) time.Time {

	OrderDate := Order.GetCreatedAt()

	cutOfftime := Carrier.GetCutOffTime()
	cutHour, cutMinute := cutOfftime.GetHourAndMinute()

	orderHour := OrderDate.Hour()
	orderMinute := OrderDate.Minute()

	if orderHour > cutHour || (orderHour == cutHour && orderMinute > cutMinute) {
		startOfDay := time.Date(
			OrderDate.Year(),
			OrderDate.Month(),
			OrderDate.Day(),
			0, 0, 0, 0,
			OrderDate.Location(),
		)

		OrderDate = startOfDay.AddDate(0, 0, 1)
	}

	return OrderDate
}
