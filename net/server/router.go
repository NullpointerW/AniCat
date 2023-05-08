package server

import (
	"strconv"

	"github.com/NullpointerW/mikanani/net/cmd"
	"github.com/NullpointerW/mikanani/subject"
	"github.com/liushuochen/gotable"
)

func route(c *cmd.Command) {
	switch c.Opt {
	case cmd.Add:
		p := subject.NewPip(c.N)
		subject.Create <- p
		c.Err = p.Error()
	case cmd.Del:
		i, err := strconv.Atoi(c.N)
		if err != nil {
			c.Err = err
			return
		}
		p := subject.NewPip(i)
		subject.Create <- p
		c.Err = p.Error()
	case cmd.Ls:
		ls := subject.Manager.List()
		tb, _ := gotable.Create("sid", "type", "name", "episode", "status", "compl")
		for _, s := range ls {
			sid := strconv.Itoa(s.SubjId)
			fin := "updating"
			compl := "N"
			epi := strconv.Itoa(s.Episode)
			if s.Finished {
				fin = "fin"
			}
			if s.Terminate {
				compl = "Y"
			}
			tb.AddRow([]string{sid, s.Typ.String(), s.Name, epi, fin, compl})
		}
		c.N = "\n" + tb.String()
	}
}
