package presenters

import "regexp"

// SoDigitos Limpa a string via express├úo regular e deixa s├│ digitos numericos
//
// PARAMETERES:
//
//     sTexto string que sera tratada
// RETURN:
//
//     String - O texto passado por par├ómetro com a express├úo regular e ou filtros aplicados
func SoDigitos(sTexto string) string {
    reg, _ := regexp.Compile(`[^0-9]`)
    return reg.ReplaceAllString(sTexto, "")
}
