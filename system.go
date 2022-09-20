package mawsgo

import "os"

// ---------------------------------------------------------------------------
// system ENV | def ->
func Env(envName, defValue string) string {
	//
	if val := os.Getenv(envName); len(val) > 0 {
		//
		return val
	}

	//
	return defValue
}

// ---------------------------------------------------------------------------
// ...
func Ifpanic(e error) {
	//
	if e != nil {
		//
		panic(e)
	}
}
