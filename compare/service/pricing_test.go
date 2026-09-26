package service

import (
	"testing"

	money "github.com/ChristianDenniss/go-data-model/money/entity"
	quoteentity "github.com/ChristianDenniss/go-data-model/quoteobs/entity"
)

func TestSumQuoteFees(t *testing.T) {
	total := sumQuoteFees(quoteentity.Observation{
		FeeLines: []quoteentity.FeeLine{
			{Amount: money.Money{AmountCents: 299}},
			{Amount: money.Money{AmountCents: 150}},
		},
	})
	if total != 449 {
		t.Fatalf("expected 449, got %d", total)
	}
}
