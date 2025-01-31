package state

import (
	"time"

	"github.com/cloudfoundry-community/carousel/credhub"
)

type Filter func(*Credential) bool

func NotFilter(fn Filter) Filter {
	return func(c *Credential) bool {
		return !fn(c)
	}
}

func OrFilter(fns ...Filter) Filter {
	return func(c *Credential) bool {
		for _, fn := range fns {
			if fn(c) {
				return true
			}
		}
		return false
	}
}

func AnyFilter(fn Collector) Filter {
	return func(c *Credential) bool {
		return len(Credentials{c}.Collect(fn)) != 0
	}
}

func AndFilter(fns ...Filter) Filter {
	return func(c *Credential) bool {
		for _, fn := range fns {
			if !fn(c) {
				return false
			}
		}
		return true
	}
}

func SelfSignedFilter() Filter {
	return func(c *Credential) bool {
		return c.SignedBy == nil
	}
}

func ActiveFilter() Filter {
	return func(c *Credential) bool {
		return c.Active()
	}
}

func LatestFilter() Filter {
	return func(c *Credential) bool {
		return c.Latest
	}
}

func SigningFilter() Filter {
	return func(c *Credential) bool {
		return c.Signing != nil && *c.Signing
	}
}

func TransitionalFilter() Filter {
	return func(c *Credential) bool {
		return c.Transitional
	}
}

func TypeFilter(types ...credhub.CredentialType) Filter {
	return func(c *Credential) bool {
		for _, t := range types {
			if c.Type == t {
				return true
			}
		}
		return false
	}
}

func DeploymentFilter(deployment string) Filter {
	return func(c *Credential) bool {
		for _, d := range c.Path.Deployments {
			if d.Name == deployment {
				return true
			}
		}
		return false
	}
}

func NameFilter(name string) Filter {
	return func(c *Credential) bool {
		return c.Name == name
	}
}

func CertificateAuthorityFilter(expected bool) Filter {
	return func(c *Credential) bool {
		return c.CertificateAuthority == expected
	}
}

func ExpiresBeforeFilter(t time.Time) Filter {
	return func(c *Credential) bool {
		return c.ExpiryDate != nil && c.ExpiryDate.Before(t)
	}
}

func OlderThanFilter(t time.Time) Filter {
	return func(c *Credential) bool {
		return c.VersionCreatedAt.Before(t)
	}
}

func SignedByFilter(name string) Filter {
	return func(c *Credential) bool {
		return c.SignedBy != nil && c.SignedBy.Name == name
	}
}

func References(c *Credential) Filter {
	return func(c *Credential) bool {
		return c.References.Includes(c)
	}
}
