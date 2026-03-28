package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {

	times := []int{}
	cmd := exec.Command("sudo", "rm", "/var/log/cwall/cwall.log")

	res, err := cmd.CombinedOutput()

	println(string(res))

	if err != nil {
		println(err.Error())
	}

	cmd = exec.Command("sudo", "systemctl", "restart", "cwalld-enforce.service")

	res, err = cmd.CombinedOutput()
	
	println(string(res))

	cmd = exec.Command("sudo", "sh", "-c", "sudo echo '' > /home/testgrounds/objects/beta_plans")
	res, _ = cmd.CombinedOutput()
	println(string(res))

	if err != nil {
		println(err.Error())
	}

	MAX := 500
	timer := 0 // time is the amount of milliseconds that the daemon will sleep for on this run
	breached := 0
	for { // loop until we get it
		newSpeedd(MAX - timer)

		time.Sleep(time.Millisecond * 600)
		
		breach := checkBeta()

		if breach {
			println("1 breach!")
			times = append(times, MAX - timer)
			if (MAX - timer) < 125 {
				fmt.Printf("Time found: %d x 10^-5!", MAX - timer)
				cmd := exec.Command("sudo", "systemctl", "stop", "cwalldspeedd.service")
				_ = cmd.Run()

				f, err := os.OpenFile("results.csv", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

				if err != nil {
					println(err.Error())
				}

				writer := bufio.NewWriter(f)

				for t := range times {
					writer.WriteString((strconv.Itoa(times[t]) + ","))
				}

				writer.Flush()

				f.Close()
				os.Exit(1)
			}
			cmd := exec.Command("sudo", "sh", "-c", "sudo echo '' > /home/testgrounds/objects/beta_plans")
			res, _ := cmd.CombinedOutput()

			println("Beta current contents below")
			
			cmd = exec.Command("sudo", "cat", "/home/testgrounds/objects/beta_plans")
			res, err = cmd.CombinedOutput()

			println(string(res))
			breached++
			breach = false
		}

		timer += 1

	}
}

func newSpeedd(num int) {
	fmt.Printf("Making new speed d with %d x 10^-5\n", num)
	cmd := exec.Command("sudo", "systemctl", "stop", "cwalldspeedd.service")

	res, err := cmd.CombinedOutput()
	// res = res

	println(string(res))

	if err != nil {
		println(err.Error())
	}

	editLine(num)

	cmd = exec.Command("gcc", "-o", "/home/cal/cs408-cwalld/cwalld_test/benchmark/cwalldspeedd/cwalldspeedd", "/home/cal/cs408-cwalld/cwalld_test/benchmark/cwalldspeedd/cwalldspeedd.c")

	res, err = cmd.CombinedOutput()

	// println(string(res)) 

	if err != nil {
		println(err.Error())
	}

	cmd = exec.Command("sudo", "cp", "/home/cal/cs408-cwalld/cwalld_test/benchmark/cwalldspeedd/cwalldspeedd", "/usr/local/bin")

	res, err = cmd.CombinedOutput()

	// println(string(res)) 

	if err != nil {
		println(err.Error())
	}

	cmd = exec.Command("sudo", "chcon", "-t", "bin_t", "/usr/local/bin/cwalldspeedd")

	res, err = cmd.CombinedOutput()

	// println(string(res)) 

	if err != nil {
		println(err.Error())
	}

	cmd = exec.Command("sudo", "systemctl", "start", "cwalldspeedd.service")

	res, err = cmd.CombinedOutput()

	println(string(res)) 

	if err != nil {
		println(err.Error())
	}

}

func editLine(num int) {
	path := "/home/cal/cs408-cwalld/cwalld_test/benchmark/cwalldspeedd/cwalldspeedd.c"
	data, err := os.ReadFile(path)

	if err != nil {
		fmt.Printf("Error %s", err.Error())
	}

	lines := strings.Split(string(data), "\n")

	content := fmt.Sprintf("      usleep(%d);", num * 10) // should change how long the daemon sleeps for

	lines[15] = content

	err = os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
}

func checkBeta() bool {
	cmd := exec.Command("sudo", "cat", "/home/testgrounds/objects/beta_plans")

	res, err := cmd.CombinedOutput()

	println(string(res))

	if bytes.Contains(res, []byte("Don't let this get to beta!")) {
		return true
	}

	if err != nil {
		println(err.Error())
	}

	return false
}
