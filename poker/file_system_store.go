package poker

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
)

type FileSystemPlayerStore struct {
	database *json.Encoder
	league   League
}

func initialisePlayerDBFile(file *os.File) error {
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		return err
	}

	if info.Size() == 0 {
		_, err = file.Write([]byte("[]"))
		if err != nil {
			return err
		}
		_, err = file.Seek(0, 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	err := initialisePlayerDBFile(file)

	if err != nil {
		return nil, fmt.Errorf("error initialising player db file %s, %q", file.Name(), err)
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
}

func FileSystemPlayerStoreFromFile(path string) (ps *FileSystemPlayerStore, closeFunc func(), err error) {
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, nil, fmt.Errorf("problem opening %s %v", path, err)
	}

	closeFunc = func() {
		err := db.Close()
		if err != nil {
			log.Fatalf("error closing database file %s", db.Name())
		}
	}

	store, err := NewFileSystemPlayerStore(db)

	if err != nil {
		return nil, nil, fmt.Errorf("problem creating file system player store, %v", err)
	}
	return store, closeFunc, nil
}

func (f *FileSystemPlayerStore) GetLeague() League {
	sort.Slice(f.league, func(i, j int) bool {
		return f.league[i].Wins > f.league[j].Wins
	})
	return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	league := f.GetLeague()

	player := league.Find(name)
	if player == nil {
		return 0
	}
	return player.Wins
}

func (f *FileSystemPlayerStore) RecordWin(name string) error {
	league := f.GetLeague()
	player := league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(league, Player{name, 1})
	}

	return f.database.Encode(f.league)
}
