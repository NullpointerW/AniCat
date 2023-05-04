package test

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/NullpointerW/mikanani/errs"
	"github.com/NullpointerW/mikanani/subject"
)

func TestOS(t *testing.T) {
	_, err := os.ReadDir("rows")
	fmt.Println(err)
}

func TestMap(t *testing.T) {
	subject.Manager.Add(&subject.Subject{})
	subject.Manager.Add(&subject.Subject{})
	subject.Manager.Add(&subject.Subject{})
}

func TestJsonSubj(t *testing.T) {
	b, _ := json.Marshal(subject.Subject{})
	fmt.Println(string(b))
	s := &subject.Subject{}
	js := `{"subjId":214,"name":"","path":"","finished":false,"episode":0,"resourceTyp":1,"resourceUrl":"","typ":1,"startTime":"","endTime":"","torrentHash":""}`
	json.Unmarshal([]byte(js), s)
	fmt.Printf("%#+v\n", s)
	fmt.Println(s.ResourceTyp == subject.Torrent)
	fmt.Println(s.ResourceTyp == 1)
}

func TestScan(t *testing.T) {
	subject.Scan()
}

func TestCreateSubj(t *testing.T) {
	err := subject.CreateSubject("轻音少女")
	errs.NoError(t, err)
}

func TestXxx(t *testing.T) {
	typ := reflect.TypeOf(subject.Extra{}.RssOption)
	fmt.Println("Type:", typ)
	typ = reflect.TypeOf(subject.Extra{})
	fmt.Println("Type:", typ)
}
