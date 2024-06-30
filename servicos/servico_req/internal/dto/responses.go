package dto

type GeneralResponse struct {
    City  string  `json:"city"`
    TempC float64 `json:"temp_c"`
    TempF float64 `json:"temp_f"`
    TempK float64 `json:"temp_K"`
}

type GeneralResponseError struct {
    Msg string `json:"msg"`
}
