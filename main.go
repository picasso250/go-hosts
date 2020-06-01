package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const hostsFileName = "c:\\windows\\System32\\drivers\\etc\\hosts"

func getURL() string {
	f, err := os.Open("url")
	if err != nil {
		return "http://xx"
	}
	defer f.Close()
	s, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Trim(string(s), "\n\r\t ")
}
func get() []byte {
	url := getURL()
	fmt.Printf("Fetching from %s\n", url)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	z, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", z)
	return z
}
func new(hosts []byte) []byte {
	f, err := os.Open(hostsFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	s, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("old hosts\n%s\n", string(s))
	return subst(s, hosts)
}

func subst(allHosts []byte, newHosts []byte) []byte {
	const begin = "#### BEGIN GO-HOSTS ####"
	const end = "#### END GO-HOSTS ####"
	const newline = "\r\n"

	s := [][]byte{[]byte(begin), newHosts, []byte(end)}
	toInsert := bytes.Join(s, []byte(newline))

	b := bytes.Index(allHosts, []byte(begin))
	if b < 0 {
		s := [][]byte{allHosts, toInsert}
		return bytes.Join(s, []byte(newline))
	}
	e := bytes.Index(allHosts, []byte(end))
	if e < 0 {
		log.Fatal("no end")
	}
	e += len([]byte(end))

	s = [][]byte{allHosts[:b], toInsert, allHosts[e:]}
	return bytes.Join(s, []byte(""))
}
func main() {
	h := get()
	hosts := new(h)
	err := ioutil.WriteFile(hostsFileName, hosts, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Wrote to %s\n", hostsFileName)
}
