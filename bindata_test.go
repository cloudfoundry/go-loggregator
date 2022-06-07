// Code generated by go-bindata.
// sources:
// test-certs/CA.crl
// test-certs/CA.crt
// test-certs/CA.key
// test-certs/client.crt
// test-certs/client.key
// test-certs/invalid-ca.crt
// test-certs/metron.csr
// test-certs/reverselogproxy.csr
// test-certs/server.crt
// test-certs/server.key
// DO NOT EDIT!

package loggregator_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _caCrl = []byte(`-----BEGIN X509 CRL-----
MIIChzBxAgEBMA0GCSqGSIb3DQEBCwUAMBgxFjAUBgNVBAMTDWxvZ2dyZWdhdG9y
Q0EXDTIxMDYxNzIxNTU1NloXDTQ2MDYxNzIxNTU1NVowAKAjMCEwHwYDVR0jBBgw
FoAUX/2QoNqcyZrxb4ktjdoI/p1IsQ8wDQYJKoZIhvcNAQELBQADggIBAE+g2muh
hzWQMptI9Apsnn7m4qUq/y6vtIFiJwt0vQuwb7zOYNiDMVOtSsnzwo/kp6rYzxOj
bsDafUUcHpmQbKl1YlIP/XjJi/SoMGIilXhQCboARTvJTjtpaAH4OFWSOzKQ3Hs8
ggO9iFgTmXXNiqeUH99SKemTcDrduFC9LXAzBH4V6oVyoDssiqzUgpk24EgHP3W0
fC0TuFnT8QBOurv4dTTBk7aVcZ8Yxx5XIQizFDDAmpLBKIuCAXMCSjrYmQuqwtxz
I1CIlvPumal00BqcLGDZTodQoy2bgTJC0KoxBJYAtZLWxwLshVgABN53Vq5vcm5G
KZMhBpxebm+xzVGvIHg1n/ZqeydBGe8I6A56w24TY0lOGZuXYE9rW29EvduRRvXX
8MEiVMKDX3ZrvrxJlRzhFSnxekVdqFCAp1o6OTZ0NeOZbQdFWDzaq0kFpcTqk9/Y
L/924VWuGnbaFAbXPx361HL+byN7CCPBVuWdVxGx2RzR1iW7AVIGVztIHRSD/iwg
qD3Ak9ded0er9XkMqjQpV/VLKqr8GDad1UgZdcqPxTfwvvFRU/XuPv4OlylY+hNt
9S2T+4t5BLRtefWu7n5VMVgYJXqu4vLegFppM+SDOGZQDxhMeQsEhNMndTUy0QBL
GLhpE2YAsO1aBAEm0ObHRizTzbLPmGkqlBSh
-----END X509 CRL-----
`)

func caCrlBytes() ([]byte, error) {
	return _caCrl, nil
}

func caCrl() (*asset, error) {
	bytes, err := caCrlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "CA.crl", size: 930, mode: os.FileMode(292), modTime: time.Unix(1623966956, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _caCrt = []byte(`-----BEGIN CERTIFICATE-----
MIIE8DCCAtigAwIBAgIBATANBgkqhkiG9w0BAQsFADAYMRYwFAYDVQQDEw1sb2dn
cmVnYXRvckNBMB4XDTIxMDYxNzIxNTU1NloXDTQ2MDYxNzIxNTU1NVowGDEWMBQG
A1UEAxMNbG9nZ3JlZ2F0b3JDQTCCAiIwDQYJKoZIhvcNAQEBBQADggIPADCCAgoC
ggIBALqfmAQDG6pe09+iFr7q623Cp0SZj/STs8EvkAuSehWh8MHProhYjiZTxXAb
MTsUL/4Ap8PXCprC9Cs6wsz1kzeZMWmLHfxUkIjvgeNMe0osO0uWD5dj6r4rrO0+
ArNSLui/kH3msWwc3upMAJBE+PTiGHtd1T8etqaEz8Rp7bIHIlf6zr0scQ+UFj1J
MJUM/Q2xvCMpTfT27gbWug9QS/dUstQEvfUeoc2rzvUmwMk1Bq3w7l00GV/aHcxV
SU0v77fyw2PdXhf0Jv/xH5x78QeslctWKpMa7fMkVkCIyka2mZaaVeo8yeSrzwHy
MadwqpHP+WikknUWxz2UxARfkJ6eJ/XpewdtGguthfV7klgvzdCWO8nsaMHffrhz
xO4+ahfiTr+iTMl8iTx5Xvo+cH3mudlTuuAA12xXsOnXnRjTTwQ/ewh+wWpgq1Va
5pg5IpphhPmXiQZnSdshIrB0ghDjIdOvKG3Ga5E4UvdSuIBeBMHT3VLD5rpuWLp1
fNNqGMqJOA7Tq/c4GNgA/ehefAw/MDXp7duyrySyLgolnGZl1kCLk2Vo70DlF3nS
kSA0PDN+kRu11Ejv+1bUSVStd1eVfIXwyiHPXr8kjx97tyiUah2U8m+9QIX445Kl
gJuWbI7kcn507Nc40ait+9bxBimMv4x+IK7B3XfLCFg3luYxAgMBAAGjRTBDMA4G
A1UdDwEB/wQEAwIBBjASBgNVHRMBAf8ECDAGAQH/AgEAMB0GA1UdDgQWBBRf/ZCg
2pzJmvFviS2N2gj+nUixDzANBgkqhkiG9w0BAQsFAAOCAgEAZmwkW8PvDCGSdxV3
PntWk1lB9xOXAnZC0Bpm3G5h0+nxP0T/mWNJL9OpFGVoadFtTD4jECmWOHR83dLr
oiW0nZdNw0bDX+uJF5jP6xeXotvT706nDs0Ye6J2oRprYADOBzKttaFn1YHfHiF5
uz7VvrUpk7X8wimRqPWJnc0apqsgCbb4tkTOk/joSuavtw6ep0tfurnLEklhq7vd
dxeG+S0vmjyLHqinmhNvEaRMbfMynoHb9ajEeI4Eyy/4gFkqRG9ZfWvV0NrQOGeA
cX/5L9Nn7HZX2TK63+phQY3/BEGLUwMo/7qf9CwaKLTWKCW+4rm8s2Bx2DoNaXfh
lC9jmNCh5AB+LnJp8FP3QtHvrr4RmqxkjdPCj+qg4Sw5F0Jdv5iv8fPo7H1OXRvx
aOaTg3HQhitHTQC7fbeiwWjaT+dUoM9pztcijeQx/mSoOm/+wYfJEckgwZyuVSSj
RPaWg5CNQC2IeJotpXg1Gv0YguU6s1EIyRj2jiLRFHpRxtYu2KNqS9cs/uJPae0A
usvDlyqw/D/dhRl+JH8rMK8Zxj1vhPD2so0aYU8NTdmeL3tXjRTK0VVEK6nVOvNH
gQcIxHigQjm00Xc0OAVhdUxzVgwbI8ADa0v+6bZjqjlKGytNogkmUwbkZwf5t1/k
LybgDAH2yTpnRmw/PTTVG4ssZcM=
-----END CERTIFICATE-----
`)

func caCrtBytes() ([]byte, error) {
	return _caCrt, nil
}

func caCrt() (*asset, error) {
	bytes, err := caCrtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "CA.crt", size: 1773, mode: os.FileMode(292), modTime: time.Unix(1623966956, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _caKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----
`)

func caKeyBytes() ([]byte, error) {
	return _caKey, nil
}

func caKey() (*asset, error) {
	bytes, err := caKeyBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "CA.key", size: 3243, mode: os.FileMode(288), modTime: time.Unix(1623966956, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _clientCrt = []byte(`-----BEGIN CERTIFICATE-----
MIIETDCCAjSgAwIBAgIRAL5EcrmG6+xxPFSHMCpnfTYwDQYJKoZIhvcNAQELBQAw
GDEWMBQGA1UEAxMNbG9nZ3JlZ2F0b3JDQTAeFw0yMTA2MTcyMTU1NTZaFw00NjA2
MTcyMTU1NTRaMBoxGDAWBgNVBAMTD3JldmVyc2Vsb2dwcm94eTCCASIwDQYJKoZI
hvcNAQEBBQADggEPADCCAQoCggEBAKOth4D/CFwI5tTgp0w+bp4Z4UTCzBopkdc6
5YVQI6tNd5xwP8+hcxdG9Cy183/P8l/8rVxitf9RFnf1v29u2+3lrljBfpfKBhBw
+NPlxltgok/HvJVWKXO9kSWpPPzBYkVxJ5zlz8CkDOalkMoM4ctkel85SwnQ+kpw
BlUi0eUjFOZftxQ0/lqDALhu9xPIA6cMC5Kn2AEi2kzRnsmGS4GJN6ZMGugubv9z
DC3xkrxx5qrg85miRSgbhAVib4Kz1ZsbpkLZadwg1qt0ASVGgE3fXMNnEQ54A87Y
fPQgkGkF+yIXovRhKQ2wmmmoKtT8fXW7LK+QCg1C6vuat6+XWScCAwEAAaOBjjCB
izAOBgNVHQ8BAf8EBAMCA7gwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMC
MB0GA1UdDgQWBBTEZS2dkL3/e80NA+HYhwYKwu0N6zAfBgNVHSMEGDAWgBRf/ZCg
2pzJmvFviS2N2gj+nUixDzAaBgNVHREEEzARgg9yZXZlcnNlbG9ncHJveHkwDQYJ
KoZIhvcNAQELBQADggIBAG3wAkX24yiHPAULT/buibbLawXOmXVfvQuY8eeENNKG
gpwo9nj6+QeM5s8m5vS7IR5E2sy4DSHKoiqOrRTbP56IcgML/4NeJG2iyMeSdywA
clbnVgeDq61qQvdtss7WLznGEJOm/A7MJMfn0apocVVp9EOWLGQbjqqWz3xNL0di
u7eKb+1h0q9j0XN4zghtjw9cpU1+3x8gehTP357T8yX4mP89UOtKsguFlRFyv4Zm
yrHzcnVwtpAePyOS+zrHS6ch0uwJ7PTQJfM283DrD3nPhHuuPm+KFt/9k2C2TC9e
NSm17MPD+X60Qc/Vkc9kvBi8Cb2NGYeMDTkyK5Ew6j/63l2Py6tIHLSNaNN4rwSB
YWIenHYtxiQnHtcY6Zskz6uqUnjbFsYfX4vAe2ft9dJfFPm0C3zje62OROClgVi3
i3eqMzrPANAxfYXioEOdIwp74QNstiAEWrrFqwYjvdBCsD0Bmr++7nRHJZnAzU50
3jXHbMeqe2fJ/H6Fz0e/wRswaEnkvHiZSksQhcsV1PCJXQrWHRKG4+eKks39a2ho
R11GcxbZuh1dqLJZDan57aREq+aa8IurUDpq265rkzqGldESAJh4R9QLo2ItRtJh
PtnI7Bi0PiLgy4asjpXWvqWzSIlvcEsQSNYcZSbztgPR/ycSyWV6mqgjp40vlgPu
-----END CERTIFICATE-----
`)

func clientCrtBytes() ([]byte, error) {
	return _clientCrt, nil
}

func clientCrt() (*asset, error) {
	bytes, err := clientCrtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "client.crt", size: 1549, mode: os.FileMode(292), modTime: time.Unix(1623966956, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _clientKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAo62HgP8IXAjm1OCnTD5unhnhRMLMGimR1zrlhVAjq013nHA/
z6FzF0b0LLXzf8/yX/ytXGK1/1EWd/W/b27b7eWuWMF+l8oGEHD40+XGW2CiT8e8
lVYpc72RJak8/MFiRXEnnOXPwKQM5qWQygzhy2R6XzlLCdD6SnAGVSLR5SMU5l+3
FDT+WoMAuG73E8gDpwwLkqfYASLaTNGeyYZLgYk3pkwa6C5u/3MMLfGSvHHmquDz
maJFKBuEBWJvgrPVmxumQtlp3CDWq3QBJUaATd9cw2cRDngDzth89CCQaQX7Ihei
9GEpDbCaaagq1Px9dbssr5AKDULq+5q3r5dZJwIDAQABAoIBAQCc/A2P6ouRAjSr
DkFHPkYIO1g1BASQPziSzvlecKuVXDqRf5OkR/xD4hjFXUoLV13nNEjt5/sYwyQi
sEAI17H2rSkYFZWebfR9Bm2RhrtF3xwuGqtELByb1sCE95n37l6brdjJgh1Nbcq1
2SxSriJOWrOH60KOSrIUfPeF5lycQA27B72POL3B3+lkXQo3Jo53xau96wtiSF8O
ePaKqt75YkkrQyuLbHOvZaxZz1fNoZSRkKcdMX+HoX36K4MjpNwI7mxnAEXlKjTS
mM/+nz+ri/pRMTyxzOE93H/9doTBaRroq6gmPNIEj9i7lpa4xADexeA2AoWST05u
+yvSSYSZAoGBAMduehfpmmzCMwiENl6DA3/wjx9mdoJ3YbIX0eEWR4CC8W9bzKBC
EhOGWoC0QSZaSrXdWW7hnuEBzsZvlpkOCz4J9/hh7IG/113OYvh9uU7y4Q+F8S5F
RlqfdPsGaz2Y1l4hpXFmqXbCU8OybPXpUcFuYnom+KVKVlvmn5fmIKMjAoGBANIa
1UBCdOI/ukZoVYSnMnax7hXkLvuhkJJuyoGXmAMDkqLif89NZwrvbo3MhY80lfOz
+/iz6/Y+Ac/cYxHDItlJFSiWo9imCva8gJ7MWQQJ97uijkcoidoL+jD1aGN7tK41
wHbSt84tyXg5+yljo36hCD+qo/rlOJgIXHSJ62QtAoGATrYRxiJS6p3zGmdkNgUr
enFta407NN79VhcEpPvwGI6Vz8fBiXbKP56FVcrO894rIoBMbfDqjg/ylmswWxQp
58BzeDrd54/Z5pIwibbFTp2ZqlDJEeQRkm7g8rpj1RnfcaOB8rH8LH0iJljjnGML
+3Yfs+pxtHsUfo0VbBRNyVUCgYAQXBSEx2fwggPQHamjBZ3RTjN6suTRpRfrvwK3
qoUknu+ZDUfkbWN7n38dPXKc3vxaGIajK+dQqi1b8Q4pwOcCwkUKfwhNA0jRQ4ZE
VycLQHdwvcyUT9zEBLC7hTBWprg/5GGTHv8+56PLX8BlzaNaZdGNm4zfKWAJvoTs
chzJFQKBgEly6GRIZtWb7UHSMf8izXa6uPeQsp5VXd2PSi0ccd9pffRACMOKnytm
iOmSCLvtA/5+PpCjxRiFUXCOxDDV59XtfHvVU1ZEcJukpRpIAlSu+vjUylgSYMeV
69Dne5L27XtBQ1xvN0cLG7p8Jom7urH4925dpcnhxq5KVSQ9SjUH
-----END RSA PRIVATE KEY-----
`)

func clientKeyBytes() ([]byte, error) {
	return _clientKey, nil
}

func clientKey() (*asset, error) {
	bytes, err := clientKeyBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "client.key", size: 1675, mode: os.FileMode(288), modTime: time.Unix(1623966956, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _invalidCaCrt = []byte(`foobar
`)

func invalidCaCrtBytes() ([]byte, error) {
	return _invalidCaCrt, nil
}

func invalidCaCrt() (*asset, error) {
	bytes, err := invalidCaCrtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "invalid-ca.crt", size: 7, mode: os.FileMode(436), modTime: time.Unix(1623966956, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _metronCsr = []byte(`-----BEGIN CERTIFICATE REQUEST-----
MIICejCCAWICAQAwETEPMA0GA1UEAxMGbWV0cm9uMIIBIjANBgkqhkiG9w0BAQEF
AAOCAQ8AMIIBCgKCAQEA2jF5+Bl4Ayh+xdef50a81gJPksWzEtj907zyJ1SUS2ts
O/SgeFcBoE8ezgFSBOzCiwXzRmy+8S4uvfOhKG0xfygHO7YQJFKnptGmbEyNVpDW
Nhz61YrVh5TK173nu38TgtZqzKakUfxbUr/0JSlV9XXSNf5m1K3SK5jZ2jE14v1l
0O7Nk6pGqI2gbqfkQ42FhsJrWZ2A4TPtMnHEAxRKcSTEJ3HbrGWEbBAxjPe+t1RA
LoiEHCGNx9jsL4Kq+04mCxtGPwNIBgxD8ya8Jwv8VHonHCWupoFvdKsl/zchkG+e
ZCEi9fSLVqrN/eJYo2+R2o/tqWPXsLbRnng3hrcZUQIDAQABoCQwIgYJKoZIhvcN
AQkOMRUwEzARBgNVHREECjAIggZtZXRyb24wDQYJKoZIhvcNAQELBQADggEBANT5
RYbf/OCDlcIo6NIDzFG3Y0S1sQcTFg0IwCaUYgubTfjWOZ3dc7xpqMgoTUXktRzw
GThO0TEnoH628m/3nP5rll/KVQTv3DU3p79l00WlQbaCJvXJ3e2Nz74jumcLuEji
n4g37CBadmCwUymK0+yMfm+94bjBJ3Hdt1XBJslVrBlhPTWtcpIfKKIT2ckCeZNc
FGT2tgob0WG/ys1MQ+bmoezwUb1TZLAcDeKQukUTx395WHL95A6JQRK2QCRm0ZZS
XPQsr4lMFYMztJVc3b4TRxE7o72pV3imnoc+QHX5xfS2eiAtg6ZPKhjCgAHKmgDy
Y+CWzgW9kggs5mfVHOM=
-----END CERTIFICATE REQUEST-----
`)

func metronCsrBytes() ([]byte, error) {
	return _metronCsr, nil
}

func metronCsr() (*asset, error) {
	bytes, err := metronCsrBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "metron.csr", size: 936, mode: os.FileMode(292), modTime: time.Unix(1623966956, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _reverselogproxyCsr = []byte(`-----BEGIN CERTIFICATE REQUEST-----
MIICjDCCAXQCAQAwGjEYMBYGA1UEAxMPcmV2ZXJzZWxvZ3Byb3h5MIIBIjANBgkq
hkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAo62HgP8IXAjm1OCnTD5unhnhRMLMGimR
1zrlhVAjq013nHA/z6FzF0b0LLXzf8/yX/ytXGK1/1EWd/W/b27b7eWuWMF+l8oG
EHD40+XGW2CiT8e8lVYpc72RJak8/MFiRXEnnOXPwKQM5qWQygzhy2R6XzlLCdD6
SnAGVSLR5SMU5l+3FDT+WoMAuG73E8gDpwwLkqfYASLaTNGeyYZLgYk3pkwa6C5u
/3MMLfGSvHHmquDzmaJFKBuEBWJvgrPVmxumQtlp3CDWq3QBJUaATd9cw2cRDngD
zth89CCQaQX7Ihei9GEpDbCaaagq1Px9dbssr5AKDULq+5q3r5dZJwIDAQABoC0w
KwYJKoZIhvcNAQkOMR4wHDAaBgNVHREEEzARgg9yZXZlcnNlbG9ncHJveHkwDQYJ
KoZIhvcNAQELBQADggEBAJ1E0OdsYteb7B44mdtzLDaErh1voCVQTQeGvhQNIPh4
ODBANBWc26VO/oWvKMN+8SddZ0UmhOCYT6CLAmWzb+U4K5/1X+9tQ5Pv9m5pyDJU
60Wiq8Pp4Oxj5zVX18StDHgZRF1kHoLEw8DF5WoD9zT+pl0rNDNPQ3FXME1ndlkD
266Qj7ubADh8CY1BX61CQtMncg30hawaKRzpG03D2cnuUiwfrOwVTFn9CBptbhK4
gy4vuWw0grPqFlG91J5mbtftlpgrpkwI3jznAKpCV0fRHAGzcagf1y59MLypJQB4
oMilDqr0XVgo8jGY7tPhDTHHaz0cCinuDXY4Jwxv9Vo=
-----END CERTIFICATE REQUEST-----
`)

func reverselogproxyCsrBytes() ([]byte, error) {
	return _reverselogproxyCsr, nil
}

func reverselogproxyCsr() (*asset, error) {
	bytes, err := reverselogproxyCsrBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "reverselogproxy.csr", size: 960, mode: os.FileMode(292), modTime: time.Unix(1623966956, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _serverCrt = []byte(`-----BEGIN CERTIFICATE-----
MIIEOjCCAiKgAwIBAgIRAOe16YJ7n4JN2Lk+qx+bd+8wDQYJKoZIhvcNAQELBQAw
GDEWMBQGA1UEAxMNbG9nZ3JlZ2F0b3JDQTAeFw0yMTA2MTcyMTU1NTZaFw00NjA2
MTcyMTU1NTRaMBExDzANBgNVBAMTBm1ldHJvbjCCASIwDQYJKoZIhvcNAQEBBQAD
ggEPADCCAQoCggEBANoxefgZeAMofsXXn+dGvNYCT5LFsxLY/dO88idUlEtrbDv0
oHhXAaBPHs4BUgTswosF80ZsvvEuLr3zoShtMX8oBzu2ECRSp6bRpmxMjVaQ1jYc
+tWK1YeUyte957t/E4LWasympFH8W1K/9CUpVfV10jX+ZtSt0iuY2doxNeL9ZdDu
zZOqRqiNoG6n5EONhYbCa1mdgOEz7TJxxAMUSnEkxCdx26xlhGwQMYz3vrdUQC6I
hBwhjcfY7C+CqvtOJgsbRj8DSAYMQ/MmvCcL/FR6JxwlrqaBb3SrJf83IZBvnmQh
IvX0i1aqzf3iWKNvkdqP7alj17C20Z54N4a3GVECAwEAAaOBhTCBgjAOBgNVHQ8B
Af8EBAMCA7gwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMB0GA1UdDgQW
BBRj6WW2oRJ5KMaM2FiOqscbWH373TAfBgNVHSMEGDAWgBRf/ZCg2pzJmvFviS2N
2gj+nUixDzARBgNVHREECjAIggZtZXRyb24wDQYJKoZIhvcNAQELBQADggIBALZR
9fJ1RaOgTLZmo7u9gavi7Kl1dwci3z0MQtEiFQLYvEG5rwQVF9VU6Lsrl9IGx/lA
JEFXYX5JAcnvreFTSOl8R2nDXNhVK01yWqk4RG0LnGFeld21Mj2Sxnc6OSy61p2k
Tk6q5LPqqpzzI9KNmAk4RPVeQZsr+qKuu04nr+bjZKX1zjzE+SHnGOIyYj9CsRpG
4AVJFZkQsx124XYMQ2gRGjnqwP73+2mO/nVKR3op7Fl9TN8HB0CJtL3+Rs5aCys2
B2DiRUNfzsGxV6fLnFDM8n0fHl729F6HdFoyUCFWtq08XZ4jhNR+0ufWG3YvvIcv
ZtEdjNzE367/Z843EPlaJGEGOhMkGZzmKqutO5z7NCYFEOBsoN3WF9EMD61IBVkB
/+ljtzPTaet6hTD32po4+bykv4H2DaJ/1oGNDwum7LmaXo3KsrCIstZOTW9Sj3rV
CaX/9w2+JeEXHXI/lmowFAvV17XNtBEyiCRJj0vToo4L6ROnosJN5HRmlc04rC0D
KFStJ3RsB8ghSA6kjiYQcugGNWx4vS5W3FC1AptUdqtd/qV92hKZCpMquHTvKQkL
sBjBTsgHA2sFf6+ybeuKgFYdlmeTycZyMQekeiT8Yn131XHK7TPl0CXx4P7iQ2gm
wFHQbuZKItL28i0keRrtsQAJjiNTlfskahfnD5us
-----END CERTIFICATE-----
`)

func serverCrtBytes() ([]byte, error) {
	return _serverCrt, nil
}

func serverCrt() (*asset, error) {
	bytes, err := serverCrtBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "server.crt", size: 1525, mode: os.FileMode(292), modTime: time.Unix(1623966956, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _serverKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEA2jF5+Bl4Ayh+xdef50a81gJPksWzEtj907zyJ1SUS2tsO/Sg
eFcBoE8ezgFSBOzCiwXzRmy+8S4uvfOhKG0xfygHO7YQJFKnptGmbEyNVpDWNhz6
1YrVh5TK173nu38TgtZqzKakUfxbUr/0JSlV9XXSNf5m1K3SK5jZ2jE14v1l0O7N
k6pGqI2gbqfkQ42FhsJrWZ2A4TPtMnHEAxRKcSTEJ3HbrGWEbBAxjPe+t1RALoiE
HCGNx9jsL4Kq+04mCxtGPwNIBgxD8ya8Jwv8VHonHCWupoFvdKsl/zchkG+eZCEi
9fSLVqrN/eJYo2+R2o/tqWPXsLbRnng3hrcZUQIDAQABAoIBABNR7DXes3S8IjWM
eTk2V0Qv4jHh1ZBmrFsSUPLQl2zyLxxD9e2R7q/uMZEWJtgOys8akMb8nA+pAjSJ
nQyCVq6msbVE0rsUoomv6KeQQ7eVjZIvELrV10BxGWDvDNTaVLTyHXTPkJ891CxI
mOqtdVINw7ZKN3q/hWuc1jVuxBG8KlhdhCs8+TNgWKOSwUaY2At3y9ub9/HvtIaf
au2skcHfsOEf8wQLacMVYQrZhyar6uXA8qEUijGJmIAhYrPOqaVa/G+pHKiZwg6P
+AxDVmbLwJVirYQfE/PgL2EhhaoR4e3C4KskZWCji4DQKc7cHvMQ1C8m9Ig1n8CF
XbBfA9ECgYEA3nQmLHCbGX1o83b7ie3aAX1bAuEZgJMiRoQAUABW8CrGIg51g+po
wE7RZ4ESniSnGRz9yQKoDXq4FcPifpZrS3I6r8Zvwtoy9S9z2ag2niGnDOuNQ5fc
Rq1CJEr7TuD3t62dlh7FJ/yyKRdsPT+3hmeA3JyNDwi0NW9vqJpSdl0CgYEA+xja
QH3XhyXek3j1D60tHOHAX7fUFvv8dJinHfAlLjp7oLGs+TgnIAc5wStS7gIbrhNE
7IHX8IbXE0+0sb7daGHhyb+T0J7tQ0fHPS+V8fpwT9Pl44AVQKWzPbsTbC4N+xwe
mIe7HzErD7sy+wING2kv1UZnL778vpZNPHveV4UCgYAgr/SEBy/jOPhY/hzMEbU9
Dsx2ydjTectJjU/2cXZU6BQhIPrHnYQy7eH7UY4Iyt365LWt+cPz5xpxqEz5yOSP
O4PAHGqDuUhPmt9tFjigV9WSInKpggEOKZtUdegjmQ8NYGeNjYvu6kTLoPN4tIol
J8RZpm9bzC2exHcl0TdYyQKBgQCB0oulJGs2uOGnJbaucD8O27l2w7ioWYhhUDu3
Qt42VI5uuu5PvDSeXp4BvcCWxghBrDzKeyeGeHDizycBb0lSGql+gcqO5lyNmKLu
g5fnEDDZVRla0nIqhoFxvTOBjx4zYop/Gk4pBmbZL1RgauMT9QKCJnBbQ0ex0kwE
pZaDcQKBgHxyz/DyXQb3PYQ1+7dYBGQjRtY2nqEKHKdmvMf/Cwj6zqcqNY2uOcmu
8ooq+4RC/uKUOQMjKepTAJcFNfhjncp7NKUHsE6gEIZa9mAADpsPlwehqC3TbqWr
NnhaI9JoDrLxus/3u6nMoWEPRu9g7W8r9HFk/fHdpnFL5kaTI0db
-----END RSA PRIVATE KEY-----
`)

func serverKeyBytes() ([]byte, error) {
	return _serverKey, nil
}

func serverKey() (*asset, error) {
	bytes, err := serverKeyBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "server.key", size: 1675, mode: os.FileMode(288), modTime: time.Unix(1623966956, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"CA.crl":              caCrl,
	"CA.crt":              caCrt,
	"CA.key":              caKey,
	"client.crt":          clientCrt,
	"client.key":          clientKey,
	"invalid-ca.crt":      invalidCaCrt,
	"metron.csr":          metronCsr,
	"reverselogproxy.csr": reverselogproxyCsr,
	"server.crt":          serverCrt,
	"server.key":          serverKey,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}

var _bintree = &bintree{nil, map[string]*bintree{
	"CA.crl":              &bintree{caCrl, map[string]*bintree{}},
	"CA.crt":              &bintree{caCrt, map[string]*bintree{}},
	"CA.key":              &bintree{caKey, map[string]*bintree{}},
	"client.crt":          &bintree{clientCrt, map[string]*bintree{}},
	"client.key":          &bintree{clientKey, map[string]*bintree{}},
	"invalid-ca.crt":      &bintree{invalidCaCrt, map[string]*bintree{}},
	"metron.csr":          &bintree{metronCsr, map[string]*bintree{}},
	"reverselogproxy.csr": &bintree{reverselogproxyCsr, map[string]*bintree{}},
	"server.crt":          &bintree{serverCrt, map[string]*bintree{}},
	"server.key":          &bintree{serverKey, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
