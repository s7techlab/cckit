# X.509 certificate for testing

generate self-signed certificate:

```
$ openssl genrsa -out private.key 3072
$ openssl req -new -x509 -key private.key -sha256 -out certificate.pem -days 730
```

You are about to be asked to enter information that will be incorporated into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.



examine the certificate:

`$ openssl x509 -in certificate.pem -text -noout`


````
Certificate:
     Data:
         Version: 1 (0x0)
         Serial Number: 13222534896082439009 (0xb77fe16e97334b61)
     Signature Algorithm: sha256WithRSAEncryption
         Issuer: C=RU, ST=Moscow, L=Moscow, O=S7Techlab, OU=Blockchain dept, CN=Victor Nosov/emailAddress=vitiko@mail.ru
         Validity
             Not Before: Apr 24 07:49:10 2018 GMT
             Not After : Jul  6 07:49:10 2018 GMT
         Subject: C=RU, ST=Moscow, L=Moscow, O=S7Techlab, OU=Blockchain dept, CN=Victor Nosov/emailAddress=vitiko@mail.ru
         Subject Public Key Info
````