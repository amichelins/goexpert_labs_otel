package request

import (
    "context"
    "crypto/tls"
    "encoding/json"
    "io"
    "log"
    "net/http"
    "net/url"
    "strings"

    "github.com/amichelins/goexpert_labs_otel/servicos/servico_req/internal/dto"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/propagation"
)

type Request struct {
    cep         string
    externalURL string
    ctx         context.Context
}

func NewRequest(sCep string, ExternalURL string, ctx context.Context) *Request {
    return &Request{cep: sCep, externalURL: ExternalURL, ctx: ctx}
}

func (r *Request) Valida() bool {
    return len(r.cep) == 8
}

func (r *Request) CallServicoOrc() (*dto.GeneralResponse, error) {
    var GeneralResp dto.GeneralResponse

    data := url.Values{}
    data.Set("cep", r.cep)

    log.Println("QQQ " + data.Encode())
    http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

    req, err := http.NewRequestWithContext(r.ctx, "POST", r.externalURL, strings.NewReader(data.Encode()))

    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    otel.GetTextMapPropagator().Inject(r.ctx, propagation.HeaderCarrier(req.Header))

    resp, err := http.DefaultClient.Do(req)

    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    Data, err := io.ReadAll(resp.Body)

    if err != nil {
        return nil, err
    }

    err = json.Unmarshal(Data, &GeneralResp)

    if err != nil {
        return nil, err
    }

    return &GeneralResp, nil
}
