package main

import (
	SOURCE "github.com/Ceald1/HTB-TUI/src"
)


func main(){
	client, err  := SOURCE.API_Token()
	if err != nil {
		panic(err)
	}
	

}
