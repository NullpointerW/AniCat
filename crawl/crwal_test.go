package crawl

import (
	"fmt"
	"testing"

	"github.com/tidwall/gjson"
)





func TestJsonArray(t *testing.T) {
	jsonstr := `[1,2,3]`
	fmt.Println(gjson.Get(jsonstr, "0").Int())
	// https://movie.douban.com/subject/4074292/?suggest=%E5%87%89%E5%AE%AB+%E6%98%A5%E6%97%A5
}




