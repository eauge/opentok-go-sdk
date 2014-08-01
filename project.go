package opentok

import (
	"errors"
	"fmt"
)

// Struct that contains all the necessary information to
// create sessions and interact with the OpenTok platform
type OpenTok struct {

	// This is the ApiKey that you get after creating
	// a project at the OpenTok Dashboard
	ApiKey int

	// This is the ApiSecret that you get after
	// creating a project with the OpenTok Dashboard
	ApiSecret string

	// This is just used internally to test the sdk
	// for the different environments
	apiUrl string
}

// We create this new type token for the tokens created
// by a session.
type Token string

// Validates that the necessary parameters in the OpenTok
// structure have been set
func validate(ot OpenTok) error {
	if ot.ApiKey == 0 {
		return errors.New("ApiKey is not set")
	}

	if len(ot.ApiSecret) == 0 {
		return errors.New("ApiSecret is not set")
	}

	return nil
}

// Retrieves an archive from the server. If the
// archive does not exist an error will be raised
func GetArchive(ot OpenTok, archiveId string) (Archive, error) {
	if archiveId == "" {
		return Archive{}, errors.New("Archive id cannot be empty")
	}

	var (
		client = newHttpClient(ot)
		url = fmt.Sprintf("v2/partner/%d/archive/%s", ot.ApiKey, archiveId)
		response, err = client.get(url, nil)
		a Archive
	)

	if err != nil {
		return Archive{}, err
	}

	a, err = decodeArchive(response)
	if err != nil {
		return Archive{}, errors.New(fmt.Sprintf("Archive could not be decoded successfully"))
	}
	return a, nil
}

// Deletes an existing archive with status available. If
// the archive is in any other state the operation will
// fail and return an error
func DeleteArchive(ot OpenTok, archiveId string) error {
	if archiveId == "" {
		return errors.New("Archive id cannot be empty")
	}

	var (
		client = newHttpClient(ot)
		url = fmt.Sprintf("v2/partner/%d/archive/%s", ot.ApiKey, archiveId)
		err = client.delete(url, nil)
	)

	return err
}

// Returns a list of archives. If Count == 0, the limit of
// the number of archives returned by the server is limited
// by the server. Otherwise it will be count. Offset is
// useful for pagination
func ListArchives(ot OpenTok, count, offset int) (as []Archive, err error) {
	if count < 0 {
		return nil, errors.New("count cannot be smaller than 1")
	}

	var (
		client = newHttpClient(ot)
		response []byte
		url = fmt.Sprintf("v2/partner/%d/archive?offset=%d", ot.ApiKey, offset)
	)

	if count > 0 {
		url = fmt.Sprintf("%s&count=%d", url, count)
	}

	response, err = client.get(url, nil)
	if err != nil {
		return nil, err
	}

	as, err = decodeArchiveList(response)
	if err != nil {
		return nil, err
	}

	return as, nil
}
