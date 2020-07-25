REM Output Files
REM ca.key: Certificate authority private key file (this should not be shared IRL)
REM ca.crt: Certificate authority trust certificate (this should be shared IRL)
REM server.key: Server private key, password protected (this shouldnt be shared)
REM server.csr: Server certificate signing request (this should be shared with the CA owner)
REM server.crt: Server certificate signed by the CA (this would be sent back by the CA owner and kept on server)
REM server.pem: Conversion of server.key into a format gRPC likes (this shouldnt be shared)

REM private files: ca.key, server.key, server.pem, server.crt
REM share files: ca.crt (needed by client), server.csr (needed by CA)

REM Step 1: Generate Certificate Authority + Trust Certificate (ca.crt)
openssl genrsa -passout pass:1111 -des3 -out ca.key 4096
REM the server is claiming to be localhost (thats its common name, would normally be some domain name)
openssl req -passin pass:1111 -new -x509 -days 365 -key ca.key -out ca.crt -subj "/CN=localhost"

REM Step 2: Generate the server private key (server.key)
openssl genrsa -passout pass:1111 -des3 -out server.key 4096

REM Step 3: Get a certificate signing request for the CA (server.csr)
openssl req -passin pass:1111 -new -key server.key -out server.csr -subj "/CN=localhost"

REM Step 4: Sign the certificate with the CA we created (its called self signing) - server.crt
openssl x509 -req -passin pass:1111 -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial -01 -out server.crt

REM Step 5: Convert the server certificate to .pem format (server.pem) - usable by gRPC
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out server.pem