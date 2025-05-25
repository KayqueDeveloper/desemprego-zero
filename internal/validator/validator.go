package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidaEmail verifica se o email é válido
func ValidaEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidaTelefone verifica se o telefone é válido (formato brasileiro)
func ValidaTelefone(telefone string) bool {
	telefoneRegex := regexp.MustCompile(`^\(\d{2}\)\s\d{5}-\d{4}$`)
	return telefoneRegex.MatchString(telefone)
}

// ValidaStruct valida uma struct usando o validator
func ValidaStruct(s interface{}) error {
	return validate.Struct(s)
}

// ErroValidacao representa um erro de validação
type ErroValidacao struct {
	Campo    string `json:"campo"`
	Mensagem string `json:"mensagem"`
}

// ErroValidacaoResponse representa a resposta de erro de validação
type ErroValidacaoResponse struct {
	Erros []ErroValidacao `json:"erros"`
}

// NovoErroValidacao cria um novo erro de validação
func NovoErroValidacao(campo, mensagem string) ErroValidacao {
	return ErroValidacao{
		Campo:    campo,
		Mensagem: mensagem,
	}
}

// ErroValidacaoParaResponse converte um erro de validação para a resposta
func ErroValidacaoParaResponse(err error) ErroValidacaoResponse {
	var erros []ErroValidacao

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			erros = append(erros, NovoErroValidacao(e.Field(), mensagemErro(e)))
		}
	} else {
		erros = append(erros, NovoErroValidacao("", err.Error()))
	}

	return ErroValidacaoResponse{Erros: erros}
}

// mensagemErro retorna uma mensagem de erro amigável para o usuário
func mensagemErro(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "Campo obrigatório"
	case "email":
		return "Email inválido"
	case "min":
		return "Valor muito curto"
	case "max":
		return "Valor muito longo"
	default:
		return "Valor inválido"
	}
}
