package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "time"

    "github.com/amichelins/goexpert_labs_otel/servicos/servico_req/internal/web"
    "github.com/spf13/viper"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func init() {
    viper.AutomaticEnv()
}

func main() {
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt)

    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
    defer cancel()

    shutdown, err := initProvider(viper.GetString("OTEL_SERVICE_NAME"), viper.GetString("OTEL_EXPORTER_OTLP_ENDPOINT"))
    if err != nil {
        log.Fatal(err, "\n", viper.GetString("OTEL_SERVICE_NAME"), viper.GetString("OTEL_EXPORTER_OTLP_ENDPOINT"))
    }
    defer func() {
        if err := shutdown(ctx); err != nil {
            log.Fatal("failed to shutdown TracerProvider: %w", err)
        }
    }()

    tracer := otel.Tracer("servico-req")

    println(viper.GetString("EXTERNAL_CALL_URL"))
    Server := web.NewWebServer(web.WebserverProperties{
        ResponseTime:    time.Duration(10),
        ExternalCallURL: viper.GetString("EXTERNAL_CALL_URL"),
        OTELTracer:      tracer,
    })

    Router := Server.CreateServer()

    go func() {
        log.Println("Starting server on port", ":8000")
        if err := http.ListenAndServe(":8000", Router); err != nil {
            log.Fatal(err)
        }
    }()

    select {
    case <-sigCh:
        log.Println("Shutting down gracefully, CTRL+C pressed...")
    case <-ctx.Done():
        log.Println("Shutting down due to other reason...")
    }

    // Create a timeout context for the graceful shutdown
    _, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer shutdownCancel()
}

func initProvider(serviceName, collectorURL string) (func(context.Context) error, error) {
    ctx := context.Background()

    res, err := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceName(serviceName),
        ),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create resource: %w", err)
    }

    ctx, cancel := context.WithTimeout(ctx, time.Second)
    defer cancel()

    conn, err := grpc.NewClient("otel-collector:4317",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    //conn, err := grpc.DialContext(ctx, collectorURL,
    //    grpc.WithTransportCredentials(insecure.NewCredentials()),
    //    grpc.WithBlock(),
    //)
    if err != nil {
        return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
    }

    traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
    if err != nil {
        return nil, fmt.Errorf("failed to create trace exporter: %w", err)
    }

    bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
    tracerProvider := sdktrace.NewTracerProvider(
        sdktrace.WithSampler(sdktrace.AlwaysSample()),
        sdktrace.WithResource(res),
        sdktrace.WithSpanProcessor(bsp),
    )
    otel.SetTracerProvider(tracerProvider)

    otel.SetTextMapPropagator(propagation.TraceContext{})

    return tracerProvider.Shutdown, nil
}
