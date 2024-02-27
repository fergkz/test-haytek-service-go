package DomainEntity

import "time"

type Order struct {
	id        string
	addressId string
	carrierId string
	createdAt time.Time
	quantity  int
}

func NewOrders(
	id string,
	addressId string,
	carrierId string,
	createdAt time.Time,
	quantity int,
) *Order {
	order := new(Order)
	order.id = id
	order.addressId = addressId
	order.carrierId = carrierId
	order.createdAt = createdAt
	order.quantity = quantity
	return order
}

func (order *Order) GetId() string {
	return order.id
}

func (order *Order) GetAddressId() string {
	return order.addressId
}

func (order *Order) GetCarrierId() string {
	return order.carrierId
}

func (order *Order) GetCreatedAt() time.Time {
	return order.createdAt
}

func (order *Order) GetQuantity() int {
	return order.quantity
}
