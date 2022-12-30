package main

import "gitlab.com/kendellfab/mtgstudio/internal"

func main() {
	mtgStudio := internal.NewMtgStudio()
	mtgStudio.Start()
}
