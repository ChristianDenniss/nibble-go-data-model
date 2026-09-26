package entity

import "testing"

func TestCanonicalizeLegacy(t *testing.T) {
	f, e := Canonicalize("delivery_3p", "")
	if f != ModeDelivery || e != ExecutorThirdParty {
		t.Fatalf("got %s %s", f, e)
	}
	f, e = Canonicalize("pickup", "")
	if f != ModePickup || e != "" {
		t.Fatalf("got %s %s", f, e)
	}
}

func TestInferDeliveryExecutor(t *testing.T) {
	if InferDeliveryExecutor(ModeDelivery, "", ChannelKindAggregator) != ExecutorThirdParty {
		t.Fatal("expected third_party from aggregator channel")
	}
	if InferDeliveryExecutor(ModeDelivery, "", "merchant_app") != ExecutorMerchant {
		t.Fatal("expected merchant default for merchant channel")
	}
}

func TestOptionMatchesFilter(t *testing.T) {
	if !OptionMatchesFilter(ModeDelivery, ExecutorThirdParty, "delivery_3p") {
		t.Fatal("legacy filter should match third party delivery")
	}
	if OptionMatchesFilter(ModeDelivery, ExecutorMerchant, "delivery_3p") {
		t.Fatal("merchant delivery should not match delivery_3p filter")
	}
	if !OptionMatchesFilter(ModeDelivery, ExecutorMerchant, ModeDelivery) {
		t.Fatal("delivery without executor filter should match any executor")
	}
}
