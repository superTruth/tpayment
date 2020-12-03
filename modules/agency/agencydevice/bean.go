package agencydevice

type DeviceBindRequest struct {
	AgencyId uint64 `json:"agency_id"`
	DeviceId uint64 `json:"device_id"`
	FileUrl  string `json:"file_url"`
}
