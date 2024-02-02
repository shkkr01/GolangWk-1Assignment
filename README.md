# Encrypt and Decrypt the given data in golang 

## Description

Key Features:

You are part of the security team at Webknot and have been given the task to transport some confidential data to a client. To prevent any loss, the data should be encrypted in transit. 
The data is stored in a text file (Approx 1GB) and each line in the file contains exactly one JSON object representing some important information. The goal is to read the file efficiently using different concepts of concurrency and parallelism, encrypt the contents using an encryption key and store the encrypted contents back to a file. The file will then be transferred to the client where he can use the encryption key to get back the original contents.

### Executing program
git pull https://github.com/shkkr01/GolangWk-1Assignment.git
Change the path of the data with the original path in MakeFile and the encryption key as well if needed 
encrypt:
	go run main.go --file imp-data.txt --key "bHB&fuw2GD*I@Bsd" --out outputpath.txt --encrypt

decrypt:
	go run main.go --file outputpath.txt --key "bHB&fuw2GD*I@Bsd" --out outputpathnew.txt

 ## Help
Iam using The Symmetric Encryption is one of these method.We are going to encrypt and decrypt a file using symmetric encryption
To encrypt a file, we are going to use crypto package that Go’s built-in package.
after we read the file we use block cipher algorithm (AES (Advanced Encryption System)) for this case 
send it to AES to create block cipher algorithm. (crypto/aes)
After created block of algorithm, we are going to use GCM (Galois/Counter Mode) mode. The GCM is a stream mode and provides data authenticity and confidentially.(crypto/cipher)
The nonce has to be unique and it changes every time when data is encrypted.
To generate random nonce, we are going to use the package crypto/rand

Encryption Time : 10.977630792s
Decryption Time : 12.306711166s
