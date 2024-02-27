package DomainEntity

type CarrierCutOffTime struct {
	hour   int
	minute int
}

func NewCarrierCutOffTime(
	hour int,
	minute int,
) *CarrierCutOffTime {
	carrierCutOffTime := new(CarrierCutOffTime)
	carrierCutOffTime.hour = hour
	carrierCutOffTime.minute = minute
	return carrierCutOffTime
}

func (carrierCutOffTime *CarrierCutOffTime) GetHourAndMinute() (Hour int, Minute int) {
	return carrierCutOffTime.hour, carrierCutOffTime.minute
}

type Carrier struct {
	id         string
	name       string
	cutOffTime CarrierCutOffTime
}

func NewCarrier(
	id string,
	name string,
	cutOffTime CarrierCutOffTime,
) *Carrier {
	carrier := new(Carrier)
	carrier.id = id
	carrier.name = name
	carrier.cutOffTime = cutOffTime
	return carrier
}

func (carrier *Carrier) GetId() string {
	return carrier.id
}

func (carrier *Carrier) GetName() string {
	return carrier.name
}

func (carrier *Carrier) GetCutOffTime() CarrierCutOffTime {
	return carrier.cutOffTime
}
