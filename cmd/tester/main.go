// Copyright 2020 ActiveState Software. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"os/signal"
	"time"

	"github.com/ActiveState/termtest/winterm"
)

var exit1 = flag.Bool("exit1", false, "exit the script with exit code 1")
var sleep = flag.Bool("sleep", false, "sleep for an hour, basically never return unless interrupted")
var consoleMode = flag.Bool("console-mode", false, "show current console mode (for windows only)")
var fillBuffer = flag.Bool("fill-buffer", false, "print a string with 100,00 characters")
var stutter = flag.Bool("stutter", false, "print 50 messages with 50 ms delays")

func main() {
	c := make(chan os.Signal, 1)
	defer close(c)
	signal.Notify(c, os.Interrupt)

	flag.Parse()

	fmt.Println("an expected string")

	if *sleep {
		/* This will listen to a ctrl-c event for up to two hours
		 * Notice: That it is *only* necessary to watch for an interrupt
		 * signal on Windows.  On Linux&MacOS the interrupt signal would
		 * always break the control flow, whereas on Windows it is not really
		 * clear when and how this is happening.
		 */
		select {
		case <-time.After(1 * time.Hour):
			fmt.Println("returning after an hour, this will never happen")
		case sig := <-c:
			fmt.Printf("received %v\n", sig)
			os.Exit(123)
		}
	}

	if *consoleMode {
		mode, err := winterm.GetStdoutConsoleMode()
		if err != nil {
			log.Fatalf("Could not get console mode: %v\n", err)
		}
		fmt.Printf("console mode: %d\n", mode)
	}

	if *fillBuffer {
		/*
		err = winterm.SetConsoleMode(uintptr(stdOutHandle), winterm.ENABLE_WRAP_AT_EOL_OUTPUT | winterm.ENABLE_VIRTUAL_TERMINAL_PROCESSING | winterm.DISABLE_NEWLINE_AUTO_RETURN)
		if err != nil {
			log.Fatalf("Could not set console mode: %v\n", err)
		}
		*/
		for i := 0; i < 300; i++ {
			os.Stdout.WriteString(fmt.Sprintf(":%03d:", i))
			for j := 5; j < 80; j++ {
				os.Stdout.WriteString(fmt.Sprintf("%d", j%10))
			}
		}
		os.Stdout.Write([]byte("\n"))
	}

	if *stutter {
		for i := 0; i < 20; i++ {
			fmt.Printf("stuttered %d times\n", i+1)
			time.Sleep(50 * time.Millisecond)
		}
	}

	if *exit1 {
		os.Exit(1)
	}
}
