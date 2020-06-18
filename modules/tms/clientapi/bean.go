package clientapi

// 请求结构体
type RequestBean struct {
	DeviceSn    string    `json:"device_sn,omitempty"`
	DeviceCsn   string    `json:"device_csn,omitempty"`
	DeviceModel string    `json:"device_model,omitempty"`
	Alias       string    `json:"alias,omitempty"`
	LocationLat string    `json:"location_lat,omitempty"`
	LocationLon string    `json:"location_lon,omitempty"`
	PushToken   string    `json:"push_token,omitempty"`
	Power       int       `json:"power", omitempty`
	StoreID     int       `json:"store_id", omitempty`
	AppInfos    []AppInfo `json:"app_infos,omitempty"`
}

// 批量设置数据请求体
type ConfigInBatchRequestBean struct {
}

// 请求
type RequestWaitingApproveDevice struct {
	AccountID uint   `json:"account_id,omitempty"`
	DeviceSn  string `json:"device_sn,omitempty"`
}

type RequestApprove struct {
	ID     uint `json:"id,omitempty"`
	Status int  `json:"status,omitempty"`
}

type RequestDecodeApkFile struct {
	ID uint `json:"id,omitempty"`
}

type RequestBathUpdate struct {
	ID uint `json:"id,omitempty"`
}

type UploadFileRequest struct {
	FileName string `json:"file_name,omitempty"`
	FileSize uint   `json:"file_size,omitempty"`
	Tag      string `json:"devicetag,omitempty"`
}

type CreateFileRequest struct {
	DeviceSn string `json:"device_sn,omitempty"`
	FileName string `json:"file_name,omitempty"`
	FileUrl  string `json:"file_url,omitempty"`
}

// 返回结构体
type ResponseBean struct {
	DeviceCsn        string    `json:"device_csn,omitempty"`
	Alias            string    `json:"alias,omitempty"`
	RebootMode       int       `json:"reboot_mode"`
	RebootTime       string    `json:"reboot_time"`
	RebootDayInWeek  int       `json:"reboot_day_in_week,omitempty"`
	RebootDayInMonth int       `json:"reboot_day_in_month,omitempty"`
	AppInfos         []AppInfo `json:"app_infos,omitempty"`
}

type AppInfo struct {
	Name        string `json:"name"`
	PackageId   string `json:"package_id"`
	VersionName string `json:"version_name"`
	VersionCode int    `json:"version_code"`
	Description string `json:"description"`

	Status   string    `json:"status"`
	FileInfo *FileInfo `json:"file_info,omitempty"`
}

type FileInfo struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type UploadFileResponse struct {
	UploadUrl   string `json:"upload_url,omitempty"`
	DownloadUrl string `json:"download_url,omitempty"`
	Exp         int64  `json:"exp,omitempty"`
}
