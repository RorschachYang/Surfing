package cdn_cf

import (
	"Surfing/util"
	"strconv"

	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/google/uuid"
)

func CreateDistribution(origin string) string {
	// 获取配置信息中的credentials
	credentialID := util.GetConfigString("Cloudfront.credential.id")
	credentialSecret := util.GetConfigString("Cloudfront.credential.secret")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("global"), // replace with your region
		Credentials: credentials.NewStaticCredentials(credentialID, credentialSecret, ""),
	})
	if err != nil {
		util.PrintLog(err.Error())
	}

	svc := cloudfront.New(sess)

	// create distribution config
	distributionConfig := &cloudfront.DistributionConfig{
		Comment:         aws.String(""),
		Enabled:         aws.Bool(true),
		CallerReference: aws.String(uuid.NewString()),
		Origins: &cloudfront.Origins{
			Quantity: aws.Int64(1),
			Items: []*cloudfront.Origin{
				{
					DomainName: aws.String(origin),
					Id:         aws.String(origin),
					CustomOriginConfig: &cloudfront.CustomOriginConfig{
						HTTPPort:             aws.Int64(80),
						HTTPSPort:            aws.Int64(443),
						OriginProtocolPolicy: aws.String("match-viewer"),
					},
				},
			},
		},
		DefaultCacheBehavior: &cloudfront.DefaultCacheBehavior{
			TargetOriginId:       aws.String(origin),
			ViewerProtocolPolicy: aws.String("allow-all"),
			Compress:             aws.Bool(false),
			AllowedMethods: &cloudfront.AllowedMethods{
				Items: []*string{
					aws.String("GET"),
					aws.String("HEAD"),
				},
				Quantity: aws.Int64(2),
			},
			ForwardedValues: &cloudfront.ForwardedValues{
				QueryString: aws.Bool(false),
				Cookies: &cloudfront.CookiePreference{
					Forward: aws.String("none"),
				},
			},
			MinTTL:     aws.Int64(0),
			DefaultTTL: aws.Int64(86400),
			MaxTTL:     aws.Int64(31536000),
		},
	}

	// create distribution input
	distributionInput := &cloudfront.CreateDistributionInput{
		DistributionConfig: distributionConfig,
	}

	// create distribution
	result, err := svc.CreateDistribution(distributionInput)
	if err != nil {
		util.PrintLog(err.Error())
	}

	distributionDomain := aws.StringValue(result.Distribution.DomainName)

	util.PrintLog("Created distribution: " + *result.Distribution.Id)

	return distributionDomain
}

func DeleteDistribution(distributionID string) {
	// 获取配置信息中的credentials
	credentialID := util.GetConfigString("Cloudfront.credential.id")
	credentialSecret := util.GetConfigString("Cloudfront.credential.secret")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("global"), // replace with your region
		Credentials: credentials.NewStaticCredentials(credentialID, credentialSecret, ""),
	})
	if err != nil {
		util.PrintLog(err.Error())
	}

	svc := cloudfront.New(sess)

	resp, err := svc.GetDistributionConfig(&cloudfront.GetDistributionConfigInput{
		Id: aws.String(distributionID),
	})
	if err != nil {
		util.PrintLog("Failed to get distribution config:" + err.Error())
	}

	_, err = svc.DeleteDistribution(&cloudfront.DeleteDistributionInput{
		Id:      aws.String(distributionID),
		IfMatch: resp.ETag,
	})
	if err != nil {
		util.PrintLog("Failed to delete distribution:" + err.Error())
	}

	util.PrintLog("Deleted distribution: " + distributionID)

}

func DisableDistribution(distributionID string) {
	// 获取配置信息中的credentials
	credentialID := util.GetConfigString("Cloudfront.credential.id")
	credentialSecret := util.GetConfigString("Cloudfront.credential.secret")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("global"), // replace with your region
		Credentials: credentials.NewStaticCredentials(credentialID, credentialSecret, ""),
	})
	if err != nil {
		util.PrintLog(err.Error())
	}

	svc := cloudfront.New(sess)

	resp, err := svc.GetDistributionConfig(&cloudfront.GetDistributionConfigInput{
		Id: aws.String(distributionID),
	})
	if err != nil {
		util.PrintLog("Failed to get distribution config:" + err.Error())
	}

	// Disable the distribution
	resp.DistributionConfig.Enabled = aws.Bool(false)

	// Update the configuration with the new state
	_, err = svc.UpdateDistribution(&cloudfront.UpdateDistributionInput{
		Id:                 aws.String(distributionID),
		IfMatch:            resp.ETag,
		DistributionConfig: resp.DistributionConfig,
	})
	if err != nil {
		util.PrintLog("Failed to update distribution:" + err.Error())
	}

	util.PrintLog("Distribution disabled:" + distributionID)

}

func GetDistributionsCreatedBefore(hours int64) []string {

	duration := int64(hours * -1 * int64(time.Hour.Abs()))

	startTime := time.Now().UTC().Add(time.Duration(duration))

	// 获取配置信息中的credentials
	credentialID := util.GetConfigString("Cloudfront.credential.id")
	credentialSecret := util.GetConfigString("Cloudfront.credential.secret")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("global"), // replace with your region
		Credentials: credentials.NewStaticCredentials(credentialID, credentialSecret, ""),
	})
	if err != nil {
		util.PrintLog(err.Error())
	}

	// Create a new CloudFront client
	svc := cloudfront.New(sess)

	// List all CloudFront distributions
	var marker *string
	var distributions []*cloudfront.DistributionSummary
	for {
		resp, err := svc.ListDistributions(&cloudfront.ListDistributionsInput{
			Marker: marker,
		})
		if err != nil {
			util.PrintLog("Failed to list distributions:" + err.Error())
		}

		distributions = append(distributions, resp.DistributionList.Items...)
		if *resp.DistributionList.IsTruncated {
			marker = resp.DistributionList.NextMarker
		} else {
			break
		}
	}

	var distributionIDs []string
	for _, dist := range distributions {
		if *dist.Status == "Deployed" && dist.LastModifiedTime.Before(startTime) {
			distributionIDs = append(distributionIDs, *dist.Id)
		}
	}

	for _, id := range distributionIDs {
		util.PrintLog("Distributions created more than " + strconv.FormatInt(hours, 10) + " hours ago:" + id)
	}

	return distributionIDs

}
