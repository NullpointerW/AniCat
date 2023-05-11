package resource
type Opt int 
const(
	Ls =1
	LSGroup =2
)
const resourcesBaseUrl = `https://mikanime.tv`

var ResourceAPIs = map[string]string{
	"search": "/Home/Search?searchstr=",
}
