package sync

import (
	"fmt"
	"log"
	"sync"

	"github.com/barelyhuman/musync/account"
	"github.com/barelyhuman/musync/utils"
	"github.com/zmb3/spotify"
)

type Syncer struct {
	source        string
	destination   string
	authenticator *account.Auth
}

type Options func(*Syncer)

func (s *Syncer) getAllTrackIds(tracks *spotify.SavedTrackPage) (result []string) {
	for page := 1; ; page++ {
		ids := utils.PickField(tracks.Tracks, func(track spotify.SavedTrack) string { return track.ID.String() })
		result = append(result, ids...)
		err := s.authenticator.Client.NextPage(tracks)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	return
}

func (s *Syncer) getAllPlaylistIds(playlist *spotify.FullPlaylist) (result []string) {
	for page := 1; ; page++ {
		ids := utils.PickField(playlist.Tracks.Tracks, func(item spotify.PlaylistTrack) string { return string(item.Track.ID) })
		result = append(result, ids...)
		err := s.authenticator.Client.NextPage(&playlist.Tracks)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	return result
}

func (s *Syncer) Sync() error {
	var syncBatchWG sync.WaitGroup

	targetPlaylistID := s.destination
	user, err := s.authenticator.Client.CurrentUser()

	if err != nil {
		return fmt.Errorf("failed to get user...\n Error: %v", err)
	}

	fmt.Println("Logged in as:", user.DisplayName)

	var sourceTrackIds []string

	if s.source == "lib" {
		tracks, err := s.authenticator.Client.CurrentUsersTracks()
		if err != nil {
			return err
		}

		sourceTrackIds = s.getAllTrackIds(tracks)
		fmt.Println("Processing User Library...")

		fmt.Println("Total Tracks in Library:", tracks.Total)
	} else {
		sourcePlaylist, err := s.authenticator.Client.GetPlaylist(spotify.ID(s.source))
		if err != nil {
			return err
		}
		fmt.Println("Processing Source Playlist:", sourcePlaylist.Name)
		sourceTrackIds = s.getAllPlaylistIds(sourcePlaylist)
	}

	targetPlaylist, err := s.authenticator.Client.GetPlaylist(spotify.ID(targetPlaylistID))
	if err != nil {
		log.Fatal("Failed while trying to get playlist, Error:", err)
	}

	fmt.Println("Comparing Library and Playlist:", targetPlaylist.Name)

	fmt.Println("Please wait...")

	trackIdsInDestinationPlaylist := s.getAllPlaylistIds(targetPlaylist)

	uniqueTrackIds := []spotify.ID{}

	for _, idFromSource := range sourceTrackIds {
		if !utils.DoesSliceContain(trackIdsInDestinationPlaylist, idFromSource) {
			uniqueTrackIds = append(uniqueTrackIds, spotify.ID(idFromSource))
		}
	}

	fmt.Printf("Total Tracks to be Moved: %d\n", len(uniqueTrackIds))

	if len(uniqueTrackIds) < 1 {
		fmt.Println("Playlists already synced")
	}

	batches := utils.Chunk(uniqueTrackIds, 100)

	for _, batch := range batches {
		syncBatchWG.Add(1)
		batch := batch
		go func() {
			defer syncBatchWG.Done()
			_, err = s.authenticator.Client.AddTracksToPlaylist(spotify.ID(targetPlaylistID), batch...)
			if err != nil {
				log.Fatal("Failed to sync library with playlist,Error:", err)
			}
		}()
	}

	syncBatchWG.Wait()

	fmt.Println("Synced!")

	return nil
}

func NewSyncer(source, destination string, opts ...Options) *Syncer {
	syncer := &Syncer{}

	syncer.source = source
	syncer.destination = destination

	for _, opt := range opts {
		opt(syncer)
	}

	return syncer
}

func WithAuthenticator(auth *account.Auth) Options {
	return func(s *Syncer) {
		s.authenticator = auth
	}
}
