package agencydevice

type DeviceBindRequest struct {
	AgencyId uint   `json:"agency_id"`
	DeviceId uint   `json:"device_id"`
	FileUrl  string `json:"file_url"`
}
