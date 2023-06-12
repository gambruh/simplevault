package helpers

import "testing"

func TestCompareTwoMaps(t *testing.T) {
	mapServer := map[string]struct{}{
		"file1":  {},
		"file2":  {},
		"file3":  {},
		"file10": {},
	}

	mapLocal := map[string]struct{}{
		"file1": {},
		"file2": {},
		"file4": {},
		"file5": {},
	}

	expectedToUpload := map[string]struct{}{
		"file4": {},
		"file5": {},
	}
	expectedToDownload := map[string]struct{}{
		"file3":  {},
		"file10": {},
	}

	toUpload, toDownload := CompareTwoMaps(mapServer, mapLocal)

	if len(toUpload) != len(expectedToUpload) {
		t.Errorf("Unexpected number of files to upload. Expected: %d, Got: %d", len(expectedToUpload), len(toUpload))
	}

	for file := range expectedToUpload {
		if _, ok := toUpload[file]; !ok {
			t.Errorf("File %s is missing from toUpload", file)
		}
	}

	if len(toDownload) != len(expectedToDownload) {
		t.Errorf("Unexpected number of files to download. Expected: %d, Got: %d", len(expectedToDownload), len(toDownload))
	}

	for file := range expectedToDownload {
		if _, ok := toDownload[file]; !ok {
			t.Errorf("File %s is missing from toDownload", file)
		}
	}
}
