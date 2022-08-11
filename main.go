package main

import (
	"fmt"

	"demees.local/r53update/utils"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"os"
	"bufio"
	"strings"
)

const (
	// AWS region
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

	// Read credentials from AWS credentials file
	file, err := os.Open("/home/ec2-user/.aws/credentials")
	
	if err != nil {
		fmt.Println(err)
		return
	}

	scanner := bufio.NewScanner(file)
	var match = false
	var access_key_id = ""
	var secret_access_key = ""

	for scanner.Scan(){
		line := scanner.Text()
		
		if(line == "[default]"){
			match=true
		}
		if(match == true && strings.HasPrefix(line, "aws_access_key_id=")){
			access_key_id = strings.Split(line, "=")[1]
		}

		if(match == true && strings.HasPrefix(line, "aws_secret_access_key=")){
			secret_access_key = strings.Split(line, "=")[1]	
			match=false // assume secret access key is the last in the config block - BAD CODE
		}
		
	}
	

	// get AWS session - we are using the NewStaticCredentials on purpose 
	session, err := session.NewSession(
		&aws.Config{
			Credentials: credentials.NewStaticCredentials(access_key_id,secret_access_key,""),
		})

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
