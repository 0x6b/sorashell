package sorashell

import "github.com/soracom/soracom-cli/generators/lib"

// SoracomCompleter returns suggestions for given go-prompt document.
type SoracomCompleter struct {
	// SORACOM CLI API definitions
	apiDef *lib.APIDefinitions
	// worker which executes a command for device or imsi completion
	worker *SoracomWorker
}

// SoracomExecutor executes given string with the shell.
type SoracomExecutor struct {
	// worker which executes a command
	worker *SoracomWorker
}

// SoracomWorker executes given string with the shell.
type SoracomWorker struct {
	// shell which executes a command (NOT_USED)
	shell        string
	profileName  string
	coverageType string
	apiKey       string
	apiToken     string
	command      chan string
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

var inventoryDevices []struct {
	DeviceID     string `json:"deviceId"`
	Endpoint     string `json:"endpoint"`
	Imei         string `json:"imei"`
	Imsi         string `json:"imsi"`
	IPAddress    string `json:"ipAddress"`
	Manufacturer string `json:"manufacturer"`
	Online       bool   `json:"online"`
	Status       string `json:"status"`
}

var sigfoxDevices []struct {
	DeviceID string `json:"deviceId"`
	Status   string `json:"status"`
	Tags     struct {
		Name string `json:"name"`
	} `json:"tags"`
}

var orders struct {
	OrderList []struct {
		Currency      string `json:"currency"`
		Email         string `json:"email"`
		OrderDateTime string `json:"orderDateTime"`
		OrderID       string `json:"orderId"`
		OrderItemList []struct {
			Product struct {
				Count            int           `json:"count"`
				Currency         string        `json:"currency"`
				Price            int           `json:"price"`
				ProductCode      string        `json:"productCode"`
				ProductImageURLs []interface{} `json:"productImageURLs"`
				ProductName      string        `json:"productName"`
				ProductType      string        `json:"productType"`
				Properties       struct {
					ContractType    string `json:"contractType"`
					SimSize         string `json:"simSize"`
					SimSubscription string `json:"simSubscription"`
				} `json:"properties"`
			} `json:"product"`
			ProductAmount int `json:"productAmount"`
			Quantity      int `json:"quantity"`
		} `json:"orderItemList"`
		OrderStatus     string `json:"orderStatus"`
		ShippingAddress struct {
			AddressLine1 string `json:"addressLine1"`
			AddressLine2 string `json:"addressLine2"`
			Building     string `json:"building"`
			City         string `json:"city"`
			CompanyName  string `json:"companyName"`
			CountryCode  string `json:"countryCode"`
			Department   string `json:"department"`
			FullName     string `json:"fullName"`
			PhoneNumber  string `json:"phoneNumber"`
			State        string `json:"state"`
			ZipCode      string `json:"zipCode"`
		} `json:"shippingAddress"`
		ShippingAddressID   string `json:"shippingAddressId"`
		ShippingCost        int    `json:"shippingCost"`
		ShippingLabelNumber string `json:"shippingLabelNumber"`
		TaxAmount           int    `json:"taxAmount"`
		TotalAmount         int    `json:"totalAmount"`
	} `json:"orderList"`
}

var groups []struct {
	Configuration    interface{} `json:"configuration"`
	CreatedAt        int64       `json:"createdAt"`
	CreatedTime      int64       `json:"createdTime"`
	GroupID          string      `json:"groupId"`
	LastModifiedAt   int64       `json:"lastModifiedAt"`
	LastModifiedTime int64       `json:"lastModifiedTime"`
	OperatorID       string      `json:"operatorId"`
	Tags             struct {
		Name string `json:"name"`
	} `json:"tags"`
}
