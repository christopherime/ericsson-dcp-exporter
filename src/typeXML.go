package main

import (
	"encoding/xml"
)

type HeaderObject struct {
	Text     string `xml:",chardata"`
	Security struct {
		Text           string              `xml:",chardata"`
		MustUnderstand string              `xml:"soapenv:mustUnderstand,attr"`
		Wsse           string              `xml:"xmlns:wsse,attr"`
		Wsu            string              `xml:"xmlns:wsu,attr"`
		UsernameToken  UsernameTokenObject `xml:"wsse:UsernameToken"`
	} `xml:"wsse:Security"`
}

type UsernameTokenObject struct {
	Text     string `xml:",chardata"`
	Username string `xml:"wsse:Username"`
	Password string `xml:"wsse:Password"`
}

type EnvelopeObjectEcho struct {
	XMLName xml.Name     `xml:"soapenv:Envelope"`
	Text    string       `xml:",chardata"`
	Soapenv string       `xml:"xmlns:soapenv,attr"`
	Apis    string       `xml:"xmlns:apis,attr"`
	Header  HeaderObject `xml:"soapenv:Header"`
	Body    struct {
		Text string `xml:",chardata"`
		Echo struct {
			Text    string `xml:",chardata"`
			Message string `xml:"message"`
		} `xml:"apis:Echo"`
	} `xml:"soapenv:Body"`
}

type EnvelopeRespEcho struct {
	XMLName xml.Name `xml:"Envelope"`
	Text    string   `xml:",chardata"`
	Env     string   `xml:"env,attr"`
	Header  string   `xml:"Header"`
	Body    struct {
		Text         string `xml:",chardata"`
		EchoResponse struct {
			Text    string `xml:",chardata"`
			Ns1     string `xml:"ns1,attr"`
			Message string `xml:"message"`
		} `xml:"EchoResponse"`
	} `xml:"Body"`
}
