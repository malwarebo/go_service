package main

import (
	"log"
	"net/http"

	"github.com/malwarebo/gopay/api"
	"github.com/malwarebo/gopay/config"
	"github.com/malwarebo/gopay/db"
	"github.com/malwarebo/gopay/providers"
	"github.com/malwarebo/gopay/repositories"
	"github.com/malwarebo/gopay/services"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := db.NewDB(cfg.GetDatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize payment providers
	stripeProvider := providers.NewStripeProvider(cfg.Stripe.Secret)
	xenditProvider := providers.NewXenditProvider(cfg.Xendit.Secret)

	// Create a provider selector that can handle multiple providers
	providerSelector := &providers.MultiProviderSelector{
		Providers: []providers.PaymentProvider{stripeProvider, xenditProvider},
	}

	// Initialize repositories
	paymentRepo := repositories.NewPaymentRepository(db)
	planRepo := repositories.NewPlanRepository(db)
	subscriptionRepo := repositories.NewSubscriptionRepository(db)
	disputeRepo := repositories.NewDisputeRepository(db.DB)

	// Initialize services
	paymentService := services.NewPaymentService(paymentRepo, providerSelector)
	subscriptionService := services.NewSubscriptionService(planRepo, subscriptionRepo, providerSelector)
	disputeService := services.NewDisputeService(disputeRepo, providerSelector)

	// Initialize handlers
	paymentHandler := api.NewPaymentHandler(paymentService)
	subscriptionHandler := api.NewSubscriptionHandler(subscriptionService)
	disputeHandler := api.NewDisputeHandler(disputeService)

	// Setup payment routes
	http.HandleFunc("/charge", paymentHandler.HandleCharge)
	http.HandleFunc("/refund", paymentHandler.HandleRefund)

	// Setup subscription routes
	http.HandleFunc("/plans", subscriptionHandler.HandlePlans)
	http.HandleFunc("/plans/", subscriptionHandler.HandlePlans)
	http.HandleFunc("/subscriptions", subscriptionHandler.HandleSubscriptions)
	http.HandleFunc("/subscriptions/", subscriptionHandler.HandleSubscriptions)

	// Setup dispute routes
	http.HandleFunc("/disputes", disputeHandler.HandleDisputes)
	http.HandleFunc("/disputes/", disputeHandler.HandleDisputes)
	http.HandleFunc("/disputes/stats", disputeHandler.HandleDisputes)

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
