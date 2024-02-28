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
	id           string
	state        string
	zipCode      AddressZipCode
	street       string
	complement   string
	neighborhood string
	city         string
}

func NewAddress(
	id string,
	state string,
	zipCode AddressZipCode,
	street string,
	complement string,
	neighborhood string,
	city string,
) *Address {
	address := new(Address)
	address.id = id
	address.state = state
	address.zipCode = zipCode
	address.street = street
	address.complement = complement
	address.neighborhood = neighborhood
	address.city = city
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

func (address *Address) GetStreet() string {
	return address.street
}

func (address *Address) GetComplement() string {
	return address.complement
}

func (address *Address) GetNeighborhood() string {
	return address.neighborhood
}

func (address *Address) GetCity() string {
	return address.city
}
