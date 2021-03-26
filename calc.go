package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"gopkg.in/yaml.v2"
)

type conf struct {
	Year                         int64 `yaml:"year"`
	Revenue                      int64 `yaml:"revenue"`
	CostOfRevenue                int64 `yaml:"costOfRevenue"`
	OperatingRevenue             int64 `yaml:"opRevenue"`
	CostOfSales                  int64 `yaml:"costOfSales"`
	CostOfAdministrative         int64 `yaml:"costOfAdmin"`
	CostOfResearchAndDevelopment int64 `yaml:"costOfR&D"`
	CostOfFinancing              int64 `yaml:"costOfFinancing"`
	OpratingProfit               int64 `yaml:"opProfit"`
	NetProfit                    int64 `yaml:"netProfit"`
	IrregularProfit              int64 `yaml:"irregularProfit"`
	TotalAssets                  int64 `yaml:"assets"`
	TotalAssetsPreviously        int64 `yaml:"assetsFromLastTime"`
	TotalLiabilities             int64 `yaml:"liabilities"`
	AccountReceivable            int64 `yaml:"receivable"`
	FixedAssets                  int64 `yaml:"fixedAssets"`
	OperatingNetCash             int64 `yaml:"opNetCash"`
	TotalEquityPreviously        int64 `yaml:"equityFromLastTime"`
}

func (c *conf) getConf(filename string) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("err %+v", err)
	}

	err = yaml.Unmarshal(content, c)
	if err != nil {
		log.Fatalf("Unmarshal: %+v", err)
	}
}

type report struct {
	OperatingRevenue                float64
	CostOfRevenue                   float64
	SalesExpense                    float64
	AdminExpense                    float64
	DevExpense                      float64
	FinancingExpense                float64
	ThreeExpenses                   float64
	NetProfit                       float64
	LiabilitiesToAssets             float64
	Receivable                      float64
	FixedAssets                     float64
	WeightedRoE                     float64
	WeightedAverageCostOfCapital    float64
	ReturnOnInvestedCapital         float64
	AssetsTurnoverRatio             float64
	OperatingNetCashToNetProfit     float64
	OperatingNetCashToRegularProfit float64
}

func (c *conf) toReport() *report {
	//Revenue Growth
	//EBIT growth
	//Leverage

	RegularProfit := c.NetProfit - c.IrregularProfit
	roe := toPercent(
		RegularProfit*2,
		c.TotalAssets-c.TotalLiabilities+c.TotalEquityPreviously)
	r := report{
		OperatingRevenue:                toPercent(c.OperatingRevenue, c.Revenue),
		CostOfRevenue:                   toPercent(c.CostOfRevenue, c.Revenue),
		SalesExpense:                    toPercent(c.CostOfSales, c.Revenue),
		AdminExpense:                    toPercent(c.CostOfAdministrative, c.Revenue),
		DevExpense:                      toPercent(c.CostOfResearchAndDevelopment, c.Revenue),
		FinancingExpense:                toPercent(c.CostOfFinancing, c.Revenue),
		ThreeExpenses:                   toPercent((c.CostOfSales + c.CostOfAdministrative + c.CostOfResearchAndDevelopment + c.CostOfFinancing), c.Revenue),
		NetProfit:                       toPercent(c.NetProfit, c.Revenue),
		LiabilitiesToAssets:             toPercent(c.TotalLiabilities, c.TotalAssets),
		Receivable:                      toPercent(c.AccountReceivable, c.TotalAssets),
		FixedAssets:                     toPercent(c.FixedAssets, c.TotalAssets),
		WeightedRoE:                     roe,
		AssetsTurnoverRatio:             toPercent(c.Revenue*2, (c.TotalAssets + c.TotalAssetsPreviously)),
		OperatingNetCashToNetProfit:     toPercent(c.OperatingNetCash, c.NetProfit),
		OperatingNetCashToRegularProfit: toPercent(c.OperatingNetCash, RegularProfit),
	}
	return &r
}

func toPercent(a, b int64) float64 {
	return float64(a) * float64(100) / float64(b)
}

func main() {
	var configFilename = flag.String("config", "", "configuration filename")
	flag.Parse()
	if *configFilename == "" {
		log.Fatalf("Please provide the configuration file\n")

	}

	var c conf
	c.getConf(*configFilename)
	//fmt.Printf("input: %+v\n", c)
	prettyPrint(reflect.ValueOf(&c), "%+v")

	r := c.toReport()
	//fmt.Printf("output : %+v\n", r)
	prettyPrint(reflect.ValueOf(r), "%.2f")
}

func prettyPrint(obj reflect.Value, control string) {
	v := obj.Elem()

	typeOf := v.Type()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		//fmt.Printf("%d ", f.Type())
		fmt.Printf("%s = "+control+"\n",
			typeOf.Field(i).Name, f.Interface())
	}

	fmt.Println()
}
