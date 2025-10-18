package bgworker

import (
	"context"
	"log"
	"time"

	"github.com/Sampath942/ecommerce/internal/user/handler"
	"github.com/Sampath942/ecommerce/internal/user/repository"
)

func RemoveExpiredTokens(ctx context.Context, h *handler.UserHandler) {
	err := repository.RemoveExpiredTokens(h.DB)
	if err != nil {
		log.Printf("The expired verification tokens clearing had a problem. Please check")
	}
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		select {
			case <- ticker.C:
				err := repository.RemoveExpiredTokens(h.DB)
				if err != nil {
					log.Printf("The expired verification tokens clearing had a problem. Please check")
				}
			case <- ctx.Done():
				log.Printf("Background worker to remove expired verifications stopped")
				return
		}
	}
}
