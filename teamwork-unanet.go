package main

import (
	"flag"
	"fmt"
	"os"

	teamworkapi "github.com/Foxtrot-Division/teamworkAPI"
	unanet "github.com/Foxtrot-Division/teamworkTime/unanent"
)

const (
	apiConf               = "api.json"
	peopleTimeDetailsConf = "peopletimedetails.json"
	defaultConfDir        = "./conf"
)

var (
	unanetUpload = flag.NewFlagSet("upload", flag.ExitOnError)

	uploadFile = unanetUpload.String("file", "", "Unanet report containing time entries to upload. (Required)")
	confDir    = unanetUpload.String("conf", defaultConfDir, "Directory containing config files.")
)

func init() {

	unanetUpload.StringVar(uploadFile, "f", "", "Unanet report containing time entries to upload. (shorthand) (Required)")
	unanetUpload.StringVar(confDir, "c", defaultConfDir, "Directory containing config files. (shorthand)")
}

func displayUsage(msg string) {

	if msg != "" {
		fmt.Println(msg)
	}
	fmt.Println("Usage: teamwork-unanet upload <options>")
	unanetUpload.PrintDefaults()
}

func main() {	

	if len(os.Args) < 2 {
		displayUsage("")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "upload":
		unanetUpload.Parse(os.Args[2:])

	default:
		displayUsage("No valid command found.  Valid commands include: upload")
		os.Exit(1)
	}

	errBuff := ""

	if unanetUpload.Parsed() {

		if *uploadFile == "" {
			if errBuff != "" {
				errBuff += "; "
			}
			errBuff += "--file (-f)"
		}

		if errBuff != "" {
			displayUsage("Missing required options: " + errBuff)
			os.Exit(1)
		}
		
		// create new Teamwork API connection from json config file
		conn, err := teamworkapi.NewConnectionFromJSON(fmt.Sprintf("%s/%s", *confDir, apiConf))
		if err != nil {
			fmt.Printf("[ERROR]\tfailed to initialize API connection: %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Println("[OK]\tsuccessfully established API connection")

		// create new Unanet TimeDetailsReport
		rep, err := unanet.NewTimeDetailsReport(conn)
		if err != nil {
			fmt.Printf("[ERROR]\tfailed to create new report: %s\n", err.Error())
			os.Exit(1)
		}
		
		// load Unanet TimeDetailsReport configuration
		err = rep.Report.LoadConfig(fmt.Sprintf("%s/%s", *confDir, peopleTimeDetailsConf))
		if err != nil {
			fmt.Printf("[ERROR]\tfailed to load report configuration: %s\n", err.Error())
			os.Exit(1)
		}

		fmt.Println("[OK]\tsuccessfully configured time details report")

		// parse Unanet time details report to create Teamwork time entries
		entries, err := rep.ParseTimeDetailsReport(*uploadFile)
		if err != nil {
			fmt.Printf("[ERROR]\tfailed to parse time details report (%s): %s\n", *uploadFile, err)
			os.Exit(1)
		}

		fmt.Printf("[OK]\tfound %d time entries in %s\n", len(entries), *uploadFile)
		
		// upload time entries to Teamwork and save Unanet time details report
		// with Teamwork Time Entry IDs
		err = rep.UploadTimeEntries()
		if err != nil {
			fmt.Printf("[ERROR]\tfailed to upload time entries: %s", err)
			os.Exit(1)
		}

		fmt.Printf("[OK]\tsuccessfully uploaded %d time entries to Teamwork\n", len(entries))
		fmt.Printf("[OK]\tsaved Unanet report %s", rep.Report.Filename)
	}
}
