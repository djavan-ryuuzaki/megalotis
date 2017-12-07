package megalotis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	ComandoEscolhido string
  Proibidos string
	iniciado bool
	modulosRegistrados = make(map[string]ComandoInterface)
	modulos *ModulosStruct
	config *configStruct
	configPadrao *configStruct
)

type configStruct struct {
	PrefixoComando string `json:"Prefixo"`
	PalavrasProibidas string `json:"proibidas"`
}

type ModulosStruct struct {
  Modulos []ModuloStruct `json:"modulos"`
}

type ModuloStruct struct {
	Nome string `json:"nome"`
	Metodos []MetodoStruct `json:"metodos"`
}

type MetodoStruct struct {
	Metodo        string `json:"metodo"`
	Parametros	 	string `json:"parametros"`
	Descricao     string `json:"descrição"`
	Info          string `json:"info"`
  Cooldown      int `json:"cooldown"`
  Funcao        string `json:"funcao"`
	Ativa         bool   `json:"ativa"`
}

type ComandoErro struct {
	s string
}

func (e ComandoErro) Erro() string {
	return e.s
}

func Inicia() error {

	fmt.Println("Lendo o arquivo de comandos...")

	file, err := ioutil.ReadFile("./comandos.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = json.Unmarshal(file, &modulos)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	errConfig := json.Unmarshal(file, &config)

	if errConfig != nil {

		fmt.Println(errConfig.Error())
		return errConfig

	}

	Proibidos = config.PalavrasProibidas

	iniciado = true

	return nil
}

func RegistrarModulo(nome string, classe ComandoInterface){
	modulosRegistrados[nome] = classe
}

func AnalisaComando(m string, objetos []interface{}) (ComandoInterface, ComandoErroInterface) {

	if iniciado == false {
		return nil, ComandoErro{"Comandos não foi iniciado!"}
	}

	existeModulo := false
	existeMetodo := false
	existeParametro := false

	if strings.HasPrefix(m, config.PrefixoComando) != true{
      return nil, nil
    }

	prefix, metod, params := QuebraStringComando(m)

	for m := range modulos.Modulos {
		//Verifica se existe o Módulo
		if modulos.Modulos[m].Nome == prefix {

			if existeModulo == false{
				existeModulo = true
			}

      for n := range modulos.Modulos[m].Metodos  {

					//Se existir o Módulo procuro pelo Método
					if modulos.Modulos[m].Metodos[n].Metodo == metod || (modulos.Modulos[m].Metodos[n].Metodo == "default" && metod == "") {

							if ( modulos.Modulos[m].Metodos[n].Ativa == false ){
								existeMetodo = false
								break
							}
							if( existeMetodo == false ){
								existeMetodo = true
								//Se o metodo não exigir parâmetros encerra a verificação
								if( modulos.Modulos[m].Metodos[n].Parametros == "" ){
									existeParametro = true
									break
								}
								//Se o número de parametros for igual ao número de parametros obrigatorios
								if( len( strings.Split(params, " ") ) == len ( strings.Split(modulos.Modulos[m].Metodos[n].Parametros, " ") ) ){
									existeParametro = true
									break
								}

							}
							break
					}
			}
			break
		}

	}

	if( existeModulo == false){
		return nil, ComandoErro{"Desculpe, mas eu não entendo esse comando!"}
	}

	if(existeMetodo == false){
		return nil, ComandoErro{"Desculpe, mas eu não tenho essa habilidade!"}
	}

	if(existeParametro == false){
		return nil, ComandoErro{"Parece que falta alguma coisa no seu pedido!"}
	}

	if existeModulo && existeMetodo && existeParametro {
		return executaComando(prefix, metod, params, objetos)
	}

	return nil, nil
}

func QuebraStringComando(m string) (prefix string, metodo string, parametros string){

	s      := strings.Split(m, " ")
	comp   := len( s )
	pref   := ""
	met    := ""
	params := ""
	switch  {
		case comp == 1:
			pref   = strings.Replace(s[0], "!", "", -1)
			met    = ""
			params = ""
		case comp == 2:
			pref   = strings.Replace(s[0], "!", "", -1)
			met    = s[1]
			params = ""
		case comp >= 3:
			pref   = strings.Replace(s[0], "!", "", -1)
			met    = s[1]
			params = m[ (len(s[0]) + len(s[1]) + 1 )  :len(m)]

	}

	return strings.Trim(pref, " "), strings.Trim(met, " "), strings.Trim(params, " ")

}

func executaComando( prefix string, metod string, params string, objetos []interface{}) (ComandoInterface, ComandoErroInterface) {

	if com, ok := modulosRegistrados[prefix]; ok {
		com.SetMetodo(metod)
		com.SetParametros(params)
		com.SetObjetos(objetos)
		return com, nil
	}

	return nil, ComandoErro{"O comando '"+prefix+"' não foi registrado ainda!"}


}
