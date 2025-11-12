package model

// SatcomDataInput represents the input for creating/updating satcom data
type SatcomDataInput struct {
	Company  string `json:"company" binding:"required"`
	Category string `json:"category" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Date     string `json:"date" binding:"required"`
	Time     string `json:"time" binding:"required"`
	DbPort   string `json:"db_port" binding:"required"`
	UiPort   string `json:"ui_port" binding:"required"`
	URL      string `json:"url" binding:"required"`
	IP       string `json:"ip" binding:"required"`
	Status   bool   `json:"status"`
}

// SatcomDataResponse represents the response model for satcom data
type SatcomDataResponse struct {
	ID       int32  `json:"id"`
	Company  string `json:"company"`
	Category string `json:"category"`
	Type     string `json:"type"`
	Date     string `json:"date"`
	Time     string `json:"time"`
	DbPort   string `json:"db_port"`
	UiPort   string `json:"ui_port"`
	URL      string `json:"url"`
	IP       string `json:"ip"`
	Status   bool   `json:"status"`
}

