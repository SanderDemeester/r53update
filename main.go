package main

import (
	"fmt"

	"demees.local/r53update/utils"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func main() {
	// get public IP address
	ip, err := utils.GetIP()
	if err != nil {
		fmt.Println(err)
		return
	}

	// get AWS session
	session, err := session.NewSession()
	if err != nil {
		fmt.Println(err)
		return
	}

	// get AWS service client
	svc := route53.New(session)

	// get hosted zone ID
	zoneId, err := utils.GetZoneID(svc, "demees.dev")
	if err != nil {
		fmt.Println(err)
		return
	}
	// update DNS record
	err = utils.UpdateRecord(svc, zoneId, "compute.demees.dev", ip)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Updated DNS record for compute.demees.dev (" + zoneId + ") -> " + ip)
}
