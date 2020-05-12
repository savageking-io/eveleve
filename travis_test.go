package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTravis_ParsePayload(t *testing.T) {

	pl := `{"id":685826671,"number":"7","config":{"language":"cpp","os":["linux"],"dist":"bionic","install":["pyenv global 3.7","pip install conan","conan user","sudo apt-get -y install libsdl2-dev libsdl2-gfx-dev libsdl2-image-dev libsdl2-mixer-dev libsdl2-net-dev libsdl2-ttf-dev","conan install --file ./.conanfile"],"script":["./configure","make"],"notifications":{"webhooks":[{"urls":["http://savageking.io:12811/travis"],"on_success":"always","on_failure":"always","on_start":"always","on_cancel":"always","on_error":"always"}]}},"type":"push","state":"errored","status":1,"result":1,"status_message":"Errored","result_message":"Errored","started_at":"2020-05-11T21:06:36Z","finished_at":"2020-05-11T21:07:36Z","duration":60,"build_url":"https://travis-ci.org/savageking-io/evelengine/builds/685826671","commit_id":209828901,"commit":"0421420c3da3ac558bda584e8dfa0c092eb68996","base_commit":null,"head_commit":null,"branch":"master","message":"Added travis notifications","compare_url":"https://github.com/savageking-io/evelengine/compare/9919622aefcc...0421420c3da3","committed_at":"2020-05-11T20:10:22Z","author_name":"vozgua","author_email":"criotos@gmail.com","committer_name":"vozgua","committer_email":"criotos@gmail.com","pull_request":false,"pull_request_number":null,"pull_request_title":null,"tag":null,"repository":{"id":28682454,"name":"evelengine","owner_name":"savageking-io","url":null},"matrix":[{"id":685826672,"repository_id":28682454,"parent_id":685826671,"number":"7.1","state":"errored","config":{"os":"linux","language":"cpp","dist":"bionic","install":["pyenv global 3.7","pip install conan","conan user","sudo apt-get -y install libsdl2-dev libsdl2-gfx-dev libsdl2-image-dev libsdl2-mixer-dev libsdl2-net-dev libsdl2-ttf-dev","conan install --file ./.conanfile"],"script":["./configure","make"]},"status":1,"result":1,"commit":"0421420c3da3ac558bda584e8dfa0c092eb68996","branch":"master","message":"Added travis notifications","compare_url":"https://github.com/savageking-io/evelengine/compare/9919622aefcc...0421420c3da3","started_at":"2020-05-11T21:06:36Z","finished_at":"2020-05-11T21:07:36Z","committed_at":"2020-05-11T20:10:22Z","author_name":"vozgua","author_email":"criotos@gmail.com","committer_name":"vozgua","committer_email":"criotos@gmail.com","allow_failure":null}]}`

	travis := new(Travis)
	res, err := travis.ParsePayload(pl)
	if err != nil {
		t.Errorf("Payload parse failed: %s", err.Error())
	}
	fmt.Printf("%+v", res)

	type fields struct {
		conf *TravisConfig
	}
	type args struct {
		pl string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TravisPacket
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Travis{
				conf: tt.fields.conf,
			}
			got, err := tr.ParsePayload(tt.args.pl)
			if (err != nil) != tt.wantErr {
				t.Errorf("Travis.ParsePayload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Travis.ParsePayload() = %v, want %v", got, tt.want)
			}
		})
	}
}
