package httppreparer

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/Azure/go-autorest/autorest"
)

// WithBaseURL returns a PrepareDecorator that populates the http.Request with a url.URL constructed
// from the supplied baseUrl.  Query parameters will be encoded as required.
func WithBaseURL(baseURL string) autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r, err := p.Prepare(r)
			if err == nil {
				var u *url.URL
				if u, err = url.Parse(baseURL); err != nil {
					return r, err
				}
				if u.Scheme == "" {
					return r, fmt.Errorf("autorest: No scheme detected in URL %s", baseURL)
				}
				if u.RawQuery != "" {
					// handle unencoded semicolons (ideally the server would send them already encoded)
					u.RawQuery = strings.Replace(u.RawQuery, ";", "%3B", -1)
					var q url.Values
					q, err = url.ParseQuery(u.RawQuery)
					if err != nil {
						return r, err
					}
					u.RawQuery = q.Encode()
				}
				if r.URL == nil {
					r.URL = u
				} else {
					OverrideURL(r.URL, u)
				}
			}
			return r, err
		})
	}
}

func OverrideURL(dst, src *url.URL) {
	if src.Scheme != "" {
		dst.Scheme = src.Scheme
	}

	if src.Opaque != "" {
		dst.Opaque = src.Opaque
	}

	if src.User != nil {
		dst.User = src.User
	}

	if src.Host != "" {
		dst.Host = src.Host
	}

	if src.Path != "" {
		dst.Path = path.Join(src.Path, dst.Path)
	}

	if src.RawPath != "" {
		dst.RawPath = src.RawPath
	}

	if src.ForceQuery {
		dst.ForceQuery = src.ForceQuery
	}

	if src.RawQuery != "" {
		dst.RawQuery = src.RawQuery
	}

	if src.Fragment != "" {
		dst.Fragment = src.Fragment
	}

	if src.RawFragment != "" {
		dst.RawFragment = src.RawFragment
	}
}
