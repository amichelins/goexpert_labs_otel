package web

import (
    "net/http"
    "strings"
    "time"

    "github.com/amichelins/goexpert_labs_otel/servicos/servico_orc/internal/dto"
    "github.com/amichelins/goexpert_labs_otel/servicos/servico_orc/internal/infra/request"
    "github.com/amichelins/goexpert_labs_otel/servicos/servico_orc/internal/presenters"
    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/trace"
)

type WebserverProperties struct {
    ResponseTime  time.Duration
    WeatherApiKey string
    OTELTracer    trace.Tracer
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

    r.Post("/temp_cep", ws.TempCepHandler)
    return r
}

func (ws *Webserver) TempCepHandler(w http.ResponseWriter, r *http.Request) {
    var sCep = r.FormValue("cep")

    carrier := propagation.HeaderCarrier(r.Header)
    ctx := r.Context()
    ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

    ctx, span := ws.OTELTracer.Start(ctx, "servico-orc")
    defer span.End()

    if strings.Trim(ws.WeatherApiKey, " ") == "" {
        http.Error(w, presenters.ToJson(dto.GeneralResponseError{Msg: "Missing Weather Api Key"}), http.StatusInternalServerError)
        return
    }

    Request := request.NewRequest(presenters.SoDigitos(sCep), ws.WeatherApiKey, ctx)

    if !Request.Valida() {
        http.Error(w, presenters.ToJson(dto.GeneralResponseError{Msg: "invalid zipcode"}), http.StatusUnprocessableEntity)
        return
    }

    err := Request.ViaCep()

    if err != nil && err != request.ErrNoCep {
        http.Error(w, presenters.ToJson(dto.GeneralResponseError{Msg: "An error has occurred" + err.Error()}), http.StatusInternalServerError)
        return
    }

    if err == request.ErrNoCep {
        http.Error(w, presenters.ToJson(dto.GeneralResponseError{Msg: "can not find zipcode"}), http.StatusNotFound)
        return
    }

    err = Request.GetTemperatura()

    if err != nil {
        http.Error(w, presenters.ToJson(dto.GeneralResponseError{Msg: "An error has occurred" + err.Error()}), http.StatusInternalServerError)
        return
    }

    _, _ = w.Write([]byte(presenters.ToJson(dto.GeneralResponse{City: Request.GetCity(), TempC: Request.GetTempC(), TempF: Request.GetTempF(), TempK: Request.GetTempK()})))
}
