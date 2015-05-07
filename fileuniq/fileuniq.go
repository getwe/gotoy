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

var fileset map[string]string
var duplicateFileList []string
var opts struct {
	Path         string `short:"p" long:"path" description:"file path" required:"true"`
	DelDuplicate bool   `short:"d" long:"delete" description:"delete duplicate file or just output the file name" default:"false"`
	Rename       bool   `short:"r" long:"rename" description:"rename use md5 value" default:"true"`
}

func main() {
	fmt.Println(logo())

	parser := flags.NewParser(&opts, flags.HelpFlag)
	_, err := parser.ParseArgs(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	duplicateFileList = make([]string, 0)
	fileset = make(map[string]string)

	err = filepath.Walk(opts.Path, pathTravelHandle)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if opts.DelDuplicate {
		fmt.Printf("Found %d Duplicate File\n", len(duplicateFileList))
		for _, f := range duplicateFileList {
			fmt.Printf("Remove file : %s\n", f)
			err := os.Remove(f)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	} else {
		fmt.Printf("Found %d Duplicate File :\n", len(duplicateFileList))
		for _, f := range duplicateFileList {
			fmt.Println(f)
		}
	}

	if opts.Rename {
		for md5Sign, fileFullPath := range fileset {
			err := renameFile(md5Sign, fileFullPath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
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

	md5Sign, err := fileMd5(fileFullPath)
	if err != nil {
		fmt.Println("ca")
		return err
	}

	_, ok := fileset[md5Sign]
	if !ok {
		fileset[md5Sign] = fileFullPath
	} else {
		duplicateFileList = append(duplicateFileList, fileFullPath)
	}
	return nil
}

func logo() string {
	str := "file uniq"
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

func renameFile(md5Sign, fileFullPath string) error {

	fDir := filepath.Dir(fileFullPath)
	fExt := filepath.Ext(fileFullPath)

	newName := fmt.Sprintf("%s/%s%s", fDir, md5Sign, fExt)

	fmt.Printf("rename from [%s] to [%s]\n", fileFullPath, newName)

	return os.Rename(fileFullPath, newName)
}
