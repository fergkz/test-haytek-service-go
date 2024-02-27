package ApplicationContractService

import DomainEntity "github.com/fergkz/test-haytek-service-go/src/Domain/Entity"

type Carrier interface {
	ListAllCarriers() ([]DomainEntity.Carrier, error)
}
