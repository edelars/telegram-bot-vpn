package iphone

import (
	"backend-vpn/internal/dto"
	"backend-vpn/pkg/controller"
	"backend-vpn/pkg/storage"
	"backend-vpn/pkg/transport"
	"context"
	"fmt"
	"github.com/rs/zerolog"
)

const (
	fileDumb = `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
 <dict>
 <key>PayloadContent</key>
 <array>
 <dict>
 <key>IKEv2</key>
 <dict>
 <key>AuthName</key>
 <string>%s</string>
 <key>AuthPassword</key>
 <string>%s</string>
 <key>AuthenticationMethod</key>
 <string>SharedSecret</string>
 <key>ChildSecurityAssociationParameters</key>
 <dict>
 <key>DiffieHellmanGroup</key>
 <integer>2</integer>
 <key>EncryptionAlgorithm</key>
 <string>AES-128</string>
 <key>IntegrityAlgorithm</key>
 <string>SHA1-96</string>
 <key>LifeTimeInMinutes</key>
 <integer>1440</integer>
 </dict>
 <key>DeadPeerDetectionRate</key>
 <string>Medium</string>
 <key>DisableMOBIKE</key>
 <integer>0</integer>
 <key>DisableRedirect</key>
 <integer>0</integer>
 <key>EnableCertificateRevocationCheck</key>
 <integer>0</integer>
 <key>EnablePFS</key>
 <integer>0</integer>
 <key>ExtendedAuthEnabled</key>
 <true/>
 <key>IKESecurityAssociationParameters</key>
 <dict>
 <key>DiffieHellmanGroup</key>
 <integer>2</integer>
 <key>EncryptionAlgorithm</key>
 <string>AES-128</string>
 <key>IntegrityAlgorithm</key>
 <string>SHA1-96</string>
 <key>LifeTimeInMinutes</key>
 <integer>1440</integer>
 </dict>
 <key>LocalIdentifier</key>
 <string>%s</string>
 <key>RemoteAddress</key>
 <string>65.108.96.44</string>
 <key>RemoteIdentifier</key>
 <string>%s</string>
 <key>SharedSecret</key>
 <string>%s</string>
 <key>UseConfigurationAttributeInternalIPSubnet</key>
 <integer>0</integer>
 </dict>
 <key>IPv4</key>
 <dict>
 <key>OverridePrimary</key>
 <integer>1</integer>
 </dict>
 <key>PayloadDescription</key>
 <string>Configures VPN settings for iphone</string>
 <key>PayloadDisplayName</key>
 <string>TutorialVPN</string>
 <key>PayloadIdentifier</key>
 <string>com.apple.vpn.managed.03d38bd2-2153-4ad9-ba85-c202276e4d64</string>
 <key>PayloadType</key>
 <string>com.apple.vpn.managed</string>
 <key>PayloadUUID</key>
 <string>9d667b63-854e-4e82-896b-38f0a23b861b</string>
 <key>PayloadVersion</key>
 <real>1</real>
 <key>Proxies</key>
 <dict>
 <key>HTTPEnable</key>
 <integer>0</integer>
 <key>HTTPSEnable</key>
 <integer>0</integer>
 <key>ProxyAutoConfigEnable</key>
 <integer>0</integer>
 <key>ProxyAutoDiscoveryEnable</key>
 <integer>0</integer>
 </dict>
 <key>UserDefinedName</key>
 <string>OnionVPN</string>
 <key>VPNType</key>
 <string>IKEv2</string>
 <key>VendorConfig</key>
 <dict/>
 </dict>
 </array>
 <key>PayloadDisplayName</key>
 <string>IKEv2</string>
 <key>PayloadIdentifier</key>
 <string>5f83468c-3b50-40d6-8518-28448b170795</string>
 <key>PayloadRemovalDisallowed</key>
 <false/>
 <key>PayloadType</key>
 <string>Configuration</string>
 <key>PayloadUUID</key>
 <string>ae2c121b-35fd-466f-9f21-f16cc021cc49</string>
 <key>PayloadVersion</key>
 <integer>1</integer>
 </dict>
</plist>
`
)

type Iphone struct {
	ctrl     controller.Controller
	logger   *zerolog.Logger
	endpoint string
	text     string
}

func NewIphone(ctrl controller.Controller, logger *zerolog.Logger, text string, endpoint string) *Iphone {
	return &Iphone{ctrl: ctrl, logger: logger, endpoint: endpoint, text: text}
}
func (p *Iphone) Data() (text, unique string) {
	return p.text, p.endpoint
}
func (p *Iphone) Handler() func(data transport.HandlerData) interface{} {
	return func(data transport.HandlerData) interface{} {

		u := storage.NewUserQuery(data.Username, data.Id, "")

		if err := p.ctrl.Exec(context.Background(), u); err != nil {
			p.logger.Debug().Err(err).Msg("fail")
			return fmt.Sprint("Ошибка, попробуйте позже")
		}
		p.logger.Debug().Msgf("user %s press %s button", u.Out.User, p.endpoint)

		fileStr := fmt.Sprintf(fileDumb, u.Out.User.Login, u.Out.User.Password, u.Out.User.Login, u.Out.User.Login, u.Out.User.Psk)

		tgf := dto.TgFile{
			File:     []byte(fileStr),
			Filename: fmt.Sprintf("%s.mobileconfig", u.Out.User.Login),
			Caption:  "Сохраните и запустите файл. После этого в настройках надо применить и активировать профиль VPN",
		}

		return tgf
	}
}
