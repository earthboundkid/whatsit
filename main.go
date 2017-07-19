package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	flag.Parse()
	if err := Run(flag.Arg(0)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func Run(filename string) (err error) {
	f := os.Stdin
	if filename != "-" && filename != "" {
		f, err = os.Open(filename)
		if err != nil {
			return err
		}
		defer DeferClose(&err, f.Close)
	}

	const magicSize = 512 // Size that DetectContentType expects
	buf := make([]byte, magicSize)

	n, err := io.ReadFull(f, buf)
	if err != nil && err != io.ErrUnexpectedEOF {
		return err
	}
	buf = buf[:n]

	ct := http.DetectContentType(buf)
	_, err = fmt.Println(ct)
	return err
}

func DeferClose(err *error, f func() error) {
	newErr := f()
	if *err == nil {
		*err = newErr
	}
}
