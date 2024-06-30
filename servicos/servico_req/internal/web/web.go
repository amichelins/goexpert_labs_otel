package web

import (
    "io"
    "net/http"
    "time"

    "github.com/amichelins/goexpert_labs_otel/servicos/servico_req/internal/dto"
    "github.com/amichelins/goexpert_labs_otel/servicos/servico_req/internal/presenters"
    "github.com/amichelins/goexpert_labs_otel/servicos/servico_req/internal/request"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
    "go.opentelemetry.io/otel/trace"
)

type WebserverProperties struct {
    ResponseTime    time.Duration
    ExternalCallURL string
    RequestNameOTEL string
    OTELTracer      trace.Tracer
}

type Webserver struct {
    WebserverProperties
}

func NewWebServer(webserverProperties WebserverProperties) *Webserver {
    return &Webserver{WebserverProperties: webserverProperties}
}

func (ws *Webserver) CreateServer() *chi.Mux {
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(middleware.RealIP)
    r.Use(middleware.RequestID)

    r.Post("/temp_cep", ws.TempCep)
    return r
}

func (ws *Webserver) TempCep(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    ctx, spanInicial := ws.OTELTracer.Start(ctx, "servico-req")
    time.Sleep(time.Second)

    spanInicial.End()

    data, err := io.ReadAll(r.Body)

    if err != nil {
        http.Error(w, presenters.ToJson(dto.GeneralResponseError{Msg: "Read body error"}), http.StatusInternalServerError)
        return
    }

    CepDto, err := presenters.FromJson(data)

    if err != nil {
        http.Error(w, presenters.ToJson(dto.GeneralResponseError{Msg: "Json parse error"}), http.StatusInternalServerError)
        return
    }

    req := request.NewRequest(presenters.SoDigitos(CepDto.Cep), ws.ExternalCallURL, ctx)

    if !req.Valida() {
        http.Error(w, presenters.ToJson(dto.GeneralResponseError{Msg: "Invalid zipcode"}), http.StatusUnprocessableEntity)
        return
    }

    Response, err := req.CallServicoOrc()

    if err != nil {
        http.Error(w, presenters.ToJson(dto.GeneralResponseError{Msg: "Error calling servico-orc"}), http.StatusInternalServerError)
        return
    }

    if Response.City == "" {
        http.Error(w, presenters.ToJson(dto.GeneralResponseError{Msg: "can not find zipcode"}), http.StatusNotFound)
        return
    }

    http.Error(w, presenters.ToJson(Response), http.StatusOK)
}
