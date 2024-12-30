package es

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/loveuer/upp/pkg/tool"
	"github.com/samber/lo"
)

// New elasticsearch client v7
// example:
//   - uri: http://127.0.0.1:9200
//   - uri: https://<username>:<password>@node1:9200,node2:9200,node3:9200
func New(ctx context.Context, uri string) (*elastic.Client, error) {
	var (
		err      error
		username string
		password string
		client   *elastic.Client
		ins      *url.URL
	)

	if ins, err = url.Parse(uri); err != nil {
		return nil, err
	}

	endpoints := lo.Map(
		strings.Split(ins.Host, ","),
		func(item string, index int) string {
			return fmt.Sprintf("%s://%s", ins.Scheme, item)
		},
	)

	if ins.User != nil {
		username = ins.User.Username()
		password, _ = ins.User.Password()
	}

	if client, err = elastic.NewClient(
		elastic.Config{
			Addresses:     endpoints,
			Username:      username,
			Password:      password,
			CACert:        nil,
			RetryOnStatus: []int{429},
			MaxRetries:    3,
			RetryBackoff:  nil,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				DialContext:     (&net.Dialer{Timeout: 10 * time.Second}).DialContext,
			},
		},
	); err != nil {
		return nil, err
	}

	var res *esapi.Response
	if res, err = client.Ping(client.Ping.WithContext(tool.TimeoutCtx(ctx, 5))); err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		err = fmt.Errorf("ping client response: %s", res.String())
		return nil, err
	}

	return client, nil
}
