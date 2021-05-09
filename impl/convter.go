package impl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	Verbose  = false
	FilePath = ""
)

func Cli() {
	if FilePath == "" {
		log.Fatalln("-f is not specify")
	}
	msg, err := readMsgFile(FilePath)
	if err != nil {
		log.Fatalln("The file format is wrong")
	}
	if len(msg) == 0 {
		log.Fatalln("The file is empty")
	}
	converter, err := NewConverter(msg)
	if err != nil {
		log.Fatalf("parse filed;\n %v\n", err)
	}
	converter.Echo()
	os.Exit(0)
}

func NewConverter(msg string) (*Converter, error) {
	if Verbose {
		log.Println(msg)
	}
	c := &Converter{
		msg:       msg,
		msgReader: bufio.NewReader(strings.NewReader(msg)),
	}
	err := c.parse()
	if err != nil {
		return nil, err
	}
	return c, nil
}

type Converter struct {
	msg       string
	msgReader *bufio.Reader
	request   *http.Request
	builder   strings.Builder
}

func (c *Converter) parse() error {
	var err error
	c.request, err = http.ReadRequest(c.msgReader)
	if err != nil {
		log.Println(err)
		return errors.New("http message is illegal")
	}
	//err = c.fix()
	return err
}

func (c *Converter) check() error {
	if c.request.Host == "" {
		return errors.New("host is required")
	}
	if c.request.Method == "" {
		return errors.New("method is required")
	}
	return nil
}

func (c *Converter) fix() error {
	realLen := c.calLen()
	if c.request.ContentLength == realLen {
		return nil
	}
	compile := regexp.MustCompile(`Content-Length: \d*`)
	findString := compile.FindString(c.msg)
	if findString != "" {
		c.msg = compile.ReplaceAllString(c.msg, fmt.Sprintf("Content-Length: %d", realLen))
	}
	return c.parse()
}

func (c *Converter) calLen() int64 {
	var start bool
	length := 0
	for {
		line, _, err := c.msgReader.ReadLine()
		if err == io.EOF {
			break
		}
		if !start && len(line) == 0 {
			start = true
		}
		if start {
			length += len(line)
		}
	}
	return int64(length)
}

func (c *Converter) Echo() {
	str, err := c.do()
	if err != nil {
		println(err)
		return
	}
	println(str)
}

func (c *Converter) do() (string, error) {
	err := c.check()
	if err != nil {
		return "", err
	}
	c.major()
	c.headers()
	err = c.body()
	if err != nil {
		return "", err
	}
	return c.builder.String(), nil
}

func (c *Converter) major() {
	c.builder.WriteString(fmt.Sprintf(`curl --location --request %s "%s%s" \`, c.request.Method, c.request.Host, c.request.RequestURI))
	c.builder.WriteString("\n")

}

func (c *Converter) headers() {
	for name, vals := range c.request.Header {
		for _, val := range vals {
			c.builder.WriteString(fmt.Sprintf(`--header '%s: %s' \`, name, val))
			c.builder.WriteString("\n")
		}
	}
}

func (c *Converter) body() error {
	all, err := ioutil.ReadAll(c.request.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("parse body failed,%s", err))
	}
	c.builder.WriteString(fmt.Sprintf(`--data-raw '%s'`, all))
	return nil
}

func (c *Converter) form() error {
	err := c.request.ParseForm()
	if err != nil {
		return errors.New(fmt.Sprintf("parse form failed,%s", err))
	}
	for name, vals := range c.request.Form {
		for _, val := range vals {
			c.builder.WriteString(fmt.Sprintf(`--form '%s="%s"'`, name, val))
			c.builder.WriteString("\n")
		}
	}
	return nil
}

func readMsgFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	return string(content), nil
}
