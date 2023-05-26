package domain

type Delivery struct {
	TrackingCode       string `json:"tracking_code"`
	SourceAddress      string `json:"source_address"`
	DestinationAddress string `json:"destination_address"`
	Status             string `json:"status"`
	Created            string `json:"created"`
	Modified           string `json:"modified"`
}
