// Package router contains all routes
package router

import (
	"github.com/gin-gonic/gin"
	gconfig "github.com/pilinux/gorest/config"
	gcontroller "github.com/pilinux/gorest/controller"
	gmiddleware "github.com/pilinux/gorest/lib/middleware"

	"apidev/controller"
)

// SetupRouter sets up all the routes
func SetupRouter(configure *gconfig.Configuration) (*gin.Engine, error) {
	if configure.Server.ServerEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// gin.Default() = gin.New() + gin.Logger() + gin.Recovery()
	r := gin.Default()

	// Which proxy to trust:
	// disable (set nil) this feature as it still fails
	// to provide the real client IP in
	// different scenarios
	err := r.SetTrustedProxies(nil)
	if err != nil {
		return r, err
	}

	// when using Cloudflare's CDN:
	// router.TrustedPlatform = gin.PlatformCloudflare
	//
	// when running on Google App Engine:
	// router.TrustedPlatform = gin.PlatformGoogleAppEngine
	//
	/*
		when using apache or nginx reverse proxy
		without Cloudflare's CDN or Google App Engine

		config for nginx:
		=================
		proxy_set_header X-Real-IP       $remote_addr;
		proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	*/
	// router.TrustedPlatform = "X-Real-Ip"
	//
	// set TrustedPlatform to get the real client IP
	trustedPlatform := configure.Security.TrustedPlatform
	if trustedPlatform == "cf" {
		trustedPlatform = gin.PlatformCloudflare
	}
	if trustedPlatform == "google" {
		trustedPlatform = gin.PlatformGoogleAppEngine
	}
	r.TrustedPlatform = trustedPlatform

	// CORS
	if configure.Security.MustCORS == gconfig.Activated {
		r.Use(gmiddleware.CORS(configure.Security.CORS))
	}

	// Sentry.io
	if configure.Logger.Activate == gconfig.Activated {
		r.Use(gmiddleware.SentryCapture(configure.Logger.SentryDsn))
	}

	// WAF
	if configure.Security.MustFW == gconfig.Activated {
		r.Use(gmiddleware.Firewall(
			configure.Security.Firewall.ListType,
			configure.Security.Firewall.IP,
		))
	}

	// API Status
	r.GET("", controller.APIStatus)

	// API:v1
	v1 := r.Group("/api/v1/")
	{
		// RDBMS
		if configure.Database.RDBMS.Activate == gconfig.Activated {
			// Register - no JWT required
			v1.POST("register", gcontroller.CreateUserAuth)

			// Login - app issues JWT
			v1.POST("login", gcontroller.Login)

			// Refresh - app issues new JWT
			rJWT := v1.Group("refresh")
			rJWT.Use(gmiddleware.RefreshJWT())
			rJWT.POST("", gcontroller.Refresh)

			// Two-factor authentication
			if configure.Security.Must2FA == gconfig.Activated {
				r2FA := v1.Group("2fa")
				r2FA.Use(gmiddleware.JWT())
				r2FA.POST("setup", gcontroller.Setup2FA)
				r2FA.POST("activate", gcontroller.Activate2FA)
				r2FA.POST("validate", gcontroller.Validate2FA)
				if configure.Security.Must2FA == gconfig.Activated {
					r2FA.Use(gmiddleware.TwoFA(
						configure.Security.TwoFA.Status.On,
						configure.Security.TwoFA.Status.Off,
						configure.Security.TwoFA.Status.Verified,
					))
				}
				// disable 2FA
				r2FA.POST("deactivate", gcontroller.Deactivate2FA)
			}

			// Update/reset password
			rPass := v1.Group("password")
			rPass.Use(gmiddleware.JWT())
			if configure.Security.Must2FA == gconfig.Activated {
				rPass.Use(gmiddleware.TwoFA(
					configure.Security.TwoFA.Status.On,
					configure.Security.TwoFA.Status.Off,
					configure.Security.TwoFA.Status.Verified,
				))
			}
			// change password while logged in
			rPass.POST("edit", gcontroller.PasswordUpdate)

			// User
			rUsers := v1.Group("users")
			rUsers.Use(gmiddleware.JWT())
			if configure.Security.Must2FA == gconfig.Activated {
				rUsers.Use(gmiddleware.TwoFA(
					configure.Security.TwoFA.Status.On,
					configure.Security.TwoFA.Status.Off,
					configure.Security.TwoFA.Status.Verified,
				))
			}
			rUsers.GET("", controller.GetUserProfile)
			rUsers.POST("", controller.CreateUserProfile)
			rUsers.PUT("", controller.UpdateUserProfile)

			// Note
			rPosts := v1.Group("notes")
			rPosts.Use(gmiddleware.JWT())
			if configure.Security.Must2FA == gconfig.Activated {
				rPosts.Use(gmiddleware.TwoFA(
					configure.Security.TwoFA.Status.On,
					configure.Security.TwoFA.Status.Off,
					configure.Security.TwoFA.Status.Verified,
				))
			}
			rPosts.GET("", controller.GetNotes)
			rPosts.GET("/:id", controller.GetNote)
			rPosts.POST("", controller.CreateNote)
			rPosts.PUT("/:id", controller.UpdateNote)
			rPosts.DELETE("/:id", controller.DeleteNote)
		}
	}

	return r, nil
}
