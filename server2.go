package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type Calificacion struct {
	Alumno       string `json:"alumno"`
	Calificacion string `json:"calificacion"`
	Materia      string `json:"materia"`
}

type Mensaje struct {
	alumno       string
	calificacion string
	materia      string
}

type AdminCalificacion struct {
	Calificaciones []Calificacion
}

func (calificaciones *AdminCalificacion) AgregarCalificacion(calificacion Calificacion) {
	calificaciones.Calificaciones = append(calificaciones.Calificaciones, calificacion)
}

var misCalificaciones AdminCalificacion

func main() {
	http.HandleFunc("/", root)
	http.HandleFunc("/agregar", viewAgregar)
	http.HandleFunc("/promA", viewPromA)
	http.HandleFunc("/promG", viewPromG)
	http.HandleFunc("/promM", viewPromM)
	http.HandleFunc("/promedioA", promedioA)
	http.HandleFunc("/promedioM", promedioM)
	http.HandleFunc("/calificacion", calificacion)
	http.HandleFunc("/respaldar", respaldar)
	http.HandleFunc("/recuperar", recuperar)
	http.HandleFunc("/vaciar", vaciar)

	fmt.Println("Corriendo servirdor calificaciones...")
	http.ListenAndServe(":9000", nil)
}

func calificacion(res http.ResponseWriter, req *http.Request) {
	var msg Mensaje
	fmt.Println(req.Method)
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		msg.alumno = req.FormValue("nombre")
		msg.materia = req.FormValue("materia")
		msg.calificacion = req.FormValue("calificacion")
		//// anexar codigo RPC
		if existeMateriaAndAlumno(msg.materia, msg.alumno) {

			err := "ya existe calificacion para ese alumno en la materia"
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				cargarHtml("error.html"),
				err,
			)
			return

		}

		var aux Calificacion

		aux.Alumno = msg.alumno
		aux.Calificacion = msg.calificacion
		aux.Materia = msg.materia
		misCalificaciones.AgregarCalificacion(aux)
		fmt.Println(misCalificaciones.Calificaciones)

		fmt.Println("Calificacion Registrada")
		//// fin codigo rpc
		m := "se agrego la calificacion de " + msg.alumno
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("respuesta.html"),
			m,
		)
	case "GET":
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("index.html"),
		)
	}
}

func promedioA(res http.ResponseWriter, req *http.Request) {

	fmt.Println(req.Method)
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		alumnoName := req.FormValue("nombre")
		//// anexar codigo RPC
		var total float64
		var promedio float64
		var contador float64
		existe := false
		for _, i := range misCalificaciones.Calificaciones {
			if i.Alumno == alumnoName {
				existe = true
				aux, _ := strconv.ParseFloat(i.Calificacion, 64)
				total += aux
				contador++
			}
		}
		promedio = total / contador
		if existe {

			m := "El promedio de " + alumnoName + " es: " + strconv.FormatFloat(promedio, 'f', 1, 64)
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				cargarHtml("respuesta.html"),
				m,
			)
		} else {

			err := "No existe ese alumno"
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				cargarHtml("error.html"),
				err,
			)
			return
		}

		//// fin codigo rpc

	case "GET":
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("index.html"),
		)
	}
}

func promedioM(res http.ResponseWriter, req *http.Request) {

	fmt.Println(req.Method)
	switch req.Method {
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, "ParseForm() error %v", err)
			return
		}
		fmt.Println(req.PostForm)
		materiaName := req.FormValue("nombre")
		//// anexar codigo RPC
		var total float64
		var promedio float64
		var contador float64
		existe := false
		for _, i := range misCalificaciones.Calificaciones {
			if i.Materia == materiaName {
				existe = true
				aux, _ := strconv.ParseFloat(i.Calificacion, 64)
				total += aux
				contador++

			}
		}
		promedio = total / contador
		if existe {

			m := "El promedio de la materia " + materiaName + " es: " + strconv.FormatFloat(promedio, 'f', 1, 64)
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				cargarHtml("respuesta.html"),
				m,
			)
		} else {

			err := "No existe esa Materia"
			res.Header().Set(
				"Content-Type",
				"text/html",
			)
			fmt.Fprintf(
				res,
				cargarHtml("error.html"),
				err,
			)
			return
		}

		//// fin codigo rpc

	case "GET":
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("index.html"),
		)
	}
}

func isvisitado(l []string, name string) bool {
	visitado := false
	for _, i := range l {
		if i == name {
			visitado = true
		}
	}
	return visitado
}

func promedioindividual(name string) float64 {
	var total float64
	var contador float64
	var promedio float64
	for _, c := range misCalificaciones.Calificaciones {

		if name == c.Alumno {
			aux, _ := strconv.ParseFloat(c.Calificacion, 64)
			total += aux
			contador++
		}

	}
	promedio = total / contador
	return promedio
}

func viewPromG(res http.ResponseWriter, req *http.Request) {
	alumnos := make([]string, 0)
	var contador float64
	var total float64
	var promedioG float64 = 0

	// crear lista de alumnos
	for _, i := range misCalificaciones.Calificaciones {
		if !isvisitado(alumnos, i.Alumno) {
			alumnos = append(alumnos, i.Alumno)
		}
	}

	// promedio de un alumno

	for _, i := range alumnos {
		total += promedioindividual(i)
		contador++
	}
	promedioG = total / contador

	if len(misCalificaciones.Calificaciones) == 0 {
		err := "No existen alumnos"
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("error.html"),
			err,
		)
		return
	} else {
		aux := strconv.FormatFloat(promedioG, 'f', 1, 64)
		res.Header().Set(
			"Content-Type",
			"text/html",
		)
		fmt.Fprintf(
			res,
			cargarHtml("promGeneral.html"),
			aux,
		)
	}

}

func existeMateriaAndAlumno(materia string, alumno string) bool {
	existe := false
	for _, i := range misCalificaciones.Calificaciones {
		if i.Alumno == alumno && i.Materia == materia {
			existe = true

		}

	}

	return existe
}

func (calificaciones *AdminCalificacion) String() string {
	var html string
	alumnos := make([]string, 0)
	for _, i := range misCalificaciones.Calificaciones {
		if !isvisitado(alumnos, i.Alumno) {
			alumnos = append(alumnos, i.Alumno)
			html +=
				"<option value=\"" + i.Alumno + "\"" + ">" + i.Alumno + "</option>"
		}
	}

	return html
}
func (calificaciones *AdminCalificacion) StringM() string {
	var html string
	materias := make([]string, 0)
	for _, i := range misCalificaciones.Calificaciones {
		if !isvisitado(materias, i.Materia) {
			materias = append(materias, i.Materia)
			html +=
				"<option value=\"" + i.Materia + "\"" + ">" + i.Materia + "</option>"
		}
	}

	return html
}

func root(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("index.html"),
	)
}

func cargarHtml(a string) string {
	html, _ := ioutil.ReadFile(a)

	return string(html)
}

func viewAgregar(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("form-agregar.html"),
	)
}
func viewPromA(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("promAlumno.html"),
		misCalificaciones.String(),
	)
}

func viewPromM(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("promMateria.html"),
		misCalificaciones.StringM(),
	)
}

func respaldar(res http.ResponseWriter, req *http.Request) {

	misCalificaciones.crearRespaldo()
	m := "Se ah respaldado la informacion Correctamente"
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("respuesta.html"),
		m,
	)

}

func vaciar(res http.ResponseWriter, req *http.Request) {

	var aux AdminCalificacion
	misCalificaciones = aux
	m := "Se ah vaciado la informacion correctamente"
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("respuesta.html"),
		m,
	)

}

func recuperar(res http.ResponseWriter, req *http.Request) {

	misCalificaciones.recuperarJson()
	m := "Se ah restaurado la informacion correctamente"
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(
		res,
		cargarHtml("respuesta.html"),
		m,
	)

}

func (calificaciones *AdminCalificacion) recuperarJson() {
	loadJSON("calificaciones.json", &misCalificaciones)
	fmt.Println(calificacion)

}

func (calificaciones *AdminCalificacion) crearRespaldo() {
	fmt.Println(calificaciones.Calificaciones)

	jsonData, err := json.MarshalIndent(misCalificaciones, "", "    ")
	if err != nil {
		fmt.Println("Error al convertir")
		return
	}
	fmt.Println(string(jsonData))
	saveJSON("calificaciones.json", misCalificaciones)

}

func saveJSON(fileName string, object interface{}) {
	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error al convertir a JSON", err.Error())
		return
	}
	err = json.NewEncoder(outFile).Encode(object)
	if err != nil {
		fmt.Println("Error al convertir a JSON", err.Error())
		return
	}
	outFile.Close()
}

func loadJSON(fileName string, object interface{}) {
	inFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error al abrir el archivo", err.Error())
		return
	}
	err = json.NewDecoder(inFile).Decode(object)
	if err != nil {
		fmt.Println("Error de conversión", err.Error())
		return
	}
	inFile.Close()
}
