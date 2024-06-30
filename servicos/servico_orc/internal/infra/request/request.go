package request

import (
    "context"
    "crypto/tls"
    "encoding/json"
    "errors"
    "io"
    "net/http"
    "net/url"
    "strings"

    "github.com/amichelins/goexpert_labs_otel/servicos/servico_orc/internal/dto"
    "github.com/valyala/fastjson"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/propagation"
)

var ErrNoCep = errors.New("Cep não foi achado")

type Request struct {
    cep         string
    cidade      string
    key         string
    temperatura float64
    ctx         context.Context
}

func NewRequest(sCep string, sKey string, ctx context.Context) *Request {
    return &Request{cep: sCep, key: sKey, ctx: ctx}
}

// ViaCep Recebe um CEP e consulta VIACEP para saber os dados do cep
//
// PARAMETERS
//
//     sCep string Cep para obter os dados
//
// RETURN
//
//     *dto.ViaCepInput Dados do CEP
//
//     error Erro ocorrido ou nil
//
func (r *Request) ViaCep() error {
    var CepBrasil dto.ViaCep

    http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
    //request, err := http.Get("https://viacep.com.br/ws/" + r.cep + "/json/")
    request, err := http.NewRequestWithContext(r.ctx, "GET", "https://viacep.com.br/ws/"+r.cep+"/json/", nil)

    if err != nil {
        return err
    }

    otel.GetTextMapPropagator().Inject(r.ctx, propagation.HeaderCarrier(request.Header))

    resp, err := http.DefaultClient.Do(request)

    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Pegamos as resposta
    data, err := io.ReadAll(resp.Body)

    if err != nil {
        return err
    }

    if strings.Contains(string(data), `"erro"`) {
        return ErrNoCep
    }

    err = json.Unmarshal(data, &CepBrasil)

    if err != nil {
        return err
    }

    r.cidade = CepBrasil.Localidade

    return nil

}

func (r *Request) Valida() bool {

    if len(r.cep) != 8 || len(r.key) == 0 {
        return false
    }
    return true
}

func (r *Request) GetTemperatura() error {
    var p fastjson.Parser

    request, err := http.NewRequestWithContext(r.ctx, "GET", "http://api.weatherapi.com/v1/forecast.json?key="+r.key+"&q="+url.QueryEscape(r.cidade)+"&days=1&aqi=no&alerts=no", nil)

    if err != nil {
        return err
    }

    otel.GetTextMapPropagator().Inject(r.ctx, propagation.HeaderCarrier(request.Header))
    request.Header.Add("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(request)

    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Pegamos as resposta
    data, err := io.ReadAll(resp.Body)

    if err != nil {
        return err
    }

    Values, err := p.ParseBytes(data)

    if err != nil {
        return err
    }

    Temp, err := Values.GetObject("current").Get("temp_c").Float64()

    if err != nil {
        return err
    }

    r.temperatura = Temp

    return nil
}

func (r *Request) GetTempC() float64 {
    return r.temperatura
}

func (r *Request) GetTempF() float64 {
    return (r.temperatura * 1.8) + 32
}

func (r *Request) GetTempK() float64 {
    return r.temperatura + 273
}

func (r *Request) GetCity() string {
    return r.cidade
}
