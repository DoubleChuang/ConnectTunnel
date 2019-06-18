package cfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type JsonConfig struct {
	ExePath    string `json:"exe_path"`
	ServerIp   string `json:"server_ip"`
	ServerPort int    `json:"server_port"`
	BusterPort int    `json:"buster_port"`
	EnableSSL  bool   `json:"ssl"`
	MemberList []struct {
		UserName string `json:"user_name"`
		UserPswd string `json:"user_password"`
		Port     int    `json:"port"`
	} `json:"member_list"`
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
	url := "http://" + res.ServerIp + ":11111/admin?cmd=servers"
	s, err := GetConfigFromServer(url)
	for i := 0; i < len(res.MemberList); i++ {
		if contains(s, res.MemberList[i].UserName) == false {
			//delete no exist list
			res.MemberList = append(res.MemberList[:i], res.MemberList[i+1:]...)
			i--
		}
	}

	return res
}
func contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}
func GetConfigFromServer(url string) ([]string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	type memberList struct {
		Code int    `json:"Code"`
		Msg  string `json:"Msg"`
	}
	var memberLists memberList
	err := json.Unmarshal([]byte(body), &memberLists)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	memberName := make([]string, 0)
	if memberLists.Code == 200 {
		type memberInfo struct {
			Default []string `json:"default"`
		}
		var memberInfos memberInfo
		err = json.Unmarshal([]byte(memberLists.Msg), &memberInfos)
		if err != nil {
			fmt.Println("error:", err)
			return nil, err
		}
		for i, v := range memberInfos.Default {
			if i%2 == 0 {
				memberName = append(memberName, v)
			}
		}
	} else {
		return nil, fmt.Errorf("Fail Connect:%d", memberLists.Code)
	}

	return memberName, nil
}




