package main

import (
	"crypto/md5"
	"fmt"
	"github.com/fatih/color"
	"github.com/getwe/figlet4go"
	flags "github.com/jessevdk/go-flags"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var opts struct {
	From   string `short:"f" long:"from" description:"copy file from path" required:"true" default:"WechatPath"`
	To     string `short:"t" long:"to" description:"copy file destination" required:"true"`
	Suffix string `short:"s" long:"suffix" description:"file suffix" required:"true" default:"mp4"`
	Move   bool   `short:"m" long:"move" description:"move file or copy" default:"false"`
}

var fileset map[string]string

func main() {
	fmt.Println(logo())

	parser := flags.NewParser(&opts, flags.HelpFlag)
	_, err := parser.ParseArgs(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fileset = make(map[string]string)

	if opts.From == "WechatPath" {
		opts.From = getWeChatPath()
	}

	err = filepath.Walk(opts.From, pathTravelHandle)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	for md5Sign, fileFullPath := range fileset {
		fExt := filepath.Ext(fileFullPath)
		newName := fmt.Sprintf("%s/%s%s", opts.To, md5Sign, fExt)
		fmt.Printf("%s => %s\n", fileFullPath, newName)
		if opts.Move {
			err = moveFile(fileFullPath, newName)
		} else {
			err = copyFile(fileFullPath, newName)
		}
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

}

func pathTravelHandle(path string, info os.FileInfo, err error) error {
	if info == nil {
		return nil
	}
	if info.IsDir() {
		return nil
	}
	if err != nil {
		return err
	}

	if strings.HasPrefix(info.Name(), ".") {
		return nil
	}

	fileFullPath := path

	fExt := filepath.Ext(fileFullPath)

	if fmt.Sprintf(".%s", opts.Suffix) != fExt {
		return nil
	}

	md5Sign, err := fileMd5(fileFullPath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, ok := fileset[md5Sign]
	if !ok {
		fileset[md5Sign] = fileFullPath
	}

	return nil
}

func logo() string {
	str := "findfile"
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

func fileMd5(fullpath string) (string, error) {
	fh, err := os.Open(fullpath)
	if err != nil {
		return "", err
	}
	m := md5.New()
	io.Copy(m, fh)
	return fmt.Sprintf("%x", m.Sum(nil)), nil
}

func copyFile(from, to string) error {
	f, err := os.Open(from)
	defer f.Close()
	if err != nil {
		return err
	}

	t, err := os.Create(to)
	defer t.Close()
	if err != nil {
		return err
	}

	_, err = io.Copy(t, f)
	if err != nil {
		return err
	}

	return nil
}

func moveFile(from, to string) error {
	return os.Rename(from, to)
}

func getWeChatPath() string {
	wx := "/Library/Containers/com.tencent.xinWeChat/Data/Library/Application Support/Wechat"
	return fmt.Sprintf("%s%s", os.Getenv("HOME"), wx)
}
