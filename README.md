# patients-poc-stable

This root directory represents the backend REST API made for operational/admin type personnel to interact with registered patient data.

Problem: Build a website for patient registration. 
Requirements: 1. The patient has to submit their name, date of birth, phone number, email, address, photo (driver license) and appointment time to register. 2. The admin should be able to view all the registered patients from the website. 

Simply run ```go run main.go```
Your serverside health check you should now find reachable unsecured over at http at http://localhost:8888/health


-------------------------------------------------------


# API endpoints:

#### Retrieve all Patients

```/patients ``` (GET)
		
Example: (inbrowser)	
```
http://localhost:8888/patients
```

-------------------------------------------------------

#### Register to add a NEW Patient to collection

```/register``` (POST)

Required request data params: 
 `id=[int64]`
 
 Example: 
 ```
curl -H "Content-Type: application/json" -X POST -d '{"Id":12345,"firstname":"heather","lastname":"spankx","address":"1231232 some place Way","state":"GA","city":"Atlanta","zip":30106,"telephone":7706958723423}' http://localhost:8888/register
  ```
  

 #### Delete an item - Authentication required
 
 ```/delete/<some_ID> ``` (DELETE)
 

Example:

```
curl -X DELETE 'http://localhost:8888/delete/12345'
```

-------------------------------------------------------


 #### Create a NEW Patient record by Id or UPDATE existing
 
 ```/edit/<some_ID> ``` (DELETE)
 

Example:

```
curl --header "Content-Type: application/json" \
 --request PUT \
 --data '{"id":12345,"firstname":"heather","lastname":"spankx","dob":"","address":"1231232 some place Way","state":"GA","city":"Atlanta","zip":30106,"email":"","telephone":7706958723423,"appointment":"09/10/2022"}' http://localhost:8888/edit/12345
```

-------------------------------------------------------


#### Find the Patient record by Id

```/find/<some_id> ``` (GET)
		
Example: (inbrowser)	Returns <b> Patient  <b>
```

http://localhost:8888/find/12345
