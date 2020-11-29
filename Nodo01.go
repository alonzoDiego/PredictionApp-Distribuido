package main

import (
    "fmt"
	"strconv"
	"bufio"
	"io"
	"net/http"
	"strings"
	"net"
)

type Data struct {
	peso float64
	altura float64	
}

func UrlToLines(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return LinesFromReader(resp.Body)
}

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func regresion_lineal(datos []Data, x float64)float64{
	xmedia := 0.
	ymedia := 0.
	for _, j := range datos{
		xmedia += j.altura
		ymedia += j.peso
	}
	xmedia /= float64(len(datos))
	ymedia /= float64(len(datos))

	usuma := 0.
	dsuma := 0.
	for _, j := range datos{
		usuma += (j.altura-xmedia)*(j.peso-ymedia)
		dsuma += (j.altura-xmedia)*(j.altura-xmedia)
	}
	a := usuma/dsuma
	b := ymedia - (a * xmedia)

	y := a * x + b
	fmt.Println("Para x =", x, "tenemos que y =", y)
	return y
}

func handle(datos []Data, con net.Conn){
	
	defer con.Close()
	r := bufio.NewReader(con)
	
	respuesta, _ := r.ReadString('\n') //Recibe los datos del usuario desde el API
	resp := strings.Split(respuesta, ",")
	val, err := strconv.ParseFloat(strings.ReplaceAll(resp[0],"\r\n",""), 64)
	if err != nil{
		fmt.Println("Error: ", err)
	}
	result := regresion_lineal(datos, val) //aplica el algoritmo de Machine Learning para el primer dato (altura)
	resStr := fmt.Sprintf("%f", result) //convierte el resultado en string
	fmt.Println("Enviando ",resStr, " a nodo02...")

	resSend := resStr + "," + resp[1] //concatena los strings

	completeResult := send(resSend)
	fmt.Print("Enviando ",completeResult, "...")
	fmt.Fprintf(con, "%d\n", completeResult) //envia todo el resultado final en un string al API
}

func main(){
	var datos []Data
	data := Data{
		peso: 0.,
		altura: 0.,
	}

	lines, err := UrlToLines("https://gitlab.com/Dunke28/concurrente/-/raw/master/dataset.txt")
	if err != nil {
		fmt.Println(err)
	}
	for _, line := range lines {
		i := 0
		pesorow := 0.
		alturarow := 0.
		arr := strings.Split(line," ")
		for _, valor := range arr {
			if i == 0 {
				alturarow, err = strconv.ParseFloat(valor, 64)
				i++
			} else {
				pesorow, err = strconv.ParseFloat(valor, 64)
			}
		}
		data.peso = pesorow
		data.altura = alturarow
		datos = append(datos, data)
	}

	ln, _ := net.Listen("tcp", "localhost:9000")
	defer ln.Close()
	
	for {
		con, _ := ln.Accept()
		go handle(datos, con)
	}
}

func send(result string)string {
	conn, _ := net.Dial("tcp", "localhost:9002")
    defer conn.Close()
	fmt.Fprintf(conn, "%d\n", result) 
	
	completeResult := returnRes(conn) 
	
	return completeResult
}

func returnRes(conn net.Conn)string{
	r := bufio.NewReader(conn)
	respuesta, _ := r.ReadString('\n') 

	return respuesta
}
