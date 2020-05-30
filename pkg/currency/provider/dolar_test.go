package provider

import (
	"reflect"
	"testing"

	"coinbani/pkg/currency"
)

func Test_addDollarPrices(t *testing.T) {
	type args struct {
		lastPrices []*currency.CurrencyPrice
		price      dollarPrice
	}
	tests := []struct {
		name string
		args args
		want []*currency.CurrencyPrice
	}{
		{
			name: "add dollar price",
			args: args{
				lastPrices: []*currency.CurrencyPrice{
					{
						Desc:          "Oficial",
						BidPrice:      float64(65.72),
						AskPrice:      float64(70.72),
						Currency:      "USD",
						PercentChange: "+0,040",
					},
					{
						Desc:          "Blue",
						BidPrice:      float64(115.00),
						AskPrice:      float64(125.00),
						Currency:      "USD",
						PercentChange: "+0,810",
					},
				},
				price: dollarPrice{
					Name:          "Dolar Ahorro",
					BidPrice:      "87,05",
					AskPrice:      "91,13",
					PercentChange: "-0,040",
				},
			},
			want: []*currency.CurrencyPrice{
				{
					Desc:          "Oficial",
					Currency:      "USD",
					BidPrice:      float64(65.72),
					AskPrice:      float64(70.72),
					PercentChange: "+0,040",
				},
				{
					Desc:          "Blue",
					Currency:      "USD",
					BidPrice:      float64(115.00),
					AskPrice:      float64(125.00),
					PercentChange: "+0,810",
				},
				{
					Desc:          "Ahorro",
					Currency:      "USD",
					BidPrice:      float64(87.05),
					AskPrice:      float64(91.13),
					PercentChange: "-0,040",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := addDollarPrices(tt.args.lastPrices, tt.args.price)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addDollarPrices() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatDollarName(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "format dollar MEP name",
			args: args{v: "Dolar Bolsa"},
			want: "MEP",
		},
		{
			name: "format dollar CCL name",
			args: args{v: "Dolar Contado con Liqui"},
			want: "CCL",
		},
		{
			name: "format other things should not replace the name",
			args: args{v: "Dolar Oficial"},
			want: "Oficial",
		},
		{
			name: "format empty string shouldn't fail",
			args: args{v: ""},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatDollarName(tt.args.v); got != tt.want {
				t.Errorf("formatDollarName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_replaceComa(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "replace comma with dot",
			args: args{value: "70,43"},
			want: "70.43",
		},
		{
			name: "replace without comma shouln´t fail",
			args: args{value: "10"},
			want: "10",
		},
		{
			name: "replace with empty string shouln´t fail",
			args: args{value: ""},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceComa(tt.args.value); got != tt.want {
				t.Errorf("replaceComa() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatPercent(t *testing.T) {
	type args struct {
		percentChange string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "add plus sign to percentage",
			args: args{percentChange: "0,810"},
			want: "+0,810",
		},
		{
			name: "skip plus sign if percentage is negative",
			args: args{percentChange: "-0,920"},
			want: "-0,920",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatPercent(tt.args.percentChange); got != tt.want {
				t.Errorf("formatPercent() = %v, want %v", got, tt.want)
			}
		})
	}
}
