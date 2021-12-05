package types

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gammtypes "github.com/osmosis-labs/osmosis/x/gamm/types"
	locktypes "github.com/osmosis-labs/osmosis/x/lockup/types"
)

func AnysToSdkMsgs(cdc codec.BinaryCodec, anys []*codectypes.Any) ([]sdk.Msg, error) {
	msgs := make([]sdk.Msg, len(anys))
	for i, any := range anys {
		var msg sdk.Msg
		err := cdc.UnpackAny(any, &msg)
		if err != nil {
			return nil, err
		}
		msgs[i] = msg
	}
	return msgs, nil
}

func SdkMsgsToAnys(cdc codec.BinaryCodec, data interface{}) ([]*codectypes.Any, error) {
	msgs := make([]sdk.Msg, 0)
	switch data := data.(type) {
	case sdk.Msg:
		msgs = append(msgs, data)
	case []sdk.Msg:
		if len(data) == 0 {
			return []*codectypes.Any{}, nil
		}
		msgs = append(msgs, data...)
	default:
		return nil, fmt.Errorf("wrong data")
	}

	anys := make([]*codectypes.Any, len(msgs))

	for i, msg := range msgs {
		var err error
		anys[i], err = codectypes.NewAnyWithValue(msg)
		if err != nil {
			return nil, err
		}
	}

	return anys, nil

}

func SdkMsgFromFile(fileDir string, msgType string) (sdk.Msg, error) {
	file, err := os.Open(fileDir)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	jsonData, _ := ioutil.ReadAll(reader)

	switch msgType {
	case "swap_and_pool":
		msg := &gammtypes.MsgJoinSwapExternAmountIn{}
		jsonErr := json.Unmarshal(jsonData, msg)
		if jsonErr != nil {
			return nil, err
		}
		return msg, nil
	case "lock":
		msg := &locktypes.MsgLockTokens{}
		jsonErr := json.Unmarshal(jsonData, msg)
		if jsonErr != nil {
			fmt.Println("Unable to map JSON at " + " to Investments")
		}
		return msg, nil
	}
	return nil, nil
}
