package bchutil

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"testing"
)

type SigHashVector struct {
	RawTx  string
	Inputs []Input
}

type Input struct {
	Pubkey    string
	Signature string
	Value     int64
}

var SigHashTestVectors = []SigHashVector{
	{
		RawTx: "0100000001d71dd1e5eab582c7b6ec624a70c9b3b515e11ede8bc86a7b11cac4e3835f4935010000006b483045022100c377c67448794dd1a49b7406b4178f4862f86d91061e402e0f1959e87456844002200c5055448cf4c270c653667006c88134b0c22b8a53182f6c5078f18fa820af3d412102bc4f7f11f7b45f4be411d4891bde9d0aedd59a5873c21a844b8dbb01c64dfb35ffffffff0200e1f505000000001976a914fb6553eda6df1a1fbf14d4fe43687b3b12aec92088ac02e6d213010000001976a9142035366682baeda80a415d80a8100f17890853d388ac00000000",
		Inputs: []Input{
			{"02bc4f7f11f7b45f4be411d4891bde9d0aedd59a5873c21a844b8dbb01c64dfb35", "3045022100c377c67448794dd1a49b7406b4178f4862f86d91061e402e0f1959e87456844002200c5055448cf4c270c653667006c88134b0c22b8a53182f6c5078f18fa820af3d41", int64(4727610477)},
		},
	},
	{
		RawTx: "0100000002770a4ae52da85244127710adee78135bf9a03bcea4062ffbeed1571221f445d8010000006b483045022100b343b9db00ba30d0db1577ae363466ca6851b6ddc2f529b7cc0a0652177510340220303cd29ab1a115aac65210b5edef8f917f12bb98400eaec754aa890e1d97a01041210210f1cd9e15cd041894192c68009b4f58171aaf01f89899f5df00b1bc205ccc72ffffffffe9aeabe253cedc1c877bf4d6a0479f17a436e9693ea8d4ab4727a1203dc4e8d8010000006a4730440220159b05368e5d3c45b55b6802456abf82d27128be7c1ae34e4be2fbb4404e9cef022058f1c642a79ccd184bd92b7cec3e3c7bc6db8c3061119907c95ed279cc3c4c33412103fe95d4ffb769c5329e6423348029112bd49792e9595cc19cdb551d4713a6dfe9ffffffff020084d717000000001976a91451fed1517527e3a5b8080cb1abe3d6167d9afa7e88ac760ecd17000000001976a9148065346932ae9d3d253f4db10c302a05d78e9e6b88ac00000000",
		Inputs: []Input{
			{"0210f1cd9e15cd041894192c68009b4f58171aaf01f89899f5df00b1bc205ccc72", "3045022100b343b9db00ba30d0db1577ae363466ca6851b6ddc2f529b7cc0a0652177510340220303cd29ab1a115aac65210b5edef8f917f12bb98400eaec754aa890e1d97a01041", int64(399420000)},
			{"03fe95d4ffb769c5329e6423348029112bd49792e9595cc19cdb551d4713a6dfe9", "30440220159b05368e5d3c45b55b6802456abf82d27128be7c1ae34e4be2fbb4404e9cef022058f1c642a79ccd184bd92b7cec3e3c7bc6db8c3061119907c95ed279cc3c4c3341", int64(400000000)},
		},
	},
	{
		RawTx: "0100000002917145d86182b8a8ed62759a45494bae9e22fa0feaa09b24a2a15abbc1d056cd010000006b4830450221009582cb61f98d96094cba3787c76a4817718b57dd002ade62126e021830f67a4002203331030e35b8d1b7fd0d3d4338c9e136e601605a9ec80d72e842a06009c26c9d412103b3ea4932c641d3419d2d51d735db8310c6c0e55228fbcbfa75f3f2a3a39190fdffffffff5938a89e87a46a15693244533f9326f35be99b5d92c248912d5b400aba3c1552010000006a47304402204a3c6b0ed55f4e8714eb63681eeddaf672b481adc1f1b7c2d5297a14a5750c42022020cc8f033eded2bfcdd8eadbd4e73798e7edb12cca985a0d44553f7ad04b32e9412102e82f6c2c29ad67232a3a27dc13d26eca7b741d60fae3aa4eba47d028a6dcd8b5ffffffff02809b60d7010000001976a9140efdf418b594e8994e2af434e8790a346c2d14fa88ac49548b10010000001976a9143c094094fb6d5d573f8d166e795bb56445495d6988ac00000000",
		Inputs: []Input{
			{"03b3ea4932c641d3419d2d51d735db8310c6c0e55228fbcbfa75f3f2a3a39190fd", "30450221009582cb61f98d96094cba3787c76a4817718b57dd002ade62126e021830f67a4002203331030e35b8d1b7fd0d3d4338c9e136e601605a9ec80d72e842a06009c26c9d41", int64(6245933502)},
			{"02e82f6c2c29ad67232a3a27dc13d26eca7b741d60fae3aa4eba47d028a6dcd8b5", "304402204a3c6b0ed55f4e8714eb63681eeddaf672b481adc1f1b7c2d5297a14a5750c42022020cc8f033eded2bfcdd8eadbd4e73798e7edb12cca985a0d44553f7ad04b32e941", int64(6235001007)},
		},
	},
	{
		RawTx: "01000000011b538d2b11c88a90b88e6b89cc36dac2e3bfb508e807837005cc0577bed1720a000000006a4730440220146b6997e53de3965ea427ebd583a29eedf963b7105c45c7a28cc93fd327223a02200fa1c990c94fea445ce6a0b427960d3ee48c2f5e2088d0f80be6022fc6df4877412102d4023f62f0b56f0912d6d568cba805802419f8d6e1078ceb701abe8a1d76c69cffffffff019caca907000000001976a9144fda757634c0d56211c64d1cb029e995a4f21ee588ac00000000",
		Inputs: []Input{
			{"02d4023f62f0b56f0912d6d568cba805802419f8d6e1078ceb701abe8a1d76c69c", "30440220146b6997e53de3965ea427ebd583a29eedf963b7105c45c7a28cc93fd327223a02200fa1c990c94fea445ce6a0b427960d3ee48c2f5e2088d0f80be6022fc6df487741", int64(128560666)},
		},
	},
	{
		RawTx: "020000000162bc145631c3f60f6bc90ddd7e3cf6aec7ebc9a49b73a0765a2a151b4b1475d4010000006a473044022069dec1655011257d8947b9e0222c28b008a76459a8b280f2a7b6daf02c8c9a7e022069093b8d0cf18f153ef460cf61408fb288e1f7c19dfe974caa02df601ab85867412102e156045831e1e93e32da05a00529599febcf988b1f3a8f41397bd7a6a0d1d380feffffff021e2c3101000000001976a914bdc4d05afd8da87e2aa03dd011279df54b95163388ac00f90295000000001976a914f9a9d317345b92f0d286245521fc8eaead8a44d488ac95ce0700",
		Inputs: []Input{
			{"02e156045831e1e93e32da05a00529599febcf988b1f3a8f41397bd7a6a0d1d380", "3044022069dec1655011257d8947b9e0222c28b008a76459a8b280f2a7b6daf02c8c9a7e022069093b8d0cf18f153ef460cf61408fb288e1f7c19dfe974caa02df601ab8586741", int64(2520000000)},
		},
	},
}

func TestCalulateSigHash(t *testing.T) {
	for i, v := range SigHashTestVectors {
		raw, err := hex.DecodeString(v.RawTx)
		if err != nil {
			t.Error(err)
			return
		}
		r := bytes.NewReader(raw)
		msgTx := wire.NewMsgTx(1)
		msgTx.BtcDecode(r, 1, wire.BaseEncoding)
		for idx, _ := range msgTx.TxIn {
			pubKeyBytes, err := hex.DecodeString(v.Inputs[idx].Pubkey)
			if err != nil {
				t.Error(err)
			}
			addr, err := NewCashAddressPubKeyHash(btcutil.Hash160(pubKeyBytes), &chaincfg.MainNetParams)
			if err != nil {
				t.Error(err)
			}
			prevScript, err := PayToAddrScript(addr)
			if err != nil {
				t.Error(err)
			}
			pubkey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
			if err != nil {
				t.Error(err)
			}
			sigBytes, err := hex.DecodeString(v.Inputs[idx].Signature)
			if err != nil {
				t.Error(err)
			}
			sig, err := btcec.ParseDERSignature(sigBytes, btcec.S256())
			if err != nil {
				t.Error(err)
			}
			hash := calcBip143SignatureHash(prevScript, txscript.NewTxSigHashes(msgTx), txscript.SigHashAll, msgTx, idx, v.Inputs[idx].Value)
			if !sig.Verify(hash, pubkey) {
				t.Errorf("Calcualted invalid hash for vector %d  input %d ", i, idx)
			}
		}
	}
}