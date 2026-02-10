package yaml

import (
	"fmt"
	"strings"

	"github.com/Ceald1/HTB-TUI/src/format"
	HTB "github.com/gubarz/gohtb"

)

func GetMachineID(name string, HTBClient HTB.Client) (id int, err error) {
	name = strings.ToLower(name)
	format.TaskResult = 0
	task := format.Task(func(a any) any {
		HTBClient, _ := a.(*HTB.Client)
		contents, _ := HTBClient.Machines.List().ByState("active").AllResults(ctx)
		Response := contents.Data
		for _, r := range Response {
			if strings.ToLower(r.Name) == name {
				id = r.Id
				return id
			}
		}
		contents, err = HTBClient.Machines.List().ByState("retired").AllResults(ctx)

		Response = contents.Data
		for _, r := range Response {
			if strings.ToLower(r.Name) == name {
				id = r.Id
				return id
			}
		}

		contents1, _ := HTBClient.Machines.List().ByState("unreleased").AllResults(ctx)

		Response1 := contents1.Data
		for _, r := range Response1 {
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

	id, _ = format.TaskResult.(int)
	if id == 0{
		err = fmt.Errorf("unable to find machine")
	}
	return

}
