/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"net"
	"os"
	"strings"
	"time"
)

func waitFor(target string, timeout uint, quiet bool) error {
	if !quiet {
		if timeout > 0 {
			fmt.Printf("Waiting for %s to start within %d seconds.\n", target, timeout)
		} else {
			fmt.Printf("Waiting for %s indefinitely ... \n", target)
		}
	}
	conn, err := net.DialTimeout("tcp", target, time.Duration(timeout)*time.Second)
	if conn != nil {
		defer conn.Close()
	}
	if err, ok := err.(*net.OpError); ok && err.Timeout() {
		return fmt.Errorf("timed out after %d seconds", timeout)
	}
	if err != nil {
		return err
	}

	return nil

}

func catchHelp() {
	if pflag.Arg(0) == "help" {
		fmt.Println("waitfor is a go cli that will wait on the availability of a host and TCP port.")
		fmt.Println("For example to manage the start order of docker containers.")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("waitfor host:port [-q quiet] [-t timeout]")
		fmt.Println("")
		pflag.PrintDefaults()
		fmt.Println("")

		os.Exit(0)
	}
}

func main() {
	timeout := pflag.UintP("timeout", "t", 0, "The timeout in seconds until the service is considered non-responsive")
	quiet := pflag.BoolP("quiet", "q", false, "Only write errors to output.")
	pflag.Parse()

	catchHelp()

	if len(pflag.Args()) != 1 {
		fmt.Println("waitfor requires exactly one argument: 'host:port'")
		fmt.Printf("You called it with: %s\n", pflag.Args())
		os.Exit(1)
	}
	if !strings.Contains(pflag.Arg(0), ":") {
		fmt.Println("The argument is not valid, please make sure you are calling waitfor with a valid argument: 'host:port'")
	}

	err := waitFor(pflag.Arg(0), *timeout, *quiet)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !*quiet {
		fmt.Printf("%s is up and running.\n", pflag.Arg(0))
	}
	os.Exit(0)
}
