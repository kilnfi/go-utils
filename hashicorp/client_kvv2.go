package hashicorp

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

type KVv2Client struct {
	*api.Client

	cfg *ClientConfig

	mountPath, basePath string

	logger logrus.FieldLogger
}

func NewKVv2Client(cfg *ClientConfig) (*KVv2Client, error) {
	hashicorpConfig, err := cfg.ToHashicorpConfig()
	if err != nil {
		return nil, err
	}

	client, err := api.NewClient(hashicorpConfig)
	if err != nil {
		return nil, err
	}

	c := &KVv2Client{
		cfg:    cfg,
		Client: client,
	}

	c.SetLogger(logrus.StandardLogger())

	return c, nil
}

func (c *KVv2Client) Logger() logrus.FieldLogger {
	return c.logger
}

func (c *KVv2Client) SetLogger(logger logrus.FieldLogger) {
	c.logger = logger.WithField("component", "hashicorp.kvv2-client")
}

func (c *KVv2Client) Init(ctx context.Context) error {
	err := c.validateAddress()
	if err != nil {
		return err
	}

	err = c.initAuth(ctx)
	if err != nil {
		return err
	}

	err = c.initKVv2(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *KVv2Client) validateAddress() (err error) {
	logger := c.logger.WithField("address", c.cfg.Address)
	logger.Info("validate Vault address")
	defer func() {
		if err != nil {
			logger.WithError(err).Error("invalid Vault address")
		}
	}()

	if c.cfg.Address == "" {
		err = fmt.Errorf("missing Vault address")
		return err
	}

	// validate vault address
	_, err = url.Parse(c.cfg.Address)
	if err != nil {
		return fmt.Errorf("invalid Vault address %q [err=%v]", c.cfg.Address, err)
	}

	return nil
}
func (c *KVv2Client) initAuth(ctx context.Context) (err error) {
	c.logger.Info("authenticate on Vault")
	defer func() {
		if err != nil {
			c.logger.WithError(err).Error("authentication failed")
		}
	}()

	if c.cfg.Auth == nil {
		err = fmt.Errorf("vault authentication credentials missing")
		return
	}

	// Vault token has been provided
	if c.cfg.Auth.Token != "" {
		c.SetToken(c.cfg.Auth.Token)
		c.logger.
			WithField("auth", "token").
			Info("authentication succeeded")
		return nil
	}

	// GitHub token has been provided
	if c.cfg.Auth.GitHubToken != "" {
		logger := c.logger.WithField("auth", "github")

		// Perform Github login
		ghSecAuth, ghErr := c.GithubLogin(ctx)
		if ghErr != nil {
			logger.WithError(ghErr).Errorf("authentication failed")
			return fmt.Errorf("GitHub authentication failed %v", ghErr)
		}

		logger.Infof("authentication succeeded")

		// Set token for subsequent requests
		c.SetToken(ghSecAuth.ClientToken)

		return nil
	}

	err = fmt.Errorf("vault authentication credentials missing")

	return
}

func (c *KVv2Client) initKVv2(ctx context.Context) error {
	c.logger.WithField("path", c.cfg.Path).Info("check Vault mount is kv-v2")
	mountPath, ok, err := c.IsKVv2(ctx, c.cfg.Path)
	if err != nil {
		c.logger.WithError(err).Errorf("failed to check Vault mount")
		return err
	}

	logger := c.logger.WithField("mount", mountPath)
	if !ok {
		logger.Errorf("Vault mount is not a kv-v2 secret engine")
		return fmt.Errorf("mount %q is not a kv-v2 secret engine", mountPath)
	}

	c.mountPath = mountPath
	c.basePath = strings.TrimPrefix(c.cfg.Path, mountPath)

	logger.WithField("base.path", "c.basePath").Infof("Vault mount is kv-v2")

	return nil
}

func (c *KVv2Client) GithubLogin(ctx context.Context) (*api.SecretAuth, error) {
	secret, err := c.Client.Logical().WriteWithContext(
		ctx,
		"auth/github/login",
		map[string]interface{}{
			"token": strings.TrimSpace(c.cfg.Auth.GitHubToken),
		},
	)

	if err != nil {
		return nil, err
	}

	if secret == nil {
		return nil, fmt.Errorf("empty response")
	}

	return secret.Auth, nil
}

func (c *KVv2Client) SetToken(token string) {
	c.Client.SetToken(token)
}

func (c *KVv2Client) LookupToken(ctx context.Context, token string) (*api.Secret, error) {
	return c.Client.Logical().WriteWithContext(
		ctx,
		"auth/token/lookup",
		map[string]interface{}{
			"token": token,
		},
	)
}

func (c *KVv2Client) UnwrapToken(ctx context.Context, token string) (*api.Secret, error) {
	return c.Client.Logical().UnwrapWithContext(ctx, token)
}

func (c *KVv2Client) HealthCheck(ctx context.Context) error {
	resp, err := c.Client.Sys().HealthWithContext(ctx)
	if err != nil {
		return err
	}

	if !resp.Initialized {
		return fmt.Errorf("hashicorp client is not initialized")
	}

	return nil
}

func (c *KVv2Client) Put(ctx context.Context, id string, data map[string]interface{}) (*api.Secret, error) {
	return c.Client.Logical().WriteWithContext(
		ctx,
		c.dataPath(id),
		map[string]interface{}{
			"data": data,
		},
	)
}

func (c *KVv2Client) Get(ctx context.Context, pth, version string) (secret *api.Secret, data, metadata map[string]interface{}, err error) {
	query := map[string][]string{}
	if version != "" {
		query["version"] = []string{version}
	}

	secret, err = c.Client.Logical().ReadWithDataWithContext(
		ctx,
		c.dataPath(pth),
		query,
	)
	if err != nil {
		return
	}

	if secret == nil || secret.Data["data"] == nil {
		return secret, nil, nil, fmt.Errorf("empty secret at path %q", pth)
	}

	var ok bool
	data, ok = secret.Data["data"].(map[string]interface{})
	if !ok {
		return secret, nil, nil, fmt.Errorf("invalid hashicorp vault response data: %v", secret.Data["data"])
	}

	metadata, ok = secret.Data["metadata"].(map[string]interface{})
	if !ok {
		return secret, nil, nil, fmt.Errorf("invalid hashicorp vault response metadata: %v", secret.Data["metadata"])
	}

	return
}

func (c *KVv2Client) List(ctx context.Context, pth string) ([]string, error) {
	secret, err := c.Client.Logical().ListWithContext(
		ctx,
		c.metadaDataPath(pth),
	)
	if err != nil {
		return nil, err
	}

	return extractListData(secret)
}

func extractListData(secret *api.Secret) ([]string, error) {
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("empty response")
	}

	keys, ok := secret.Data["keys"]
	if !ok {
		return nil, fmt.Errorf("invalid response body missing \"keys\" field")
	}

	ikeys, ok := keys.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid \"keys\" format: %v", keys)
	}

	skeys := make([]string, len(ikeys))
	for i, key := range ikeys {
		skeys[i] = fmt.Sprintf("%v", key)
	}

	return skeys, nil
}

func (c *KVv2Client) Delete(ctx context.Context, id string) error {
	return fmt.Errorf("not implemented error")
}

func (c *KVv2Client) dataPath(id string) string {
	return path.Join(c.mountPath, "data", c.basePath, id)
}

func (c *KVv2Client) metadaDataPath(id string) string {
	return path.Join(c.mountPath, "metadata", c.basePath, id)
}

func (c *KVv2Client) kvPreflightVersionRequest(ctx context.Context, pth string) (mountPath string, version int, err error) {
	// We don't want to use a wrapping call here so save any custom value and
	// restore after
	currentWrappingLookupFunc := c.Client.CurrentWrappingLookupFunc()
	c.Client.SetWrappingLookupFunc(nil)
	defer c.Client.SetWrappingLookupFunc(currentWrappingLookupFunc)
	currentOutputCurlString := c.Client.OutputCurlString()
	c.Client.SetOutputCurlString(false)
	defer c.Client.SetOutputCurlString(currentOutputCurlString)

	r := c.Client.NewRequest("GET", path.Join("/v1/sys/internal/ui/mounts", pth))

	resp, err := c.Client.RawRequestWithContext(ctx, r) //nolint
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		// If we get a 404 we are using an older version of vault, default to
		// version 1
		if resp != nil && resp.StatusCode == 404 {
			return "", 1, nil
		}

		return "", 0, err
	}

	secret, err := api.ParseSecret(resp.Body)
	if err != nil {
		return "", 0, err
	}
	if secret == nil {
		return "", 0, fmt.Errorf("nil response from pre-flight request")
	}

	if mountPathRaw, ok := secret.Data["path"]; ok {
		mountPath = mountPathRaw.(string)
	}
	options := secret.Data["options"]
	if options == nil {
		return mountPath, 1, nil
	}
	versionRaw := options.(map[string]interface{})["version"]
	if versionRaw == nil {
		return mountPath, 1, nil
	}

	switch versionRaw.(string) {
	case "", "1":
		return mountPath, 1, nil
	case "2":
		return mountPath, 2, nil
	}

	return mountPath, 1, nil
}

func (c *KVv2Client) IsKVv2(ctx context.Context, pth string) (mount string, isKVv2 bool, err error) {
	mountPath, version, err := c.kvPreflightVersionRequest(ctx, pth)
	if err != nil {
		return "", false, err
	}

	return mountPath, version == 2, nil
}
