package conf

import "testing"

// This is correct struct
type TestArgs struct {
	Configfile string `name:"config" default:"/etc/daemon.conf" description:"Конфигурационный файл"`
	Daemon     bool   `name:"daemon" default:"false" description:"Запуск приложения в режиме daemon"`
	Test       uint32 `name:"test" default:"200" description:"Test field"`
}

// This struct has incorrect fields
type TestArgsF struct {
	Configfile string `required:"false" default:"/etc/daemon.conf"`
	Daemon     bool   `required:"true" name:"daemon" default:"false" description:"Запуск приложения в режиме daemon"`
	Test       uint32 `required:"false" name:"test" default:"200" description:"Test field"`
}

func TestGetArguments(t *testing.T) {
	correct := TestArgs{}
	err := GetArguments(&correct)
	if err != nil {
		t.Errorf("Must be ok, got %-v", err)
	}

	fail := TestArgsF{}
	err = GetArguments(&fail)
	if err == nil {
		t.Errorf("Must be error, got %-v", err)
	}
}
