package template

type PeerConfigTmplData struct {
	PeerPublicKey           string
	PeerPresharedKey        string
	PeerEndpoint            string
	PeerAllowedIPs          []string
	PeerPersistentKeepalive int
}

type ConfigTmplData struct {
	InterfacePrivateKey string
	InterfaceAddress    []string
	InterfacePort       string
	InterfaceDNS        string
	InterfaceMTU        int
	InterfaceTable      string
	InterfaceFwMark     int
	InterfacePreUp      string
	InterfacePostUp     string
	InterfacePreDown    string
	InterfacePostDown   string
	InterfacePeers      []PeerConfigTmplData
	SaveConfig          bool
}

var ConfigTemplate = `[Interface]
PrivateKey = {{ .InterfacePrivateKey }}
Address = {{ StringsJoin .InterfaceAddress ", " }}
{{ if ne .InterfacePort "" -}}
ListenPort = {{ .InterfacePort }}
{{- end}}
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
Table = {{ .InterfaceTable }}
{{- end}}
{{ if ne .InterfacePreUp "" -}}
PreUp = {{ .InterfacePreUp }}
{{- end}}
{{ if ne .InterfacePreDown "" -}}
PreDown = {{ .InterfacePreDown }}
{{- end}}
{{ if ne .InterfacePostUp "" -}}
PostUp = {{ .InterfacePostUp }}
{{- end}}
{{ if ne .InterfacePostDown "" -}}
PostDown = {{ .InterfacePostDown }}
{{- end}}
{{ if .SaveConfig -}}
SaveConfig = true
{{- end}}

{{range .InterfacePeers}}
[Peer]
PublicKey = {{ .PeerPublicKey }}
{{ if ne .PeerPresharedKey "" -}}
PresharedKey = {{ .PeerPresharedKey }}
{{- end}}
{{ if ne .PeerEndpoint "" -}}
Endpoint = {{ .PeerEndpoint }}
{{- end}}
AllowedIPs = {{ StringsJoin .PeerAllowedIPs ", " }}
{{ if ne .PeerPersistentKeepalive 0 -}}
PersistentKeepalive = {{ .PeerPersistentKeepalive }}
{{end}}
{{end}}
`
