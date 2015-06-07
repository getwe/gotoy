package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/getwe/figlet4go"
	flags "github.com/jessevdk/go-flags"
	"github.com/tealeg/xlsx"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println(logo())

	var opts struct {
		InFile []string `short:"i" long:"in" description:"tab隔开的输入文件" required:"true"`

		OutFile string `short:"o" long:"out" description:"结果文件" required:"true"`
	}
	parser := flags.NewParser(&opts, flags.HelpFlag)
	_, err := parser.ParseArgs(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(opts.InFile) == 0 || opts.OutFile == "" {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	xlsxFile := xlsx.NewFile()

	for _, inFileName := range opts.InFile {
		sheet := xlsxFile.AddSheet(inFileName)

		inFh, err := os.OpenFile(inFileName, os.O_RDONLY, 0755)
		if err != nil {
			fmt.Printf("open file %s fail , error : %s", inFileName, err.Error())
			continue
		}

		scanner := bufio.NewScanner(inFh)
		for scanner.Scan() {
			arr := strings.Split(scanner.Text(), "\t")
			row := sheet.AddRow()
			for _, s := range arr {
				cell := row.AddCell()
				cell.Value = s
			}
		}
	}

	err = xlsxFile.Save(opts.OutFile)
	if err != nil {
		fmt.Printf("save xlsx file fail , error : %s", err.Error())
		return
	}
}

func logo() string {
	str := "tab2xlsx"
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
