package imds

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
)

const (
	IPv4Mode        = "ipv4"
	IPv6Mode        = "ipv6"
	spotITNPath     = "spot/termination-time"
	scheduledEvents = "events/maintenance/scheduled"
)

type IMDS struct {
	client *imds.Client
}

type ScheduledEventDetail struct {
	NotBefore   string `json:"NotBefore"`
	Code        string `json:"Code"`
	Description string `json:"Description"`
	EventID     string `json:"EventId"`
	NotAfter    string `json:"NotAfter"`
	State       string `json:"State"`
}

type InstanceAction struct {
	Action string `json:"action"`
	Time   string `json:"time"`
}

type RebalanceRecommendation struct {
	NoticeTime string `json:"noticeTime"`
}

func NewClient(ctx context.Context, endpoint string, ipMode string) (*IMDS, error) {
	cfg, err := config.LoadDefaultConfig(ctx, withIMDSEndpoint(endpoint), withIPMode(ipMode))
	if err != nil {
		return nil, err
	}
	return &IMDS{
		client: imds.NewFromConfig(cfg),
	}, nil
}

func withIMDSEndpoint(imdsEndpoint string) func(*config.LoadOptions) error {
	return func(lo *config.LoadOptions) error {
		lo.EC2IMDSEndpoint = imdsEndpoint
		return nil
	}
}

func withIPMode(ipMode string) func(*config.LoadOptions) error {
	return func(lo *config.LoadOptions) error {
		if ipMode == IPv6Mode {
			lo.EC2IMDSEndpointMode = imds.EndpointModeStateIPv6
		} else if ipMode == IPv4Mode {
			lo.EC2IMDSEndpointMode = imds.EndpointModeStateIPv4
		} else {
			return fmt.Errorf("invalid IMDS IP Mode \"%s\"", ipMode)
		}
		return nil
	}
}

// TODO: use spot/instance-action instead
func (i IMDS) GetSpotInterruptionNotification(ctx context.Context) (*time.Time, bool, error) {
	output, err := i.client.GetMetadata(ctx, &imds.GetMetadataInput{Path: spotITNPath})
	if err != nil {
		return nil, false, fmt.Errorf("IMDS Failed to get \"%s\": %w", spotITNPath, err)
	}
	termTimeBytes := new(bytes.Buffer)
	termTimeBytes.ReadFrom(output.Content)
	termTime, err := time.Parse("2006-01-02T15:04:05Z", termTimeBytes.String())
	if err != nil {
		return nil, true, fmt.Errorf("invalid time received from \"%s\": %w", spotITNPath, err)
	}
	return &termTime, true, nil
}

//TODO: Make this work
// func (i IMDS) GetMaintenanceEvent(ctx context.Context) (bool, error) {
// 	output, err := i.client.GetMetadata(ctx, &imds.GetMetadataInput{Path: scheduledEvents})
// 	if err != nil {
// 		return false, fmt.Errorf("IMDS Failed to get \"%s\": %w", scheduledEvents, err)
// 	}
// 	return true, nil
// }
