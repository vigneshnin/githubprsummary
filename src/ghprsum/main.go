package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/go-github/v53/github"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var (
	TimeNow         = time.Now()
	LastWeekTime    = time.Now().AddDate(0, 0, -7)
	Repository      = os.Getenv("GITHUB_REPOSITORY")
	RepositoryOwner = os.Getenv("GITHUB_REPOSITORY_OWNER")
)

func handler(ctx context.Context, event events.CloudWatchEvent) {
	log.Println("Starting...")
	// github authentication - none in this case
	client := github.NewClient(nil)

	if RepositoryOwner == "" {
		log.Fatalln("GITHUB_REPOSITORY_OWNER not set")
	}
	if Repository == "" {
		log.Fatalln("GITHUB_REPOSITORY not set")
	}
	log.Println("Getting PRs from the respository: ", RepositoryOwner+"/"+Repository)

	//list PRs of last week from the repository
	prs, _, err := client.PullRequests.List(context.Background(), RepositoryOwner, Repository, &github.PullRequestListOptions{
		State:     "all",
		Direction: "desc",
		Sort:      "updated",
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    1,
		},
	})
	if err != nil {
		log.Fatalf("unable to list PRs, %v", err)
	}

	//log.Println(prs)

	UpdatedPrSlice := []github.PullRequest{}

	for _, pr := range prs {
		if pr.UpdatedAt.Compare(LastWeekTime) == 1 || pr.CreatedAt.Compare(LastWeekTime) == 1 {
			UpdatedPrSlice = append(UpdatedPrSlice, *pr)
			log.Println(pr.UpdatedAt.Format("2006-01-02"))
		}
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Using the Config value, create the SNS client
	svc := sns.NewFromConfig(cfg)

	// Parse PR response for last week and create notification body
	var NotificationBody string

	if UpdatedPrSlice == nil {
		log.Println("No PRs updated or created in the last week")
		NotificationBody = "No PRs updated or created in the last week"
	} else {
		NotificationBody += "PRs updated or created in the last week" + "\n" + "URL: " + "https://github.com/" + RepositoryOwner + "/" + Repository + "/pulls" + "\n\n"
		for _, updatedprs := range UpdatedPrSlice {
			NotificationBody = NotificationBody + "PR: " + *updatedprs.Title + " | " + *updatedprs.HTMLURL + "\n" +
				"--------------------------" + "\n" +
				"Summary: " + *updatedprs.Body + "\n" +
				"#######################" + "\n\n"
		}
	}

	// Build the request with its input parameters
	resp, err := svc.Publish(context.TODO(), &sns.PublishInput{
		Message:  aws.String(NotificationBody),
		Subject:  aws.String("Github PR Weekly Summary | " + RepositoryOwner + "/" + Repository + " | " + TimeNow.Format("2006-01-02")),
		TopicArn: aws.String(os.Getenv("SNS_ARN")),
	})
	if err != nil {
		log.Fatalf("unable to publish, %v", err)
	}
	log.Println(resp)
}

func main() {
	lambda.Start(handler)
}
