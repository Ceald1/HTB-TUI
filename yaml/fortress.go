package yaml

import (
	"fmt"
	"strings"

	"github.com/Ceald1/HTB-TUI/src/format"
	HTB "github.com/gubarz/gohtb"

)



func GetFortressID(name string, HTBClient HTB.Client) (id int, err error) {
	name = strings.ToLower(name)
	format.TaskResult = 0
	task := format.Task(func(a any) any {
		contents, _ := HTBClient.Fortresses.List(ctx)

		Response := contents.Data
		for _, r := range Response {
			if strings.ToLower(r.Name) == name {
				id = r.Id
				return id
			}
		}
	return 0
	})
	err = format.RunLoading(task, &HTBClient)
	if err != nil {
		return
	}
	if format.TaskResult == 0 {
		err = fmt.Errorf("unable to find challenge")
		return
	}
	id, _ = format.TaskResult.(int)
	return

}