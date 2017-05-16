package proxmox

import (
	"github.com/hashicorp/terraform/helper/schema"
	"net/http"
	"crypto/tls"
	"encoding/json"

	"io/ioutil"
	"net/url"
)

type ProxMoxAuthData struct {
	Data struct {
		CSRFPreventionToken string `json:"CSRFPreventionToken"`
		Ticket string `json:"ticket"`
		Username string `json:"username"`
		Cap struct {
			Storage struct { } `json:"storage"`
			Dc struct { } `json:"dc"`
			Vms struct { } `json:"vms"`
			Nodes struct { } `json:"nodes"`
			Access struct { } `json:"access"` } `json:"cap"` } `json:"data"`
	Errors struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"errors"`

}

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"node": &schema.Schema{
				Type:	schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	println("Test 0")
	address := d.Get("address").(string)
	if !ping(address) {
		auth(address)
	}


	auth(address)
	//d.SetId(address)
	return nil
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	address := d.Get("address").(string)
	auth(address)



	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func ping(address string) bool {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(address + "/api2/json/version")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func auth(address string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	user := url.QueryEscape("khazrak@pve")
	pass := url.QueryEscape("Carnagea1987!")

	resp, err := client.Post(address + "/api2/json/access/ticket?username=" + user +"&password=" + pass, "text/plain", nil)

	responseData,err := ioutil.ReadAll(resp.Body)
	responseString := string(responseData)
	println(responseString)
	println(resp.StatusCode)

	res := ProxMoxAuthData{}
	json.Unmarshal([]byte(responseString), &res)

	if resp.StatusCode < 400 {
		println(res.Data.Ticket)

	} else {
		println("Error!")
		println("Username: " + res.Errors.Username)
		println("Password: " + res.Errors.Password)
		if err != nil {
			panic(err)
		}
	}

	defer resp.Body.Close()
}


