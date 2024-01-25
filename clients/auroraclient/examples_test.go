package auroraclient_test

import (
	"context"
	"fmt"
	"time"

	"github.com/diamcircle/go/clients/auroraclient"
	"github.com/diamcircle/go/keypair"
	"github.com/diamcircle/go/network"
	hProtocol "github.com/diamcircle/go/protocols/aurora"
	"github.com/diamcircle/go/protocols/aurora/effects"
	"github.com/diamcircle/go/protocols/aurora/operations"
	"github.com/diamcircle/go/txnbuild"
)

func ExampleClient_Accounts() {
	client := auroraclient.DefaultPublicNetClient
	accountsRequest := auroraclient.AccountsRequest{Signer: "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU"}

	account, err := client.Accounts(accountsRequest)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print(account)
}

func ExampleClient_AccountDetail() {
	client := auroraclient.DefaultPublicNetClient
	accountRequest := auroraclient.AccountRequest{AccountID: "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU"}

	account, err := client.AccountDetail(accountRequest)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print(account)
}

func ExampleClient_Assets() {
	client := auroraclient.DefaultPublicNetClient
	// assets for asset issuer
	assetRequest := auroraclient.AssetRequest{ForAssetIssuer: "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU"}
	asset, err := client.Assets(assetRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(asset)

	// all assets
	assetRequest = auroraclient.AssetRequest{}
	asset, err = client.Assets(assetRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(asset)
}

func ExampleClient_Effects() {
	client := auroraclient.DefaultPublicNetClient
	// effects for an account
	effectRequest := auroraclient.EffectRequest{ForAccount: "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU"}
	effect, err := client.Effects(effectRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(effect)

	// all effects
	effectRequest = auroraclient.EffectRequest{}
	effect, err = client.Effects(effectRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	records := effect.Embedded.Records
	if records[0].GetType() == "account_created" {
		acc, ok := records[0].(effects.AccountCreated)
		if ok {
			fmt.Print(acc.Account)
			fmt.Print(acc.StartingBalance)
		}
	}
}

func ExampleClient_FeeStats() {
	client := auroraclient.DefaultPublicNetClient
	// aurora fees
	fees, err := client.FeeStats()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(fees)

}

func ExampleClient_Fund() {
	client := auroraclient.DefaultTestNetClient
	// fund an account
	resp, err := client.Fund("GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(resp)
}

func ExampleClient_LedgerDetail() {
	client := auroraclient.DefaultPublicNetClient
	// details for a ledger
	sequence := uint32(12345)
	ledger, err := client.LedgerDetail(sequence)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(ledger)

}

func ExampleClient_NextAccountsPage() {
	client := auroraclient.DefaultPublicNetClient
	// accounts with signer
	accountsRequest := auroraclient.AccountsRequest{Signer: "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU",
		Limit: 20}
	accounts, err := client.Accounts(accountsRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Page 1:")
	for _, a := range accounts.Embedded.Records {
		fmt.Println(a.ID)
	}

	// next page
	accounts2, err := client.NextAccountsPage(accounts)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Page 2:")
	for _, a := range accounts2.Embedded.Records {
		fmt.Println(a.ID)
	}
}

func ExampleClient_NextAssetsPage() {
	client := auroraclient.DefaultPublicNetClient
	// assets for asset issuer
	assetRequest := auroraclient.AssetRequest{ForAssetIssuer: "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU",
		Limit: 20}
	asset, err := client.Assets(assetRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(asset)

	// all assets
	assetRequest = auroraclient.AssetRequest{}
	asset, err = client.Assets(assetRequest)
	if err != nil {
		fmt.Println(err)
		return
	}

	// next page
	nextPage, err := client.NextAssetsPage(asset)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(nextPage)
}

func ExampleClient_NextEffectsPage() {
	client := auroraclient.DefaultPublicNetClient
	// all effects
	effectRequest := auroraclient.EffectRequest{Limit: 20}
	efp, err := client.Effects(effectRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(efp)

	// get next pages.
	recordsFound := false
	if len(efp.Embedded.Records) > 0 {
		recordsFound = true
	}
	page := efp
	// get the next page of records if recordsFound is true
	for recordsFound {
		// next page
		nextPage, err := client.NextEffectsPage(page)
		if err != nil {
			fmt.Println(err)
			return
		}

		page = nextPage
		if len(nextPage.Embedded.Records) == 0 {
			recordsFound = false
		}
		fmt.Println(nextPage)
	}
}

func ExampleClient_NextLedgersPage() {
	client := auroraclient.DefaultPublicNetClient
	// all ledgers
	ledgerRequest := auroraclient.LedgerRequest{Limit: 20}
	ledgers, err := client.Ledgers(ledgerRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(ledgers)

	// get next pages.
	recordsFound := false
	if len(ledgers.Embedded.Records) > 0 {
		recordsFound = true
	}
	page := ledgers
	// get the next page of records if recordsFound is true
	for recordsFound {
		// next page
		nextPage, err := client.NextLedgersPage(page)
		if err != nil {
			fmt.Println(err)
			return
		}

		page = nextPage
		if len(nextPage.Embedded.Records) == 0 {
			recordsFound = false
		}
		fmt.Println(nextPage)
	}
}

func ExampleClient_NextOffersPage() {
	client := auroraclient.DefaultPublicNetClient
	// all offers
	offerRequest := auroraclient.OfferRequest{ForAccount: "GAQHWQYBBW272OOXNQMMLCA5WY2XAZPODGB7Q3S5OKKIXVESKO55ZQ7C", Limit: 20}
	offers, err := client.Offers(offerRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(offers)

	// get next pages.
	recordsFound := false
	if len(offers.Embedded.Records) > 0 {
		recordsFound = true
	}
	page := offers
	// get the next page of records if recordsFound is true
	for recordsFound {
		// next page
		nextPage, err := client.NextOffersPage(page)
		if err != nil {
			fmt.Println(err)
			return
		}

		page = nextPage
		if len(nextPage.Embedded.Records) == 0 {
			recordsFound = false
		}
		fmt.Println(nextPage)
	}
}
func ExampleClient_NextOperationsPage() {
	client := auroraclient.DefaultPublicNetClient
	// all operations
	operationRequest := auroraclient.OperationRequest{Limit: 20}
	ops, err := client.Operations(operationRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(ops)

	// get next pages.
	recordsFound := false
	if len(ops.Embedded.Records) > 0 {
		recordsFound = true
	}
	page := ops
	// get the next page of records if recordsFound is true
	for recordsFound {
		// next page
		nextPage, err := client.NextOperationsPage(page)
		if err != nil {
			fmt.Println(err)
			return
		}

		page = nextPage
		if len(nextPage.Embedded.Records) == 0 {
			recordsFound = false
		}
		fmt.Println(nextPage)
	}
}

func ExampleClient_NextTradeAggregationsPage() {
	client := auroraclient.DefaultPublicNetClient
	testTime := time.Unix(int64(1517521726), int64(0))
	// Find trade aggregations
	ta := auroraclient.TradeAggregationRequest{
		StartTime:          testTime,
		EndTime:            testTime,
		Resolution:         auroraclient.FiveMinuteResolution,
		BaseAssetType:      auroraclient.AssetTypeNative,
		CounterAssetType:   auroraclient.AssetType4,
		CounterAssetCode:   "SLT",
		CounterAssetIssuer: "GCKA6K5PCQ6PNF5RQBF7PQDJWRHO6UOGFMRLK3DYHDOI244V47XKQ4GP",
		Order:              auroraclient.OrderDesc,
	}
	tradeAggs, err := client.TradeAggregations(ta)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(tradeAggs)

	// get next pages.
	recordsFound := false
	if len(tradeAggs.Embedded.Records) > 0 {
		recordsFound = true
	}
	page := tradeAggs
	// get the next page of records if recordsFound is true
	for recordsFound {
		// next page
		nextPage, err := client.NextTradeAggregationsPage(page)
		if err != nil {
			fmt.Println(err)
			return
		}

		page = nextPage
		if len(nextPage.Embedded.Records) == 0 {
			recordsFound = false
		}
		fmt.Println(nextPage)
	}
}

func ExampleClient_NextTradesPage() {
	client := auroraclient.DefaultPublicNetClient
	// all trades
	tradeRequest := auroraclient.TradeRequest{Cursor: "123456", Limit: 30, Order: auroraclient.OrderAsc}
	trades, err := client.Trades(tradeRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(trades)

	// get next pages.
	recordsFound := false
	if len(trades.Embedded.Records) > 0 {
		recordsFound = true
	}
	page := trades
	// get the next page of records if recordsFound is true
	for recordsFound {
		// next page
		nextPage, err := client.NextTradesPage(page)
		if err != nil {
			fmt.Println(err)
			return
		}

		page = nextPage
		if len(nextPage.Embedded.Records) == 0 {
			recordsFound = false
		}
		fmt.Println(nextPage)
	}
}

func ExampleClient_NextTransactionsPage() {
	client := auroraclient.DefaultPublicNetClient
	// all transactions
	transactionRequest := auroraclient.TransactionRequest{Limit: 20}
	transactions, err := client.Transactions(transactionRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(transactions)

	// get next pages.
	recordsFound := false
	if len(transactions.Embedded.Records) > 0 {
		recordsFound = true
	}
	page := transactions
	// get the next page of records if recordsFound is true
	for recordsFound {
		// next page
		nextPage, err := client.NextTransactionsPage(page)
		if err != nil {
			fmt.Println(err)
			return
		}

		page = nextPage
		if len(nextPage.Embedded.Records) == 0 {
			recordsFound = false
		}
		fmt.Println(nextPage)
	}
}

func ExampleClient_OfferDetails() {
	client := auroraclient.DefaultPublicNetClient
	offer, err := client.OfferDetails("2")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print(offer)
}

func ExampleClient_Offers() {
	client := auroraclient.DefaultPublicNetClient
	offerRequest := auroraclient.OfferRequest{
		ForAccount: "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU",
		Cursor:     "now",
		Order:      auroraclient.OrderDesc,
	}
	offers, err := client.Offers(offerRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(offers)

	offerRequest = auroraclient.OfferRequest{
		Seller:  "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU",
		Selling: "COP:GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU",
		Buying:  "EUR:GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU",
		Order:   auroraclient.OrderDesc,
	}

	offers, err = client.Offers(offerRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(offers)
}

func ExampleClient_OperationDetail() {
	client := auroraclient.DefaultPublicNetClient
	opID := "123456"
	// operation details for an id
	ops, err := client.OperationDetail(opID)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(ops)
}

func ExampleClient_Operations() {
	client := auroraclient.DefaultPublicNetClient
	// operations for an account
	opRequest := auroraclient.OperationRequest{ForAccount: "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU"}
	ops, err := client.Operations(opRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(ops)

	// all operations
	opRequest = auroraclient.OperationRequest{Cursor: "now"}
	ops, err = client.Operations(opRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(ops)
	records := ops.Embedded.Records

	for _, value := range records {
		// prints the type
		fmt.Print(value.GetType())
		// for example if the type is change_trust
		c, ok := value.(operations.ChangeTrust)
		if ok {
			// access ChangeTrust fields
			fmt.Print(c.Trustee)
		}

	}
}

func ExampleClient_OrderBook() {
	client := auroraclient.DefaultPublicNetClient
	// orderbook for an asset pair, e.g XLM/NGN
	obRequest := auroraclient.OrderBookRequest{
		BuyingAssetType:    auroraclient.AssetTypeNative,
		SellingAssetCode:   "USD",
		SellingAssetType:   auroraclient.AssetType4,
		SellingAssetIssuer: "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU",
	}
	obs, err := client.OrderBook(obRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(obs)
}

func ExampleClient_Paths() {
	client := auroraclient.DefaultPublicNetClient
	// Find paths for XLM->NGN
	pr := auroraclient.PathsRequest{
		DestinationAccount:     "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU",
		DestinationAmount:      "100",
		DestinationAssetCode:   "NGN",
		DestinationAssetIssuer: "GDZST3XVCDTUJ76ZAV2HA72KYQODXXZ5PTMAPZGDHZ6CS7RO7MGG3DBM",
		DestinationAssetType:   auroraclient.AssetType4,
		SourceAccount:          "GDZST3XVCDTUJ76ZAV2HA72KYQODXXZ5PTMAPZGDHZ6CS7RO7MGG3DBM",
	}
	paths, err := client.StrictReceivePaths(pr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(paths)
}

func ExampleClient_StrictSendPaths() {
	client := auroraclient.DefaultPublicNetClient
	// Find paths for USD->EUR
	pr := auroraclient.StrictSendPathsRequest{
		SourceAmount:      "20",
		SourceAssetCode:   "USD",
		SourceAssetIssuer: "GDUKMGUGDZQK6YHYA5Z6AY2G4XDSZPSZ3SW5UN3ARVMO6QSRDWP5YLEX",
		SourceAssetType:   auroraclient.AssetType4,
		DestinationAssets: "EURT:GAP5LETOV6YIE62YAM56STDANPRDO7ZFDBGSNHJQIYGGKSMOZAHOOS2S",
	}
	paths, err := client.StrictSendPaths(pr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(paths)
}

func ExampleClient_Payments() {
	client := auroraclient.DefaultPublicNetClient
	// payments for an account
	opRequest := auroraclient.OperationRequest{ForAccount: "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU"}
	ops, err := client.Payments(opRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(ops)

	// all payments
	opRequest = auroraclient.OperationRequest{Cursor: "now"}
	ops, err = client.Payments(opRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(ops)
	records := ops.Embedded.Records

	for _, value := range records {
		// prints the type
		fmt.Print(value.GetType())
		// for example if the type is create_account
		c, ok := value.(operations.CreateAccount)
		if ok {
			// access create_account fields
			fmt.Print(c.StartingBalance)
		}

	}
}

func ExampleClient_PrevAssetsPage() {
	client := auroraclient.DefaultPublicNetClient
	// assets for asset issuer
	assetRequest := auroraclient.AssetRequest{ForAssetIssuer: "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU",
		Limit: 20}
	asset, err := client.Assets(assetRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(asset)

	// all assets
	assetRequest = auroraclient.AssetRequest{}
	asset, err = client.Assets(assetRequest)
	if err != nil {
		fmt.Println(err)
		return
	}

	// next page
	prevPage, err := client.PrevAssetsPage(asset)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(prevPage)
}

func ExampleClient_PrevEffectsPage() {
	client := auroraclient.DefaultPublicNetClient
	// all effects
	effectRequest := auroraclient.EffectRequest{Limit: 20}
	efp, err := client.Effects(effectRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(efp)

	// get prev pages.
	recordsFound := false
	if len(efp.Embedded.Records) > 0 {
		recordsFound = true
	}
	page := efp
	// get the prev page of records if recordsFound is true
	for recordsFound {
		// prev page
		prevPage, err := client.PrevEffectsPage(page)
		if err != nil {
			fmt.Println(err)
			return
		}

		page = prevPage
		if len(prevPage.Embedded.Records) == 0 {
			recordsFound = false
		}
		fmt.Println(prevPage)
	}
}

func ExampleClient_PrevLedgersPage() {
	client := auroraclient.DefaultPublicNetClient
	// all ledgers
	ledgerRequest := auroraclient.LedgerRequest{Limit: 20}
	ledgers, err := client.Ledgers(ledgerRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(ledgers)

	// get prev pages.
	recordsFound := false
	if len(ledgers.Embedded.Records) > 0 {
		recordsFound = true
	}
	page := ledgers
	// get the prev page of records if recordsFound is true
	for recordsFound {
		// prev page
		prevPage, err := client.PrevLedgersPage(page)
		if err != nil {
			fmt.Println(err)
			return
		}

		page = prevPage
		if len(prevPage.Embedded.Records) == 0 {
			recordsFound = false
		}
		fmt.Println(prevPage)
	}
}

func ExampleClient_PrevOffersPage() {
	client := auroraclient.DefaultPublicNetClient
	// all offers
	offerRequest := auroraclient.OfferRequest{ForAccount: "GAQHWQYBBW272OOXNQMMLCA5WY2XAZPODGB7Q3S5OKKIXVESKO55ZQ7C", Limit: 20}
	offers, err := client.Offers(offerRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(offers)

	// get prev pages.
	recordsFound := false
	if len(offers.Embedded.Records) > 0 {
		recordsFound = true
	}
	page := offers
	// get the prev page of records if recordsFound is true
	for recordsFound {
		// prev page
		prevPage, err := client.PrevOffersPage(page)
		if err != nil {
			fmt.Println(err)
			return
		}

		page = prevPage
		if len(prevPage.Embedded.Records) == 0 {
			recordsFound = false
		}
		fmt.Println(prevPage)
	}
}

func ExampleClient_PrevOperationsPage() {
	client := auroraclient.DefaultPublicNetClient
	// all operations
	operationRequest := auroraclient.OperationRequest{Limit: 20}
	ops, err := client.Operations(operationRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(ops)

	// get prev pages.
	recordsFound := false
	if len(ops.Embedded.Records) > 0 {
		recordsFound = true
	}
	page := ops
	// get the prev page of records if recordsFound is true
	for recordsFound {
		// prev page
		prevPage, err := client.PrevOperationsPage(page)
		if err != nil {
			fmt.Println(err)
			return
		}

		page = prevPage
		if len(prevPage.Embedded.Records) == 0 {
			recordsFound = false
		}
		fmt.Println(prevPage)
	}
}

func ExampleClient_PrevTradeAggregationsPage() {
	client := auroraclient.DefaultPublicNetClient
	testTime := time.Unix(int64(1517521726), int64(0))
	// Find trade aggregations
	ta := auroraclient.TradeAggregationRequest{
		StartTime:          testTime,
		EndTime:            testTime,
		Resolution:         auroraclient.FiveMinuteResolution,
		BaseAssetType:      auroraclient.AssetTypeNative,
		CounterAssetType:   auroraclient.AssetType4,
		CounterAssetCode:   "SLT",
		CounterAssetIssuer: "GCKA6K5PCQ6PNF5RQBF7PQDJWRHO6UOGFMRLK3DYHDOI244V47XKQ4GP",
		Order:              auroraclient.OrderDesc,
	}
	tradeAggs, err := client.TradeAggregations(ta)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(tradeAggs)

	// get prev pages.
	recordsFound := false
	if len(tradeAggs.Embedded.Records) > 0 {
		recordsFound = true
	}
	page := tradeAggs
	// get the prev page of records if recordsFound is true
	for recordsFound {
		// prev page
		prevPage, err := client.PrevTradeAggregationsPage(page)
		if err != nil {
			fmt.Println(err)
			return
		}

		page = prevPage
		if len(prevPage.Embedded.Records) == 0 {
			recordsFound = false
		}
		fmt.Println(prevPage)
	}
}

func ExampleClient_PrevTradesPage() {
	client := auroraclient.DefaultPublicNetClient
	// all trades
	tradeRequest := auroraclient.TradeRequest{Cursor: "123456", Limit: 30, Order: auroraclient.OrderAsc}
	trades, err := client.Trades(tradeRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(trades)

	// get prev pages.
	recordsFound := false
	if len(trades.Embedded.Records) > 0 {
		recordsFound = true
	}
	page := trades
	// get the prev page of records if recordsFound is true
	for recordsFound {
		// prev page
		prevPage, err := client.PrevTradesPage(page)
		if err != nil {
			fmt.Println(err)
			return
		}

		page = prevPage
		if len(prevPage.Embedded.Records) == 0 {
			recordsFound = false
		}
		fmt.Println(prevPage)
	}
}

func ExampleClient_PrevTransactionsPage() {
	client := auroraclient.DefaultPublicNetClient
	// all transactions
	transactionRequest := auroraclient.TransactionRequest{Limit: 20}
	transactions, err := client.Transactions(transactionRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(transactions)

	// get prev pages.
	recordsFound := false
	if len(transactions.Embedded.Records) > 0 {
		recordsFound = true
	}
	page := transactions
	// get the prev page of records if recordsFound is true
	for recordsFound {
		// prev page
		prevPage, err := client.PrevTransactionsPage(page)
		if err != nil {
			fmt.Println(err)
			return
		}

		page = prevPage
		if len(prevPage.Embedded.Records) == 0 {
			recordsFound = false
		}
		fmt.Println(prevPage)
	}
}

func ExampleClient_Root() {
	client := auroraclient.DefaultTestNetClient
	root, err := client.Root()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(root)
}

func ExampleClient_SetAuroraTimeout() {
	client := auroraclient.DefaultTestNetClient

	// https://www.diamcircle.org/laboratory/#xdr-viewer?input=AAAAABB90WssODNIgi6BHveqzxTRmIpvAFRyVNM%2BHm2GVuCcAAAAZAAABD0AAuV%2FAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAyTBGxOgfSApppsTnb%2FYRr6gOR8WT0LZNrhLh4y3FCgoAAAAXSHboAAAAAAAAAAABhlbgnAAAAEAivKe977CQCxMOKTuj%2BcWTFqc2OOJU8qGr9afrgu2zDmQaX5Q0cNshc3PiBwe0qw%2F%2BD%2FqJk5QqM5dYeSUGeDQP&type=TransactionEnvelope&network=test
	txXdr := `AAAAABB90WssODNIgi6BHveqzxTRmIpvAFRyVNM+Hm2GVuCcAAAAZAAABD0AAuV/AAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAyTBGxOgfSApppsTnb/YRr6gOR8WT0LZNrhLh4y3FCgoAAAAXSHboAAAAAAAAAAABhlbgnAAAAEAivKe977CQCxMOKTuj+cWTFqc2OOJU8qGr9afrgu2zDmQaX5Q0cNshc3PiBwe0qw/+D/qJk5QqM5dYeSUGeDQP`

	// test user timeout
	client = client.SetAuroraTimeout(30 * time.Second)
	resp, err := client.SubmitTransactionXDR(txXdr)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print(resp)
}

func ExampleClient_StreamEffects() {
	client := auroraclient.DefaultTestNetClient
	// all effects
	effectRequest := auroraclient.EffectRequest{Cursor: "760209215489"}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Stop streaming after 60 seconds.
		time.Sleep(60 * time.Second)
		cancel()
	}()

	printHandler := func(e effects.Effect) {
		fmt.Println(e)
	}
	err := client.StreamEffects(ctx, effectRequest, printHandler)
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleClient_StreamLedgers() {
	client := auroraclient.DefaultTestNetClient
	// all ledgers from now
	ledgerRequest := auroraclient.LedgerRequest{}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Stop streaming after 60 seconds.
		time.Sleep(60 * time.Second)
		cancel()
	}()

	printHandler := func(ledger hProtocol.Ledger) {
		fmt.Println(ledger)
	}
	err := client.StreamLedgers(ctx, ledgerRequest, printHandler)
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleClient_StreamOffers() {
	client := auroraclient.DefaultTestNetClient
	// offers for account
	offerRequest := auroraclient.OfferRequest{ForAccount: "GAQHWQYBBW272OOXNQMMLCA5WY2XAZPODGB7Q3S5OKKIXVESKO55ZQ7C", Cursor: "1"}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Stop streaming after 60 seconds.
		time.Sleep(60 * time.Second)
		cancel()
	}()

	printHandler := func(offer hProtocol.Offer) {
		fmt.Println(offer)
	}
	err := client.StreamOffers(ctx, offerRequest, printHandler)
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleClient_StreamOperations() {
	client := auroraclient.DefaultTestNetClient
	// operations for an account
	opRequest := auroraclient.OperationRequest{ForAccount: "GAIH3ULLFQ4DGSECF2AR555KZ4KNDGEKN4AFI4SU2M7B43MGK3QJZNSR", Cursor: "760209215489"}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Stop streaming after 60 seconds.
		time.Sleep(60 * time.Second)
		cancel()
	}()

	printHandler := func(op operations.Operation) {
		fmt.Println(op)
	}
	err := client.StreamOperations(ctx, opRequest, printHandler)
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleClient_StreamOrderBooks() {
	client := auroraclient.DefaultTestNetClient
	orderbookRequest := auroraclient.OrderBookRequest{
		SellingAssetType:  auroraclient.AssetTypeNative,
		BuyingAssetType:   auroraclient.AssetType4,
		BuyingAssetCode:   "ABC",
		BuyingAssetIssuer: "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU",
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Stop streaming after 60 seconds.
		time.Sleep(60 * time.Second)
		cancel()
	}()

	printHandler := func(orderbook hProtocol.OrderBookSummary) {
		fmt.Println(orderbook)
	}
	err := client.StreamOrderBooks(ctx, orderbookRequest, printHandler)
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleClient_StreamPayments() {
	client := auroraclient.DefaultTestNetClient
	// all payments
	opRequest := auroraclient.OperationRequest{Cursor: "760209215489"}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Stop streaming after 60 seconds.
		time.Sleep(60 * time.Second)
		cancel()
	}()

	printHandler := func(op operations.Operation) {
		fmt.Println(op)
	}
	err := client.StreamPayments(ctx, opRequest, printHandler)
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleClient_StreamTrades() {
	client := auroraclient.DefaultTestNetClient
	// all trades
	tradeRequest := auroraclient.TradeRequest{Cursor: "760209215489"}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Stop streaming after 60 seconds.
		time.Sleep(60 * time.Second)
		cancel()
	}()

	printHandler := func(tr hProtocol.Trade) {
		fmt.Println(tr)
	}
	err := client.StreamTrades(ctx, tradeRequest, printHandler)

	if err != nil {
		fmt.Println(err)
	}
}

func ExampleClient_StreamTransactions() {
	client := auroraclient.DefaultTestNetClient
	// all transactions
	transactionRequest := auroraclient.TransactionRequest{Cursor: "760209215489"}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Stop streaming after 60 seconds.
		time.Sleep(60 * time.Second)
		cancel()
	}()

	printHandler := func(tr hProtocol.Transaction) {
		fmt.Println(tr)
	}
	err := client.StreamTransactions(ctx, transactionRequest, printHandler)
	if err != nil {
		fmt.Println(err)
	}
}

func ExampleClient_SubmitFeeBumpTransaction() {
	kp := keypair.MustParseFull("SDQQUZMIPUP5TSDWH3UJYAKUOP55IJ4KTBXTY7RCOMEFRQGYA6GIR3OD")
	client := auroraclient.DefaultTestNetClient
	ar := auroraclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		return
	}

	op := txnbuild.Payment{
		Destination: kp.Address(),
		Amount:      "1",
		Asset:       txnbuild.NativeAsset{},
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []txnbuild.Operation{&op},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(), // Use a real timeout in production!
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	tx, err = tx.Sign(network.TestNetworkPassphrase, kp)
	if err != nil {
		fmt.Println(err)
		return
	}

	feeBumpKP := keypair.MustParseFull("SA5ZEFDVFZ52GRU7YUGR6EDPBNRU2WLA6IQFQ7S2IH2DG3VFV3DOMV2Q")
	feeBumpTx, err := txnbuild.NewFeeBumpTransaction(txnbuild.FeeBumpTransactionParams{
		Inner:      tx,
		FeeAccount: feeBumpKP.Address(),
		BaseFee:    txnbuild.MinBaseFee * 2,
	})
	feeBumpTx, err = feeBumpTx.Sign(network.TestNetworkPassphrase, feeBumpKP)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := client.SubmitFeeBumpTransaction(feeBumpTx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}

func ExampleClient_SubmitFeeBumpTransactionWithOptions() {
	kp := keypair.MustParseFull("SDQQUZMIPUP5TSDWH3UJYAKUOP55IJ4KTBXTY7RCOMEFRQGYA6GIR3OD")
	client := auroraclient.DefaultTestNetClient
	ar := auroraclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		return
	}

	op := txnbuild.Payment{
		Destination: kp.Address(),
		Amount:      "1",
		Asset:       txnbuild.NativeAsset{},
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []txnbuild.Operation{&op},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(), // Use a real timeout in production!
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	tx, err = tx.Sign(network.TestNetworkPassphrase, kp)
	if err != nil {
		fmt.Println(err)
		return
	}

	feeBumpKP := keypair.MustParseFull("SA5ZEFDVFZ52GRU7YUGR6EDPBNRU2WLA6IQFQ7S2IH2DG3VFV3DOMV2Q")
	feeBumpTx, err := txnbuild.NewFeeBumpTransaction(txnbuild.FeeBumpTransactionParams{
		Inner:      tx,
		FeeAccount: feeBumpKP.Address(),
		BaseFee:    txnbuild.MinBaseFee * 2,
	})
	feeBumpTx, err = feeBumpTx.Sign(network.TestNetworkPassphrase, feeBumpKP)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := client.SubmitFeeBumpTransactionWithOptions(
		feeBumpTx,
		auroraclient.SubmitTxOpts{SkipMemoRequiredCheck: true},
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}

func ExampleClient_SubmitTransaction() {
	kp := keypair.MustParseFull("SDQQUZMIPUP5TSDWH3UJYAKUOP55IJ4KTBXTY7RCOMEFRQGYA6GIR3OD")
	client := auroraclient.DefaultTestNetClient
	ar := auroraclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		return
	}

	op := txnbuild.Payment{
		Destination: kp.Address(),
		Amount:      "1",
		Asset:       txnbuild.NativeAsset{},
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []txnbuild.Operation{&op},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(), // Use a real timeout in production!
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	tx, err = tx.Sign(network.TestNetworkPassphrase, kp)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := client.SubmitTransaction(tx)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}

func ExampleClient_SubmitTransactionWithOptions() {
	kp := keypair.MustParseFull("SDQQUZMIPUP5TSDWH3UJYAKUOP55IJ4KTBXTY7RCOMEFRQGYA6GIR3OD")
	client := auroraclient.DefaultTestNetClient
	ar := auroraclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		return
	}

	op := txnbuild.Payment{
		Destination: kp.Address(),
		Amount:      "1",
		Asset:       txnbuild.NativeAsset{},
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []txnbuild.Operation{&op},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(), // Use a real timeout in production!
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	tx, err = tx.Sign(network.TestNetworkPassphrase, kp)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := client.SubmitTransactionWithOptions(tx, auroraclient.SubmitTxOpts{SkipMemoRequiredCheck: true})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}

func ExampleClient_SubmitTransactionWithOptions_skip_memo_required_check() {
	kp := keypair.MustParseFull("SDQQUZMIPUP5TSDWH3UJYAKUOP55IJ4KTBXTY7RCOMEFRQGYA6GIR3OD")
	client := auroraclient.DefaultTestNetClient
	ar := auroraclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		return
	}

	op := txnbuild.Payment{
		Destination: kp.Address(),
		Amount:      "1",
		Asset:       txnbuild.NativeAsset{},
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: false,
			Operations:           []txnbuild.Operation{&op},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(), // Use a real timeout in production!
		},
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	tx, err = tx.Sign(network.TestNetworkPassphrase, kp)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := client.SubmitTransactionWithOptions(tx, auroraclient.SubmitTxOpts{
		SkipMemoRequiredCheck: true,
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}

func ExampleClient_SubmitTransactionXDR() {
	client := auroraclient.DefaultPublicNetClient
	// https://www.diamcircle.org/laboratory/#xdr-viewer?input=AAAAAOoS%2F5V%2BBiCPXRiVcz8YsnkDdODufq%2Bg7xdqTdIXN8vyAAAE4gFiW0YAAALxAAAAAQAAAAAAAAAAAAAAAFyuBUcAAAABAAAABzIyMjgyNDUAAAAAAQAAAAEAAAAALhsY%2FFdAHXllTmb025DtCVBw06WDSQjq6I9NrCQHOV8AAAABAAAAAHT8zKV7bRQzuGTpk9AO3gjWJ9jVxBXTgguFORkxHVIKAAAAAAAAAAAAOnDwAAAAAAAAAAIkBzlfAAAAQPefqlsOvni6xX1g3AqddvOp1GOM88JYzayGZodbzTfV5toyhxZvL1ZggY3prFsvrereugEpj1kyPJ67z6gcRg0XN8vyAAAAQGwmoTssW49gaze8iQkz%2FUA2E2N%2BBOo%2B6v7YdOSsvIcZnMc37KmXH920nLosKpDLqkNChVztSZFcbVUlHhjbQgA%3D&type=TransactionEnvelope&network=public
	txXdr := `AAAAAOoS/5V+BiCPXRiVcz8YsnkDdODufq+g7xdqTdIXN8vyAAAE4gFiW0YAAALxAAAAAQAAAAAAAAAAAAAAAFyuBUcAAAABAAAABzIyMjgyNDUAAAAAAQAAAAEAAAAALhsY/FdAHXllTmb025DtCVBw06WDSQjq6I9NrCQHOV8AAAABAAAAAHT8zKV7bRQzuGTpk9AO3gjWJ9jVxBXTgguFORkxHVIKAAAAAAAAAAAAOnDwAAAAAAAAAAIkBzlfAAAAQPefqlsOvni6xX1g3AqddvOp1GOM88JYzayGZodbzTfV5toyhxZvL1ZggY3prFsvrereugEpj1kyPJ67z6gcRg0XN8vyAAAAQGwmoTssW49gaze8iQkz/UA2E2N+BOo+6v7YdOSsvIcZnMc37KmXH920nLosKpDLqkNChVztSZFcbVUlHhjbQgA=`

	// submit transaction
	resp, err := client.SubmitTransactionXDR(txXdr)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print(resp)
}

func ExampleClient_TradeAggregations() {
	client := auroraclient.DefaultPublicNetClient
	testTime := time.Unix(int64(1517521726), int64(0))
	// Find trade aggregations
	ta := auroraclient.TradeAggregationRequest{
		StartTime:          testTime,
		EndTime:            testTime,
		Resolution:         auroraclient.FiveMinuteResolution,
		BaseAssetType:      auroraclient.AssetTypeNative,
		CounterAssetType:   auroraclient.AssetType4,
		CounterAssetCode:   "SLT",
		CounterAssetIssuer: "GCKA6K5PCQ6PNF5RQBF7PQDJWRHO6UOGFMRLK3DYHDOI244V47XKQ4GP",
		Order:              auroraclient.OrderDesc,
	}
	tradeAggs, err := client.TradeAggregations(ta)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(tradeAggs)
}

func ExampleClient_Trades() {
	client := auroraclient.DefaultPublicNetClient
	// Find all trades
	tr := auroraclient.TradeRequest{Cursor: "123456", Limit: 30, Order: auroraclient.OrderAsc}
	trades, err := client.Trades(tr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(trades)
}

func ExampleClient_Transactions() {
	client := auroraclient.DefaultPublicNetClient
	// transactions for an account
	txRequest := auroraclient.TransactionRequest{ForAccount: "GCLWGQPMKXQSPF776IU33AH4PZNOOWNAWGGKVTBQMIC5IMKUNP3E6NVU"}
	txs, err := client.Transactions(txRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(txs)

	// all transactions
	txRequest = auroraclient.TransactionRequest{Cursor: "now", Order: auroraclient.OrderDesc}
	txs, err = client.Transactions(txRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print(txs)
	records := txs.Embedded.Records

	for _, tx := range records {
		fmt.Print(tx)
	}
}
