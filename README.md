# rei3key
Generate keys for R3 low-code application.

Requirements:

R3 Low-code application must be compiled with an RSA public key.
You must have the corresponding private key pair to generate a valid license key.

To compile:

go build generate_prod.go

To use:

./generate_prod <license_parameter_filename>
Without the open/close brackets. Copy the output to a file.
The "license_parameter_filename" must be in the json format. 
See file license_param.json for the format.

or

./generate_prod <license_parameter_filename> > mylicense.lic

Open your R3 instance, go to System menu --> Activate License, upload mylicense.lic 
to activate the license.

How to generate public/private key pair:
# private key
openssl genrsa -out privkey.pem 2048

# public key
openssl rsa -in privkey.pem -pubout -out pubkey.pem

Convert to RSA format:
openssl rsa -in public.pem -pubin -RSAPublicKey_out -out rsa_public.pem

# Modify number of concurrent licenses (1:N)
1. Edit config/config.go on line: 93: return license.LoginCount * N
2. Edit www/stores/store.js on line 41: loginLimitedFactor:N, where N is the factor * logincount(unlimited users)
