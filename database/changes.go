package database

import (
	"log"

	"github.com/yuin/gopher-lua"
)

type viewRecalcQueue map[string]map[string]interface{}

func (vrq viewRecalcQueue) queue(key string) {
	// trigger an update on this key parent

	luacode, err := db.GetValue("/_view")
	if err != nil {
		log.Print("couldn't find view function at " + parentPath)
		return
	}

	vrq[parentKey] = luacode
}

func (vrq viewRecalcQueue) trigger() {
	for def, doc := range vrq {
		go viewRecalc(def, doc)
	}
}

func viewRecalc(definition string, doc map[string]interface{}) {
	L := lua.NewState(lua.Options{
		IncludeGoStackTrace: true,
		SkipOpenLibs:        true,
	})
	defer L.Close()
	for _, pair := range []struct {
		n string
		f lua.LGFunction
	}{
		{lua.LoadLibName, lua.OpenPackage},
		{lua.BaseLibName, lua.OpenBase},
		{lua.TabLibName, lua.OpenTable},
		{lua.MathLibName, lua.OpenMath},
		{lua.DebugLibName, lua.OpenDebug},
	} {
		if err := L.CallByParam(lua.P{
			Fn:      L.NewFunction(pair.f),
			NRet:    0,
			Protect: true,
		}, lua.LString(pair.n)); err != nil {
			log.Print("failed to load lua libraries")
		}
	}

	if err := L.DoString(definition); err != nil {
		log.Print("failed to execute view function")
		return
	}
}
