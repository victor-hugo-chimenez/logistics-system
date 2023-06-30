package driver

type Driver struct {
	ID                  int `db:"id"`
	Name                string
	VehicleModel        string `db:"vehicle_model"`
	VehicleLicensePlate string `db:"vehicle_license_plate"`
	LicenseNumber       int    `db:"license_number"`
}
