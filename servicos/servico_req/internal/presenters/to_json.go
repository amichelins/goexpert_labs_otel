package presenters

import (
    "encoding/json"

    "github.com/amichelins/goexpert_labs_otel/servicos/servico_req/internal/dto"
)

func ToJson(Message any) string {

    RespMessage, _ := json.Marshal(Message)

    return string(RespMessage)
}

func FromJson(jsonData []byte) (*dto.Request, error) {
    var Req dto.Request
    err := json.Unmarshal(jsonData, &Req)

    if err != nil {
        return nil, err
    }

    return &Req, nil
}
