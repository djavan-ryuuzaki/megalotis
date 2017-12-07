package megalotis


type ComandoHandler struct {
  Ci ComandoInterface
}

func(ch *ComandoHandler) Executar() string {
  return ch.Ci.Executar()
}

func(ch *ComandoHandler) SetHandler(ci ComandoInterface) {
  ch.Ci = ci
}
