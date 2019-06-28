package main

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"gopkg.in/resty.v1"
)

const (
	deviceEnd = "ap/devices/%s"
	keysEnd   = "/api/devices/%s/activate"
)

// Hostname refers to the lora-app-server address.
var Hostname string

type device struct {
	DevEUI            string  `csv:"dev_eui" json:"devEUI"`
	Name              string  `csv:"name" json:"name"`
	ApplicationID     int64   `csv:"application_id" json:"applicationID"`
	Description       string  `csv:"description" json:"description"`
	DeviceProfileID   int64   `csv:"device_profile_id" json:"deviceProfileID"`
	SkipFCntCheck     bool    `csv:"skip_f_cnt_check" json:"skip_f_cnt_check"`
	ReferenceAltitude float64 `csv:"reference_altitude" json:"reference_altitude"`
	DevAddr           string  `csv:"dev_addr" json:"devAddr"`
	JoinEUI           string  `csv:"join_eui" json:"joinEUI"`
	FNwkSIntKey       string  `csv:"f_nwk_s_int_key" json:"f_nwk_s_int_key"`
	SNwkSIntKey       string  `csv:"s_nwk_s_int_key" json:"s_nwk_s_int_key"`
	NwkSEncKey        string  `csv:"nwk_s_enc_key" json:"nwk_s_enc_key"`
	Activation        string  `csv:"activation"`
}

// Load takes a filename and tries to retrieve devices from it.
func Load(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	devices := []*device{}

	if err := gocsv.UnmarshalFile(file, &devices); err != nil { // Load clients from file
		return err
	}

	Provision(devices)

	return nil
}

// Provision loops through devices to provision them.
func Provision(devices []*device) {
	var result interface{}
	var err error
	for i := 0; i < len(devices); i++ {
		resp, err := resty.R().
			SetBody(*devices[i]).
			SetResult(&result).
			SetError(&err).
			Post(Hostname + fmt.Sprintf(deviceEnd, devices[i].DevEUI))
		if err != nil {
			fmt.Printf("couldn't provision device %s: %s\n", devices[i].DevEUI, err)
		} else if resp.StatusCode() != 200 {
			fmt.Printf("couldn't provision device %s: error code %d\n", devices[i].DevEUI, resp.StatusCode())
		} else {
			if devices[i].Activation == "ABP" {
				resp, err = resty.R().
					SetBody(*devices[i]).
					SetResult(&result).
					SetError(&err).
					Post(Hostname + fmt.Sprintf(keysEnd, devices[i].DevEUI))
				if err != nil {
					fmt.Printf("couldn't provision device %s: %s\n", devices[i].DevEUI, err)
				} else if resp.StatusCode() != 200 {
					fmt.Printf("couldn't provision device %s: error code %d\n", devices[i].DevEUI, resp.StatusCode())
				}
			}
		}
	}
}
