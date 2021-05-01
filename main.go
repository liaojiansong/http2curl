package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var Verbose bool

func main() {
	var httpMsg string
	flag.StringVar(&httpMsg, "m", "", "standard http msg")
	flag.BoolVar(&Verbose, "v", false, "show verbose")
	flag.Parse()
	if httpMsg == "" {
		log.Fatalf("please fill the 'm' optionl")
	}
	newConverter := NewConverter(httpMsg)
	newConverter.echo()

}
func NewConverter(msg string) *converter {
	if Verbose {
		println(msg)
	}
	c := &converter{
		originalMsg: msg,
		msgReader:   bufio.NewReader(strings.NewReader(msg)),
	}
	var err error
	c.request, err = http.ReadRequest(c.msgReader)
	if err != nil {
		fmt.Printf("===> %#v <===\n", "http message is illegal")
		os.Exit(-1)
	}
	return c
}

type converter struct {
	originalMsg string
	msgReader   *bufio.Reader
	request     *http.Request
	cBuilder    strings.Builder
}

func (c *converter) check() {
	if c.request.Host == "" {
		_, _ = fmt.Fprintf(os.Stdout, "%s ", "Header->host is required")
		os.Exit(-1)
	}
	if c.request.Method == "" {
		_, _ = fmt.Fprintf(os.Stdout, "%s ", "request method is required")
		os.Exit(-1)
	}
}

func (c *converter) echo() {
	c.major()
	c.headers()
	c.body()
	println(c.cBuilder.String())

}

func (c *converter) major() {
	c.cBuilder.WriteString(fmt.Sprintf(`curl --location --request %s "%s%s" \`, c.request.Method, c.request.Host, c.request.RequestURI))
	c.cBuilder.WriteString("\n")

}
func (c *converter) headers() {
	for name, vals := range c.request.Header {
		for _, val := range vals {
			c.cBuilder.WriteString(fmt.Sprintf(`--header '%s: %s' \`, name, val))
			c.cBuilder.WriteString("\n")
		}
	}
}
func (c *converter) body() {
	all, err := ioutil.ReadAll(c.request.Body)
	if err != nil {
		log.Fatalf("parse body failed,%s", err)
	}
	c.cBuilder.WriteString(fmt.Sprintf(`--data-raw '%s'`, all))
}
func (c *converter) form() {
	err := c.request.ParseForm()
	if err != nil {
		log.Fatalf("parse form failed,%s", err)
	}
	for name, vals := range c.request.Form {
		for _, val := range vals {
			c.cBuilder.WriteString(fmt.Sprintf(`--form '%s="%s"'`, name, val))
			c.cBuilder.WriteString("\n")
		}
	}
}
