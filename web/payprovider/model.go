package payprovider

import (
	"github.com/stripe/stripe-go/v71"
	"strings"
)

type Currency struct {
	Value       string
	ZeroDecimal bool
}

func (c Currency) CalculatePayment(paymentAmount float64, serviceFee float64) (int64, int64) {
	if c.ZeroDecimal {
		amount := int64(paymentAmount)
		commission := int64(paymentAmount * serviceFee)
		return amount, commission
	}
	return int64(paymentAmount * 100), int64(paymentAmount * 100 * serviceFee)
}

var m = make(map[string]Currency)

func AllCurrencies() map[string]Currency {
	if len(m) == 0 {
		m["aed"] = Currency{strings.ToUpper(string(stripe.CurrencyAED)), false}
		m["afn"] = Currency{strings.ToUpper(string(stripe.CurrencyAFN)), false}
		m["all"] = Currency{strings.ToUpper(string(stripe.CurrencyALL)), false}
		m["amd"] = Currency{strings.ToUpper(string(stripe.CurrencyAMD)), false}
		m["ang"] = Currency{strings.ToUpper(string(stripe.CurrencyANG)), false}
		m["aoa"] = Currency{strings.ToUpper(string(stripe.CurrencyAOA)), false}
		m["ars"] = Currency{strings.ToUpper(string(stripe.CurrencyARS)), false}
		m["aud"] = Currency{strings.ToUpper(string(stripe.CurrencyAUD)), false}
		m["awg"] = Currency{strings.ToUpper(string(stripe.CurrencyAWG)), false}
		m["azn"] = Currency{strings.ToUpper(string(stripe.CurrencyAZN)), false}
		m["bam"] = Currency{strings.ToUpper(string(stripe.CurrencyBAM)), false}
		m["bbd"] = Currency{strings.ToUpper(string(stripe.CurrencyBBD)), false}
		m["bdt"] = Currency{strings.ToUpper(string(stripe.CurrencyBDT)), false}
		m["bgn"] = Currency{strings.ToUpper(string(stripe.CurrencyBGN)), false}
		m["bif"] = Currency{strings.ToUpper(string(stripe.CurrencyBIF)), true}
		m["bmd"] = Currency{strings.ToUpper(string(stripe.CurrencyBMD)), false}
		m["bnd"] = Currency{strings.ToUpper(string(stripe.CurrencyBND)), false}
		m["bob"] = Currency{strings.ToUpper(string(stripe.CurrencyBOB)), false}
		m["brl"] = Currency{strings.ToUpper(string(stripe.CurrencyBRL)), false}
		m["bsd"] = Currency{strings.ToUpper(string(stripe.CurrencyBSD)), false}
		m["bwp"] = Currency{strings.ToUpper(string(stripe.CurrencyBWP)), false}
		m["bzd"] = Currency{strings.ToUpper(string(stripe.CurrencyBZD)), false}
		m["cad"] = Currency{strings.ToUpper(string(stripe.CurrencyCAD)), false}
		m["cdf"] = Currency{strings.ToUpper(string(stripe.CurrencyCDF)), false}
		m["chf"] = Currency{strings.ToUpper(string(stripe.CurrencyCHF)), false}
		m["clp"] = Currency{strings.ToUpper(string(stripe.CurrencyCLP)), true}
		m["cny"] = Currency{strings.ToUpper(string(stripe.CurrencyCNY)), false}
		m["cop"] = Currency{strings.ToUpper(string(stripe.CurrencyCOP)), false}
		m["crc"] = Currency{strings.ToUpper(string(stripe.CurrencyCRC)), false}
		m["cve"] = Currency{strings.ToUpper(string(stripe.CurrencyCVE)), false}
		m["czk"] = Currency{strings.ToUpper(string(stripe.CurrencyCZK)), false}
		m["djf"] = Currency{strings.ToUpper(string(stripe.CurrencyDJF)), true}
		m["dkk"] = Currency{strings.ToUpper(string(stripe.CurrencyDKK)), false}
		m["dop"] = Currency{strings.ToUpper(string(stripe.CurrencyDOP)), false}
		m["dzd"] = Currency{strings.ToUpper(string(stripe.CurrencyDZD)), false}
		m["eek"] = Currency{strings.ToUpper(string(stripe.CurrencyEEK)), false}
		m["egp"] = Currency{strings.ToUpper(string(stripe.CurrencyEGP)), false}
		m["etb"] = Currency{strings.ToUpper(string(stripe.CurrencyETB)), false}
		m["eur"] = Currency{strings.ToUpper(string(stripe.CurrencyEUR)), false}
		m["fjd"] = Currency{strings.ToUpper(string(stripe.CurrencyFJD)), false}
		m["fkp"] = Currency{strings.ToUpper(string(stripe.CurrencyFKP)), false}
		m["gbp"] = Currency{strings.ToUpper(string(stripe.CurrencyGBP)), false}
		m["gel"] = Currency{strings.ToUpper(string(stripe.CurrencyGEL)), false}
		m["gip"] = Currency{strings.ToUpper(string(stripe.CurrencyGIP)), false}
		m["gmd"] = Currency{strings.ToUpper(string(stripe.CurrencyGMD)), false}
		m["gnf"] = Currency{strings.ToUpper(string(stripe.CurrencyGNF)), true}
		m["gtq"] = Currency{strings.ToUpper(string(stripe.CurrencyGTQ)), false}
		m["gyd"] = Currency{strings.ToUpper(string(stripe.CurrencyGYD)), false}
		m["hkd"] = Currency{strings.ToUpper(string(stripe.CurrencyHKD)), false}
		m["hnl"] = Currency{strings.ToUpper(string(stripe.CurrencyHNL)), false}
		m["hrk"] = Currency{strings.ToUpper(string(stripe.CurrencyHRK)), false}
		m["htg"] = Currency{strings.ToUpper(string(stripe.CurrencyHTG)), false}
		m["huf"] = Currency{strings.ToUpper(string(stripe.CurrencyHUF)), false}
		m["idr"] = Currency{strings.ToUpper(string(stripe.CurrencyIDR)), false}
		m["ils"] = Currency{strings.ToUpper(string(stripe.CurrencyILS)), false}
		m["inr"] = Currency{strings.ToUpper(string(stripe.CurrencyINR)), false}
		m["isk"] = Currency{strings.ToUpper(string(stripe.CurrencyISK)), false}
		m["jmd"] = Currency{strings.ToUpper(string(stripe.CurrencyJMD)), false}
		m["jpy"] = Currency{strings.ToUpper(string(stripe.CurrencyJPY)), true}
		m["kes"] = Currency{strings.ToUpper(string(stripe.CurrencyKES)), false}
		m["kgs"] = Currency{strings.ToUpper(string(stripe.CurrencyKGS)), false}
		m["khr"] = Currency{strings.ToUpper(string(stripe.CurrencyKHR)), false}
		m["kmf"] = Currency{strings.ToUpper(string(stripe.CurrencyKMF)), true}
		m["krw"] = Currency{strings.ToUpper(string(stripe.CurrencyKRW)), true}
		m["kyd"] = Currency{strings.ToUpper(string(stripe.CurrencyKYD)), false}
		m["kzt"] = Currency{strings.ToUpper(string(stripe.CurrencyKZT)), false}
		m["lak"] = Currency{strings.ToUpper(string(stripe.CurrencyLAK)), false}
		m["lbp"] = Currency{strings.ToUpper(string(stripe.CurrencyLBP)), false}
		m["lkr"] = Currency{strings.ToUpper(string(stripe.CurrencyLKR)), false}
		m["lrd"] = Currency{strings.ToUpper(string(stripe.CurrencyLRD)), false}
		m["lsl"] = Currency{strings.ToUpper(string(stripe.CurrencyLSL)), false}
		m["ltl"] = Currency{strings.ToUpper(string(stripe.CurrencyLTL)), false}
		m["lvl"] = Currency{strings.ToUpper(string(stripe.CurrencyLVL)), false}
		m["mad"] = Currency{strings.ToUpper(string(stripe.CurrencyMAD)), false}
		m["mdl"] = Currency{strings.ToUpper(string(stripe.CurrencyMDL)), false}
		m["mga"] = Currency{strings.ToUpper(string(stripe.CurrencyMGA)), true}
		m["mkd"] = Currency{strings.ToUpper(string(stripe.CurrencyMKD)), false}
		m["mnt"] = Currency{strings.ToUpper(string(stripe.CurrencyMNT)), false}
		m["mop"] = Currency{strings.ToUpper(string(stripe.CurrencyMOP)), false}
		m["mro"] = Currency{strings.ToUpper(string(stripe.CurrencyMRO)), false}
		m["mur"] = Currency{strings.ToUpper(string(stripe.CurrencyMUR)), false}
		m["mvr"] = Currency{strings.ToUpper(string(stripe.CurrencyMVR)), false}
		m["mwk"] = Currency{strings.ToUpper(string(stripe.CurrencyMWK)), false}
		m["mxn"] = Currency{strings.ToUpper(string(stripe.CurrencyMXN)), false}
		m["myr"] = Currency{strings.ToUpper(string(stripe.CurrencyMYR)), false}
		m["mzn"] = Currency{strings.ToUpper(string(stripe.CurrencyMZN)), false}
		m["nad"] = Currency{strings.ToUpper(string(stripe.CurrencyNAD)), false}
		m["ngn"] = Currency{strings.ToUpper(string(stripe.CurrencyNGN)), false}
		m["nio"] = Currency{strings.ToUpper(string(stripe.CurrencyNIO)), false}
		m["nok"] = Currency{strings.ToUpper(string(stripe.CurrencyNOK)), false}
		m["npr"] = Currency{strings.ToUpper(string(stripe.CurrencyNPR)), false}
		m["nzd"] = Currency{strings.ToUpper(string(stripe.CurrencyNZD)), false}
		m["pab"] = Currency{strings.ToUpper(string(stripe.CurrencyPAB)), false}
		m["pen"] = Currency{strings.ToUpper(string(stripe.CurrencyPEN)), false}
		m["pgk"] = Currency{strings.ToUpper(string(stripe.CurrencyPGK)), false}
		m["php"] = Currency{strings.ToUpper(string(stripe.CurrencyPHP)), false}
		m["pkr"] = Currency{strings.ToUpper(string(stripe.CurrencyPKR)), false}
		m["pln"] = Currency{strings.ToUpper(string(stripe.CurrencyPLN)), false}
		m["pyg"] = Currency{strings.ToUpper(string(stripe.CurrencyPYG)), true}
		m["qar"] = Currency{strings.ToUpper(string(stripe.CurrencyQAR)), false}
		m["ron"] = Currency{strings.ToUpper(string(stripe.CurrencyRON)), false}
		m["rsd"] = Currency{strings.ToUpper(string(stripe.CurrencyRSD)), false}
		m["rub"] = Currency{strings.ToUpper(string(stripe.CurrencyRUB)), false}
		m["rwf"] = Currency{strings.ToUpper(string(stripe.CurrencyRWF)), true}
		m["sar"] = Currency{strings.ToUpper(string(stripe.CurrencySAR)), false}
		m["sbd"] = Currency{strings.ToUpper(string(stripe.CurrencySBD)), false}
		m["scr"] = Currency{strings.ToUpper(string(stripe.CurrencySCR)), false}
		m["sek"] = Currency{strings.ToUpper(string(stripe.CurrencySEK)), false}
		m["sgd"] = Currency{strings.ToUpper(string(stripe.CurrencySGD)), false}
		m["shp"] = Currency{strings.ToUpper(string(stripe.CurrencySHP)), false}
		m["sll"] = Currency{strings.ToUpper(string(stripe.CurrencySLL)), false}
		m["sos"] = Currency{strings.ToUpper(string(stripe.CurrencySOS)), false}
		m["srd"] = Currency{strings.ToUpper(string(stripe.CurrencySRD)), false}
		m["std"] = Currency{strings.ToUpper(string(stripe.CurrencySTD)), false}
		m["svc"] = Currency{strings.ToUpper(string(stripe.CurrencySVC)), false}
		m["szl"] = Currency{strings.ToUpper(string(stripe.CurrencySZL)), false}
		m["thb"] = Currency{strings.ToUpper(string(stripe.CurrencyTHB)), false}
		m["tjs"] = Currency{strings.ToUpper(string(stripe.CurrencyTJS)), false}
		m["top"] = Currency{strings.ToUpper(string(stripe.CurrencyTOP)), false}
		m["try"] = Currency{strings.ToUpper(string(stripe.CurrencyTRY)), false}
		m["ttd"] = Currency{strings.ToUpper(string(stripe.CurrencyTTD)), false}
		m["twd"] = Currency{strings.ToUpper(string(stripe.CurrencyTWD)), false}
		m["tzs"] = Currency{strings.ToUpper(string(stripe.CurrencyTZS)), false}
		m["uah"] = Currency{strings.ToUpper(string(stripe.CurrencyUAH)), false}
		m["ugx"] = Currency{strings.ToUpper(string(stripe.CurrencyUGX)), true}
		m["usd"] = Currency{strings.ToUpper(string(stripe.CurrencyUSD)), false}
		m["uyu"] = Currency{strings.ToUpper(string(stripe.CurrencyUYU)), false}
		m["uzs"] = Currency{strings.ToUpper(string(stripe.CurrencyUZS)), false}
		m["vef"] = Currency{strings.ToUpper(string(stripe.CurrencyVEF)), false}
		m["vnd"] = Currency{strings.ToUpper(string(stripe.CurrencyVND)), true}
		m["vuv"] = Currency{strings.ToUpper(string(stripe.CurrencyVUV)), true}
		m["wst"] = Currency{strings.ToUpper(string(stripe.CurrencyWST)), false}
		m["xaf"] = Currency{strings.ToUpper(string(stripe.CurrencyXAF)), true}
		m["xcd"] = Currency{strings.ToUpper(string(stripe.CurrencyXCD)), false}
		m["xof"] = Currency{strings.ToUpper(string(stripe.CurrencyXOF)), true}
		m["xpf"] = Currency{strings.ToUpper(string(stripe.CurrencyXPF)), true}
		m["yer"] = Currency{strings.ToUpper(string(stripe.CurrencyYER)), false}
		m["zar"] = Currency{strings.ToUpper(string(stripe.CurrencyZAR)), false}
		m["zmw"] = Currency{strings.ToUpper(string(stripe.CurrencyZMW)), false}
	}
	return m
}
