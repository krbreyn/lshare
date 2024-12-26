package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/krbreyn/sendto"
)

//encode port information in url

func main() {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		var filename string
		fmt.Println("data is being piped to stdin")
		flag.StringVar(&filename, "file", "", "name of input file")
		flag.Parse()
		if filename == "" {
			fmt.Println("must provide -file filename argument")
		} else {
			fmt.Println(filename)
			if path.Ext(filename) == "" {
				fmt.Println("must provide extension")
				return
			}
			return
		}
	}

	if len(os.Args) == 1 {
		fmt.Println("please route a file into stdin or type in your client keys")
		return
	}

	if len(os.Args) == 3 && os.Args[1] == "send" {
		fmt.Println("sending", os.Args[2])
		if path.Ext(os.Args[2]) == "" {
			fmt.Println("must provide extension")
			return
		}

		file, err := os.Open(os.Args[2])
		if err != nil {
			panic(err)
		}

		data, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}

		server := sendto.NewFileServer()

		url := strings.TrimSuffix(os.Args[2], path.Ext(os.Args[2]))

		server.RegisterEndpoint(url, os.Args[2], data)

		ip, _ := sendto.GetLocalIP()
		fmt.Printf("serving at %s:%s/%s\n", ip, "8000", url)

		server.StartServer(":8000")
		//todo press enter to close and quit
		return
	}

	if len(os.Args) == 5 && os.Args[1] == "client" {
		ip := "192.168." + os.Args[2]
		port := os.Args[3]
		path := os.Args[4]

		url := fmt.Sprintf("http://%s:%s/%s", ip, port, path)
		//fmt.Println(url)
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		filename := resp.Header.Get("Content-Disposition")
		filename = strings.TrimPrefix(filename, "attachment; filename=")
		if resp.Header.Get("Content-Disposition") == "" {
			fmt.Println("error:", filename)
			return
		}
		fmt.Println(filename)

		//todo prompt to actually accept download

		out, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer out.Close()

		_, err = out.ReadFrom(resp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println("saved", filename)
	}

}
