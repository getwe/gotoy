package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/getwe/figlet4go"
	flags "github.com/jessevdk/go-flags"
	"math/rand"
	"os"
	"time"
)

func main() {
	fmt.Println(logo())

	var opts struct {
		// build mode
		Sum float64 `short:"s" long:"sum" description:"贷款总额"`

		Months float64 `short:"m" long:"months" description:"分期月数"`

		Rate float64 `short:"r" long:"rate" description:"手续费(总利率),如8表示8%"`
	}
	parser := flags.NewParser(&opts, flags.HelpFlag)
	_, err := parser.ParseArgs(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if opts.Sum == 0 || opts.Rate == 0 || opts.Months == 0 {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	calRate(opts.Sum, opts.Months, opts.Rate/100)
}

func calRate(sum, months, rate float64) {

	allInterest := sum * rate
	InterestPerMonth := allInterest / months

	// 算法:假如贷款,全部钱sum拿去投资,每个月从里面扣钱去还贷款
	// 每个月的月利率达到多少,最终利息收益能够覆盖贷款的利息
	InterestRateMonth := (allInterest) * 2 * months / (sum * (1 + months) * months)
	money := sum
	earn := float64(0)

	// 投资收益扣除贷款利息剩的钱
	earnSave := float64(0)

	for i := 0; i < int(months); i++ {
		// 做贷款还欠银行的本金
		loan := sum * float64(int(months)-i) / months
		fmt.Printf("第%2d个月\t欠款%.3f\t月利率%.3f%%\t年化利率%3.2f%%\t利息%.3f\n",
			i+1,
			loan,
			InterestPerMonth/loan*100,
			InterestPerMonth/loan*100*12,
			InterestPerMonth)

		earnSave += (money*InterestRateMonth - InterestPerMonth)

		fmt.Printf("第%2d个月\t投资%.3f\t月利率%.3f%%\t年化利率%3.2f%%\t收益%.3f"+
			"\t收益扣除利息剩%.3f\t活钱%3.2f\n",
			i+1,
			money,
			InterestRateMonth*100,
			InterestRateMonth*100*12,
			money*InterestRateMonth,
			money*InterestRateMonth-InterestPerMonth,
			earnSave)
		earn += money * InterestRateMonth
		// 扣掉一些钱去还贷款本金
		money = money - sum/months
	}

	fmt.Printf("贷款总偿还利息%.3f\t投资总收益利息%.3f\n",
		allInterest, earn)
}

func logo() string {
	str := "Loan Rate"
	ascii := figlet4go.NewAsciiRender()
	// change the font color
	colors := [...]color.Attribute{
		color.FgMagenta,
		color.FgYellow,
		color.FgBlue,
		color.FgCyan,
		color.FgRed,
		color.FgWhite,
		color.FgGreen,
	}
	rand.Seed(time.Now().UnixNano())
	options := figlet4go.NewRenderOptions()
	options.FontColor = make([]color.Attribute, len(str))
	for i := range options.FontColor {
		options.FontColor[i] = colors[rand.Int()%len(colors)]
	}
	renderStr, _ := ascii.RenderOpts(str, options)
	return renderStr
}
