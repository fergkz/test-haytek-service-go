package DomainEntity

type AddressZipCode struct {
	code string
}

func NewAddressZipCode(
	code string,
) *AddressZipCode {
	addressZipCode := new(AddressZipCode)
	addressZipCode.code = code
	return addressZipCode
}

func (addressZipCode *AddressZipCode) GetCode() string {
	return addressZipCode.code
}

func (addressZipCode *AddressZipCode) SetCode(code string) {
	addressZipCode.code = code
}

type Address struct {
	id      string
	state   string
	zipCode AddressZipCode
}

func NewAddress(
	id string,
	state string,
	zipCode AddressZipCode,
) *Address {
	address := new(Address)
	address.id = id
	address.state = state
	address.zipCode = zipCode
	return address
}

func (address *Address) GetId() string {
	return address.id
}

func (address *Address) GetState() string {
	return address.state
}

func (address *Address) GetZipCode() AddressZipCode {
	return address.zipCode
}
