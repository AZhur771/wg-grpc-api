package template

type TmplData struct {
	InterfacePrivateKey string
	InterfaceAddress    []string
	InterfaceDNS        []string
	InterfaceMTU        int

	PeerPublicKey       string
	PeerPresharedKey    string
	PeerEndpoint        string
	PeerAllowedIPs      []string
	PersistentKeepalive int
}

type TmplDevice struct {
	Endpoint            string
	AllowedIPs          []string
	PersistentKeepalive int
}

var ClientConfigTemplate = `[Interface]
PrivateKey = {{ .InterfacePrivateKey }}
Address = {{ StringsJoin .InterfaceAddress ", " }}
DNS = {{ StringsJoin .InterfaceDNS ", " }}
MTU = {{ .InterfaceMTU }}

[Peer]
PublicKey = {{ .PeerPublicKey }}
{{ if ne .PeerPresharedKey "" -}}
PresharedKey = {{ .PeerPresharedKey }}
{{- end}}
Endpoint = {{ .PeerEndpoint }}
AllowedIPs = {{ StringsJoin .PeerAllowedIPs ", " }}
{{ if ne .PersistentKeepalive 0 -}}
PersistentKeepalive = {{ .PersistentKeepalive }}
{{- end}}
`
