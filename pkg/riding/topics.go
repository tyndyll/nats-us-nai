package riding

const (
	RideRequestTopic         = "riding.request"
	RideTopic                = "riding.ride"
	DriverNotificationsTopic = "driver.notifications"
	DriverInterestedTopic    = "driver.interested"
	DriverTopic              = "driver"
	CustomerTopic            = "customer"

	ActionRideAccepted = "accepted"
)

func CustomerRideAcceptedTopic(customerID string) string {
	return CustomerTopic + "." + customerID + "." + ActionRideAccepted
}

func DriverRideAcceptedTopic(driverID string) string {
	return DriverTopic + "." + driverID + "." + ActionRideAccepted
}
