[![Go Report Card](https://goreportcard.com/badge/github.com/Foxtrot-Division/teamwork-unanet)](https://goreportcard.com/report/github.com/Foxtrot-Division/teamwork-unanet)
# teamwork-unanet
Command-line utility to upload Unanet time logs to Teamwork Projects.
Currently, this app supports parsing of the Unanet People Time Details report to
generate Teamwork time entry for each Unanet time record in the report.  Two
configuration files are required:
- conf/api.json
- conf/peopletimedetails.json

## Configuring Teamwork API Connection
In the *conf/api.json* file, edit the "apiKey" and "siteName" appropriately:
```
{
    "apiKey" : "<Teamwork API key>",
    "siteName" : "<Teamwork Sitename>",
    "dataPreference" : "json"
}
```

## Configuring the Unanet People Time Details Report
You'll need to provide some mapping between the Unanet data and your Teamwork
data.  In the *conf/peopletimedetails.json* file, edit the "companyMappings",
"projectMappings", "taskMappings", and "userMappings" accordingly:
```
{
    "reportName": "PeopleTimeDetails",
    "columns": [
        "PersonOrganization",
        "Person",
        "ProjectOrganization",
        "ProjectCode",
        "TaskNumber",
        "Task",
        "LaborCategory",
        "Location",
        "ProjectType",
        "PayCode",
        "Reference",
        "Date",
        "ADJPostedDate",
        "FinancialPostedDate",
        "Hours"
    ],
    "companyMappings": {
        "<Unanet Company Name>": "<Teamwork Company ID>"
    },
    "projectMappings": {
        "<Unanet project ID>": "<Teamwork Project ID>"
    },
    "taskMappings": { 
            "<Unanet Task ID>":"<Teamwork Task ID>"
    },
    "userMappings": {
        "<Unanet Person Name 1>": "<Teamwork Person ID>",
        "<Unanet Person Name 2>": "<Teamwork Person ID>",
        "<Unanet Person Name n>": "<Teamwork Person ID>"
    }
}
```
## Downloading the Unanet People Time Details Report
This app works with Unanet v20.6.3.  To download a Unanet People Time Details report, navigate
to Reports > Dashboard.  In the **Detail Reports** pane, click on *Time Details*
link.  Once you've selected the criteria you want, click the *Run Report* link
and review the results.  An example of a properly-formed Unanet People Time Details report is available
in *reports/example-report.csv*. 

Some things to verify:
- make sure all users in the report are appropriately mapped to a Teamwork
  Person ID in *conf/peopletimedetails.json*
- make sure all task numbers are appropriately mapped to a Teamwork Task ID in
  *conf/peopletimedetails.json*


## Uploading Time Entries to Teamwork
First, make sure you build the teamwork-unanet source code to create the executable.  Navigate to the
teamwork-unanet directory and run:
```
go build 
```
Once you've downloaded the Unanet People Time Details report, run the following:

```
./teamwork-unanet upload -c ./conf -f <path to Unanet report>
```
The -c option specifies the location of the configuration files (*api.json* and
*peopletimedetails.json*).

The -f option specifies the location of the Unanet People Time Details report
that contains the Unanet time logs.

Example:
```
./teamwork-unanet upload -c ./conf -f ./reports_private/report.csv
```
If successfull, you should see the following output:
```
[OK]    successfully established API connection
[OK]    successfully configured time details report
[OK]    found 34 time entries in ./reports_private/report.csv
[OK]    successfully uploaded 34 time entries to Teamwork
[OK]    saved Unanet report time-details-report_20201201-20201215.csv
```
A new file is created to include all the original Unanet records, along with a
new column **TeamworkTimeEntryID** that maps the Unanet record to its
corresponding Teamwork Time Entry ID.  If an error occurred with a specific
Unanet record, the error will be recorded in the **TeamworkTimeEntryID** column.
