package sc

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type ABIResponse struct {
	Message *string `json:"message"`
	Result  *struct {
		Output struct {
			Abi []struct {
				Inputs []struct {
					InternalType string `json:"internalType"`
					Name         string `json:"name"`
					Type         string `json:"type"`
				} `json:"inputs,omitempty"`
				Payable         bool   `json:"payable"`
				StateMutability string `json:"stateMutability"`
				Type            string `json:"type"`
				Constant        bool   `json:"constant,omitempty"`
				Name            string `json:"name,omitempty"`
				Outputs         []struct {
					InternalType string `json:"internalType"`
					Name         string `json:"name"`
					Type         string `json:"type"`
				} `json:"outputs,omitempty"`
			} `json:"abi"`
			Devdoc struct {
				Methods map[string]struct {
					Details string `json:"details"`
				} `json:"methods"`
			} `json:"devdoc"`
		} `json:"output"`
	} `json:"result"`
	Status *bool `json:"status"`
}

// Import resty into your code and refer it as `resty`.

func GetContractABI(address string) (*ABIResponse, error) {
	client := resty.New()
	abi_resp := &ABIResponse{}
	resp, err := client.R().SetResult(abi_resp).
		Get(fmt.Sprintf("https://explorer-kintsugi.roninchain.com/v2/2020/contract/%s/abi", address))

	if err != nil {
		return nil, err
	}
	if !resp.IsSuccess() {
		return nil, fmt.Errorf(fmt.Sprintf("Get contract raw abi failed: %s\n", resp))
	}
	if !*abi_resp.Status {
		return nil, fmt.Errorf("get contract raw abi failed")
	}

	return abi_resp, nil
}
