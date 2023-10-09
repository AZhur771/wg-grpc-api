package template

type ClientConfigTmplData struct {
	InterfacePrivateKey string
	InterfaceAddress    []string
	InterfaceDNS        string
	InterfaceMTU        int

	PeerPublicKey           string
	PeerPresharedKey        string
	PeerEndpoint            string
	PeerAllowedIPs          []string
	PeerPersistentKeepalive int
}

type ServerConfigTmplData struct {
	InterfacePrivateKey          string
	InterfaceAddress             string
	InterfacePort                string
	InterfaceDNS                 string
	InterfaceMTU                 int
	InterfaceTable               string
	InterfaceFwMark              int
	InterfacePreUp               string
	InterfacePostUp              string
	InterfacePreDown             string
	InterfacePostDown            string
	InterfacePersistentKeepAlive int
}

var ClientConfigTemplate = `[Interface]
PrivateKey = {{ .InterfacePrivateKey }}
Address = {{ StringsJoin .InterfaceAddress ", " }}
{{ if ne .InterfaceDNS "" -}}
DNS = {{ .InterfaceDNS }}
{{- end}}
{{ if ne .InterfaceMTU 0 -}}
MTU = {{ .InterfaceMTU }}
{{- end}}

[Peer]
PublicKey = {{ .PeerPublicKey }}
{{ if ne .PeerPresharedKey "" -}}
PresharedKey = {{ .PeerPresharedKey }}
{{- end}}
Endpoint = {{ .PeerEndpoint }}
AllowedIPs = {{ StringsJoin .PeerAllowedIPs ", " }}
{{ if ne .PeerPersistentKeepalive 0 -}}
PersistentKeepalive = {{ .PeerPersistentKeepalive }}
{{- end}}
`

var ServerConfigTemplate = `[Interface]
PrivateKey = {{ .InterfacePrivateKey }}
Address = {{ .InterfaceAddress }}
ListenPort = {{ .InterfacePort }}
{{ if ne .InterfaceMTU 0 -}}
MTU = {{ .InterfaceMTU }}
{{- end}}
{{ if ne .InterfaceDNS "" -}}
DNS = {{ .InterfaceDNS }}
{{- end}}
{{ if ne .InterfaceFwMark 0 -}}
FwMark = {{ .InterfaceFwMark }}
{{- end}}
{{ if ne .InterfaceTable "" -}}
InterfaceTable = {{ .InterfaceTable }}
{{- end}}
{{ if ne .InterfacePersistentKeepAlive 0 -}}
PersistentKeepalive = {{ .InterfacePersistentKeepAlive }}
{{- end}}
{{ if ne .InterfacePreUp "" -}}
PreUp = {{ .InterfacePreUp }}
{{- end}}
{{ if ne .InterfacePostUp "" -}}
PostUp = {{ .InterfacePostUp }}
{{- end}}
{{ if ne .InterfacePreDown "" -}}
PreDown = {{ .InterfacePreDown }}
{{- end}}
{{ if ne .InterfacePostDown "" -}}
PostDown = {{ .InterfacePostDown }}
{{- end}}
SaveConfig = true

`
