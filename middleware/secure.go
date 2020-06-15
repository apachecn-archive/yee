package middleware

import (
	"fmt"
	"github.com/cookieY/yee"
)

type (
	SecureConfig struct {
		XSSProtection string `yaml:"xss_protection"`

		ContentTypeNosniff string `yaml:"content_type_nosniff"`

		XFrameOptions string `yaml:"x_frame_options"`

		HSTSMaxAge int `yaml:"hsts_max_age"`

		HSTSExcludeSubdomains bool `yaml:"hsts_exclude_subdomains"`

		ContentSecurityPolicy string `yaml:"content_security_policy"`

		CSPReportOnly bool `yaml:"csp_report_only"`

		HSTSPreloadEnabled bool `yaml:"hsts_preload_enabled"`

		ReferrerPolicy string `yaml:"referrer_policy"`
	}
)

var DefaultSecureConfig = SecureConfig{
	XSSProtection:      "1; mode=block",
	ContentTypeNosniff: "nosniff",
	XFrameOptions:      "SAMEORIGIN",
	HSTSPreloadEnabled: false,
}

func Secure() yee.HandlerFunc {
	return SecureWithConfig(DefaultSecureConfig)
}

func SecureWithConfig(config SecureConfig) yee.HandlerFunc {
	return yee.HandlerFunc{
		Func: func(c yee.Context) (err error) {

			if config.XSSProtection != "" {
				c.SetHeader(yee.HeaderXXSSProtection, config.XSSProtection)
			}

			if config.ContentTypeNosniff != "" {
				c.SetHeader(yee.HeaderXContentTypeOptions, config.ContentTypeNosniff)
			}

			if config.XFrameOptions != "" {
				c.SetHeader(yee.HeaderXFrameOptions, config.XFrameOptions)
			}

			if (c.IsTls() || (c.GetHeader(yee.HeaderXForwardedProto) == "https")) && config.HSTSMaxAge != 0 {
				subdomains := ""
				if !config.HSTSExcludeSubdomains {
					subdomains = "; includeSubdomains"
				}
				if config.HSTSPreloadEnabled {
					subdomains = fmt.Sprintf("%s; preload", subdomains)
				}
				c.SetHeader(yee.HeaderStrictTransportSecurity, fmt.Sprintf("max-age=%d%s", config.HSTSMaxAge, subdomains))
			}
			// CSP
			// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Security-Policy-Report-Only
			// https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/Content_Security_Policy
			if config.ContentSecurityPolicy != "" {
				if config.CSPReportOnly {
					c.SetHeader(yee.HeaderContentSecurityPolicyReportOnly, config.ContentSecurityPolicy)
				} else {
					c.SetHeader(yee.HeaderContentSecurityPolicy, config.ContentSecurityPolicy)
				}
			}

			// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Referrer-Policy
			if config.ReferrerPolicy != "" {
				c.SetHeader(yee.HeaderReferrerPolicy, config.ReferrerPolicy)
			}
			return
		},
		IsMiddleware: true,
	}
}
