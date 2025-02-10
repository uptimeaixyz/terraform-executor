package awsclient

import (
	"context"
	"fmt"
	"time"

	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// Required environment variables for AWS authentication:
// AWS_ACCESS_KEY_ID - AWS access key
// AWS_SECRET_ACCESS_KEY - AWS secret key
// AWS_REGION - AWS region (e.g., us-west-2)
// Optional:
// AWS_SESSION_TOKEN - Required only when using temporary credentials
// AWS_PROFILE - AWS credential profile name
type AWSClient struct {
	cfg       aws.Config
	S3Client  *s3.Client
	IAMClient *iam.Client
	STSClient *sts.Client
}

// NewAWSClient creates a new AWS client using credentials from environment variables
// or AWS credential file (~/.aws/credentials)
func NewAWSClient(ctx context.Context) (*AWSClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := &AWSClient{
		cfg: cfg,
	}

	if err := client.initServices(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *AWSClient) initServices() error {
	c.S3Client = s3.NewFromConfig(c.cfg)
	c.IAMClient = iam.NewFromConfig(c.cfg)
	c.STSClient = sts.NewFromConfig(c.cfg)
	return nil
}

const (
	RolePath            = "/app/uptimeai/"
	s3BoundaryPolicyArn = "arn:aws:iam::%s:policy/s3-boundary" // Changed to include account ID placeholder
)

// GetAccountID returns the AWS account ID for the current credentials
func (c *AWSClient) GetAccountID(ctx context.Context) (string, error) {
	result, err := c.STSClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", fmt.Errorf("failed to get account ID: %w", err)
	}
	return *result.Account, nil
}

// AttachRolePolicy attaches the specified policy to an IAM role
// Required IAM permissions: iam:AttachRolePolicy
func (c *AWSClient) AttachRolePolicy(ctx context.Context, roleName string, policyArn string) error {
	input := &iam.AttachRolePolicyInput{
		RoleName:  aws.String(roleName),
		PolicyArn: aws.String(policyArn),
	}

	_, err := c.IAMClient.AttachRolePolicy(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to attach policy to role: %w", err)
	}

	return nil
}

// CreateRole creates an IAM role with specified name under the /app/uptimeai/ path
func (c *AWSClient) CreateRole(ctx context.Context, roleName string, trustPolicy string) (*types.Role, error) {
	accountID, err := c.GetAccountID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS account ID: %w", err)
	}

	boundaryArn := fmt.Sprintf(s3BoundaryPolicyArn, accountID)

	input := &iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		Path:                     aws.String(RolePath),
		AssumeRolePolicyDocument: aws.String(trustPolicy),
		PermissionsBoundary:      aws.String(boundaryArn),
		Tags: []types.Tag{
			{
				Key:   aws.String("UserId"), // Add UserId tag required by the boundary policy
				Value: aws.String(roleName), // Use role name as UserId
			},
			{
				Key:   aws.String("CreatedBy"),
				Value: aws.String("terraform-executor"),
			},
			{
				Key:   aws.String("CreatedAt"),
				Value: aws.String(time.Now().Format(time.RFC3339)),
			},
		},
	}

	result, err := c.IAMClient.CreateRole(ctx, input)
	if err != nil {
		var alreadyExists *types.EntityAlreadyExistsException
		if errors.As(err, &alreadyExists) {
			return nil, fmt.Errorf("role already exists: %w", err)
		}
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	// Update to use the correct boundary ARN
	if err := c.AttachRolePolicy(ctx, roleName, boundaryArn); err != nil {
		// If policy attachment fails, attempt to clean up by deleting the role
		deleteInput := &iam.DeleteRoleInput{
			RoleName: aws.String(roleName),
		}
		_, deleteErr := c.IAMClient.DeleteRole(ctx, deleteInput)
		if deleteErr != nil {
			return nil, fmt.Errorf("failed to attach policy and cleanup role: %v (cleanup error: %v)", err, deleteErr)
		}
		return nil, fmt.Errorf("failed to attach policy: %w", err)
	}

	return result.Role, nil
}

// AssumeRole assumes the specified IAM role and returns temporary credentials
func (c *AWSClient) AssumeRole(ctx context.Context, roleArn, sessionName string) (*sts.AssumeRoleOutput, error) {
	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(sessionName),
		DurationSeconds: aws.Int32(3600), // 1 hour
	}

	return c.STSClient.AssumeRole(ctx, input)
}

// GetTemporaryCredentials is a helper function that returns formatted temporary credentials
func (c *AWSClient) GetTemporaryCredentials(ctx context.Context, roleArn, sessionName string) (accessKey, secretKey, sessionToken string, expiration time.Time, err error) {
	result, err := c.AssumeRole(ctx, roleArn, sessionName)
	if err != nil {
		return "", "", "", time.Time{}, fmt.Errorf("failed to assume role: %w", err)
	}

	if result.Credentials == nil {
		return "", "", "", time.Time{}, fmt.Errorf("no credentials returned")
	}

	return *result.Credentials.AccessKeyId,
		*result.Credentials.SecretAccessKey,
		*result.Credentials.SessionToken,
		*result.Credentials.Expiration,
		nil
}

// RoleExists checks if the specified IAM role exists
// Required IAM permissions: iam:GetRole
func (c *AWSClient) RoleExists(ctx context.Context, roleName string) (bool, error) {
	input := &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	}

	_, err := c.IAMClient.GetRole(ctx, input)
	if err != nil {
		// Check if the error is because the role doesn't exist
		var noSuchEntity *types.NoSuchEntityException
		if errors.As(err, &noSuchEntity) {
			return false, nil
		}
		// For any other error, return it
		return false, fmt.Errorf("failed to check role existence: %w", err)
	}

	return true, nil
}
