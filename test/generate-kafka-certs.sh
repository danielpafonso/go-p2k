# kafka-dev-with-docker/part-09/generate.sh
#!/usr/bin/env bash
 
set -eu

CN="${CN:-kafka-admin}"
PASSWORD="${PASSWORD:-password}"
TO_GENERATE_PEM="${CITY:-yes}"

VALIDITY_IN_DAYS=3650
CA_WORKING_DIRECTORY="certs/certificate-authority"
TRUSTSTORE_WORKING_DIRECTORY="certs/truststore"
KEYSTORE_WORKING_DIRECTORY="certs/keystore"
PEM_WORKING_DIRECTORY="certs/pem"
CA_KEY_FILE="ca-key"
CA_CERT_FILE="ca-cert"
DEFAULT_TRUSTSTORE_FILE="kafka.truststore.jks"
KEYSTORE_SIGN_REQUEST="cert-file"
KEYSTORE_SIGN_REQUEST_SRL="ca-cert.srl"
KEYSTORE_SIGNED_CERT="cert-signed"
 
echo "Creating our certificate authority"
if [ -f "$CA_WORKING_DIRECTORY/$CA_KEY_FILE" ] && [ -f "$CA_WORKING_DIRECTORY/$CA_CERT_FILE" ]; then
  echo "  already created"
else
  rm -rf $CA_WORKING_DIRECTORY && mkdir -p $CA_WORKING_DIRECTORY
  openssl req -new -newkey rsa:4096 -days $VALIDITY_IN_DAYS -x509 -subj "/CN=$CN" \
    -keyout $CA_WORKING_DIRECTORY/$CA_KEY_FILE -out $CA_WORKING_DIRECTORY/$CA_CERT_FILE -nodes \
    2> /dev/null
fi



echo "Generate Keystores"
rm -rf $KEYSTORE_WORKING_DIRECTORY && mkdir -p $KEYSTORE_WORKING_DIRECTORY

echo "  kafka"
DNAME="CN=kafka"
KEY_STORE_FILE_NAME="kafka.keystore.jks"

# generate keystore file
keytool -genkey -keystore $KEYSTORE_WORKING_DIRECTORY/"$KEY_STORE_FILE_NAME" \
  -alias localhost -validity $VALIDITY_IN_DAYS -keyalg RSA \
  -noprompt -dname $DNAME -keypass $PASSWORD -storepass $PASSWORD \
  2> /dev/null
 
# Now a certificate signing request will be made to the keystore."
keytool -certreq -keystore $KEYSTORE_WORKING_DIRECTORY/"$KEY_STORE_FILE_NAME" \
  -alias localhost -file $KEYSTORE_SIGN_REQUEST -keypass $PASSWORD -storepass $PASSWORD \
  2> /dev/null

# sign keystore with CA
openssl x509 -req -CA $CA_WORKING_DIRECTORY/$CA_CERT_FILE \
  -CAkey $CA_WORKING_DIRECTORY/$CA_KEY_FILE \
  -in $KEYSTORE_SIGN_REQUEST -out $KEYSTORE_SIGNED_CERT \
  -days $VALIDITY_IN_DAYS -CAcreateserial \
  2> /dev/null
  # creates $CA_WORKING_DIRECTORY/$KEYSTORE_SIGN_REQUEST_SRL which is never used or needed.
 
# Add CA to keystore
keytool -keystore $KEYSTORE_WORKING_DIRECTORY/"$KEY_STORE_FILE_NAME" -alias CARoot \
  -import -file $CA_WORKING_DIRECTORY/$CA_CERT_FILE -keypass $PASSWORD -storepass $PASSWORD -noprompt \
  2> /dev/null
 
# echo "Now the keystore's signed certificate will be imported back into the keystore."
keytool -keystore $KEYSTORE_WORKING_DIRECTORY/"$KEY_STORE_FILE_NAME" -alias localhost \
  -import -file $KEYSTORE_SIGNED_CERT -keypass $PASSWORD -storepass $PASSWORD \
  2> /dev/null

# clear intermediate files
rm -f $CA_WORKING_DIRECTORY/$KEYSTORE_SIGN_REQUEST_SRL $KEYSTORE_SIGN_REQUEST $KEYSTORE_SIGNED_CERT


echo "  client"
DNAME="CN=client"
KEY_STORE_FILE_NAME="client.keystore.jks"

# generate keystore file
keytool -genkey -keystore $KEYSTORE_WORKING_DIRECTORY/"$KEY_STORE_FILE_NAME" \
  -alias localhost -validity $VALIDITY_IN_DAYS -keyalg RSA \
  -noprompt -dname $DNAME -keypass $PASSWORD -storepass $PASSWORD \
  2> /dev/null
 
# Now a certificate signing request will be made to the keystore."
keytool -certreq -keystore $KEYSTORE_WORKING_DIRECTORY/"$KEY_STORE_FILE_NAME" \
  -alias localhost -file $KEYSTORE_SIGN_REQUEST -keypass $PASSWORD -storepass $PASSWORD \
  2> /dev/null

# sign keystore with CA
openssl x509 -req -CA $CA_WORKING_DIRECTORY/$CA_CERT_FILE \
  -CAkey $CA_WORKING_DIRECTORY/$CA_KEY_FILE \
  -in $KEYSTORE_SIGN_REQUEST -out $KEYSTORE_SIGNED_CERT \
  -days $VALIDITY_IN_DAYS -CAcreateserial \
  2> /dev/null
  # creates $CA_WORKING_DIRECTORY/$KEYSTORE_SIGN_REQUEST_SRL which is never used or needed.
 
# Add CA to keystore
keytool -keystore $KEYSTORE_WORKING_DIRECTORY/"$KEY_STORE_FILE_NAME" -alias CARoot \
  -import -file $CA_WORKING_DIRECTORY/$CA_CERT_FILE -keypass $PASSWORD -storepass $PASSWORD -noprompt \
  2> /dev/null
 
# echo "Now the keystore's signed certificate will be imported back into the keystore."
keytool -keystore $KEYSTORE_WORKING_DIRECTORY/"$KEY_STORE_FILE_NAME" -alias localhost \
  -import -file $KEYSTORE_SIGNED_CERT -keypass $PASSWORD -storepass $PASSWORD \
  2> /dev/null

# clear intermediate files
rm -f $CA_WORKING_DIRECTORY/$KEYSTORE_SIGN_REQUEST_SRL $KEYSTORE_SIGN_REQUEST $KEYSTORE_SIGNED_CERT


echo "Generate trust store"
# Generate trust store
rm -rf $TRUSTSTORE_WORKING_DIRECTORY && mkdir -p $TRUSTSTORE_WORKING_DIRECTORY
keytool -keystore $TRUSTSTORE_WORKING_DIRECTORY/$DEFAULT_TRUSTSTORE_FILE \
  -alias CARoot -import -file $CA_WORKING_DIRECTORY/$CA_CERT_FILE \
  -noprompt -dname "CN=$CN" -keypass $PASSWORD -storepass $PASSWORD \
  2> /dev/null

if [ $TO_GENERATE_PEM == "yes" ]; then
  echo "Generate SSL files"
  # Generate SSL files for non-java client
  rm -rf $PEM_WORKING_DIRECTORY && mkdir -p $PEM_WORKING_DIRECTORY

  keytool -exportcert -alias CARoot -keystore $KEYSTORE_WORKING_DIRECTORY/client.keystore.jks \
    -rfc -file $PEM_WORKING_DIRECTORY/ca-root.pem -storepass $PASSWORD \
    2> /dev/null

  keytool -exportcert -alias localhost -keystore $KEYSTORE_WORKING_DIRECTORY/client.keystore.jks \
    -rfc -file $PEM_WORKING_DIRECTORY/client-certificate.pem -storepass $PASSWORD \
    2> /dev/null

  keytool -importkeystore -srcalias localhost -srckeystore $KEYSTORE_WORKING_DIRECTORY/client.keystore.jks \
    -destkeystore cert_and_key.p12 -deststoretype PKCS12 -srcstorepass $PASSWORD -deststorepass $PASSWORD \
    2> /dev/null

  openssl pkcs12 -in cert_and_key.p12 -nocerts -nodes -password pass:$PASSWORD \
    | awk '/-----BEGIN PRIVATE KEY-----/,/-----END PRIVATE KEY-----/' > $PEM_WORKING_DIRECTORY/client-private-key.pem \
    2> /dev/null
  rm -f cert_and_key.p12
fi
