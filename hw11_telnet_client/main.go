package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
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
		os.Exit(1)
	}
	tc := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	err = tc.Connect()
	if err != nil {
		fmt.Printf("connection error: %v\n", err)
		os.Exit(1)
	}
	_, _ = fmt.Fprintf(os.Stderr, "Connected to: %s\n", address)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go send(&tc, &wg)
	go receive(&tc)
	go handleSignal(&wg)

	wg.Wait()
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

func send(tc *TelnetClient, wg *sync.WaitGroup) {
	if err := (*tc).Send(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "bye-bye\n")
	}
	_ = (*tc).Close()
	wg.Done()
}

func receive(tc *TelnetClient) {
	if err := (*tc).Receive(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "Connection closed by peer\n")
	}
}

func handleSignal(wg *sync.WaitGroup) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT)
	<-c
	wg.Done()
}
