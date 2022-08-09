package main

import (
	"encoding/json"
	"fmt"
	"github.com/aviplot/go-finance-math/financial"
	"math"
	"math/rand"
	"os"
	"time"
)

type Calculations struct {
	Xirr float64 `json:"xirr"`
	Xnpv float64 `json:"xnpv"`
}

type fileRecord struct {
	Id    string               `json:"Id"`
	Flows []financial.CashFlow `json:"Flows"`
	Calc  Calculations
}

//func (f fileRecord) String() string {
//	return fmt.Sprintf("Date: %v | flow: %v", f.Id, f.Flows)
//}

func getNextIdFunc() (f func() int) {
	var i int
	f = func() int {
		i++
		return i
	}
	return
}

func randomizeFloat64(from, width float64) float64 {
	return from + math.Round(rand.Float64()*width)
}
func randomizeInt64(from, width int64) int64 {
	return from + rand.Int63n(width)
}
func randomizeDate() string {
	y := int(randomizeInt64(1990, 50))
	m := int(randomizeInt64(1, 10))
	d := 10
	t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
	return fmt.Sprintf("%v-%v-%v", t.Day(), t.Month(), t.Year())
}

func calc(cf financial.CashFlowTab) (c Calculations) {
	var e error
	c.Xirr, e = financial.Xirr(cf)
	if e != nil {
		c.Xirr = 0
	}
	c.Xnpv, e = financial.Xnpv(0.15, cf)
	if e != nil {
		c.Xnpv = 0
	}
	return
}

func getRandomFileRecord(id int) fileRecord {
	flows := financial.NewCashFlowTab(-randomizeFloat64(30000, 100000), randomizeDate(), int(randomizeInt64(10, 200)),
		float64(randomizeInt64(1000, 3000)), randomizeDate())
	f := new(fileRecord)
	f.Id = fmt.Sprintf("%v", id)
	f.Flows = flows
	f.Calc = calc(flows)
	return *f
}

func recordsToFile(fn string, d []fileRecord) {
	e, err := json.MarshalIndent(d,"","	")
	if err != nil {
		fmt.Println(err)
		return
	}
	s := string(e)

	f, err := os.Create("flows.js")
	if err != nil {
		panic("File writing error")
	}
	defer f.Close()
	f.WriteString(s)
	f.Sync()
}

func main() {
	amount := 100
	f := make([]fileRecord, amount)
	rand.Seed(time.Now().UnixNano())
	getNextId := getNextIdFunc()

	for i := 0; i < amount; i++ {
		f[i] = getRandomFileRecord(getNextId())
	}
	//x, _ := financial.Xirr(c)
	//fmt.Println(f)
	recordsToFile("", f)
}
