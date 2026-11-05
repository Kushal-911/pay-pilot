package payments

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"paypilot/internal/paypal"
)

func providerCaptureMatches(amount paypal.Money, expectedMinor int64, expectedCurrency string) bool {
	if expectedMinor <= 0 || amount.Value == "" || !strings.EqualFold(amount.CurrencyCode, expectedCurrency) {
		return false
	}
	actualMinor, err := parseMinorAmount(amount.Value)
	return err == nil && actualMinor > 0 && actualMinor == expectedMinor
}

func markPaymentReviewRequired(ctx context.Context, db *pgxpool.Pool, attemptID string) error {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Exec(ctx, `UPDATE payment_attempts SET status='unknown',updated_at=now() WHERE id=$1`, attemptID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `UPDATE orders SET status='review_required',updated_at=now() WHERE id=(SELECT order_id FROM payment_attempts WHERE id=$1)`, attemptID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}
