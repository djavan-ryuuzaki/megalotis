package megalotis

type ComandoInterface interface {
  Executar() string
  SetMetodo(metodo string)
  SetParametros(parametros string)
  SetObjetos(objetos []interface{})
}


type ComandoErroInterface interface {
  Erro() string
}
