package test

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"

	// "strings"

	// "strings"

	"testing"

	"github.com/NullpointerW/anicat/errs"
	"github.com/NullpointerW/anicat/subject"
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
	_, err := subject.CreateSubject("未闻花名", nil)
	errs.NoError(t, err)
}

func TestXxx(t *testing.T) {
	typ := reflect.TypeOf(subject.Extra{}.RssOption)
	fmt.Println("Type:", typ)
	typ = reflect.TypeOf(subject.Extra{})
	fmt.Println("Type:", typ)

	s := subject.Subject{}
	s.Pushed = make(map[string]string)
	s.Pushed["v1"] = "struct{}{}"
	b, _ := json.Marshal(s)
	fmt.Println(string(b))
	err := json.Unmarshal(b, &s)
	if err != nil {
		t.FailNow()
	}
}

func TestGetSeason(t *testing.T) {
	var s subject.Subject
	s.Name = "小林家的龙女仆 第二季"
	subject.GetSeason(&s)
	t.Log(s.Season)
}

func TestExtra(t *testing.T) {
	re, err := regexp.Compile(`\[[^\]]*\d+(\.\d+)\s*[Gg][Bb][^\]]*\]`)
	if err != nil {
		fmt.Println("正则表达式编译失败：", err)
		return
	}

	tests := []string{"【恶魔岛字幕组】★4月新番【吹响！上低音号_Hibike! Euphonium】[01-13][GB_BIG5][720P][MKV][全][ 4.7 GB]", "3.5gb", "3.5 GB", "3.5gb ", "3.5GB", "3.5gB", "3gb ", "3.0 GB", "3.0gb"}
	for _, test := range tests {
		if re.MatchString(test) {
			fmt.Println(test, "符合要求")
		} else {
			fmt.Println(test, "不符合要求")
		}
	}
}

func TestFoundLastS(t *testing.T) {
	l, err := subject.FindLastSeason(`D:\anicat\吹响！悠风号 (201504)`)
	if err != nil {
		t.Error(l)
		t.FailNow()
	}
	t.Log(l)
}
func TestReg(t *testing.T) {
	str := "这是第一季的节目"
	re := regexp.MustCompile(`第(.)季`)
	match := re.FindStringSubmatch(str)
	if len(match) > 1 {
		season := match[1]
		fmt.Printf("匹配到的季：%s\n", season)
	} else {
		fmt.Println("未匹配到季")
	}
}

func TestBuildFilterReg(t *testing.T) {
	reg := subject.BuildFilterPerlReg([]string{"1080p,1080x1920p", "简体中文,CHS", "v2"})
	fmt.Println(reg)
}

func TestFilterReg(t *testing.T) {
	cs := subject.BuildFilterRegs([]string{"1080p,1080x1920p", "简体中文,CHS"})
	fmt.Println(cs)
	cls := subject.BuildFilterRegs([]string{"外挂"})
	fmt.Println(cls)
	ok := subject.FilterWithRegs("[Lilith-Raws] Okashi na Tensei - 06 [Baha][WebDL 1080p AVC AAC][CHS]", cs, cls)
	fmt.Println(ok)
}

func TestSubsReg(t *testing.T) {
	CHSReg, err := regexp.Compile(subject.CHSSubStationReg)
	if err != nil {
		t.Error(err)
	}
	ok := CHSReg.MatchString("[Lilith-Raws] Okashi na Tensei - 06 [Baha][WebDL 1080p AVC AAC][]Chs.ast")
	fmt.Println(ok)

	CHSReg, err = regexp.Compile(subject.CHTSubStationReg)
	if err != nil {
		t.Error(err)
	}
	ok = CHSReg.MatchString("[Lilith-Raws] Okashi na Tensei - 06 [Baha][WebDL 1080p AVC AAC][]繁体中文.ast")
	fmt.Println(ok)
}
