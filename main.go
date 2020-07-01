package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	// 🐒 patching of "database/sql".
	_ "github.com/go-sql-driver/mysql"
)

const (
	query       = "SELECT * FROM sys.gr_member_routing_candidate_status"
	mysqlParams = "collation=utf8mb4_0900_ai_ci"
)

var (
	// Version version.
	Version = "DEV"

	// We pass credentials in this env var as there is no better way of doing this from haproxy.
	mysqlCredentials   = os.Getenv("PATH")
	mysqlAddress       = os.Getenv("HAPROXY_SERVER_ADDR") + ":" + os.Getenv("HAPROXY_SERVER_PORT")
	haproxyBackendName = os.Getenv("HAPROXY_PROXY_NAME")
)

type eventRow struct {
	ViableCandidate    string
	ReadOnly           string
	TransactionsBehind string
	TransactionsToCert string
}

func debugMsg(isDebug bool, msg string) {
	if isDebug {
		fmt.Println(msg)
	}
}

func main() {
	var versionFlag, debugFlag bool
	flag.BoolVar(&versionFlag, "v", false, "show version")
	flag.BoolVar(&debugFlag, "d", false, "enable debug output")
	flag.Parse()
	if versionFlag {
		fmt.Println("Version", Version)
		os.Exit(0)
	}

	if !strings.HasSuffix(haproxyBackendName, "_primary") && !strings.HasSuffix(haproxyBackendName, "_secondary") {
		debugMsg(debugFlag, "Haproxy backend name does not end with either _primary or _secondary.")
		os.Exit(1)
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s@tcp(%s)/?%s", mysqlCredentials, mysqlAddress, mysqlParams))
	if err != nil {
		fmt.Println("Error connecting to MySQL", err)
		os.Exit(1)
	}

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error selecting from MySQL table:", err)
		os.Exit(1)
	}

	var row eventRow
	for rows.Next() {
		rows.Scan(&row.ViableCandidate, &row.ReadOnly, &row.TransactionsBehind, &row.TransactionsToCert)
		debugMsg(debugFlag, fmt.Sprintf("MySQL query result: %+v\n", row))
		break
	}

	if row.ViableCandidate != "YES" {
		debugMsg(debugFlag, "GR member is not viable candidate.")
		os.Exit(1)
	}

	if strings.HasSuffix(haproxyBackendName, "_primary") && row.ReadOnly == "NO" {
		debugMsg(debugFlag, "HEALTHCHECK PRIMARY - OK")
		return
	} else if strings.HasSuffix(haproxyBackendName, "_secondary") && row.ReadOnly == "YES" {
		debugMsg(debugFlag, "HEALTHCHECK SECONDARY - OK")
		return
	}

	debugMsg(debugFlag, "HEALTHCHECK - NOT OK")
	os.Exit(1)
}
