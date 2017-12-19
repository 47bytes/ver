package main

import "github.com/47bytes/ver/cmd"

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		//fmt.Println("Unexpected error. " + err.Error())
		return
	}
}
