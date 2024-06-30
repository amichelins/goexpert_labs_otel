package request

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestRequestInvalidCepAndKey(t *testing.T) {
    ctx := context.Background()
    Request := NewRequest("910403", "", ctx)

    bCep := Request.Valida()

    // Cep Invalido e key invalido
    assert.Equal(t, false, bCep)
}

func TestRequestValidCepAndInvalidKey(t *testing.T) {
    ctx := context.Background()
    Request := NewRequest("91040300", "", ctx)

    bCep := Request.Valida()

    // Cep Valido e key invalido
    assert.Equal(t, false, bCep)
}

func TestRequestValidCepAndKey(t *testing.T) {
    ctx := context.Background()
    Request := NewRequest("91040300", "alsdjladadlkajdkljalkdjalksjdkl", ctx)

    bCep := Request.Valida()

    // Cep Valido e Key Valido
    assert.Equal(t, true, bCep)
}

func TestRequestViaCepInvalido(t *testing.T) {
    ctx := context.Background()
    Request := NewRequest("00000000", "alsdjladadlkajdkljalkdjalksjdkl", ctx)

    bCep := Request.Valida()

    // Cep Valido e Key Valido
    assert.Equal(t, true, bCep)

    err := Request.ViaCep()

    assert.Equal(t, ErrNoCep, err)

}

func TestRequestTemperatura(t *testing.T) {
    ctx := context.Background()
    Request := NewRequest("91040300", "545c605410b74c09a2921907241506", ctx)

    bCep := Request.Valida()

    // Cep Valido e Key Valido
    assert.Equal(t, true, bCep)

    err := Request.ViaCep()

    assert.Equal(t, nil, err)

    err = Request.GetTemperatura()

    assert.Equal(t, nil, err)

    TempC := Request.GetTempC()

    assert.LessOrEqual(t, 0.0, TempC)

    TempF := Request.GetTempF()

    assert.LessOrEqual(t, 0.0, TempF)

    TempK := Request.GetTempK()

    assert.LessOrEqual(t, 0.0, TempK)
}
