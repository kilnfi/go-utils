package types

import (
	"encoding/json"
	"fmt"
	"strconv"

	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
	beaconphase0 "github.com/protolambda/zrnt/eth2/beacon/phase0"
)

type Validator struct {
	Index     beaconcommon.ValidatorIndex `json:"index" yaml:"index"`
	Status    string                      `json:"status" yaml:"status"`
	Balance   beaconcommon.Gwei           `json:"balance" yaml:"balance"`
	Validator *beaconphase0.Validator     `json:"validator" yaml:"validator"`
}

type ValidatorBalance struct {
	Index   beaconcommon.ValidatorIndex `json:"index" yaml:"index"`
	Balance beaconcommon.Gwei           `json:"balance" yaml:"balance"`
}

func (val Validator) MarshalCSV() ([]string, error) {
	record := []string{
		strconv.FormatInt(int64(val.Index), 10),
		val.Status,
		val.Balance.String(),
	}

	if val.Validator != nil {
		record = append(
			record,
			val.Validator.Pubkey.String(),
			val.Validator.WithdrawalCredentials.String(),
			val.Validator.EffectiveBalance.String(),
			strconv.FormatBool(val.Validator.Slashed),
			val.Validator.ActivationEligibilityEpoch.String(),
			val.Validator.ActivationEpoch.String(),
			val.Validator.ExitEpoch.String(),
			val.Validator.WithdrawableEpoch.String(),
		)
	} else {
		record = append(record, "", "", "", "", "", "", "", "")
	}

	return record, nil
}

func (val *Validator) UnmarshalCSV(record []string) error {
	if len(record) != 11 {
		return fmt.Errorf("invalid csv record with %d fields (%d fields expected)", len(record), 11)
	}

	slashed, err := strconv.ParseBool(record[6])
	if err != nil {
		return err
	}

	b, _ := json.Marshal(map[string]interface{}{
		"index":   record[0],
		"status":  record[1],
		"balance": record[2],
		"validator": map[string]interface{}{
			"pubkey":                       record[3],
			"withdrawal_credentials":       record[4],
			"effective_balance":            record[5],
			"slashed":                      slashed,
			"activation_eligibility_epoch": record[7],
			"activation_epoch":             record[8],
			"exit_epoch":                   record[9],
			"withdrawable_epoch":           record[10],
		},
	})

	return json.Unmarshal(b, val)
}
