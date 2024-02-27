package ApplicationContractService

import DomainEntity "github.com/fergkz/test-haytek-service-go/src/Domain/Entity"

type Address interface {
	ListAllAddresses() ([]DomainEntity.Address, error)
}
