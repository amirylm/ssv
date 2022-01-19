package eth1

import (
	"encoding/hex"
	"encoding/json"
	"github.com/bloxapp/ssv/utils/logex"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"strings"
	"testing"
)

func TestParseOperatorAddedEvent(t *testing.T) {
	OldRawOperatorAdded := `{
  "address": "0x9573c41f0ed8b72f3bd6a9ba6e3e15426a0aa65b",
  "topics": [
	"0x39b34f12d0a1eb39d220d2acd5e293c894753a36ac66da43b832c9f1fdb8254e"
  ],
  "data": "0x000000000000000000000000000000000000000000000000000000000000006000000000000000000000000067ce5c69260bd819b4e0ad13f4b873074d47981100000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000005617364617300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002c0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000002644c5330744c5331435255644a54694253553045675546564354456c4449457446575330744c533074436b314a53554a4a616b464f516d64726357687261556335647a424351564646526b464254304e425554684254556c4a516b4e6e53304e4255555642623364464e303946596e643554477432636c6f7756465530616d6f4b6232393553555a34546e5a6e636c6b34526d6f7256334e736556705562486c714f4656455a6b5a7957576731565734796454525a545752425a53746a5547597857457372515339514f5668594e3039434e47356d4d51705062306457516a5a33636b4d76616d684d596e5a50534459314d484a3556566c766347565a6147785457486848626b5130646d4e3256485a6a6355784d516974315a54497661586c546546464d634670534c7a5a57436e4e554d325a47636b5676626e704756484675526b4e33513059794f476c51626b7057516d70594e6c517653474e55536a553153555272596e52766447467956545a6a6433644f543068755347743656334a324e326b4b64486c5161314930523255784d576874566b633555577053543351314e6d566f57475a4763305a764e55317855335a7863466c776246687253533936565535744f476f76624846465a465577556c6856636a517854416f7961486c4c57533977566d707a5a32316c56484e4f4e79396163554644613068355a546c47596d74574f565976566d4a556144646f56315a4d5648464855326733516c6b765244646e643039335a6e564c61584579436c52335355524255554643436930744c5330745255354549464a545153425156554a4d53554d675330565a4c5330414c53304b00000000000000000000000000000000000000000000000000000000",
  "blockNumber": "0x49f59c",
  "transactionHash": "0x097d9a621ace2ca0c78d115d833edc1901bfe75f107a7b3f427663ea308c12ca",
  "transactionIndex": "0xf",
  "blockHash": "0x9542ecebe9d541e2575cb5577dfd4b73c9b0c3ab634fcac4ce0ff319249c90e4",
  "logIndex": "0xf",
  "removed": false
}`
	rawOperatorAdded := `{
   "address":"0xd594c1ef4845713e86658cb42227a811625a285b",
   "topics":[
      "0x39b34f12d0a1eb39d220d2acd5e293c894753a36ac66da43b832c9f1fdb8254e",
      "0x000000000000000000000000a5cfd290965372553efd5fdaeb91c335207b76e2"
   ],
   "data":"0x00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000f546573744f70657261746f72383838000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002c0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000002644c533038393536343955644a54694253553045675546564354456c4449457446575330744c533074436b314a53554a4a616b464f516d64726357687261556335647a424351564646526b464254304e425554684254556c4a516b4e6e53304e42555556424f4852585247307862544e7459573552613078776556704c4d7a634b4d474e4852476f79646c42545753745257564642643342574f585a705754684b566c677a54324a30566a4e4c4c3234784e79397065475a325645783561475a4b636b677a5953747053314e4963446c3557455534635170364e3252684f546c61567a5534527a4179654446305a6e7075563152454d6d4670626b6c704d444177646a5135526a4654647a6c594f55747451556735567a4e47646a426152457061647a5a4b5646643352305a69436d5a69546d4d326347567654473575636e6c6c576c56586230395a516d733054566732556d395156325a584e554a456155526165484671566a6476624656335a6e46424d57354f65553936525846434d45746b5357384b624578535a4641344f445a424e464a725a47706a55446335615764724d30526a565664434d4468705a6c4d3453466c76533031325a555a72656b30795232646d4f47354c526e466d536e46594e7a6c796246523463417053546e6c6865555a4f5958685a57455934656e42424d486c5952474648513049315469747a5a314e32596a6731574441796457564361314e61644646554d554d7954474d78576c5a6b624552465a5670474e464e6c436b68335355524255554643436930744c5330745255354549464a545153425156554a4d53554d675330565a4c5330744c53304b00000000000000000000000000000000000000000000000000000000",
   "blockNumber":"0x5b5d76",
   "transactionHash":"0xe353261f2d9c94b08769bc1cccf719e687cce72e6ad192ad186033bd96cc94c8",
   "transactionIndex":"0x0",
   "blockHash":"0x1065e001029f348db9945b0e747e62f9dbe4aa6eb154d602026f833078424c70",
   "logIndex":"0x3",
   "removed":false
}`

	logger := logex.Build("test", zap.DebugLevel, nil)
	t.Run("legacy operator added", func(t *testing.T) {
		legacyLogOperatorAdded, legacyContractAbi := unmarshalLog(t, OldRawOperatorAdded, Legacy)
		abiParser := NewParser(logger, Legacy)
		parsed, isEventBelongsToOperator, err := abiParser.ParseOperatorAddedEvent(nil, legacyLogOperatorAdded.Data, nil, legacyContractAbi)
		require.NoError(t, err)
		require.NotNil(t, legacyContractAbi)
		require.False(t, isEventBelongsToOperator)
		require.NotNil(t, parsed)
		require.Equal(t, "asdas", parsed.Name)
		require.Equal(t, "0x67Ce5c69260bd819B4e0AD13f4b873074D479811", parsed.OwnerAddress.Hex())
	})

	t.Run("v2 operator added", func(t *testing.T) {
		LogOperatorAdded, contractAbi := unmarshalLog(t, rawOperatorAdded, V2)
		abiParser := NewParser(logger, V2)
		parsed, isEventBelongsToOperator, err := abiParser.ParseOperatorAddedEvent(nil, LogOperatorAdded.Data, LogOperatorAdded.Topics, contractAbi)
		require.NoError(t, err)
		require.NotNil(t, contractAbi)
		require.False(t, isEventBelongsToOperator)
		require.NotNil(t, parsed)
		require.Equal(t, "TestOperator888", parsed.Name)
		require.Equal(t, "LS0895649UdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBOHRXRG0xbTNtYW5Ra0xweVpLMzcKMGNHRGoydlBTWStRWVFBd3BWOXZpWThKVlgzT2J0VjNLL24xNy9peGZ2VEx5aGZKckgzYStpS1NIcDl5WEU4cQp6N2RhOTlaVzU4RzAyeDF0ZnpuV1REMmFpbklpMDAwdjQ5RjFTdzlYOUttQUg5VzNGdjBaREpadzZKVFd3R0ZiCmZiTmM2cGVvTG5ucnllWlVXb09ZQms0TVg2Um9QV2ZXNUJEaURaeHFqVjdvbFV3ZnFBMW5OeU96RXFCMEtkSW8KbExSZFA4ODZBNFJrZGpjUDc5aWdrM0RjVVdCMDhpZlM4SFlvS012ZUZrek0yR2dmOG5LRnFmSnFYNzlybFR4cApSTnlheUZOYXhZWEY4enBBMHlYRGFHQ0I1TitzZ1N2Yjg1WDAydWVCa1NadFFUMUMyTGMxWlZkbERFZVpGNFNlCkh3SURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K",
			string(parsed.PublicKey))
		require.Equal(t, "0xa5cfD290965372553Efd5fDaeB91C335207b76E2", parsed.OwnerAddress.Hex())
	})
}

func TestParseValidatorAddedEvent(t *testing.T) {
	legacyRawValidatorAdded := `{
  "address": "0x9573c41f0ed8b72f3bd6a9ba6e3e15426a0aa65b",
  "topics": [
	"0x8674c0b4bd63a0814bf1ae6d64d71cf4886880a8bdbd3d7c1eca89a37d1e9271"
  ],
  "data": "0x000000000000000000000000feedb14d8b2c76fdf808c29818b06b830e8c2c0e000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000003091db3a13ab428a6c9c20e7104488cb6961abeab60e56cf4ba199eed3b5f6e7ced670ecb066c9704dc2fa93133792381c000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000b80000000000000000000000000000000000000000000000000000000000000110000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000036000000000000000000000000000000000000000000000000000000000000003c000000000000000000000000000000000000000000000000000000000000002c0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000002644c5330744c5331435255644a54694253553045675546564354456c4449457446575330744c533074436b314a53554a4a616b464f516d64726357687261556335647a424351564646526b464254304e425554684254556c4a516b4e6e53304e42555556424d455a336545684c647a4a355a7a466a555731486446464b4d6d344b52437474646d64475331564c536e6c34515468365543394e4f446379596c70735233684b62305a4a624735684d44564d634868436545466e543035705a4556584d335a3165485647533245344f4467355558425957517070565751354d316848636b4e5a574546794f557074576a686964456c5554304e3464445634613346336557746855455a784d585647646e46774f4852574f54426f596b3536536c70354d6d786f553156794f485268436b6c755130467963304670616e7058533267314e6d705751314a4d51305253623363324d324a436456424c6545646d5a5735425a5564586448466b595842325531645861485a4a63584278636b7050566b4e576356594b63334e594d3278304f57517a656b453254325a4c4d54425455323478645545764d5539584e6a4e32596e4670534446775a475a73527a427952466b7262554e566269397453566f7a4d533876556a4e5a5a7a6856595170706246685054545a32555642685a44466c527a525253544e4761465a4d61456f77616a4d3462555a6855306446616a46435257687a52314e61536e5258526c6c684f555a4b566c4e4464587035516e705654556876436c68525355524255554643436930744c5330745255354549464a545153425156554a4d53554d675330565a4c5330744c53304b00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003092aa9de41ee36d746a37c2696816103a052bdcf03af3a9d0bf517fb9ef3c30501fcf34a73a57c78f3d3ca23da1aec4580000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000158782b6f4b3052687365507254432f4e77486766625870365a6b343679676375412b6551313733546343303771427239484a714261574636536d552f6c3833663067445474554b35434a364c7032445867437033574b476c6c513775704b72467a77517a46725252503132425968416c746a5a645a347068504f6e4343687637716e7362706d633976494f78654b4963416433696c4e51594857484c302b5758546c4c392f556343565461514d38685268656c6636504a67654e7a5145384f66413555624131667735595a51423865647932446d34734c4642617542386a4548554d4e526458413974356868633345617a7972695845736d48613047494a736e3537747a4c324a375455363566626f6d43766e57517743635a534639344142494477324b54324e6755524c34705932674870756a452f514a53573733537478593343444162654f39754979566448524b735941644238773d3d000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000036000000000000000000000000000000000000000000000000000000000000003c000000000000000000000000000000000000000000000000000000000000002c0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000002644c5330744c5331435255644a54694253553045675546564354456c4449457446575330744c533074436b314a53554a4a616b464f516d64726357687261556335647a424351564646526b464254304e425554684254556c4a516b4e6e53304e42555556424d465135526d63724d324a6f5a6b394d63546849524445324d6b4d4b5a4468704d6c5a44576a64475430647756484e4f4d6a52344f4578345754524e526d354b5a56464a656b3030633152594f546c4756324e4c5232644b62307842537a566b65545272546e465a57554675533051765451706a4e3255316348684b53465a4965484a595456567a5955744e55304a6d5a33706f556a5a6e63444e584d47704e53305a3354477836564670686231526a616c56455557356d4d336c514f584577626a426f64556879436e6c465444526862586b354f47313363453158576c70564f57644e576a4a4952544673634846305557396d516a417959574d79533146764f544e46546a4252546d6330656d6857526d466a645656355955703662326f4b6448646c574652435a574a7157453173645749315a584e5855334a45626a64555a465646654745725530566b56444261576d39526347704b596b5275646c526a643039766333557a545374505a44524e5555564f63517078556a52594e5751306133427956557377626d395864455658626d7733555445314d486c6b4f55355254305a5961304e6d57537472646c5259567a4277597a46776232466c5a336c4c64555a43566d704a52446478436b4e335355524255554643436930744c5330745255354549464a545153425156554a4d53554d675330565a4c5330744c53304b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000030a1000cd6aaffd95ebe6baff24fe4d9544b7e0ad6727b54ea36cf8117bd34b19fbb586e1ba16ce84084c0ba0a3fa8f3620000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000158435a6d6c4d4b42612f43527568346c4d57764d5a342b5a476476675373556a33416e386b573062544f77626862574d626b656d6667644651326e50344458304d374861703446584c3251774d464c6865474353534f6b7649784b42474c4c5439523253737453706358446e7955634d51464244764b5871365673537167644f61375048485a6874566b6775457037414b59566452554d345258657234536a4e6b4968474d6f5a506f775935665067784c33355864473537736c412b35745869596a4b326b4f6663685562306f6e33312f66356c365a386c31377536427679484d5a756658313932554b597a45566241505a64657a674e6e69506c35495452784f2b49423644415145323233754a6a2b48357857377570524d6a58656943796a34567769437354512b483933787151594c62553051316e4c566a3571525567534f674d52434a6133536432536537777531497633694d773d3d000000000000000000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000036000000000000000000000000000000000000000000000000000000000000003c000000000000000000000000000000000000000000000000000000000000002c0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000002644c5330744c5331435255644a54694253553045675546564354456c4449457446575330744c533074436b314a53554a4a616b464f516d64726357687261556335647a424351564646526b464254304e425554684254556c4a516b4e6e53304e42555556424e32705863457872656d643254586476527a684e64455679556a494b524768554d6b313164456c6d59556430566d784d654456574b326734616d7772646e6c7854315976636d784b5245566c517939484d7a567056304d3057455533526e464b55566331516d707651575a315458685165677052517a5a364d45453162314933656e52755748553263305633546b684a534668335245464954486c54645664514d334247596c6f30516e63356231465a54554a6d62564e734c33685852307379566e4e336156686b436b4e4663555a4b526d644e55466b334e6c4a5159306f325232646b545763725756525257565646616d6c52546a4670646d4a4b5a6a5257615570435254637262564e7465465a4e4e54417a566d6c7951575a6e646b494b656e426e64544e7a64485a496448705256315a3265484a304e545230526d39444d48526d5745315252584e53553056745456526f566b686f63566f725a544a434f43396b545751325231466f646e45355a58523152517068516b786f536c704655586c704d6b6c7055553032556c6732613031765a476447556d6376656d747454465a5851305649547a457a61465635526b6f78616e67314c304d3562454979553256454e57396a64316834436d4a525355524255554643436930744c5330745255354549464a545153425156554a4d53554d675330565a4c5330744c53304b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000030887c4fbcaf5dafaa60ea73f1944be4a3eaaf55f15fac8d2e717e10d6fcdb4c82cf000305acd3fe6eb43b51cd615df2670000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000001587468336178485a73434f75687268724c4d4c4e394874717977697533494b454479394c4d77504f6b32494d3066622f576744465671744b4d2b626b6c69327163585273346b6c39382f506e6d7378725630664a786871415343376746424955536739674a356864757a68676d6c7672443639554c376c627847545037676f337a32345741335372646254434f32675553496537675a437349796f4b617a466c4d46477342464c436647676264585555324f4e6f334937703564634271707072324e516a6f745475496146507568684a4545706436685a6b2b49766179325246655448796674704a79394569583241397478364459386b456b4269305841442f796d553230707053307a477463674b7a2b312f645a695963315347444a7a6f42424431777a737848642f6c353752766f6533372b52764d424d3853744c42436a3053324f46485673634b347545375070654368623758773d3d000000000000000000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000036000000000000000000000000000000000000000000000000000000000000003c000000000000000000000000000000000000000000000000000000000000002c0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000002644c5330744c5331435255644a54694253553045675546564354456c4449457446575330744c533074436b314a53554a4a616b464f516d64726357687261556335647a424351564646526b464254304e425554684254556c4a516b4e6e53304e4255555642623364464e303946596e643554477432636c6f7756465530616d6f4b6232393553555a34546e5a6e636c6b34526d6f7256334e736556705562486c714f4656455a6b5a7957576731565734796454525a545752425a53746a5547597857457372515339514f5668594e3039434e47356d4d51705062306457516a5a33636b4d76616d684d596e5a50534459314d484a3556566c766347565a6147785457486848626b5130646d4e3256485a6a6355784d516974315a54497661586c546546464d634670534c7a5a57436e4e554d325a47636b5676626e704756484675526b4e33513059794f476c51626b7057516d70594e6c513153474e55536a553153555272596e52766447467956545a6a6433644f543068755347743656334a324e326b4b64486c5161314930523255784d576874566b633555577053543351314e6d566f57475a4763305a764e55317855335a7863466c776246687253533936565535744f476f76624846465a465577556c6856636a517854416f7961486c4c57533977566d707a5a32316c56484e4f4e79396163554644613068355a546c47596d74574f565976566d4a556144646f56315a4d5648464855326733516c6b765244646e643039335a6e564c61584579436c52335355524255554643436930744c5330745255354549464a545153425156554a4d53554d675330565a4c5330744c53304b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000030b860fe4b61a3f9295a08b53f6676443fb3ba19ed4502540a18c9d30268c4f7018b3b039d6528f06f08a63e4cb0ca67d30000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000001586859562b4378525131614d4f4d686b4d5174706d4650316c524732433454545a51414135654866746a63484a6b6c745872364b7877614c6f4d704f367170313333592f637a325a362b5553565554702b797655312b4e49306152504835534362557058376e7057497148306f4756656f674a4550514650374e65336c376a493761583876556a6b46596a7a636d52496e696e5a7370723455593344626953582f476a327a795644462f2b5941335a2f2f49366a77746561386c6c31576c3978412b794f70427a356b2b4b726134714a787452367a7668714d646532594836305a424a515342796934722b4c5035642f3279496b4f39344f7451326f4c41576d2f756778436f4d747152314e37444e4e633467614b444869722f7658676b514f475178424e396265353352345179386b69354869655953594e6e4b2f6d6e387374304e307250756442394a334b346278796657574f41413d3d0000000000000000",
  "blockNumber": "0x4a3a2e",
  "transactionHash": "0x20b673d0be280a38daa4f636ec6ad1108c0635dcb35c603f8e401a4120a2b506",
  "transactionIndex": "0x3",
  "blockHash": "0x579a98700bc9f9b1dc6ea3d00f9fd43bf28bd795f615210fd138fe724b8654d4",
  "logIndex": "0x2",
  "removed": false
}`
	rawValidatorAdded := `{
   "address":"0xd594c1ef4845713e86658cb42227a811625a285b",
   "topics":[
      "0x088097840a21a2c763dd9bd97cc2b0b27628bb6a42124a398260fac7f31ff571"
   ],
   "data":"0x0000000000000000000000004e409db090a71d14d32adbfbc0a22b1b06dde7de00000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000d200000000000000000000000000000000000000000000000000000000000000f4000000000000000000000000000000000000000000000000000000000000000308687eb8b88ff9c39e659c47b7bb76665fabfc4fc02c4246caca49700242fa9260a145969ede608b10c711ef2d57d0da1000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000003600000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000092000000000000000000000000000000000000000000000000000000000000002c0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000002644c5330744c5331435255644a54694253553045675546564354456c4449457446575330744c533074436b314a53554a4a616b464f516d64726357687261556335647a424351564646526b464254304e425554684254556c4a516b4e6e53304e42555556424e32705863457872656d643254586476527a684e64455679556a494b524768554d6b313164456c6d59556430566d784d654456574b326734616d7772646e6c7854315976636d784b5245566c517939484d7a567056304d3057455533526e464b55566331516d707651575a315458685165677052517a5a364d45453162314933656e52755748553263305633546b684a534668335245464954486c54645664514d334247596c6f30516e63356231465a54554a6d62564e734c33685852307379566e4e336156686b436b4e4663555a4b526d644e55466b334e6c4a5159306f325232646b545763725756525257565646616d6c52546a4670646d4a4b5a6a5257615570435254637262564e7465465a4e4e54417a566d6c7951575a6e646b494b656e426e64544e7a64485a496448705256315a3265484a304e545230526d39444d48526d5745315252584e53553056745456526f566b686f63566f725a544a434f43396b545751325231466f646e45355a58523152517068516b786f536c704655586c704d6b6c7055553032556c6732613031765a476447556d6376656d747454465a5851305649547a457a61465635526b6f78616e67314c304d3562454979553256454e57396a64316834436d4a525355524255554643436930744c5330745255354549464a545153425156554a4d53554d675330565a4c5330744c53304b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002c0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000002644c5330744c5331435255644a54694253553045675546564354456c4449457446575330744c533074436b314a53554a4a616b464f516d64726357687261556335647a424351564646526b464254304e425554684254556c4a516b4e6e53304e4255555642623370566147467a534739486545315953337055627a67725348634b57475630656e4a7457454e594f546474655842706148686a4c32777853456c6c53565677563256334e6b464e4d7a6c5064314a515a3256564d465a33516d51324e485a68627a5a7354544e615157785464565a6c4d677061626c4e305430314a636b4a5457475673596b633062314272524735785a6b4e4e62474a6d6131524e526c685856466f776445314964474a77566b55334e326f3061457078615549335a553133596974774e585578436c6f764e6d5678576a5a6d5257526e4f4449354d7a4e335a55686856574e7a64325a4a516d68594e6c4e61556a4e6c4d6b4a7652554a3262486c6a4e4535454e45466f4e5646615a6a4d7252577078536974356448594b63336869526d354d4e55704c5757686a536c6334596d7443647a4e6f4d3256726555597959324932655545334d336473547a5a68576b6c6152574a34516b453057446c34576a684d534642614e484a59574739476277706f4d564643643149784f555668656d463562306831546d4a6b574770426255396863315669543074744e464a42646b3979613146775a31493453304a344e474d7a637a6b304f466c696454424a526b745162304e49436b4a335355524255554643436930744c5330745255354549464a545153425156554a4d53554d675330565a4c5330744c53304b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002c0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000002644c5330744c5331435255644a54694253553045675546564354456c4449457446575330744c533074436b314a53554a4a616b464f516d64726357687261556335647a424351564646526b464254304e425554684254556c4a516b4e6e53304e4255555642636a5a586330396b4d7a4a5a5653745065566f7756565a74556c594b516b68455245744c4d3255314f545270557a56326448524c4d564a694d6c5659643359774e475a4b634764344c314e51576d6c71556d45306546646d63335a7361544d7865486731633273724d6c68364f544a3156516f35546c45344f47526c4c305978656d4a74616e51774d323577576a686153323533636d314c4f585a55524539505a4659344d3152694d554e59547a466862334a3265564d314d4552695a546c536248453253474e44436e567554545261516e6b30534864765a3270425a6a5932595446436330383565477832526a63305545677252544a3051316b305a5659774c314d3456466448626a6834523064495457354754306c31556d524d5554414b656d4d7651307050566a42494b316461534556455a5463794e5538775231417754585630516d4e485a57453152334134636b5a7757486b764d444642646d6c58616a426e4d4464714d4652314d30685a4e30646c53776f765a564e544c3168574f474a55524734344d305a516245353457486479566d6c33637a6c306347787a54464d78655578534e30787854324e5959566c344e48524c59334672565451305546686d656d395565433942436d68335355524255554643436930744c5330745255354549464a545153425156554a4d53554d675330565a4c5330744c53304b0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002c0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000002644c5330744c5331435255644a54694253553045675546564354456c4449457446575330744c533074436b314a53554a4a616b464f516d64726357687261556335647a424351564646526b464254304e425554684254556c4a516b4e6e53304e42555556426456497a5630686d55316c68576c45304e6a6b78656e52306154594b5a6c4242634578716132394c6379737251533930515764535658644862456858596d35694e6a4a5056553472613074545555353356576c4e4d4652775747644f564856534e47706a6457644b61314e54526c5253524170355745777653585270627a6c465a48453361456852513342455130784356464e59526c4e744d6a4a724e6c4e52626c6c476557733355564e6e646e6f795157396d4f584a3659566442516d566d556b5a5064557335436e465754303072627a686e526e467763586c51526e524a527939435653394662316c324d30464e5531413555574a4354585258536b4976635464325153745a4d5546725a454a6959554e756147466b4b3146555747774b5931566b537a526162485a314e566446576b784c6443394f4d6c5531524751776146683452584275526c6f334c3031534e56526e52566c324e466c336155704865574e795254464b5747565355324d724d3231445751704b656b567a596a4a50575442545a453833596a424d635764714d326856613052746345645653324e6f516c5179614777304e574a35616b3476616c5a6a555731726232396c5955677a53437432523249764e7a6856436b56335355524255554643436930744c5330745255354549464a545153425156554a4d53554d675330565a4c5330744c53304b000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000001a00000000000000000000000000000000000000000000000000000000000000030adb6d42245eaf4b00909679642964d6d5c12c4c550eaffcee499a12ea731c5f101f43a3880b9363daf873ae455fa7aa6000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000030b6de3081ad9a8becd37676827afb46386eeaa4cd7ebf8711a37505d3c5d3a7a3c1e167e3031e98094ed5262ec65ff205000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000030ad4754bd8ca755db23a0701d0dd5488403f9092912b09bcef95b8f70b380b528effd395fb3f06f92c515acf618f2cfa900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003097fceae9c1eeaeb5f9c8ebf875b7bc8248c514fe0c847cb2a15e662595ec6e214ebe3351f9b79629185008be0a1d1f5000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000240000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000005c000000000000000000000000000000000000000000000000000000000000001a0000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000001584d6a4e563955484666326f37644a76584d4c6867517866486430443945764a4f324e714c6c7753706a4c534b39497076663065514c49645a7a396e2f454941647a544734726c4f344a614332336270634946665036422b742f387a3379552b45425238785a35654235424c316f506730376672544d582f3951325a48305a7a6966716e535a372b672f396243483678675a495a634f5574687a30595141476752542f636c4d466b6162687a6d6a377172794972592f4577424a7a5164335363554d53586b4b4f65466d42496e78494241485238506e69495161597559734e314f4e48353571764b5833433452554a4175502b3675584b3949746d737353716b4249726851786a6f76696f6c4d51776b646b515038396a4c6d635835467062506973355771563856682f516b74376f4c734170306d5851546f67456d47566657426d4d736853464563384170446c6f5968344d547245413d3d000000000000000000000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000158575633707254546c593632746152366b42774f32773431544950613474796e574761596663444c672b46644b414c3061686d396136314a784276643676714c78634a7072346d472f61635053446e67657462624761425074494a423871556c4e6144796c744b7051675947526c394e657746736d69546570706c785769644d6a523445643663344154627a495a346c74486a6658683868582b772f7a7850704f4b5648344d7334414a6d50595a6835434c7057426e65554d436f4c412f6849556b71586a586d4a6c2b316d456a7270314a64526c6977762b37586f467379565570744839617767714167416d45415077454c4f575454526b6536482f4e2b774f334d526c6c4663726d476f555a73756d567a38452f523947744b32452f6864573356616934506c686a727552717a3965696d4f66564c764469774b505370546f5479675548717a72745a5047486e2b58716f64496b673d3d000000000000000000000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000158466a2b55417371307536524d74632f2b6f61446e5136386568427252745744424e37716b694c6e46526744634a6c454c426d6f6e37482f5056754a4279666775565639624c3968366b4b645844556955575873736a543048306c424e3454785475737a464d5149317158336851576f474735594a644b75596b3268362b3862572f77526c6c667033734544534b2f6e4a4936676c316e427a6a71414361794c505044466d47442b746e46767a6870765152476736334d54475969346c336939744d706f564e574573586249715773716d6f61383747695831354e59435a75397947476c66387567644d5a4b54324a66306345476c6b7957745856676433716b6b6e4d756d504c34746d2b305349636a3177434a456654737478757a614b546d44434e64315a4d72792b6d70457744414f684a6c4e36444b433348486e524d6a57672f4b376f30536d44654b504949644942344e7078673d3d000000000000000000000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000158475331774963466545787a5369587a3471384f4e6833734e765577542b493055517a70684773386a6638484b524d36615047617170384c4c65365568737941577659336373485863656144624356536178652f3556785061302f31574f34394a456c4c3346386f4975305649323769572b514d7a7735715944615a4b427476595351663379346f3036686436346c6d5a7a6c3855512b7570616161725746724f464e312b5572716565547130672b713747444b44536a614a7a524a5952566546572b7763456f5063662b515343314a5a6c5a653370314e643472524a7852304a41307477477978634f4f6a4e576b494479565866624a34766b6e72646e7539524b695a5056514d344a6a69796762684b302b516d43744a463436304935304750745247534a58756a4b3163786f6c56366c73536a557a4e51386e7141426838726278693651356b585478774155526d6b44456f2b52773d3d0000000000000000",
   "blockNumber":"0x5b5dc0",
   "transactionHash":"0x39fc924907817a759b41abd98353d3f94b9b1c159a796ad2f5339cbd2ed24dbd",
   "transactionIndex":"0x0",
   "blockHash":"0x021be90e25602cddc56386db2b690d427c05c9de288ca7c39f389157ba08c903",
   "logIndex":"0x2",
   "removed":false
}`

	t.Run("legacy validator added", func(t *testing.T) {
		vLogValidatorAdded, contractAbi := unmarshalLog(t, legacyRawValidatorAdded, Legacy)
		abiParser := NewParser(logex.Build("test", zap.InfoLevel, nil), Legacy)
		parsed, isEventBelongsToOperator, err := abiParser.ParseValidatorAddedEvent(nil, vLogValidatorAdded.Data, contractAbi)
		require.NoError(t, err)
		require.NotNil(t, contractAbi)
		require.False(t, isEventBelongsToOperator)
		require.NotNil(t, parsed)
		require.Equal(t, "91db3a13ab428a6c9c20e7104488cb6961abeab60e56cf4ba199eed3b5f6e7ced670ecb066c9704dc2fa93133792381c", hex.EncodeToString(parsed.PublicKey))
	})

	t.Run("v2 validator added", func(t *testing.T) {
		vLogValidatorAdded, contractAbi := unmarshalLog(t, rawValidatorAdded, V2)
		abiParser := NewParser(logex.Build("test", zap.InfoLevel, nil), V2)
		parsed, isEventBelongsToOperator, err := abiParser.ParseValidatorAddedEvent(nil, vLogValidatorAdded.Data, contractAbi)
		require.NoError(t, err)
		require.NotNil(t, contractAbi)
		require.False(t, isEventBelongsToOperator)
		require.NotNil(t, parsed)
		require.Equal(t, "8687eb8b88ff9c39e659c47b7bb76665fabfc4fc02c4246caca49700242fa9260a145969ede608b10c711ef2d57d0da1", hex.EncodeToString(parsed.PublicKey))
		require.Equal(t, "0x4e409dB090a71D14d32AdBFbC0A22B1B06dde7dE", parsed.OwnerAddress.Hex())
		operators := []string{"LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBN2pXcExremd2TXdvRzhNdEVyUjIKRGhUMk11dElmYUd0VmxMeDVWK2g4amwrdnlxT1YvcmxKREVlQy9HMzVpV0M0WEU3RnFKUVc1QmpvQWZ1TXhQegpRQzZ6MEE1b1I3enRuWHU2c0V3TkhJSFh3REFITHlTdVdQM3BGYlo0Qnc5b1FZTUJmbVNsL3hXR0syVnN3aVhkCkNFcUZKRmdNUFk3NlJQY0o2R2dkTWcrWVRRWVVFamlRTjFpdmJKZjRWaUpCRTcrbVNteFZNNTAzVmlyQWZndkIKenBndTNzdHZIdHpRV1Z2eHJ0NTR0Rm9DMHRmWE1RRXNSU0VtTVRoVkhocVorZTJCOC9kTWQ2R1FodnE5ZXR1RQphQkxoSlpFUXlpMklpUU02Ulg2a01vZGdGUmcvemttTFZXQ0VITzEzaFV5Rkoxang1L0M5bEIyU2VENW9jd1h4CmJRSURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K", "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBb3pVaGFzSG9HeE1YS3pUbzgrSHcKWGV0enJtWENYOTdteXBpaHhjL2wxSEllSVVwV2V3NkFNMzlPd1JQZ2VVMFZ3QmQ2NHZhbzZsTTNaQWxTdVZlMgpablN0T01JckJTWGVsYkc0b1BrRG5xZkNNbGJma1RNRlhXVFowdE1IdGJwVkU3N2o0aEpxaUI3ZU13YitwNXUxClovNmVxWjZmRWRnODI5MzN3ZUhhVWNzd2ZJQmhYNlNaUjNlMkJvRUJ2bHljNE5ENEFoNVFaZjMrRWpxSit5dHYKc3hiRm5MNUpLWWhjSlc4YmtCdzNoM2VreUYyY2I2eUE3M3dsTzZhWklaRWJ4QkE0WDl4WjhMSFBaNHJYWG9GbwpoMVFCd1IxOUVhemF5b0h1TmJkWGpBbU9hc1ViT0ttNFJBdk9ya1FwZ1I4S0J4NGMzczk0OFlidTBJRktQb0NICkJ3SURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K", "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBcjZXc09kMzJZVStPeVowVVZtUlYKQkhEREtLM2U1OTRpUzV2dHRLMVJiMlVYd3YwNGZKcGd4L1NQWmlqUmE0eFdmc3ZsaTMxeHg1c2srMlh6OTJ1VQo5TlE4OGRlL0YxemJtanQwM25wWjhaS253cm1LOXZURE9PZFY4M1RiMUNYTzFhb3J2eVM1MERiZTlSbHE2SGNDCnVuTTRaQnk0SHdvZ2pBZjY2YTFCc085eGx2Rjc0UEgrRTJ0Q1k0ZVYwL1M4VFdHbjh4R0dITW5GT0l1UmRMUTAKemMvQ0pPVjBIK1daSEVEZTcyNU8wR1AwTXV0QmNHZWE1R3A4ckZwWHkvMDFBdmlXajBnMDdqMFR1M0hZN0dlSwovZVNTL1hWOGJURG44M0ZQbE54WHdyVml3czl0cGxzTFMxeUxSN0xxT2NYYVl4NHRLY3FrVTQ0UFhmem9UeC9BCmh3SURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K", "LS0tLS1CRUdJTiBSU0EgUFVCTElDIEtFWS0tLS0tCk1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBdVIzV0hmU1lhWlE0NjkxenR0aTYKZlBBcExqa29LcysrQS90QWdSVXdHbEhXYm5iNjJPVU4ra0tTUU53VWlNMFRwWGdOVHVSNGpjdWdKa1NTRlRSRAp5WEwvSXRpbzlFZHE3aEhRQ3BEQ0xCVFNYRlNtMjJrNlNRbllGeWs3UVNndnoyQW9mOXJ6YVdBQmVmUkZPdUs5CnFWT00rbzhnRnFwcXlQRnRJRy9CVS9Fb1l2M0FNU1A5UWJCTXRXSkIvcTd2QStZMUFrZEJiYUNuaGFkK1FUWGwKY1VkSzRabHZ1NVdFWkxLdC9OMlU1RGQwaFh4RXBuRlo3L01SNVRnRVl2NFl3aUpHeWNyRTFKWGVSU2MrM21DWQpKekVzYjJPWTBTZE83YjBMcWdqM2hVa0RtcEdVS2NoQlQyaGw0NWJ5ak4valZjUW1rb29lYUgzSCt2R2IvNzhVCkV3SURBUUFCCi0tLS0tRU5EIFJTQSBQVUJMSUMgS0VZLS0tLS0K"}
		for i, pk := range parsed.OperatorPublicKeys {
			require.Equal(t, operators[i], string(pk))
		}
		shares := []string{"adb6d42245eaf4b00909679642964d6d5c12c4c550eaffcee499a12ea731c5f101f43a3880b9363daf873ae455fa7aa6", "b6de3081ad9a8becd37676827afb46386eeaa4cd7ebf8711a37505d3c5d3a7a3c1e167e3031e98094ed5262ec65ff205", "ad4754bd8ca755db23a0701d0dd5488403f9092912b09bcef95b8f70b380b528effd395fb3f06f92c515acf618f2cfa9", "97fceae9c1eeaeb5f9c8ebf875b7bc8248c514fe0c847cb2a15e662595ec6e214ebe3351f9b79629185008be0a1d1f50"}
		for i, pk := range parsed.SharesPublicKeys {
			require.Equal(t, shares[i], hex.EncodeToString(pk))
		}
	})
}

func unmarshalLog(t *testing.T, rawOperatorAdded string, abiVersion Version) (*types.Log, abi.ABI) {
	var vLogOperatorAdded types.Log
	err := json.Unmarshal([]byte(rawOperatorAdded), &vLogOperatorAdded)
	require.NoError(t, err)
	contractAbi, err := abi.JSON(strings.NewReader(ContractABI(abiVersion)))
	require.NoError(t, err)
	require.NotNil(t, contractAbi)
	return &vLogOperatorAdded, contractAbi
}