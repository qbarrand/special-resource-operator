package registry

import (
	"crypto/x509"
	"fmt"
	"os"

	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/image/v5/pkg/sysregistriesv2"
	"github.com/containers/image/v5/types"
)

type MirrorResolver interface {
	GetPullSourcesForImageReference(string) ([]sysregistriesv2.PullSource, error)
}

type mirrorResolver struct {
	sysCtx *types.SystemContext
}

func NewMirrorResolver(registriesConfPath string) MirrorResolver {
	return &mirrorResolver{
		sysCtx: &types.SystemContext{
			SystemRegistriesConfPath: registriesConfPath,
		},
	}
}

func (mr *mirrorResolver) GetPullSourcesForImageReference(imageName string) ([]sysregistriesv2.PullSource, error) {
	r, err := sysregistriesv2.FindRegistry(mr.sysCtx, imageName)
	if err != nil {
		return nil, fmt.Errorf("could not find registry for image %q: %v", imageName, err)
	}

	n, err := reference.ParseNamed(imageName)
	if err != nil {
		return nil, fmt.Errorf("could not parse image name %q: %v", imageName, err)
	}

	return r.PullSourcesFromReference(n)
}

type CertPoolGetter interface {
	SystemAndHostCertPool() (*x509.CertPool, error)
}

type certPoolGetter struct {
	hostBundlePath string
}

func NewCertPoolGetter(hostBundlePath string) CertPoolGetter {
	return &certPoolGetter{hostBundlePath: hostBundlePath}
}

func (cpg *certPoolGetter) SystemAndHostCertPool() (*x509.CertPool, error) {
	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("could not access the system certificate pool: %v", err)
	}

	b, err := os.ReadFile(cpg.hostBundlePath)
	if err != nil {
		return nil, fmt.Errorf("could not open the host's certificate bundle at %q: %v", cpg.hostBundlePath, err)
	}

	if !pool.AppendCertsFromPEM(b) {
		return nil, fmt.Errorf("could not append host certificates to the pool: %v", err)
	}

	return pool, nil
}
