package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
)

func GetIP() (string, error) {
	// Create HTTP GET Request
	resp, err := http.Get("https://api.ipify.org/?format=text")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Print response body
	return string(body), nil
}

func UpdateRecord(svc *route53.Route53, zoneId string, recordName string, recordvalue string) error {
	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(recordName),
						Type: aws.String("A"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(recordvalue),
							},
						},
						TTL: aws.Int64(1), // Set TTL to 1 second to the value is never cached
					},
				},
			},
		},
		HostedZoneId: aws.String(zoneId),
	}
	_, err := svc.ChangeResourceRecordSets(params)
	if err != nil {
		return err
	}

	return nil
}
func GetZoneID(svc *route53.Route53, domainName string) (string, error) {
	params := &route53.ListHostedZonesInput{
		MaxItems: aws.String("100"),
	}
	resp, err := svc.ListHostedZones(params)
	if err != nil {
		return "", err
	}

	for _, zone := range resp.HostedZones {
		if *zone.Name == domainName+"." {
			return *zone.Id, nil
		}
		fmt.Println(*zone.Name)
	}
	return "", fmt.Errorf("zone not found")
}
