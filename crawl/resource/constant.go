package resource
type LsTyp int 
const(
	Ls =1
	LSGroup =2
)
func (t LsTyp)String()string{
	if t==Ls{
		return "ls"
	}
	return "ls group"
}
const resourcesBaseUrl = `https://mikanime.tv`

var ResourceAPIs = map[string]string{
	"search": "/Home/Search?searchstr=",
}
