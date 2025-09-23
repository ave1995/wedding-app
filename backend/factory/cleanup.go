package factory

import "wedding-app/utils"

func (f *Factory) Cleanup() {
	if err := f.googleCloudClient.Close(); err != nil {
		// RADEK: Jetbrains kontrolují typos: Cliet :D
		f.Logger().Error("failed to close Google Cloud Storage Cliet:", utils.ErrAttr(err))
	}
}
