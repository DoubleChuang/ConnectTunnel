package cfg

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Member struct {
	UserName string `json:"user_name"`
	UserPswd string `json:"user_password"`
	Port     int    `json:"port"`
}

type JsonConfig struct {
	ExePath    string   `json:"exe_path"`
	ServerIp   string   `json:"server_ip"`
	ServerPort int      `json:"server_port"`
	EnableSSL  bool     `json:"ssl"`
	MemberList []Member `json:"member_list"`
}

func ReadFile2Buf(filePath string) []byte {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	return b
}

func GetConfig(filePath string) JsonConfig {
	byt := ReadFile2Buf(filePath)
	var res JsonConfig
	err := json.Unmarshal([]byte(byt), &res)
	if err != nil {
		panic(err)
	}
	//fmt.Println(res)
	return res
}

func main() {
	GetConfig("config.json")
}
