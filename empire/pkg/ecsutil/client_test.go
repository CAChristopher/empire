package ecsutil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/remind101/empire/empire/pkg/awsutil"
)

func TestListAppServices(t *testing.T) {
	h := awsutil.NewHandler([]awsutil.Cycle{
		awsutil.Cycle{
			Request: awsutil.Request{
				RequestURI: "/",
				Operation:  "AmazonEC2ContainerServiceV20141113.ListServices",
				Body:       `{"cluster":"cluster"}`,
			},
			Response: awsutil.Response{
				StatusCode: 200,
				Body:       `{"serviceArns":["arn:aws:ecs:us-east-1:249285743859:service/ae69bb4c-3903-4844-82fe-548ac5b74570--web"]}`,
			},
		},
	})
	m, s := newTestClient(h)
	defer s.Close()

	resp, err := m.ListAppServices(context.Background(), "ae69bb4c-3903-4844-82fe-548ac5b74570", &ecs.ListServicesInput{
		Cluster: aws.String("cluster"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if got := len(resp.ServiceARNs); got != 1 {
		t.Fatalf("Expected 1 service returned; got %d", got)
	}
}

func TestListAppServices_Pagination(t *testing.T) {
	h := awsutil.NewHandler([]awsutil.Cycle{
		awsutil.Cycle{
			Request: awsutil.Request{
				RequestURI: "/",
				Operation:  "AmazonEC2ContainerServiceV20141113.ListServices",
				Body:       `{"cluster":"cluster"}`,
			},
			Response: awsutil.Response{
				StatusCode: 200,
				Body:       `{"serviceArns":["arn:aws:ecs:us-east-1:249285743859:service/ae69bb4c-3903-4844-82fe-548ac5b74570--web"],"nextToken":"1234"}`,
			},
		},

		awsutil.Cycle{
			Request: awsutil.Request{
				RequestURI: "/",
				Operation:  "AmazonEC2ContainerServiceV20141113.ListServices",
				Body:       `{"cluster":"cluster","nextToken":"1234"}`,
			},
			Response: awsutil.Response{
				StatusCode: 200,
				Body:       `{"serviceArns":["arn:aws:ecs:us-east-1:249285743859:service/ae69bb4c-3903-4844-82fe-548ac5b74570--web"]}`,
			},
		},
	})
	m, s := newTestClient(h)
	defer s.Close()

	resp, err := m.ListAppServices(context.Background(), "ae69bb4c-3903-4844-82fe-548ac5b74570", &ecs.ListServicesInput{
		Cluster: aws.String("cluster"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if got := len(resp.ServiceARNs); got != 2 {
		t.Fatalf("Expected 2 services returned; got %d", got)
	}
}

func TestListAppTasks(t *testing.T) {
	h := awsutil.NewHandler([]awsutil.Cycle{
		awsutil.Cycle{
			Request: awsutil.Request{
				RequestURI: "/",
				Operation:  "AmazonEC2ContainerServiceV20141113.ListServices",
				Body:       `{"cluster":"cluster"}`,
			},
			Response: awsutil.Response{
				StatusCode: 200,
				Body:       `{"serviceArns":["arn:aws:ecs:us-east-1:249285743859:service/ae69bb4c-3903-4844-82fe-548ac5b74570--web"]}`,
			},
		},

		awsutil.Cycle{
			Request: awsutil.Request{
				RequestURI: "/",
				Operation:  "AmazonEC2ContainerServiceV20141113.ListTasks",
				Body:       `{"cluster":"cluster","serviceName":"ae69bb4c-3903-4844-82fe-548ac5b74570--web"}`,
			},
			Response: awsutil.Response{
				StatusCode: 200,
				Body:       `{"taskArns":["arn:aws:ecs:us-east-1:249285743859:task/ae69bb4c-3903-4844-82fe-548ac5b74570"]}`,
			},
		},
	})
	m, s := newTestClient(h)
	defer s.Close()

	resp, err := m.ListAppTasks(context.Background(), "ae69bb4c-3903-4844-82fe-548ac5b74570", &ecs.ListTasksInput{
		Cluster: aws.String("cluster"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if got := len(resp.TaskARNs); got != 1 {
		t.Fatalf("Expected 1 tasks returned; got %d", got)
	}
}

func TestListAppTasks_Paginate(t *testing.T) {
	h := awsutil.NewHandler([]awsutil.Cycle{
		awsutil.Cycle{
			Request: awsutil.Request{
				RequestURI: "/",
				Operation:  "AmazonEC2ContainerServiceV20141113.ListServices",
				Body:       `{"cluster":"cluster"}`,
			},
			Response: awsutil.Response{
				StatusCode: 200,
				Body:       `{"serviceArns":["arn:aws:ecs:us-east-1:249285743859:service/ae69bb4c-3903-4844-82fe-548ac5b74570--web"]}`,
			},
		},

		awsutil.Cycle{
			Request: awsutil.Request{
				RequestURI: "/",
				Operation:  "AmazonEC2ContainerServiceV20141113.ListTasks",
				Body:       `{"cluster":"cluster","serviceName":"ae69bb4c-3903-4844-82fe-548ac5b74570--web"}`,
			},
			Response: awsutil.Response{
				StatusCode: 200,
				Body:       `{"taskArns":["arn:aws:ecs:us-east-1:249285743859:task/ae69bb4c-3903-4844-82fe-548ac5b74570"],"nextToken":"1234"}`,
			},
		},

		awsutil.Cycle{
			Request: awsutil.Request{
				RequestURI: "/",
				Operation:  "AmazonEC2ContainerServiceV20141113.ListTasks",
				Body:       `{"cluster":"cluster","serviceName":"ae69bb4c-3903-4844-82fe-548ac5b74570--web","nextToken":"1234"}`,
			},
			Response: awsutil.Response{
				StatusCode: 200,
				Body:       `{"taskArns":["arn:aws:ecs:us-east-1:249285743859:task/ae69bb4c-3903-4844-82fe-548ac5b74570"]}`,
			},
		},
	})
	m, s := newTestClient(h)
	defer s.Close()

	resp, err := m.ListAppTasks(context.Background(), "ae69bb4c-3903-4844-82fe-548ac5b74570", &ecs.ListTasksInput{
		Cluster: aws.String("cluster"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if got := len(resp.TaskARNs); got != 2 {
		t.Fatalf("Expected 2 tasks returned; got %d", got)
	}
}

func TestListAppTaskDefinitions(t *testing.T) {
	h := awsutil.NewHandler([]awsutil.Cycle{
		awsutil.Cycle{
			Request: awsutil.Request{
				RequestURI: "/",
				Operation:  "AmazonEC2ContainerServiceV20141113.ListTaskDefinitions",
				Body:       `{}`,
			},
			Response: awsutil.Response{
				StatusCode: 200,
				Body:       `{"taskDefinitionArns":["arn:aws:ecs:us-east-1:249285743859:task-definition/ae69bb4c-3903-4844-82fe-548ac5b74570--web"]}`,
			},
		},
	})
	m, s := newTestClient(h)
	defer s.Close()

	resp, err := m.ListAppTaskDefinitions(context.Background(), "ae69bb4c-3903-4844-82fe-548ac5b74570", &ecs.ListTaskDefinitionsInput{})
	if err != nil {
		t.Fatal(err)
	}

	if got := len(resp.TaskDefinitionARNs); got != 1 {
		t.Fatalf("Expected 1 tasks returned; got %d", got)
	}
}

func newTestClient(h http.Handler) (*Client, *httptest.Server) {
	s := httptest.NewServer(h)

	return NewClient(aws.DefaultConfig.Merge(&aws.Config{
		Credentials: credentials.NewStaticCredentials(" ", " ", " "),
		Endpoint:    s.URL,
		Region:      "localhost",
	})), s
}
