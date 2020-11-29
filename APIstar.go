package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"io/ioutil"
	"fmt"
	"strconv"
	"bufio"
	"strings"
	"net"
)

//estructura
type Data struct {
	Altura float64 `json:"altura"`
	Horas  float64 `json:"horas"`
}

var datos []Data
var altura string
var hora string

func main() {

	datos = []Data{
		{120.3, 3000.2},
		{150.0, 5000},
		{130.5, 4000}}

	handleResquest()
	
}

func handleResquest() {
	http.HandleFunc("/results", getAll)
	http.HandleFunc("/pushresult", pushResult)
	http.HandleFunc("/getdata", getResult)

	log.Fatal(http.ListenAndServe(":9001", nil))
}

func getAll(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type", "application/json")
	//serializacion
	jsonBytes, _ := json.MarshalIndent(datos, "", " ")
	io.WriteString(res, string(jsonBytes))
}

func pushResult(res http.ResponseWriter, req *http.Request) {
	var consulta Data

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil{
		fmt.Fprintf(res, "Inserte datos validos")
	}

	json.Unmarshal(reqBody, &consulta)
	datos = append(datos, consulta)

	res.Header().Set("Content-Type", "application/json")

}

func getResult(res http.ResponseWriter, req *http.Request) {
	//http://localhost:3000/getdata?alt=155&hour=8555

	altura = req.FormValue("alt") 
	hora = req.FormValue("hour")

	fmt.Println(altura)
	fmt.Println(hora)

	conn, _ := net.Dial("tcp", "localhost:9000")
	defer conn.Close()
	allDate := altura + "," + hora
	fmt.Fprintf(conn, "%s\n", allDate)
	
	obtenerData(conn)

	res.Header().Set("Content-Type", "application/json")
}

func obtenerData(con net.Conn){
 	r := bufio.NewReader(con)
	respuesta, _ := r.ReadString('\n')
	resp := strings.Split(respuesta, ",")

	muestra1 := strings.Split(resp[0], "=")
	muestra2 := strings.Split(resp[1], ")")

	horaF, _ := strconv.ParseFloat(strings.ReplaceAll(hora,"\r\n",""), 64)
	calcMin := horaF / 60
 
	 fmt.Println("Con una altura de ", altura, " el peso ideal es ", muestra1[3], "kg")
	 fmt.Println("Haciendo caminata por ", calcMin, " minutos, se quema ", muestra2[0], " cal.")
}