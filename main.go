package main

import (
	"fmt"

	"demees.local/r53update/utils"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

const (
	DOMAIN     = "demees.dev"
	RECORDNAME = "compute.demees.dev"
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
	zoneId, err := utils.GetZoneID(svc, DOMAIN)
	if err != nil {
		fmt.Println(err)
		return
	}
	// update DNS record
	err = utils.UpdateRecord(svc, zoneId, RECORDNAME, ip)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Updated DNS record for " + RECORDNAME + " (" + zoneId + ") -> " + ip)
}
