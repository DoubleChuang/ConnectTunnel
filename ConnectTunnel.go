package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"./cfg"
)

type tunnel struct {
	user string
	pswd string
	port string
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

func getCmd(c cfg.JsonConfig, num int) []string {
	var bash string
	var cmdHead, cmdWin, cmdTail string

	if runtime.GOOS == "windows" {
		bash = "cmd.exe"
		cmdHead = "/C"
		cmdWin = "Start "
	} else if runtime.GOOS == "linux" {
		bash = "bash"
		cmdHead = "-c"
		cmdTail = "&"
	}

	arg := []string{
		bash,
		cmdHead,
		cmdWin + c.ExePath +
			" -link " + c.MemberList[num].UserName +
			" -clientkey " + c.MemberList[num].UserPswd +
			" -local :" + strconv.Itoa(c.MemberList[num].Port) +
			" -remote=" + c.ServerIp + ":" + strconv.Itoa(c.ServerPort) +
			" -ssl=" + If(c.EnableSSL, "true", "false").(string) + cmdTail +
			" -buster " + c.ServerIp + ":" + strconv.Itoa(c.BusterPort)}

	return arg
}

func callConnect(c cfg.JsonConfig, num int) {

	/*arg := []string{
	"-c",
		c.ExePath +
		" -link "      + c.MemberList[num].UserName +
		" -clientkey " + c.MemberList[num].UserPswd +
		" -local :"    + strconv.Itoa(c.MemberList[num].Port) +
		" -remote="    + c.ServerIp + ":" + strconv.Itoa(c.ServerPort) +
		" -ssl="       + If(c.EnableSSL, "true", "false").(string)+ "&"}*/
	arg := getCmd(c, num)

	cmd := exec.Command(arg[0], arg[1], arg[2])
	if err := cmd.Run(); err != nil {
		log.Println("Error:", err)
	}

	/*output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("callConnect:\n%v\n\n%v\n\n%v", string(output), cmd.Stdout, cmd.Stderr)*/
}

func main() {
	var selected int
	var m_json string

	if len(os.Args) == 2 {
		m_json = os.Args[1]
	} else {
		m_json = "config.json"
	}

	fmt.Printf("\nLoading %s file: %s",
		If(len(os.Args) == 2, "custom", "default").(string), m_json)

	m_cfg := cfg.GetConfig(m_json)
	fmt.Printf("\n+---+-----------+------------+------+\n")
	for i, v := range m_cfg.MemberList {
		fmt.Printf("|%2d |%10s |%10s |%6d |\n",
			i+1, v.UserName, v.UserPswd, v.Port)
		fmt.Printf("+---+-----------+------------+------+\n")
	}

	fmt.Println("\nPlease enter num: ")
	fmt.Scanln(&selected)

	if selected > len(m_cfg.MemberList) || selected <= 0 {
		fmt.Printf("Over Range:1..%d", len(m_cfg.MemberList))
		return
	}
	fmt.Printf("\n+-----------------------------------+\n")
	fmt.Println(" Connecting to", m_cfg.MemberList[selected-1].UserName)
	fmt.Printf("+-----------------------------------+\n")
	callConnect(m_cfg, selected-1)
}
