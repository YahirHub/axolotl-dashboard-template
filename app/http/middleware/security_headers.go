package middleware

import "github.com/goravel/framework/contracts/http"

type SecurityHeaders struct{}

func NewSecurityHeaders() *SecurityHeaders {
	return &SecurityHeaders{}
}

func (r *SecurityHeaders) Handle(ctx http.Context) {
	ctx.Response().
		Header("X-Content-Type-Options", "nosniff").
		Header("X-Frame-Options", "SAMEORIGIN").
		Header("Referrer-Policy", "strict-origin-when-cross-origin").
		Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")

	ctx.Request().Next()
}

func (r *SecurityHeaders) Signature() string {
	return "axolotl.security-headers"
}
