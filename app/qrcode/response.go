package qrcode

type Response struct {
	Turma     int    `json:"turma"`
	Matricula int    `json:"matricula"`
	Estudante string `json:"estudante"`
}
