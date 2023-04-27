package test

import (
	"fmt"
	"github.com/NullpointerW/mikanani/errs"
	"github.com/NullpointerW/mikanani/subject"
	"os"
	"testing"
)

func TestOS(t *testing.T) {
	_, err := os.ReadDir("rows")
	fmt.Println(err)
}

func TestMap(t *testing.T) {
	subject.Manager.Add(8848, &subject.Subject{}, true)
	subject.Manager.Add(8849, &subject.Subject{}, true)
	subject.Manager.Add(8850, &subject.Subject{}, false)
	subject.Manager.Move(8848, false)

}

func TestScan(t *testing.T) {
	subject.Scan()
}

func TestCreateSubj(t *testing.T) {
	err := subject.CreateSubject("轻音少女")
	errs.NoError(t, err)
}
