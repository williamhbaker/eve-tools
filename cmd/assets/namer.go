package main

import (
	"net/http"

	"github.com/wbaker85/eve-tools/pkg/lib"
)

type nameable interface {
	addName(string)
	id() int
}

func namer(n []nameable) {
	ids := []int{}

	for _, val := range n {
		ids = append(ids, val.id())
	}

	esi := lib.Esi{
		Client: http.DefaultClient,
	}

	nameList := esi.ItemNameList(ids)

	for _, val := range n {
		val.addName(nameList[val.id()])
	}
}
