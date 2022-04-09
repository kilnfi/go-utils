package staking

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
)

type KeystoreManager struct {
	scryptEncryptor, pbkdf2Encryptor *keystorev4.Encryptor
}

func NewKeystoreManager() *KeystoreManager {
	return &KeystoreManager{
		scryptEncryptor: keystorev4.New(keystorev4.WithCipher("scrypt")),
		pbkdf2Encryptor: keystorev4.New(keystorev4.WithCipher("pbkdf2")),
	}
}

func (mngr *KeystoreManager) GenerateValidatorKeys(mnemonic string, count int, storeMnemo bool, cb func(string) error) (keys []*ValidatorKey, err error) {
	return GenerateValidatorKeys(mnemonic, "", count, storeMnemo, cb)
}

func (mngr *KeystoreManager) EncryptToScryptKeystore(vKey *ValidatorKey, pwd string) (map[string]interface{}, error) {
	cryptoKs, err := mngr.scryptEncryptor.Encrypt(vKey.PrivKey.Marshal(), pwd)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"crypto":      cryptoKs,
		"pubkey":      vKey.Pubkey,
		"path":        vKey.Path,
		"uuid":        vKey.UUID,
		"version":     4,
		"description": "",
	}, nil
}

func (mngr *KeystoreManager) EncryptToPbkdf2Keystore(vKey *ValidatorKey, pwd string) (map[string]interface{}, error) {
	cryptoKs, err := mngr.pbkdf2Encryptor.Encrypt(vKey.PrivKey.Marshal(), pwd)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"crypto":      cryptoKs,
		"pubkey":      vKey.Pubkey,
		"path":        vKey.Path,
		"uuid":        vKey.UUID,
		"version":     4,
		"description": "",
	}, nil
}

func (mngr *KeystoreManager) DecryptFromKeystore(ks map[string]interface{}, pwd string) (*ValidatorKey, error) {
	version, ok := ks["version"]
	if ok {
		v, ok := version.(int)
		if ok && v != 4 {
			return nil, fmt.Errorf("invalid keystore version %v (version 4 expected)", v)
		}
	}

	cryptoKs, ok := ks["crypto"]
	if !ok {
		return nil, fmt.Errorf("invalid keystore missing \"crypto\" field")
	}

	iKs, ok := cryptoKs.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid \"crypto\" keystore format")
	}

	bytes, err := mngr.scryptEncryptor.Decrypt(iKs, pwd)
	if err != nil {
		return nil, err
	}

	vkey, err := ValidatorKeyFromBytes(bytes)
	if err != nil {
		return nil, err
	}

	if pth, ok := ks["path"]; ok {
		vkey.Path, _ = pth.(string)
	}

	if uuid, ok := ks["uuid"]; ok {
		vkey.UUID, _ = uuid.(string)
	}

	if desc, ok := ks["description"]; ok {
		vkey.Desc, _ = desc.(string)
	}

	return vkey, nil
}

type accountsStore struct {
	PrivateKeys [][]byte `json:"private_keys"`
	PublicKeys  [][]byte `json:"public_keys"`
}

func (mngr *KeystoreManager) EncryptToPrysmKeystore(vkeys []*ValidatorKey, pwd string) (map[string]interface{}, error) {
	accStore := new(accountsStore)

	exists := make(map[string]bool)
	for i := 0; i < len(vkeys); i++ {
		if exists[string(vkeys[i].PrivKey.Marshal())] {
			continue
		}
		accStore.PublicKeys = append(accStore.PublicKeys, vkeys[i].PrivKey.PublicKey().Marshal())
		accStore.PrivateKeys = append(accStore.PrivateKeys, vkeys[i].PrivKey.Marshal())
		exists[string(vkeys[i].PrivKey.Marshal())] = true
	}

	encodedStore, err := json.MarshalIndent(accStore, "", "\t")
	if err != nil {
		return nil, err
	}
	cryptoKs, err := mngr.pbkdf2Encryptor.Encrypt(encodedStore, pwd)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"crypto":  cryptoKs,
		"name":    "keystore",
		"version": 4,
		"uuid":    uuid.New().String(),
	}, nil
}
