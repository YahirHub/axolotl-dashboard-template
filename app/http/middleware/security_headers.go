package middleware

import "github.com/goravel/framework/contracts/http"

const contentSecurityPolicy = "default-src 'self'; " +
	"base-uri 'self'; " +
	"form-action 'self'; " +
	"frame-ancestors 'self'; " +
	"object-src 'none'; " +
	"script-src 'self' https://cdn.jsdelivr.net https://cdnjs.cloudflare.com; " +
	"style-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net https://cdnjs.cloudflare.com; " +
	"font-src 'self' data: https://cdnjs.cloudflare.com; " +
	"img-src 'self' data:; " +
	"connect-src 'self' https://cdn.jsdelivr.net https://cdnjs.cloudflare.com; " +
	"media-src 'self'; " +
	"worker-src 'self'; " +
	"manifest-src 'self'"

type SecurityHeaders struct{}

func NewSecurityHeaders() *SecurityHeaders {
	return &SecurityHeaders{}
}

func (r *SecurityHeaders) Handle(ctx http.Context) {
	response := ctx.Response().
		Header("X-Content-Type-Options", "nosniff").
		Header("X-Frame-Options", "SAMEORIGIN").
		Header("Referrer-Policy", "strict-origin-when-cross-origin").
		Header("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=(), usb=()").
		Header("Content-Security-Policy", contentSecurityPolicy).
		Header("X-Permitted-Cross-Domain-Policies", "none").
		Header("Origin-Agent-Cluster", "?1")

	if request := ctx.Request().Origin(); request != nil && request.TLS != nil {
		response.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	}

	ctx.Request().Next()
}

func (r *SecurityHeaders) Signature() string {
	return "axolotl.security-headers"
}
