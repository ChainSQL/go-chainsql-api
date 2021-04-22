package keymanage

//SignData is used to sign data, verify signed data, etc
type SignData struct {
}

func (o *SignData) SignTransaction(sPrivateKey string, sTx string, pSignedData *[]byte, pHash *[]byte) bool {
	return true
}

func (o *SignData) SignPlainData(sPrivateKey string, sPlain string, pSignedData *[]byte) bool {
	return true
}

func (o *SignData) getSerialData(sJson string) bool {
	jsonparser.ObjectEach([]byte(sJson), func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
		fmt.Printf("Key: '%s'\n Value: '%s'\n Type: %s\n", string(key), string(value), dataType)
		return nil})


	BOOST_FOREACH(ptree::value_type &v1, objJson)
	{
		auto &fieldInfo = g_map_field.at(v1.first);
		int iKey = fieldInfo.enumFildID;
		iKey <<= 16;
		iKey += fieldInfo.iOrder;
		mapField.insert(std::make_pair(iKey, v1.first));
	}
	return true
}
