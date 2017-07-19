package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	if err := Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func Run() (err error) {
	flag.Parse()
	filename := flag.Arg(0)
	f := os.Stdin
	if filename != "-" && filename != "" {
		f, err = os.Open(filename)
		if err != nil {
			return err
		}
		defer DeferClose(&err, f.Close)
	}

	if ct, err := DetectContentType(f); err != nil {
		return err
	} else {
		_, err = fmt.Println(ct)
		return err
	}
}

func DetectContentType(r io.Reader) (string, error) {
	const magicSize = 512 // Size that DetectContentType expects
	buf := make([]byte, magicSize)

	n, err := io.ReadFull(r, buf)
	if err != nil && err != io.ErrUnexpectedEOF {
		return "", err
	}
	buf = buf[:n]

	return http.DetectContentType(buf), nil
}

func DeferClose(err *error, f func() error) {
	newErr := f()
	if *err == nil {
		*err = newErr
	}
}
