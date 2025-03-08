package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type TemplateData struct {
	TaskOne     *TaskOneData
	TaskTwo     *TaskTwoData
	ActiveTab   string
	Task1Result bool
	Task2Result bool
}

type TaskOneData struct {
	// Вхідні значення
	H, C, S, N, O, W, A float64
	// Розраховані значення
	RawToDry          float64
	RawToCombustible  float64
	HDry              float64
	CDry              float64
	SDry              float64
	NDry              float64
	ODry              float64
	ADry              float64
	HCombustible      float64
	CCombustible      float64
	SCombustible      float64
	NCombustible      float64
	OCombustible      float64
	HeatCombustionRaw float64
	HeatCombustionDry float64
	HeatCombustionCom float64
}

type TaskTwoData struct {
	// Вхідні значення
	H, C, S, V, O, W, A, Q float64
	// Розраховані значення
	HRaw               float64
	CRaw               float64
	SRaw               float64
	VRaw               float64
	ORaw               float64
	ARaw               float64
	HeatCombustibleRaw float64
}

func calculateTaskOne(h, c, s, n, o, w, a float64) *TaskOneData {
	rawToDry := math.Round((100/(100-w))*100) / 100
	rawToCombustible := 100 / (100 - w - a)

	hDry := h * rawToDry
	cDry := c * rawToDry
	sDry := s * rawToDry
	nDry := n * rawToDry
	oDry := o * rawToDry
	aDry := a * rawToDry

	hCombustible := h * rawToCombustible
	cCombustible := c * rawToCombustible
	sCombustible := s * rawToCombustible
	nCombustible := n * rawToCombustible
	oCombustible := o * rawToCombustible

	lowerHeatOfCombustion := (339.0*c + 1030.0*h - 108.8*(o-s) - 25.0*w) / 1000
	dryMass := (lowerHeatOfCombustion + 0.025*w) * (100 / (100 - w))
	combustibleMass := (lowerHeatOfCombustion + 0.025*w) * (100 / (100 - w - a))

	return &TaskOneData{
		H: h, C: c, S: s, N: n, O: o, W: w, A: a,

		RawToDry:          rawToDry,
		RawToCombustible:  rawToCombustible,
		HDry:              hDry,
		CDry:              cDry,
		SDry:              sDry,
		NDry:              nDry,
		ODry:              oDry,
		ADry:              aDry,
		HCombustible:      hCombustible,
		CCombustible:      cCombustible,
		SCombustible:      sCombustible,
		NCombustible:      nCombustible,
		OCombustible:      oCombustible,
		HeatCombustionRaw: lowerHeatOfCombustion,
		HeatCombustionDry: dryMass,
		HeatCombustionCom: combustibleMass,
	}
}

func calculateTaskTwo(h, c, s, v, o, w, a, q float64) *TaskTwoData {
	hRaw := h * (100 - w - a) / 100
	cRaw := c * (100 - w - a) / 100
	sRaw := s * (100 - w - a) / 100
	oRaw := o * (100 - w - a) / 100
	vRaw := v * (100 - w) / 100
	aRaw := a * (100 - w) / 100

	heatCombustibleRaw := q*((100-w-a)/100) - 0.025*w

	return &TaskTwoData{
		H: h, C: c, S: s, V: v, O: o, W: w, A: a, Q: q,

		HRaw:               hRaw,
		CRaw:               cRaw,
		SRaw:               sRaw,
		VRaw:               vRaw,
		ORaw:               oRaw,
		ARaw:               aRaw,
		HeatCombustibleRaw: heatCombustibleRaw,
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	data := &TemplateData{
		ActiveTab: "task1", // Активна вкладка за замовчуванням
	}

	if r.Method == http.MethodPost {
		r.ParseForm()

		// Оновлення активної вкладки, якщо необхідно
		if tab := r.FormValue("activeTab"); tab != "" {
			data.ActiveTab = tab
		}

		if r.FormValue("task") == "1" {
			h, _ := strconv.ParseFloat(r.FormValue("h"), 64)
			c, _ := strconv.ParseFloat(r.FormValue("c"), 64)
			s, _ := strconv.ParseFloat(r.FormValue("s"), 64)
			n, _ := strconv.ParseFloat(r.FormValue("n"), 64)
			o, _ := strconv.ParseFloat(r.FormValue("o"), 64)
			w, _ := strconv.ParseFloat(r.FormValue("w"), 64)
			a, _ := strconv.ParseFloat(r.FormValue("a"), 64)

			data.TaskOne = calculateTaskOne(h, c, s, n, o, w, a)
			data.Task1Result = true
		}

		if r.FormValue("task") == "2" {
			h, _ := strconv.ParseFloat(r.FormValue("h"), 64)
			c, _ := strconv.ParseFloat(r.FormValue("c"), 64)
			s, _ := strconv.ParseFloat(r.FormValue("s"), 64)
			v, _ := strconv.ParseFloat(r.FormValue("v"), 64)
			o, _ := strconv.ParseFloat(r.FormValue("o"), 64)
			w, _ := strconv.ParseFloat(r.FormValue("w"), 64)
			a, _ := strconv.ParseFloat(r.FormValue("a"), 64)
			q, _ := strconv.ParseFloat(r.FormValue("q"), 64)

			data.TaskTwo = calculateTaskTwo(h, c, s, v, o, w, a, q)
			data.Task2Result = true
		}
	}

	tmpl.Execute(w, data)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
