package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration

	flag.DurationVar(&timeout, "timeout", time.Second*10, "")
	flag.Parse()

	address, err := getAddress()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Usage %s [--timeout] host port\n", os.Args[0])
		flag.PrintDefaults()
		return
	}
	tc := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	err = tc.Connect()
	if err != nil {
		fmt.Printf("connection error: %v\n", err)
		os.Exit(1)
	}
	_, _ = fmt.Fprintf(os.Stderr, "Connected to: %s\n", address)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go send(&tc, cancel)
	go receive(&tc, cancel)

	<-ctx.Done()
	err = tc.Close()
	if err != nil {
		fmt.Println("connection close error: " + err.Error())
	}
}

func getAddress() (string, error) {
	if len(flag.Args()) < 2 {
		return "", fmt.Errorf("host and port should be provided")
	}
	host := flag.Arg(0)
	port, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		return "", fmt.Errorf("can not parse port: %w", err)
	}
	return fmt.Sprintf("%s:%d", host, port), nil
}

func send(tc *TelnetClient, cancelFunc context.CancelFunc) {
	if err := (*tc).Send(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "bye-bye\n")
	}
	cancelFunc()
}

func receive(tc *TelnetClient, cancelFunc context.CancelFunc) {
	if err := (*tc).Receive(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "Connection closed by peer\n")
	}
	cancelFunc()
}
