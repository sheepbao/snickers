package downloaders

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/flavioribeiro/gonfig"
	"github.com/snickers/snickers/db"
	"github.com/snickers/snickers/types"
)

var _ = Describe("HTTP downloader", func() {
	var (
		logger     *lagertest.TestLogger
		dbInstance db.Storage
		cfg        gonfig.Gonfig
	)

	BeforeEach(func() {
		logger = lagertest.NewTestLogger("http-download")
		currentDir, _ := os.Getwd()
		cfg, _ = gonfig.FromJsonFile(currentDir + "/../fixtures/config.json")
		dbInstance, _ = db.GetDatabase(cfg)
		dbInstance.ClearDatabase()
	})

	Context("HTTP Downloader", func() {
		It("should return an error if source couldn't be fetched", func() {
			exampleJob := types.Job{
				ID:          "123",
				Source:      "http://source_here.mp4",
				Destination: "s3://user@pass:/bucket/",
				Preset:      types.Preset{Name: "presetHere", Container: "mp4"},
				Status:      types.JobCreated,
				Details:     "",
			}
			dbInstance.StoreJob(exampleJob)
			err := HTTPDownload(logger, cfg, dbInstance, exampleJob.ID)
			Expect(err.Error()).To(SatisfyAny(ContainSubstring("no such host"), ContainSubstring("No filename could be determined")))
		})

		It("Should set the local source and local destination on Job", func() {
			exampleJob := types.Job{
				ID:          "123",
				Source:      "http://flv.io/source_here.mp4",
				Destination: "s3://user@pass:/bucket/",
				Preset:      types.Preset{Name: "240p", Container: "mp4"},
				Status:      types.JobCreated,
				Details:     "",
			}
			dbInstance.StoreJob(exampleJob)
			HTTPDownload(logger, cfg, dbInstance, exampleJob.ID)
			changedJob, _ := dbInstance.RetrieveJob("123")

			swapDir, _ := cfg.GetString("SWAP_DIRECTORY", "")

			sourceExpected := swapDir + "123/src/source_here.mp4"
			Expect(changedJob.LocalSource).To(Equal(sourceExpected))

			destinationExpected := swapDir + "123/dst/source_here_240p.mp4"
			Expect(changedJob.LocalDestination).To(Equal(destinationExpected))
		})
	})

})
