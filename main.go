package main

import (
	"fmt"
	"os"

	SOURCE "github.com/Ceald1/HTB-TUI/src"
	YAML "github.com/Ceald1/HTB-TUI/yaml"
)


func main(){

	// TODO: Add yaml automation support (now make it pretty)
	args := os.Args[1:]
	if len(args) == 0{
		instance, err  := SOURCE.API_Token()
		if err != nil {
			panic(err)
		}
		SOURCE.MainMenu(instance)
	}else{
		fmt.Println("running in automation mode....")
	}
	for _, arg := range args{
		YAML.RunAutomation(arg)
	}

	

}
