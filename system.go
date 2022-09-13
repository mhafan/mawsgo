package mawsgo

import "os"

// ---------------------------------------------------------------------------
// system ENV | def ->
func Env(envName, defValue string) string {
	//
	val := os.Getenv(envName)

	//
	if len(val) <= 0 {
		//
		return defValue
	}

	//
	return val
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
