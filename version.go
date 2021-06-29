package main

import "fmt"

const Version = "0.0.1"


var banner = []string{
	"    ____  __________________  ____  _   __    ",
	"   / __ \\/  _/ ____/  _/ __ \\/ __ \\/ | / / ",
	"  / /_/ // // / __ / // /_/ / / / /  |/ /     ",
	" / ____// // /_/ // // _, _/ /_/ / /|  /      ",
	"/_/   /___/\\____/___/_/ |_|\\____/_/ |_/     "} 



func printBanner() {
	for _, line := range banner {
		fmt.Println(line)
	}
	fmt.Println()
	fmt.Printf("Version %s\n", Version)
	fmt.Println()
}
		
	
