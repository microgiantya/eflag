package testtypes

type MyExportedTypeTypedString string

const (
	OneEx   MyExportedTypeTypedString = "one"
	TwoEx   MyExportedTypeTypedString = "two"
	ThreeEx MyExportedTypeTypedString = "three"
	fourEx  MyExportedTypeTypedString = "four"
)

type myUnexportedTypeTypedString string

const (
	OneUn   myUnexportedTypeTypedString = "one"
	TwoUn   myUnexportedTypeTypedString = "two"
	ThreeUn myUnexportedTypeTypedString = "three"
	fourUn  myUnexportedTypeTypedString = "four"
)

const (
	One   = "one"
	Two   = "two"
	Three = "three"
	four  = "four"
)
