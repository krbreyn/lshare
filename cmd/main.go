package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/krbreyn/lshare"
)

type ProgressTracker struct {
	Total uint64
}

func (pt *ProgressTracker) Write(p []byte) (int, error) {
	n := len(p)
	pt.Total += uint64(n)
	pt.PrintProgress()
	return n, nil
}

func (pt *ProgressTracker) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	fmt.Printf("\rDownloading... %s", humanize.Bytes(pt.Total))
}

//look into generating a keypass by encoding ip/port/url into a string
//do base64 encoding and have the option be set to a flag

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
				fmt.Println("no extension provided")
				//return
			}
			data, err := io.ReadAll(os.Stdin)
			must(err)
			server := lshare.NewFileServer()

			url := os.Args[2]
			server.RegisterEndpoint(url, os.Args[2], data)

			ip, _ := lshare.GetLocalIP()
			fmt.Printf("serving at %s:%s/%s\n", ip, "8000", url)

			go server.StartServer(":8000")
			fmt.Printf("type sendto client %s %s %s on your client\n", strings.TrimPrefix(ip, "192.168."), "8000", url)
			fmt.Println("press ctrl+c to quit")
			select {}
			//_, _ = fmt.Scanln("\n") // doesnt work for stdin
			// return
		}
	}

	if len(os.Args) == 1 {
		fmt.Println("please use send, stdin, or type in your client keys")
		return
	}

	if len(os.Args) == 3 && os.Args[1] == "send" {
		fmt.Println("sending", os.Args[2])
		if path.Ext(os.Args[2]) == "" {
			fmt.Println("no extension provided")
			//return
		}

		file, err := os.Open(os.Args[2])
		must(err)

		data, err := io.ReadAll(file)
		must(err)

		server := lshare.NewFileServer()

		url := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(os.Args[2], "(", ""), ")", ""), " ", "")
		server.RegisterEndpoint(url, os.Args[2], data)

		ip, _ := lshare.GetLocalIP()
		fmt.Printf("serving at %s:%s/%s\n", ip, "8000", url)

		go server.StartServer(":8000")
		fmt.Printf("type sendto client %s %s %s on your client\n", strings.TrimPrefix(ip, "192.168."), "8000", url)
		fmt.Println("press enter to quit")
		_, _ = fmt.Scanln() // wait
		return
	}

	if len(os.Args) == 5 && os.Args[1] == "client" {
		ip := "192.168." + os.Args[2]
		port := os.Args[3]
		path := os.Args[4]

		url := fmt.Sprintf("http://%s:%s/%s", ip, port, path)
		resp, err := http.Get(url)
		must(err)

		defer resp.Body.Close()

		filename := resp.Header.Get("Content-Disposition")
		filename = strings.TrimPrefix(filename, "attachment; filename=")
		if resp.Header.Get("Content-Disposition") == "" {
			fmt.Println("error:", filename)
			return
		}
		fmt.Println(filename)

		//TODO prompt to actually accept download

		out, err := os.Create(filename)
		must(err)

		defer out.Close()

		tracker := &ProgressTracker{}
		_, err = io.Copy(out, io.TeeReader(resp.Body, tracker))
		must(err)

		fmt.Println("\nsaved", filename)
		return
	}

	fmt.Println("invalid inputs")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
