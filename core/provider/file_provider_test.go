package provider

import (
	"testing"
)



func Test_File_provider_1(t *testing.T) {
	fr := NewFileProvider("not_exist_file_1", '\n')

	if fr == nil {
		t.Errorf("fr is nil")
		t.FailNow()
	}
	omsg, err, hasNext := fr.Read()
	if msg, ok := omsg.(string); !ok || len(msg) > 0 || err == nil || hasNext == true {
		t.Errorf("get from not_exist_file_1 %s %v %v", msg, err, hasNext)
		t.FailNow()
	}
}

func Test_File_provider_2(t *testing.T) {
	fr := NewFileProvider("file_provider_test.dat", '\n')

	if fr == nil {
		t.Errorf("fr is nil")
		t.FailNow()
	}
	omsg, err, hasNext := fr.Read()
	if msg, ok := omsg.(string); !ok || len(msg) == 0 || msg != "12345" || err != nil || hasNext != true {
		t.Errorf("get from file_provider_test.dat 1 %s %v %v", msg, err, hasNext)
		t.FailNow()
	}

	omsg, err, hasNext = fr.Read()
	if msg, ok := omsg.(string); !ok || len(msg) == 0 || msg != "qwerty" || err != nil || hasNext != true {
		t.Errorf("get from file_provider_test.dat 2 %s %v %v", msg, err, hasNext)
		t.FailNow()
	}

	omsg, err, hasNext = fr.Read()
	if msg, ok := omsg.(string); !ok || len(msg) != 0 || msg != "" || err != nil || hasNext != true {
		t.Errorf("get from file_provider_test.dat 3 %s %v %v", msg, err, hasNext)
		t.FailNow()
	}

	omsg, err, hasNext = fr.Read()
	if msg, ok := omsg.(string); !ok || len(msg) == 0 || msg != "zxcvbnnm" || err != nil || hasNext != true {
		t.Errorf("get from file_provider_test.dat 4 %s %v %v", msg, err, hasNext)
		t.FailNow()
	}

	omsg, err, hasNext = fr.Read()
	if msg, ok := omsg.(string); !ok || len(msg) != 0 || msg != "" || err != nil || hasNext == true {
		t.Errorf("get from file_provider_test.dat 5 %s %v %v", msg, err, hasNext)
		t.FailNow()
	}

	//openand read again
	omsg, err, hasNext = fr.Read()
	if msg, ok := omsg.(string);  !ok || len(msg) == 0 || msg != "12345" || err != nil || hasNext != true {
		t.Errorf("get from file_provider_test.dat 1 %s %v %v", msg, err, hasNext)
		t.FailNow()
	}

}
