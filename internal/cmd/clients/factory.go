package clients

import (
	"io"

	"github.com/spotinst/spotinst-cli/internal/cloud"
	"github.com/spotinst/spotinst-cli/internal/dep"
	"github.com/spotinst/spotinst-cli/internal/log"
	"github.com/spotinst/spotinst-cli/internal/spotinst"
	"github.com/spotinst/spotinst-cli/internal/survey"
	"github.com/spotinst/spotinst-cli/internal/thirdparty"
	"github.com/spotinst/spotinst-cli/internal/writer"
)

type factory struct {
	in       io.Reader
	out, err io.Writer
}

func NewFactory(in io.Reader, out, err io.Writer) Factory {
	return &factory{
		in:  in,
		out: out,
		err: err,
	}
}

func (x *factory) NewSpotinst(options ...spotinst.ClientOption) (spotinst.Interface, error) {
	log.Debugf("Instantiating new spotinst client")
	return spotinst.New(options...), nil
}

func (x *factory) NewCloud(name cloud.ProviderName) (cloud.Interface, error) {
	log.Debugf("Instantiating new cloud (%s)", name)
	return cloud.GetInstance(name)
}

func (x *factory) NewCommand(name thirdparty.CommandName) (thirdparty.Command, error) {
	log.Debugf("Instantiating new command (%s)", name)
	return thirdparty.GetInstance(name, thirdparty.WithStdio(x.in, x.out, x.err))
}

func (x *factory) NewSurvey() (survey.Interface, error) {
	log.Debugf("Instantiating new survey")
	return survey.New(x.in, x.out, x.err), nil
}

func (x *factory) NewDep() (dep.Interface, error) {
	log.Debugf("Instantiating new dependency manager")
	return dep.New(survey.New(x.in, x.out, x.err)), nil
}

func (x *factory) NewWriter(format writer.Format) (writer.Writer, error) {
	log.Debugf("Instantiating new writer (%s)", format)
	return writer.GetInstance(format, x.out)
}