#pragma once

#ifdef __cplusplus
extern "C" {
#endif

#include <stdbool.h>
	/*
	* return: sige data's length
	* param@out_data: signed data
	*
	* CAUTION: param@out_data needs to be free manually.
	*/
    unsigned long sign_common_data(const unsigned char* private_key, const unsigned long private_key_len,
        const unsigned char* in_data, const unsigned long in_data_len,
        unsigned char** out_data);
	bool verify_common_data(const unsigned char* public_key, const unsigned long public_key_len,
		const unsigned char* pPlainData, const unsigned long plainDataLen,
		const unsigned char* pSignedData, const unsigned long signedDataLen);

	/*
	* return: sige data's length
	* param@out_data: signed data
	*
	* CAUTION: param@out_data needs to be free manually.
	*/
	unsigned long sign_for_common_transaction(const unsigned char* private_key, const unsigned long private_key_len,
		const unsigned char* in_data, const unsigned long in_data_len,
		unsigned char** out_data);

	/*
	* return: sige data's length
	* param@out_data: signed data
	* param@out_hash: tx's hash
	* 
	* CAUTION: param@out_data and param@out_hash both need to be free manually.
	*/
	unsigned long sign_for_common_transaction_ex(const unsigned char* private_key, const unsigned long private_key_len,
		const unsigned char* in_data, const unsigned long in_data_len,
		unsigned char** out_data, unsigned char** out_hash, unsigned long *out_hash_len);
	/*
	* return: sige data's length
	* param@out_data: signed data
	*
	* CAUTION: param@out_data needs to be free manually.
	*/
	unsigned long sign_for_mulsinger_transaction(const unsigned char* private_key, const unsigned long private_key_len,
		const unsigned char* in_data, const unsigned long in_data_len,
		unsigned char** out_data);
	/*
	* return: sige data's length
	* param@out_data: signed data
	*
	* CAUTION: param@out_data needs to be free manually.
	*/
	unsigned long encrypt_data_AES(const unsigned char* key, const unsigned long key_len,
		const unsigned char* in_data, const unsigned long in_data_len,
		unsigned char** out_data);
	/*
	* return: sige data's length
	* param@out_data: signed data
	*
	* CAUTION: param@out_data needs to be free manually.
	*/
	unsigned long decrypt_data_AES(const unsigned char* key, const unsigned long key_len,
		const unsigned char* in_data, const unsigned long in_data_len,
		unsigned char** out_data);

	/*
	* return: sige data's length
	* param@out_data: signed data
	*
	* CAUTION: param@out_data needs to be free manually.
	*/
	unsigned long  encrypt_data_asym(const unsigned char* public_key, const unsigned long public_key_len,
		const unsigned char* in_data, const unsigned long in_data_len,
		unsigned char** out_data);
	/*
	* return: sige data's length
	* param@out_data: signed data
	*
	* CAUTION: param@out_data needs to be free manually.
	*/
	unsigned long  decrypt_data_asym(const unsigned char* private_key, const unsigned long private_key_len,
		const unsigned char* in_data, const unsigned long in_data_len,
		unsigned char** out_data);

	/*
	* return: success or fail
	*
	* CAUTION: param@seed and param@public_key both need to be free manually.
	*/
	bool validation_create(unsigned char** seed, unsigned long *seed_len,unsigned char** public_key, unsigned long* public_key_len);

	/*
	* return: success or fail
	*
	* CAUTION: param@account, param@public_key, param@public_key_hex and param@private_key all need to be free manually.
	*/
	bool get_valid_address(const unsigned char* seed, const unsigned long len_seed, unsigned char** account, unsigned long* len_accout,
		unsigned char** public_key, unsigned long* len_public_key,
		unsigned char** public_key_hex, unsigned long* len_public_key_hex,
		unsigned char** private_key, unsigned long* len_private_key);


#ifdef __cplusplus
}
#endif
