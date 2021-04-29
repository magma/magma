package akatataipx

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/modules/eap/methods"
	"fbc/cwf/radius/modules/eap/methods/common"
	eap "fbc/cwf/radius/modules/eap/packet"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"

	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

// EapAkaTataIpxMethod Implementation ofthe EAP-AKA method impl with
// TATA IPX integration
type EapAkaTataIpxMethod struct {
	config Config
	client *TataClient
}

// Config the aka-magma configuration
type Config struct {
	IpxEndpoint string `json:"IpxEndpoint"`
	Username    string `json:"Username"`
	Password    string `json:"Password"`
}

// Create ...
func Create(config methods.MethodConfig) (methods.EapMethod, error) {
	// Parse config
	var akaConfig Config
	err := mapstructure.Decode(config, &akaConfig)
	if err != nil {
		return nil, errors.New("failed to parse AKA TATA configuration")
	}

	// Create TATA IPX client
	client, err := NewClient(akaConfig.Username, akaConfig.Password, akaConfig.IpxEndpoint)
	if err != nil {
		return nil, fmt.Errorf("failed creating TATA IPX client: %s", err.Error())
	}

	return &EapAkaTataIpxMethod{
		config: akaConfig,
		client: &client,
	}, nil
}

// Handle ...
// TODO: func is very long, should be broken into smaller
func (m EapAkaTataIpxMethod) Handle(
	c *modules.RequestContext,
	p *eap.Packet,
	state string,
	r *radius.Request,
) (*methods.HandlerResponse, error) {
	// Prepare result
	result := &methods.HandlerResponse{
		Packet: &eap.Packet{
			Code:       eap.CodeREQUEST,
			EAPType:    eap.EAPTypeAKA,
			Identifier: p.Identifier + 1,
			Data:       nil,
		},
		RadiusCode:       radius.CodeAccessChallenge,
		NewProtocolState: state,
		ExtraAttributes:  make(radius.Attributes),
	}

	// Handle EAP-Identity request
	if p.EAPType == eap.EAPTypeIDENTITY {
		c.Logger.Info("[aka] formatting IDENTITY response")
		result.Packet.Data = m.getAkaIdentityResponse().Bytes()
		return result, nil
	}

	// Verify we got correct EAP type
	if p.EAPType != eap.EAPTypeAKA {
		return nil, errors.New("invalid EAP packet type")
	}

	// Parse incoming packet
	c.Logger.Info("parsing AKA packet")
	akaPacket, err := NewAkaPacket(p.Data)
	if err != nil {
		return nil, err
	}

	// Deserialize the current state
	currentState := DeserializeState(state)

	// TODO: Add verification we are in got the expected akaPacket.Subtype
	// and fail if not (security!)

	// Handle the packet
	if akaPacket.Subtype == AkaIdentity {
		akaRequest, newState, err := m.handleAkaIdentity(c, p, akaPacket)
		if err != nil {
			return nil, err
		}
		result.Packet.Data = akaRequest.Bytes()
		result.NewProtocolState = *newState
	}

	if akaPacket.Subtype == AkaChallenge {
		// Log packet values
		res := akaPacket.GetFirst(AT_RES)
		mac := akaPacket.GetFirst(AT_MAC)
		c.Logger.Info(
			"got aka client challenge",
			zap.Any("AT_RES", res.Value),
			zap.Any("AT_MAC", mac.Value),
		)

		// TODO: verify challenge against TATA API (not needed for POC)

		// Prepare response
		result.RadiusCode = radius.CodeAccessAccept
		result.Packet.Code = eap.CodeSUCCESS
		result.Packet.Identifier-- // For EAP-Success we use the same ID as last challenge

		// Add MPPE keys
		keyingMaterialAttrs, err := common.GetKeyingAttributes(
			currentState.MSK,
			r.Secret,
			r.Authenticator[:],
		)
		if err != nil {
			return nil, err
		}
		result.ExtraAttributes[rfc2865.VendorSpecific_Type] = keyingMaterialAttrs

		// Add required User-Name
		result.ExtraAttributes[rfc2865.UserName_Type] =
			[]radius.Attribute{
				radius.Attribute(currentState.Identity),
			}
	}

	if akaPacket.Subtype == AkaClientError {
		c.Logger.Error(
			"got AKA error, failing authorization",
			zap.Any("packet", akaPacket),
			zap.Error(err),
		)
		result.RadiusCode = radius.CodeAccessReject
		result.Packet.Code = eap.CodeFAILURE
	}

	// We're done
	return result, nil
}

func (m EapAkaTataIpxMethod) getAkaIdentityResponse() *AkaPacket {
	return &AkaPacket{
		Subtype:  AkaIdentity,
		Reserved: 0x0000,
		Attributes: []Attribute{
			{
				Type:  AT_PERMANENT_ID_REQ,
				Value: []byte{0x00, 0x00},
			},
		},
	}
}

func (m EapAkaTataIpxMethod) handleAkaIdentity(c *modules.RequestContext, eapPacket *eap.Packet, p *AkaPacket) (*AkaPacket, *string, error) {
	identityAttr := p.GetFirst(AT_IDENTITY)
	if identityAttr == nil {
		return nil, nil, errors.New("aka identity packet is missing the AT_IDENTITY attribute")
	}

	if len(identityAttr.Value) < 16 {
		return nil, nil, errors.New("invalid identity value")
	}

	actualLength := (uint16(identityAttr.Value[0]) << 8) | uint16(identityAttr.Value[1])
	if actualLength > uint16(len(identityAttr.Value)) {
		return nil, nil, fmt.Errorf(
			"invalid actual length field %d (packet has only %d bytes)",
			actualLength,
			len(identityAttr.Value),
		)
	}

	if identityAttr.Value[2] != '0' {
		return nil, nil, fmt.Errorf(
			"identity string starts with 0x%x and not with 0x30 (the char '0')",
			identityAttr.Value[2],
		)
	}

	identity := string(identityAttr.Value[2 : actualLength+2])
	loc := strings.Index(identity, "@")
	if loc < 1 {
		return nil, nil, errors.New("identity string must be of format imsi@domain")
	}
	imsi := identity[1:loc]

	// Make the MAPSAI call to TATA IPX API to get auth vectors
	c.Logger.Info(
		"Calling TATA API to get auth vectors",
		zap.String("imsi", imsi),
		zap.String("identity", identity),
	)
	response, err := m.client.Mapsi(imsi, fmt.Sprintf("%d", c.RequestID))
	if err != nil {
		return nil, nil, err
	}

	if response.ErrorCode != 0 && response.ErrorCode != "0" {
		return nil, nil, errors.New(
			"got error response from TATA: " + response.ErrorString,
		)
	}

	// Decode vectors
	rand, err := hex.DecodeString(response.AuthQuintuplets.RAND)
	if err != nil {
		return nil, nil, err
	}

	autn, err := hex.DecodeString(response.AuthQuintuplets.AUTN)
	if err != nil {
		return nil, nil, err
	}

	ik, err := hex.DecodeString(response.AuthQuintuplets.IK)
	if err != nil {
		return nil, nil, err
	}

	ck, err := hex.DecodeString(response.AuthQuintuplets.CK)
	if err != nil {
		return nil, nil, err
	}

	c.Logger.Info(
		"auth vectors received",
		zap.String("RAND", response.AuthQuintuplets.RAND),
		zap.String("AUTN", response.AuthQuintuplets.AUTN),
		zap.String("IK", response.AuthQuintuplets.IK),
		zap.String("CK", response.AuthQuintuplets.CK),
	)

	// Derive K_aut & MSK
	_, kaut, msk, _ := MakeAKAKeys([]byte(identity), ik, ck)

	// Build packet
	result := &AkaPacket{
		Subtype:  AkaChallenge,
		Reserved: 0x0000,
		Attributes: []Attribute{
			{
				Type:  AT_RAND,
				Value: append(bytes.Repeat([]byte{0x00}, 2), rand...),
			},
			{
				Type:  AT_AUTN,
				Value: append(bytes.Repeat([]byte{0x00}, 2), autn...),
			},
		},
	}
	result.AppendMac(eap.CodeREQUEST, eapPacket.Identifier+1, kaut)
	newState := getAkaState(identity, msk).Serialize()
	return result, &newState, nil
}
