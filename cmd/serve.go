package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"wabot/frontend"
	"wabot/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the WebUI",
	Long:  `This command will run HTTP server for APIs and WebUI.`,
	Run: func(_ *cobra.Command, _ []string) {
		serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve() {
	appCfg := config.AppCfg()

	// Define Fiber config & app.
	fiberCfg := fiberConfig()
	app := fiber.New(fiberCfg)

	// ======================================
	// GLOBAL MIDDLEWARE
	// ======================================

	// recover
	app.Use(recover.New(recover.Config{
		EnableStackTrace: appCfg.Debug,
	}))

	// logger middleware
	if appCfg.Debug {
		app.Use(logger.New(logger.Config{
			Format: "[${time}] ${status} - ${latency} ${method} ${path} ${ip} ${ua}\n",
		}))
	}

	// ======================================
	// ROUTES
	// ======================================
	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.JSON(fiber.Map{
	// 		"status":  "ok",
	// 		"message": "WAT?",
	// 		"data":    nil,
	// 	})
	// })

	// ping check
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "pong",
			"data":    nil,
		})
	})

	// route.V1Routes(app)

	app.Use("/", filesystem.New(filesystem.Config{
        Root:         frontend.SvelteKitHandler(),
        // NotFoundFile: "index.html",
    }))

	// signal channel to capture system calls
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	// start shutdown goroutine
	go func() {
		// capture sigterm and other system call here
		<-sigCh
		fmt.Println("Shutting down server...")
		_ = app.Shutdown()
	}()

	// start http server
	serverAddr := fmt.Sprintf("%s:%d", appCfg.Host, appCfg.Port)
	if err := app.Listen(serverAddr); err != nil {
		fmt.Printf("Oops... server is not running! error: %v", err)
	}
}

func fiberConfig() fiber.Config {
	// Return Fiber configuration.
	return fiber.Config{
		Prefork:     config.AppCfg().Prefork,
		ProxyHeader: config.AppCfg().ProxyHeader,
		AppName:     "wabot HTTP server " + AppVer,
	}
}
