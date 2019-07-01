package main

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

const (
	deviceEnd     = "/api/devices"
	activationEnd = "/api/devices/%s/activate"
	keysEnd       = "/api/devices/%s/keys"
	loginEnd      = "/api/internal/login"
)

type deviceRequest struct {
	Device device `json:"device"`
}

type device struct {
	DevEUI            string  `csv:"dev_eui" json:"devEUI"`
	Name              string  `csv:"name" json:"name"`
	ApplicationID     int64   `csv:"application_id" json:"applicationID"`
	Description       string  `csv:"description" json:"description"`
	DeviceProfileID   string  `csv:"device_profile_id" json:"deviceProfileID"`
	SkipFCntCheck     bool    `csv:"skip_f_cnt_check" json:"skip_f_cnt_check"`
	ReferenceAltitude float64 `csv:"reference_altitude" json:"reference_altitude"`
	DevAddr           string  `csv:"dev_addr" json:"devAddr"`
	NwkKey            string  `csv:"nwk_key" json:"nwk_key"`
	AppKey            string  `csv:"app_key" json:"app_key"`
	GenAppKey         string  `csv:"gen_app_key" json:"gen_app_key"`
	AppSKey           string  `csv:"app_s_key" json:"app_s_key"`
	FNwkSIntKey       string  `csv:"f_nwk_s_int_key" json:"f_nwk_s_int_key"`
	SNwkSIntKey       string  `csv:"s_nwk_s_int_key" json:"s_nwk_s_int_key"`
	NwkSEncKey        string  `csv:"nwk_s_enc_key" json:"nwk_s_enc_key"`
	Activation        string  `csv:"activation"`
}

type deviceKeysRequest struct {
	Keys deviceKeys `json:"device_keys"`
}

type deviceKeys struct {
	DevEUI    string `json:"devEUI"`
	NwkKey    string `json:"nwk_key"`
	AppKey    string `json:"app_key"`
	GenAppKey string `json:"gen_app_key"`
}

type deviceActivationRequest struct {
	Activation deviceActivation `json:"device_activation"`
}

type deviceActivation struct {
	DevEUI      string `csv:"dev_eui" json:"devEUI"`
	DevAddr     string `csv:"dev_addr" json:"devAddr"`
	AppSKey     string `csv:"app_s_key" json:"app_s_key"`
	FNwkSIntKey string `csv:"f_nwk_s_int_key" json:"f_nwk_s_int_key"`
	SNwkSIntKey string `csv:"s_nwk_s_int_key" json:"s_nwk_s_int_key"`
	NwkSEncKey  string `csv:"nwk_s_enc_key" json:"nwk_s_enc_key"`
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"jwt"`
}

// Load takes a filename and tries to retrieve devices from it.
func Load(filename, hostname, username, password string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	devices := []*device{}

	if err := gocsv.UnmarshalFile(file, &devices); err != nil { // Load clients from file
		return err
	}

	fmt.Printf("got %d devices from file\n", len(devices))

	//Now try to login to get a jwt token.
	req := loginRequest{
		Username: username,
		Password: password,
	}

	var lr loginResponse
	var rErr interface{}

	resp, err := resty.R().
		SetBody(req).
		SetResult(&lr).
		SetError(&rErr).
		Post(hostname + loginEnd)

	if err != nil {
		return errors.Wrap(err, "couldn't login")
	} else if rErr != nil {
		fmt.Printf("unexpected error: %+v\n", rErr)
		return errors.New("error at login")
	} else if resp.StatusCode() != 200 {
		return errors.New("got status code different from 200")
	}

	Provision(devices, hostname, fmt.Sprintf("Bearer %s", lr.Token))

	return nil
}

// Provision loops through devices to provision them.
func Provision(devices []*device, hostname, jwtString string) {
	var result interface{}
	var rErr interface{}
	for i := 0; i < len(devices); i++ {
		device := *devices[i]
		dr := deviceRequest{
			Device: device,
		}
		resp, err := resty.R().
			SetBody(dr).
			SetHeader("Authorization", jwtString).
			SetResult(&result).
			SetError(&rErr).
			Post(hostname + deviceEnd)
		if err != nil {
			fmt.Printf("couldn't provision device %s: %s\n", devices[i].DevEUI, err)
		} else if rErr != nil {
			fmt.Printf("couldn't provision device %s: %+v\n", devices[i].DevEUI, rErr)
		} else if resp.StatusCode() != 200 {
			fmt.Printf("couldn't provision device %s: error code %d\n", devices[i].DevEUI, resp.StatusCode())
		} else {
			if device.Activation == "ABP" {
				da := deviceActivation{
					DevEUI:      device.DevEUI,
					DevAddr:     device.DevAddr,
					AppSKey:     device.AppSKey,
					FNwkSIntKey: device.FNwkSIntKey,
					SNwkSIntKey: device.SNwkSIntKey,
					NwkSEncKey:  device.NwkSEncKey,
				}
				dar := deviceActivationRequest{
					Activation: da,
				}
				resp, err = resty.R().
					SetBody(dar).
					SetHeader("Authorization", jwtString).
					SetResult(&result).
					SetError(&rErr).
					Post(hostname + fmt.Sprintf(activationEnd, devices[i].DevEUI))
				if err != nil {
					fmt.Printf("couldn't activate device %s: %s\n", devices[i].DevEUI, err)
				} else if rErr != nil {
					fmt.Printf("couldn't activate device %s: %+v\n", devices[i].DevEUI, rErr)
				} else if resp.StatusCode() != 200 {
					fmt.Printf("couldn't activate device %s: error code %d\n", device.DevEUI, resp.StatusCode())
				}
			} else {
				keys := deviceKeys{
					DevEUI:    device.DevEUI,
					NwkKey:    device.NwkKey,
					AppKey:    device.AppKey,
					GenAppKey: device.GenAppKey,
				}

				dkr := deviceKeysRequest{
					Keys: keys,
				}

				resp, err = resty.R().
					SetBody(dkr).
					SetHeader("Authorization", jwtString).
					SetResult(&result).
					SetError(&rErr).
					Post(hostname + fmt.Sprintf(keysEnd, device.DevEUI))
				if err != nil {
					fmt.Printf("couldn't set keys for device %s: %s\n", devices[i].DevEUI, err)
				} else if rErr != nil {
					fmt.Printf("couldn't set keys for device %s: %+v\n", devices[i].DevEUI, rErr)
				} else if resp.StatusCode() != 200 {
					fmt.Printf("couldn't set keys for device %s: error code %d\n", devices[i].DevEUI, resp.StatusCode())
				}
			}
		}
	}
}
