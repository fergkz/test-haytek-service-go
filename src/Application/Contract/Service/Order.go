package ApplicationContractService

import DomainEntity "github.com/fergkz/test-haytek-service-go/src/Domain/Entity"

type Order interface {
	ListAllOrders() ([]DomainEntity.Order, error)
}
