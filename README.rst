===============================
Go SDK for the OpenTok Platform
===============================

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
  
  import opentok "github.com/eauge/opentok-go-sdk"

And run the code below wrapped in main ::
  
	var (
		apiKey    = 44603982
		apiSecret = "1963c1345671419ea659b0feaa47b1206471463b"
		ot        = opentok.OpenTok{ApiKey: apiKey, ApiSecret: apiSecret}
		session   *opentok.Session
		token     *opentok.Token
		err       error
	)

	if session, err = opentok.NewSession(ot, opentok.SessionOptions{}); err != nil {
		log.Fatal("Error, session object could not be created: ", err)
	}
	if err := session.Create(); err != nil {
		log.Fatal("Error, session could not be created: ", err)
	}
	if token, err = session.Token(opentok.TokenProperties{}); err != nil {
		log.Fatal("Error in creating token ", err)
	}

	// We print the session and the token created
	fmt.Println("Session created", session.Id)
	fmt.Println("Token created", token.Value())
  
	
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

Sample Applications:
----------------
Find below links to sample applications developed using this sdk:

- `Sample <https://github.com/eauge/opentok-go-sample/>`_ showcasing the most basic functionality of the sdk.

- `Archiving sample <https://github.com/eauge/opentok-go-archiving/>`_ showcasing session creation and archiving.

What Comes Next:
----------------
The next step is to use the Session and the Token that you have created and
create an app. Visit https://tokbox.com/opentok/ to learn more about opentok,
how it works and what you can do with it.

  
