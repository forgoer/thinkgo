package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func ParseAddr(addrs ...string) string {
	var addr = "0.0.0.0"
	var port = "9011"
	switch len(addrs) {
	case 0:
		if a := os.Getenv("THINKGO_ADDR"); a != "" {
			addr = a
		}
		if p := os.Getenv("THINKGO_PORT"); p != "" {
			port = p
		}
	case 1:
		strs := strings.Split(addrs[0], ":")
		if len(strs) > 0 && strs[0] != "" {
			addr = strs[0]
		}
		if len(strs) > 1 && strs[1] != "" {
			port = strs[1]
		}
	default:

	}
	return addr + ":" + port
}

func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	suffix = strings.ToUpper(suffix)
	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, strings.TrimRight(dirPth, "/")+"/"+fi.Name())
		}
	}
	return files, nil
}

func MapGet(m map[string]interface{}, key string, parms ...interface{}) interface{} {
	if value, ok := m[key]; ok {
		return value
	}
	// database.mysql.host
	s := strings.Split(key, ".")
	i := 0
	for _, segment := range s {
		i++
		if _, ok := m[segment]; !ok {
			break
		}

		b, err := json.Marshal(m[segment])
		if err != nil {
			fmt.Println(err)
		}

		if i == len(s) {
			return m[segment]
		}

		vv := make(map[string]interface{})
		err = json.Unmarshal(b, &vv)
		if err != nil {
			break
		}
		m = vv
	}
	if len(parms) == 1 {
		return parms[0]
	}
	return nil
}

func FileExists(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}
