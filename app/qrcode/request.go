package qrcode

type Request struct {
	UE             int    `json:"UE"`
	MatriculaAluno int    `json:"MatriculaAluno"`
	Ano            int    `json:"Ano"`
	Visitante      int    `json:"Visitante"`
	Hash           string `json:"Hash"`
}
