package main

import (
	"os"
	"math"
	"fmt"
	"log"
	"strconv"
	"net/http"
	"github.com/gin-gonic/gin"
)


type Prestacao struct{
	Numero int32 `json:"n"`
	Amortizacao float64 `json:"amortizacao"`
	Juros float64 `json:"juros"`
	Saldo float64 `json:"saldo_devedor"`
}

func erro(c *gin.Context, msg string){
	c.String(http.StatusBadRequest, msg)
}

// um efeito colateral pra ver quantas vezes a mesma invocação do programinha é usada
// até ser fechada.
var contador = 0

func arredonda(n float64) float64{
	return math.Floor( n*100 ) / 100
}

// valor parcelas taxa
func getEmprestimo(c *gin.Context){
	contador++

	valor,err := strconv.ParseFloat( c.Param("valor"), 64)
	if err != nil {
		//erro(c,"O valor deve ser um número.")
		erro(c,"O valor deve ser um número")
		return 
	}

	parcelas , err :=   strconv.ParseInt( c.Param("parcelas"), 10,64)  
	if err != nil {
		erro(c,"As parcelas devem ser um número inteiro.")
		return 
	}

	if parcelas < 2{
		erro(c,"Você precisa contratar pelo menos 2 parcelas.")
		return
	}
	if parcelas > 240{
		erro(c,"Seu prazo é muito comprido. Escolha no máximo 240 parcelas.")
		return
	}

	taxa,err := strconv.ParseFloat( c.Param("taxa") ,64)
	if err != nil{
		erro(c,"A taxa de juros deve ser um número.")
		return 
	}
	if taxa < 0 {
		erro(c,"A taxa deve ser maior ou igual a zero.")
		return
	}

	
	var cronograma = make([]Prestacao, 0,parcelas)
	
	var saldo_devedor = valor
	var juros = taxa/100
	var prestacao float64
	var juros_pagos float64
	if taxa > 0 {
		prestacao = valor * (math.Pow( 1 + juros, float64(parcelas) ) * juros) / (math.Pow(1 + juros, float64(parcelas) ) -1)
	}else{
		prestacao = valor / float64(parcelas)
	}

	resp := fmt.Sprintf("Você pediu %g em %d parcelinhas a %g %% ao mês. Isso dá uma parcelinha de %g",valor,parcelas,taxa, arredonda(prestacao))

	for i := int64(0); i < parcelas; i++ {
		amortizacao := prestacao - saldo_devedor*juros
		saldo_devedor -= (prestacao - saldo_devedor*juros)
		juros_pagos += saldo_devedor*juros
		if saldo_devedor <= 0.00{
			saldo_devedor=0
		}

		cronograma = append(cronograma, Prestacao{
			Numero : int32(i),
			Juros  : arredonda(saldo_devedor*juros),
			Amortizacao : arredonda(amortizacao),
			Saldo  : arredonda(saldo_devedor),
		})
		
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"id" : contador,
		"mensagem" : resp,
		"parcelas" : parcelas,
		"prestacao" : arredonda(prestacao),
		"valor_solicitado" : valor,
		"juros_totais" : arredonda(juros_pagos),
		"cronograma": cronograma,
		"taxa" :taxa})
}	

func main() {
	router := gin.Default()
	router.GET("/emprestimo/:valor/:parcelas/:taxa",getEmprestimo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Não achei variável de ambiente PORT. Usando a porta padrão %s",port)
	}

	router.Run(":"+port)
}