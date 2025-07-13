package main

import (
	SOURCE "github.com/Ceald1/HTB-TUI/src"
)


func main(){
	instance, err  := SOURCE.API_Token()
	if err != nil {
		panic(err)
	}
	SOURCE.MainMenu(instance)
	
	

}
