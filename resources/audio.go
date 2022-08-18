package resources

import (
	"bytes"
	"sync"

	"github.com/gohugoio/hugo/common/hugio"
	"github.com/gohugoio/hugo/resources/audio"
	"github.com/gohugoio/hugo/resources/images"
	"github.com/gohugoio/hugo/resources/page"
)

type audioResource struct {
	audio     *audio.Audio
	audioTags *audio.AudioTags
	metaInit  sync.Once
	baseResource
}

func (a *audioResource) AudioTags() *audio.AudioTags {
	a.metaInit.Do(func() {
		preTransformationTags := a.audio.AudioTags()
		if preTransformationTags == nil {
			a.audioTags = nil
			return
		}
		if preTransformationTags.Art == nil {
			a.audioTags = &audio.AudioTags{
				BaseAudioTags: preTransformationTags.BaseAudioTags,
			}
			return
		}

		art := preTransformationTags.Art
		_, ok := images.ImageFormatFromMediaSubType(art.MimeType.SubType)
		if ok {
			imageResource, err := a.getSpec().New(ResourceSourceDescriptor{
				RelTargetFilename: a.Title() + ".art." + art.Extension,
				TargetPaths: func() page.TargetPaths {
					a.Name()
					targetPaths := a.getResourcePaths().targetPathBuilder()
					targetPaths.TargetFilename = targetPaths.SubResourceBaseTarget + "/" + a.Title() + ".art." + art.Extension
					return targetPaths
				},
				LazyPublish: true,
				MediaType:   preTransformationTags.Art.MimeType,
				OpenReadSeekCloser: func() (hugio.ReadSeekCloser, error) {
					return hugio.NewReadSeekerNoOpCloser(bytes.NewReader(art.Bytes)), nil
				},
			})
			if err != nil {
				a.audioTags = &audio.AudioTags{
					BaseAudioTags: preTransformationTags.BaseAudioTags,
				}
				return
			} else {
				a.audioTags = &audio.AudioTags{
					BaseAudioTags: preTransformationTags.BaseAudioTags,
					Art:           &imageResource,
				}
			}
		} else {
			a.audioTags = &audio.AudioTags{
				BaseAudioTags: preTransformationTags.BaseAudioTags,
				Art:           nil,
			}
		}
	})
	return a.audioTags
}
