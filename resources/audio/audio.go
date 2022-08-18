package audio

import (
	"log"

	"github.com/dhowden/tag"
	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/media"
	"github.com/gohugoio/hugo/resources/resource"
)

type Format int

const (
	MP3 Format = iota + 1
	AAC
	FLAC
	OGG
)

var (
	audioFormatsBySubType = map[string]Format{
		media.MP3Type.SubType:      MP3,
		media.AACType.SubType:      AAC,
		media.FLACType.SubType:     FLAC,
		media.OGGAudioType.SubType: OGG,
	}
)

func AudioFormatFromMediaSubType(sub string) (Format, bool) {
	f, found := audioFormatsBySubType[sub]
	return f, found
}

type Spec interface {
	ReadSeekCloser() (hugio.ReadSeekCloser, error)
}

func NewAudio(format Format, reader Spec) *Audio {
	return &Audio{format: format, reader: reader}
}

type Audio struct {
	format Format
	reader Spec
}

func (a Audio) AudioTags() *PreResourceTransformationTags {
	readSeeker, err := a.reader.ReadSeekCloser()
	if err != nil {
		log.Fatal(err)
	}
	m, err := tag.ReadFrom(readSeeker)
	if err != nil {
		log.Fatal(err)
	}
	track, discTrackCount := m.Track()
	disc, discCount := m.Disc()
	return &PreResourceTransformationTags{
		BaseAudioTags{
			Title:          m.Title(),
			Album:          m.Album(),
			Artist:         m.Artist(),
			AlbumArtist:    m.AlbumArtist(),
			Composer:       m.Composer(),
			Genre:          m.Genre(),
			Year:           m.Year(),
			Track:          track,
			DiscTrackCount: discTrackCount,
			Disc:           disc,
			DiscCount:      discCount,
			Lyrics:         m.Lyrics(),
			Comment:        m.Comment(),
		},
		PictureTag{
			Art: newArt(m.Picture()),
		},
	}
}

type AudioResource interface {
	resource.Resource
	AudioResourceOps
}

type Art struct {
	MimeType  media.Type
	Bytes     []byte
	Extension string
}

func newArt(picture *tag.Picture) *Art {
	if picture == nil {
		return nil
	}
	mimeType := media.FromContent(media.DefaultTypes, []string{picture.Ext}, picture.Data)
	if mimeType.MainType != "image" {
		return nil
	}
	return &Art{
		MimeType:  mimeType,
		Bytes:     picture.Data,
		Extension: picture.Ext,
	}
}

type BaseAudioTags struct {
	Title          string
	Album          string
	Artist         string
	AlbumArtist    string
	Composer       string
	Genre          string
	Year           int
	Track          int
	DiscTrackCount int
	Disc           int
	DiscCount      int
	Lyrics         string
	Comment        string
}

type PictureTag struct {
	Art *Art
}

type PreResourceTransformationTags struct {
	BaseAudioTags
	PictureTag
}

type AudioTags struct {
	BaseAudioTags
	Art *resource.Resource
}

type AudioResourceOps interface {
	AudioTags() *AudioTags
}
