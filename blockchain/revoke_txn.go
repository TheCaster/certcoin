package blockchain

import (
	"github.com/cfromknecht/certcoin/crypto"
)

const (
	REVOKE_FEE = uint64(100)
)

func NewRevokeTxn(onlineSecret, offlineSecret crypto.CertcoinSecretKey,
	source crypto.CertcoinPublicKey,
	identity Identity) Txn {

	fullName := identity.FullNameStr()
	return Txn{
		Type: Revoke,
		Inputs: []Input{
			Input{
				PrevHash:  identity.Domain,
				PublicKey: onlineSecret.PublicKey,
				Signature: crypto.Sign(fullName, onlineSecret),
			},
			Input{
				PrevHash:  identity.Subdomain,
				PublicKey: offlineSecret.PublicKey,
				Signature: crypto.Sign(fullName, offlineSecret),
			},
			Input{
				PrevHash:  crypto.SHA256Sum{},
				PublicKey: source,
				Signature: crypto.CertcoinSignature{},
			},
		},
		Outputs: []Output{
			Output{
				Address: crypto.SHA256Sum{},
				Value:   REVOKE_FEE,
			},
		},
	}
}

func (bc *Blockchain) ValidRevokeTxn(t Txn) bool {
	if len(t.Inputs) < 3 || !(len(t.Outputs) == 1 && len(t.Outputs) == 2) {
		return false
	}

	identity, err := NewIdentity(string(t.Inputs[0].PrevHash[:]), string(t.Inputs[1].PrevHash[:]))
	if err != nil {
		return false
	}
	// offlinePK := lookup from database
	return t.Type == Revoke &&
		t.Outputs[0].Value >= REVOKE_FEE &&
		crypto.Verify(identity.FullNameStr(), t.Inputs[0].Signature, t.Inputs[0].PublicKey)
	//Verify(t.Body.Hash(), t.OfflineSignature, offlinePK)
}
