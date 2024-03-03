package utility_test

import (
	"strconv"
	"testing"

	"github.com/heaven-chp/common-library-go/utility"
)

func TestWhetherCidrContainsIp(t *testing.T) {
	if contain, err := utility.WhetherCidrContainsIp("192.168.1.0/24", "192.168.1.111"); err != nil {
		t.Fatal(err)
	} else if contain == false {
		t.Fatal("invalid")
	}

	if contain, err := utility.WhetherCidrContainsIp("192.168.1.0/24", "192.168.2.0"); err != nil {
		t.Fatal(err)
	} else if contain {
		t.Fatal("invalid")
	}
}

func TestGetAllIpsOfCidr(t *testing.T) {
	if ips, err := utility.GetAllIpsOfCidr("192.168.1.0/24"); err != nil {
		t.Fatal(err)
	} else {
		if len(ips) != 254 {
			t.Fatal("invalid -", len(ips))
		}

		for index, ip := range ips {
			answer := "192.168.1." + strconv.Itoa(index+1)
			if ip != answer {
				t.Fatal("invalid -", ip, ",", answer)
			}
		}
	}
}
