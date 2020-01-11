// app.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gopkg.in/ini.v1"
)

// App - Application structure
type App struct {
	Router *mux.Router
	DB     *sql.DB
	cfg    *ini.File
}

//Initialize - Initializes the database connection and routes
func (a *App) Initialize(inifile string) {
        var err error

        a.cfg, err = ini.Load(inifile)
        if err != nil {
                fmt.Printf("Failed to read configuration file: %v \n", err)
                os.Exit(1)
        }

	apiProtocol := a.cfg.Section("API").Key("protocol").In("http", []string{"http", "https"})
	apiPort := a.cfg.Section("API").Key("api_port").String()
	fmt.Println("API Protocol:", apiProtocol)
	fmt.Println("API Port:", apiPort)

	dbType := a.cfg.Section("DB").Key("type").In("mysql", []string{"mysql", "mssql", "oracle"})
	dbAddress := a.cfg.Section("DB").Key("db_address").String()
	dbPort := a.cfg.Section("DB").Key("db_port").String()
	dbUser := a.cfg.Section("DB").Key("user_id").String()
	dbPwd := a.cfg.Section("DB").Key("user_password").String()
	dbName := a.cfg.Section("DB").Key("db_name").String()

	fmt.Println("DB Type:", dbType)
	fmt.Println("DB Address:", dbAddress)
	fmt.Println("DB Port:", dbPort)
	fmt.Println("DB User:", dbUser)
	fmt.Println("DB Pwd:", dbPwd)
	fmt.Println("Db Name:", dbName)

	fragmentEntryPoint := a.cfg.Section("Fragment").Key("entry_point").String()
	fmt.Println("Fragment Entry Point: ", fragmentEntryPoint)

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPwd, dbAddress, dbPort, dbName)

	//	var err error
	a.DB, err = sql.Open(dbType, connectionString)

	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	} else {
		a.DB.Ping()
		if err != nil {
			log.Fatal(err)
			fmt.Println(err)
		} else {
			fmt.Println("Database connection established ...")
		}
	}
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	fmt.Println("API server initialized ...")
	fmt.Println("Accepting requests on port: ", apiPort)
}

//Run - Runs the listener
func (a *App) Run() {
	if a.cfg.Section("API").Key("api_cert_file").String() != "" {
		log.Fatal(http.ListenAndServeTLS(a.cfg.Section("API").Key("api_port").String(), a.cfg.Section("API").Key("api_cert_file").String(), a.cfg.Section("API").Key("api_key_file").String(), a.Router))
	} else {
		log.Fatal(http.ListenAndServe(a.cfg.Section("API").Key("api_port").String(), a.Router))
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) renderFragments(fragmentID int) string {
	fragmentChildren, err := getSubfragments(a.DB, fragmentID)
	if err != nil {
		return "renderFragments 109: " + err.Error()
	}
	var outString string
	if len(fragmentChildren) > 0 {
		//debug fmt.Println(fragmentChildren)
		//render fragment code stub; needs to be recursively structured
		for i, subfragID := range fragmentChildren {
			subfragmentOutput, err := getFragment(a.DB, subfragID)
			if err != nil {
				if err == sql.ErrNoRows {
					//do nothing, fall through to end recursion
				} else {
					return "renderFragments 120: i: " + strconv.Itoa(i) + " ID:" + strconv.Itoa(subfragID) + outString + " :: " + err.Error()
				}
			}
			if subfragmentOutput.content.String == "" {
				outString = outString + subfragmentOutput.pre.String + "\n" + a.renderFragments(subfragmentOutput.ID) + "\n" + subfragmentOutput.post.String + "\n"
			} else {
				outString = outString + subfragmentOutput.pre.String + "\n" + subfragmentOutput.content.String + "\n" + subfragmentOutput.post.String + "\n"
			}
		}
	}
	return outString
}

func (a *App) serveSPA() string {
	fragmentOutput, err := getFragment(a.DB, a.cfg.Section("Fragment").Key("entry_point").MustInt(0))
	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		} else {
			return "serveSPA 136: " + err.Error()
		}
	}

	outString := fragmentOutput.pre.String + "\n" + a.renderFragments(fragmentOutput.ID) + "\n" + fragmentOutput.post.String + "\n"
	return outString
}

func (a *App) initializeRoutes() {
	//	a.Router.HandleFunc("/cashtransactions", a.getCashTransactions).Methods("GET")
	a.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, a.serveSPA())
	})
	a.Router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

}
