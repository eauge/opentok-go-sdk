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

Test your API_KEY
-----------------
You can use the cmd utility to test that your API_KEY and your API_SECRET work.
The following command will generate a sessionId and a token with your API_KEY and your API_SECRET
  $ API_KEY=API_KEY API_SECRET=API_SECRET go run cmd/cmd.go

Testing That Everything Works:
------------------------------
Once you get your API_KEY and API_SECRET, you can make sure that everything
works by running the tests that come with the SDK (the tests will not run with
go < 1.4 since we use testing.M)::
  go test github.com/eauge/opentok/

The Go OpenTok SDK works in combination with an OpenTok client. A developer
that wants to create a web application will need to add an OpenTok Server SDK
to her project and use the web client to write the other half of the application.
The web client documentation can be found: https://tokbox.com/opentok/libraries/client/js/


How To Get Started:
--------------------
You need to import::

  import opentok "github.com/eauge/opentok-go-sdk"

And run the code below wrapped in main ::

	apiKey := 123456
	apiSecret := "API_KEY"
	ot := opentok.New(apiKey, apiSecret)

	s, err := ot.Session(nil)
	if err != nil {
		panic(err)
	}

	t, err := ot.Token(s, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("session: ", s.ID)
	fmt.Println("token: ", t)


How Archiving Works:
--------------------
Create An Archive::

  archive, err := ot.ArchiveStart(session.ID, nil)

Stop An Archive::

  err := ot.ArchiveStop(archiveId)

Delete An Archive::

  err := ot.ArchiveDelete(archiveId)

Get An Archive::

  archive, err := ot.ArchiveGet(archiveId)

List All Archives linked to you API_KEY::

  archives, err := ot.ListArchives(0, 0)

What Comes Next:
----------------
The next step is to use the Session and the Token that you have created and
create an app. Visit https://tokbox.com/opentok/ to learn more about opentok,
how it works and what you can do with it.


