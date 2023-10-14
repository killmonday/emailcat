package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	emailverifier "github.com/AfterShip/email-verifier"
)

var (
	verifier = emailverifier.
		NewVerifier().
		EnableSMTPCheck().DisableCatchAllCheck()
)


var queue = make(chan string, 1000)
var wg_verfiy sync.WaitGroup
var wg_read sync.WaitGroup

func read_email_addr(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		email_addr := strings.TrimSpace(scanner.Text())
		if email_addr != "" {
			queue <- email_addr
		}
	}
	return err
}

func verify_email(queue_email chan string, sport string) error {
	for {
		select {
		case email_addr, ok := <-queue_email:
			if ok {
				username_domain := strings.Split(email_addr, "@")
				username := username_domain[0]
				domain := username_domain[1]
				ret, err := verifier.CheckSMTP(domain, username, sport)
				if err != nil {
					// fmt.Println("check smtp failed: ", err)
					continue
				}
				if ret.Deliverable == true {
					fmt.Println(email_addr)
				}

			} else {
				return nil
			}
		}
	}

}

func main() {
	var mail_path string
	var thread_num int
	var smtp_port string
	flag.StringVar(&mail_path, "f", "input.txt", "email address file")
	flag.IntVar(&thread_num, "t", 2, "thread number")
	flag.StringVar(&smtp_port, "p", "25", "smtp port")
	flag.Parse()

	wg_read.Add(1)
	go func() {
		defer wg_read.Done()
		read_email_addr(mail_path)
	}()

	for i := 0; i < thread_num; i++ {
		wg_verfiy.Add(1)
		go func() {
			defer wg_verfiy.Done()
			verify_email(queue, smtp_port)
		}()
	}
	wg_read.Wait()
	close(queue)
	wg_verfiy.Wait()
}
