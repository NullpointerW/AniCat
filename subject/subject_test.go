package subject

import (
	"fmt"
	"os"
	"testing"
)

func TestOS(t *testing.T) {
	_,err:=os.ReadDir("rows")
	fmt.Println(err)
}

func TestMap(t *testing.T){
	Manager.Add(8848,&Subject{},true)
	Manager.Add(8849,&Subject{},true)
	Manager.Add(8850,&Subject{},false)
	Manager.Move(8848,false)
	fmt.Println(len(Manager.finished))
	fmt.Println(len(Manager.unfinished))


}