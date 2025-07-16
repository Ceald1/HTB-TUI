package yaml

import (
	"fmt"
	"strings"

	"github.com/Ceald1/HTB-TUI/src/format"
	HTB "github.com/gubarz/gohtb"

)



func GetChallengeID(name string, HTBClient HTB.Client) (id int, err error) {
	name = strings.ToLower(name)
	format.TaskResult = 0
	task := format.Task(func(a any) any {
		HTBClient, _ := a.(*HTB.Client)
		contents, _ := HTBClient.Challenges.List().AllResults(ctx)
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
	if format.TaskResult == 0 {
		err = fmt.Errorf("unable to find challenge")
	}
	id, _ = format.TaskResult.(int)
	return 

}