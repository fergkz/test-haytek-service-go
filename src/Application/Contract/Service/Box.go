package ApplicationContractService

import DomainEntity "github.com/fergkz/test-haytek-service-go/src/Domain/Entity"

type Box interface {
	ListAllBoxes() ([]DomainEntity.Box, error)
}
