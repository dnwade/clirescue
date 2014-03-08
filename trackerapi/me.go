package trackerapi

import (
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
        u "os/user"

        "github.com/dnwade/clirescue/cmdutil"
        "github.com/dnwade/clirescue/user"
)

var (
        URL          string     = "https://www.pivotaltracker.com/services/v5/me"
        FileLocation string     = homeDir() + "/.tracker"
        currentUser  *user.User = user.New()
        Stdout       *os.File   = os.Stdout
)

func getAPIToken() bool {
        data, err := ioutil.ReadFile(FileLocation)
        fmt.Printf("tracker not exists? %v  len of file data: %v\n", err, len(data))
        if len(data) > 0 {
                currentUser.APIToken = string(data)
                fmt.Println(currentUser.APIToken)
        }
        return err == nil && len(data) > 0
}

func Me() {
        if getAPIToken() {
        } else {
                setCredentials()
        }
        parse(makeRequest())
        ioutil.WriteFile(FileLocation, []byte(currentUser.APIToken), 0644)
}

func makeRequest() []byte {
        client := &http.Client{}
        req, err := http.NewRequest("GET", URL, nil)
        if currentUser.APIToken == "" {
                req.SetBasicAuth(currentUser.Username, currentUser.Password)
        } else {
                req.Header.Set("X-TrackerToken", currentUser.APIToken)
        }
        resp, err := client.Do(req)
        // fmt.Println("%#v", resp)
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
                fmt.Print(err)
        }
        fmt.Printf("\n****\nAPI response: \n%s\n", string(body))
        return body
}

func parse(body []byte) {
        var meResp = new(MeResponse)
        err := json.Unmarshal(body, &meResp)
        if err != nil {
                fmt.Println("error:", err)
        }

        currentUser.APIToken = meResp.APIToken
}

// pull in data from ~/.tracker
func setCredentials() {
        fmt.Fprint(Stdout, "Username: ")
        var username = cmdutil.ReadLine()
        cmdutil.Silence()
        fmt.Fprint(Stdout, "Password: ")

        var password = cmdutil.ReadLine()
        currentUser.Login(username, password)
        cmdutil.Unsilence()
}

func homeDir() string {
        usr, _ := u.Current()
        return usr.HomeDir
}

type MeResponse struct {
        APIToken string  `json:"api_token"`
        Username string  `json:"username"`
        Name     string  `json:"name"`
        Email    string  `json:"email"`
        Initials string  `json:"initials"`
        Timezone struct {
                Kind      string  `json:"kind"`
                Offset    string  `json:"offset"`
                OlsonName string  `json:"olson_name"`
        }       `json:"time_zone"`
}
