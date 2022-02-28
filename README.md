# Blaze Mailer
Blaze Mailer is simple implementation of SMTP Go package. 

It provides an HTTP endpoint to send emails.

Plus boilerplate code to handle your credentials with env vars.

### Mac Installation 

* Clone repo

      git clone git@github.com:nixmaldonado/blazeMailer.git   

* Configure 3 env vars
      
      export SMTP_PASSWORD="YOUR_PASSWORD"
      export SMPT_HOST="SMPT_HOST" //eg: "smtp.gmail.com"
      export SMTP_PORT="SMTP_PORT" //eg: "587" 
    
* Run program from root folder
    
      go run cmd/main.go
        
* Perform HTTP request to send mail:
 
      curl --request GET 'http://localhost:3000/email/send' \
      --header 'Content-Type: application/json' \
      --data-raw '{
          "from": "sender@mail.com",
          "to": ["recipient@mail.com"],
          "subject":"Thi is Blaze Mailer",
          "body": "Hello World!"
      }'
      
 