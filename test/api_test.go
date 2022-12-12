package test

import (
	"sync"
	"testing"
	"web.misaki.world/FinalExam/xmly"
)

func TestGet(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "case 1",
			args: args{
				url: "https://www.ximalaya.com/revision/album/v1/simple?albumId=20403413",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := xmly.Get(tt.args.url)
			t.Log(int(res["ret"].(float64)))
		})
	}
}

func TestObtainDetailForAlbumId(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "case 1",
			args: args{
				id: 20403413,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg := sync.WaitGroup{}
			wg.Add(3)
			got := xmly.ObtainDetailForAlbumId(tt.args.id, &wg)
			go func() {
				got.WriteFile("20403413")
				wg.Done()
			}()
			go func() {
				got.DownLoadAudio("20403413/audio")
				wg.Done()
			}()
			wg.Wait()
		})
	}
}
