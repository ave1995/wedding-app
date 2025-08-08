package factory

import "wedding-app/utils"

func (f *Factory) Cleanup() {
	if err := f.googleCloudClient.Close(); err != nil {
		f.Logger().Error("failed to close Google Cloud Storage Cliet:", utils.ErrAttr(err))
	}
}
