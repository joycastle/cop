package config

import (
	"testing"
	"time"

	"github.com/joycastle/cop/util"
)

type Mysql struct {
	User           string        `yaml:"user"`
	Password       string        `yaml:"mypassword"`
	Host           string        `yaml:"host"`
	Port           int           `yaml:"port"`
	DbName         string        `yaml:"dbname"`
	ConnectTimeout time.Duration `yaml:"connectTimeout"`
}

type yamlConfig struct {
	M Mysql `yaml:"mysql"`
}

func Test_YamlConfig(t *testing.T) {
	content := `mysql:
 user: root
 password: mypassword
 host: 192.168.1.1
 port: 3306
 dbname: mydb1
 connectTimeout: 500ms`
	if err := util.CreateFileWithContent("./test.yaml", []byte(content)); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := util.DeleteFile("./test.yaml"); err != nil {
			t.Fatal(err)
		}
	}()

	var v yamlConfig

	if err := ReadYmalFromFile("test.yaml", &v); err != nil {
		t.Fatal(err)
	} else {
		if v.M.User != "root" || v.M.Port != 3306 || v.M.DbName != "mydb1" || v.M.ConnectTimeout != time.Millisecond*500 {
			t.Fatal(v)
		}
	}

}
