package obj

//
//-- CLASS: CtxKey
//

/*
Simple wrapper around a string that allows for suppression of "should not
use built-in type string as key for value" warnings.
*/
type CtxKey struct{ S string }
