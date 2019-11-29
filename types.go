package shell

import "github.com/soracom/soracom-cli/generators/lib"

// SoracomCompleter returns suggestions for given go-prompt document.
type SoracomCompleter struct {
	// SORACOM CLI API definitions
	apiDef                *lib.APIDefinitions
	specifiedProfileName  string
	specifiedCoverageType string
	providedAPIKey        string
	providedAPIToken      string
}

// SoracomExecutor executes given string with the shell.
type SoracomExecutor struct {
	// shell which executes a command
	shell                 string
	specifiedProfileName  string
	specifiedCoverageType string
	providedAPIKey        string
	providedAPIToken      string
}

type flag struct {
	name  string
	value string
}

type param struct {
	name        string
	required    bool
	description string
	paramType   string
	enum        []string
}

var subscribers []struct {
	Apn              string      `json:"apn"`
	CreatedAt        int64       `json:"createdAt"`
	CreatedTime      int64       `json:"createdTime"`
	ExpiredAt        interface{} `json:"expiredAt"`
	ExpiryAction     interface{} `json:"expiryAction"`
	ExpiryTime       interface{} `json:"expiryTime"`
	GroupID          string      `json:"groupId"`
	Iccid            string      `json:"iccid"`
	ImeiLock         interface{} `json:"imeiLock"`
	Imsi             string      `json:"imsi"`
	IPAddress        string      `json:"ipAddress"`
	LastModifiedAt   int64       `json:"lastModifiedAt"`
	LastModifiedTime int64       `json:"lastModifiedTime"`
	ModuleType       string      `json:"moduleType"`
	Msisdn           string      `json:"msisdn"`
	OperatorID       string      `json:"operatorId"`
	Plan             int         `json:"plan"`
	RegisteredTime   int64       `json:"registeredTime"`
	SerialNumber     string      `json:"serialNumber"`
	SessionStatus    struct {
		Cell struct {
			Eci       int    `json:"eci"`
			Mcc       int    `json:"mcc"`
			Mnc       int    `json:"mnc"`
			RadioType string `json:"radioType"`
			Tac       int    `json:"tac"`
		} `json:"cell"`
		DNSServers    []string    `json:"dnsServers"`
		Imei          string      `json:"imei"`
		LastUpdatedAt int64       `json:"lastUpdatedAt"`
		Location      interface{} `json:"location"`
		Online        bool        `json:"online"`
		UeIPAddress   string      `json:"ueIpAddress"`
	} `json:"sessionStatus"`
	SpeedClass   string `json:"speedClass"`
	Status       string `json:"status"`
	Subscription string `json:"subscription"`
	Tags         struct {
		Name string `json:"name"`
	} `json:"tags"`
	TerminationEnabled bool   `json:"terminationEnabled"`
	Type               string `json:"type"`
}

var devices []struct {
	DeviceId     string `json:"deviceId"`
	Endpoint     string `json:"endpoint"`
	Imei         string `json:"imei"`
	Imsi         string `json:"imsi"`
	IPAddress    string `json:"ipAddress"`
	Manufacturer string `json:"manufacturer"`
	Online       bool   `json:"online"`
	Status       string `json:"status"`
}
