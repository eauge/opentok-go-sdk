===============================
Go SDK for the OpenTok Platform
===============================

Warning :
--------------------
The code is not ready for building web applications. There's a a set
of tests that pass successfully but I want to add a couple of
sample apps before removing this warning. However, adventoruous developers
are more than welcome to try it out!

Besides, keep in mind that this is not an official Server SDK from 
TokBox. This is a community effort so developers that want to use Go 
can use it to build OpenTok powered web and mobile applications. 

How To Get the Code:
--------------------
::

  $ go get github.com/eauge/opentok

Before doing anything:
----------------------
The first you must do, is create an OpenTok account at TokBox website:
https://tokbox.com/
  
After signing up, you'll be able to create a project where you'll get an 
API_KEY and an API_SECRET that you'll be able to use to run the code
in the section below. Don't run the code below without an API_KEY and
API_SECRET because it's not going to work. 


Testing That Everything Works:
------------------------------
Once you get your API_KEY and API_SECRET, you can make sure that everything
works by running the tests that come with the SDK::

  API_KEY="*** YOUR API_KEY ***" API_SECRET="*** YOUR API_SECRET ***" go test github.com/eauge/opentok/

The Go OpenTok SDK works in combination with an OpenTok client. A developer
that wants to create a web application will need to add an OpenTok Server SDK
to her project and use the web client to write the other half of the application. 
The web client documentation can be found: https://tokbox.com/opentok/libraries/client/js/


How To Get Started:
--------------------
You need to import::
  
  import "github.com/eauge/opentok"

And run the code below wrapped in main ::
  
	var (
		apiKey  = 0
		apiSecret = "*** Your API SECRET ***"
		ot      = opentok.OpenTok{ApiKey: apiKey, ApiSecret: apiSecret}
		sp      = opentok.SessionOptions{}
	)

	session, err := opentok.NewSession(ot, sp)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error, session object could not be created: %s", err))
		return
	}

	err = session.Create()
	if err != nil {
		fmt.Println(fmt.Sprintf("Error, session could not be created: %s", err))
		return
	}
  
  fmt.Println("Session created", session.Id)

	var token opentok.Token
	token, err = session.GenerateToken(opentok.TokenProperties{})
	fmt.Println("Token created", token)
	
How To Interact With Archiving:
-------------------------------
Create An Archive::

  archive, err := session.StartArchive("my archive")

Stop An Archive::

  archive, err := session.StopArchive(archive.Id)

Delete An Archive::

  err := opentok.DeleteArchive(archive.Id)

Get An Archive::

  err := opentok.GetArchive(archive.Id)

List All Archives linked to you API_KEY::

  archives, err := opentok.ListArchives(archive.Id, 0, 0)

What Comes Next:
----------------
The next step is to use the Session and the Token that you have created and
create an app. Visit https://tokbox.com/opentok/ to learn more about opentok,
how it works and what you can do with it.

  
