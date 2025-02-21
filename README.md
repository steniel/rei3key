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

