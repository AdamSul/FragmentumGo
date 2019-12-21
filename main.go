package main

func main() {
	//instanciate application
	a := App{}
	//initialize the application using the intraday.ini file at the specified dir ("./" means same as binary)
	a.Initialize("./")
	// use publc key stored in ini along with licensee name to decode license key.
	// payload from key includes what options are available and expiration of key
	//	if a.license.API {
	a.Run()
	//	}
}
