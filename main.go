package main

import (
        "flag"
)

func main() {
        //instanciate application
        a := App{}
        //initialize the application using the local ./fragment.ini file by default of the one specified in the "inifile" command line parameter

        var inifile string
        flag.StringVar(&inifile, "inifile", "fragments.ini", "File name for initialization parameters.")
        flag.Parse()

        a.Initialize(inifile)
        // use publc key stored in ini along with licensee name to decode license key.
        // payload from key includes what options are available and expiration of key
//        if a.license.API {
                a.Run()
//        }
}

