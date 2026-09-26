// Generates TypeScript interfaces from go-data-model entities.
// Run from the module root: go run ./cmd/gentypes
package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	account "github.com/ChristianDenniss/go-data-model/account/entity"
	brand "github.com/ChristianDenniss/go-data-model/brand/entity"
	cart "github.com/ChristianDenniss/go-data-model/cart/entity"
	channel "github.com/ChristianDenniss/go-data-model/channel/entity"
	dish "github.com/ChristianDenniss/go-data-model/dish/entity"
	place "github.com/ChristianDenniss/go-data-model/place/entity"
	source "github.com/ChristianDenniss/go-data-model/source/entity"
	user "github.com/ChristianDenniss/go-data-model/user/entity"
	category "github.com/ChristianDenniss/go-data-model/category/entity"
	cuisine "github.com/ChristianDenniss/go-data-model/cuisine/entity"
	location "github.com/ChristianDenniss/go-data-model/location/entity"
	menu "github.com/ChristianDenniss/go-data-model/menu/entity"
	money "github.com/ChristianDenniss/go-data-model/money/entity"
	observation "github.com/ChristianDenniss/go-data-model/observation/entity"
	offer "github.com/ChristianDenniss/go-data-model/offer/entity"
	order "github.com/ChristianDenniss/go-data-model/order/entity"
	provider "github.com/ChristianDenniss/go-data-model/provider/entity"
	restaurant "github.com/ChristianDenniss/go-data-model/restaurant/entity"
)

type named struct {
	name string
	v    any
}

func main() {
	typesToEmit := []named{
		{name: "Money", v: money.Money{}},
		{name: "Location", v: location.Location{}},
		{name: "Rating", v: restaurant.Rating{}},
		{name: "Provider", v: provider.Provider{}},
		{name: "Channel", v: channel.Channel{}},
		{name: "Brand", v: brand.Brand{}},
		{name: "Place", v: place.Place{}},
		{name: "Dish", v: dish.Dish{}},
		{name: "SourceStore", v: source.Store{}},
		{name: "SourceItem", v: source.Item{}},
		{name: "User", v: user.User{}},
		{name: "ComparePrefs", v: user.ComparePrefs{}},
		{name: "Category", v: category.Category{}},
		{name: "Cuisine", v: cuisine.Cuisine{}},
		{name: "Restaurant", v: restaurant.Restaurant{}},
		{name: "Item", v: menu.Item{}},
		{name: "Offer", v: offer.Offer{}},
		{name: "Observation", v: observation.Observation{}},
		{name: "SavedAddress", v: account.SavedAddress{}},
		{name: "PaymentMethod", v: account.PaymentMethod{}},
		{name: "Account", v: account.Account{}},
		{name: "CartLine", v: cart.CartLine{}},
		{name: "Cart", v: cart.Cart{}},
		{name: "OrderLine", v: order.OrderLine{}},
		{name: "Order", v: order.Order{}},
	}

	var body bytes.Buffer
	body.WriteString("/**\n")
	body.WriteString(" * Generated from go-data-model. Do not edit.\n")
	body.WriteString(" * Regenerate: go run ./cmd/gentypes (from go-data-model).\n")
	body.WriteString(" */\n\n")
	body.WriteString("export type OrderStatus = 'pending' | 'completed' | 'cancelled'\n\n")

	for _, item := range typesToEmit {
		emitInterface(&body, item.name, reflect.TypeOf(item.v))
	}

	root, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if filepath.Base(root) == "gentypes" {
		root = filepath.Join(root, "..", "..")
	}

	targets := []string{
		filepath.Join(root, "gen", "typescript", "data-model.ts"),
		filepath.Join(root, "..", "nibble-web-platform", "src", "generated", "data-model.ts"),
	}
	for _, target := range targets {
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			panic(err)
		}
		if err := os.WriteFile(target, body.Bytes(), 0o644); err != nil {
			panic(err)
		}
		fmt.Printf("wrote %s\n", target)
	}
}

func emitInterface(buf *bytes.Buffer, name string, t reflect.Type) {
	fmt.Fprintf(buf, "export interface %s {\n", name)
	for _, f := range collectFields(t) {
		fmt.Fprintf(buf, "  %s: %s\n", f.tsName, f.tsType)
	}
	buf.WriteString("}\n\n")
}

type fieldOut struct {
	tsName string
	tsType string
}

func collectFields(t reflect.Type) []fieldOut {
	out := make([]fieldOut, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		out = append(out, fieldOut{
			tsName: jsonName(f.Name),
			tsType: tsTypeOf(f.Type),
		})
	}
	return out
}

func jsonName(name string) string {
	switch {
	case name == "ID":
		return "id"
	case name == "IDs":
		return "ids"
	case strings.HasSuffix(name, "IDs"):
		return lowerFirst(strings.TrimSuffix(name, "IDs")) + "Ids"
	case strings.HasSuffix(name, "ID"):
		return lowerFirst(strings.TrimSuffix(name, "ID")) + "Id"
	default:
		return lowerFirst(name)
	}
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func tsTypeOf(t reflect.Type) string {
	if t == reflect.TypeOf(time.Time{}) {
		return "string"
	}
	if t.PkgPath() == "github.com/ChristianDenniss/go-data-model/order/entity" && t.Name() == "Status" {
		return "OrderStatus"
	}
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Slice:
		return tsTypeOf(t.Elem()) + "[]"
	case reflect.Struct:
		if t.Name() != "" {
			return t.Name()
		}
		return "Record<string, unknown>"
	case reflect.Pointer:
		return tsTypeOf(t.Elem()) + " | null"
	default:
		return "unknown"
	}
}
