package versionbundle

import (
	"encoding/json"
	"reflect"
	"sort"

	"github.com/coreos/go-semver/semver"
	"github.com/giantswarm/microerror"
)

// Bundles is a plain validation type for a list of version bundles. A
// list of version bundles is exposed by authorities. Lists of version bundles
// of multiple authorities are aggregated and grouped to reflect releases.
type Bundles []Bundle

func (b Bundles) Contain(item Bundle) bool {
	for _, bundle := range b {
		if reflect.DeepEqual(bundle, item) {
			return true
		}
	}

	return false
}

func (b Bundles) Validate() error {
	if len(b) == 0 {
		return microerror.Maskf(invalidBundlesError, "version bundles must not be empty")
	}

	if b.hasDuplicatedVersions() {
		return microerror.Maskf(invalidBundlesError, "version bundle versions must be unique")
	}

	perProviderVersions := make(map[string][]Bundle, 0)
	for _, v := range b {
		perProviderVersions[v.Provider] = append(perProviderVersions[v.Provider], v)
	}

	for _, v := range perProviderVersions {
		for _, bundle := range v {
			// Get all bundles from the same series.
			f, err := filterSameSeriesBundles(bundle, v)
			if err != nil {
				return microerror.Maskf(invalidBundlesError, err.Error())
			}

			// Check that version number increments over time.
			err = validateIncrementsOverTime(f)
			if err != nil {
				return microerror.Maskf(invalidBundlesError, err.Error())
			}
		}
	}

	for _, bundle := range b {
		err := bundle.Validate()
		if err != nil {
			return microerror.Maskf(invalidBundlesError, err.Error())
		}
	}

	bundleName := b[0].Name
	for _, bundle := range b {
		if bundle.Name != bundleName {
			return microerror.Maskf(invalidBundlesError, "name must be the same for all version bundles")
		}
	}

	for _, bundle := range b {
		if bundle.Deprecated && bundle.WIP {
			return microerror.Maskf(invalidBundlesError, "version bundles must not be deprecated and WIP")
		}
	}

	return nil
}

func (b Bundles) hasDuplicatedVersions() bool {
	for _, b1 := range b {
		var seen int

		for _, b2 := range b {
			if b1.Version == b2.Version && b1.Provider == b2.Provider {
				seen++

				if seen >= 2 {
					return true
				}
			}
		}
	}

	return false
}

func CopyBundles(bundles []Bundle) []Bundle {
	raw, err := json.Marshal(bundles)
	if err != nil {
		panic(err)
	}

	var copy []Bundle
	err = json.Unmarshal(raw, &copy)
	if err != nil {
		panic(err)
	}

	return copy
}

// filterSameSeriesBundles filters bundles that are from the same semver series as
// input bundle. For major version bundle (e.g. 3.0.0) this function will return
// all available major bundles. For minor all minor bundles from the same major bundle.
// For patch all patch versions from the same minor bundle.
func filterSameSeriesBundles(bundle Bundle, bundles []Bundle) ([]Bundle, error) {
	if len(bundles) == 0 {
		return []Bundle{}, microerror.Maskf(executionFailedError, "bundles must not be empty")
	}

	filtered := []Bundle{}

	// Get semver for input bundle.
	bv := semver.New(bundle.Version)

	// For the case, when input bundle is a major version.
	if bv.Minor == 0 && bv.Patch == 0 {
		for _, b := range bundles {
			v := semver.New(b.Version)
			// Filter only bundles that have minor and
			// patch numbers equal to 0.
			if v.Minor == 0 && v.Patch == 0 {
				filtered = append(filtered, b)
			}
		}
		return filtered, nil
	}

	// For the case, when input bundle is a minor version.
	if bv.Patch == 0 {
		for _, b := range bundles {
			v := semver.New(b.Version)
			// Filter only bundles that have equal major number
			// and patch number is 0.
			if v.Major == bv.Major && v.Patch == 0 {
				filtered = append(filtered, b)
			}
		}
		return filtered, nil
	}

	// Finally for the case, when input bundle is a patch version.
	for _, b := range bundles {
		v := semver.New(b.Version)
		// Filter only bundles that have equal major and minor numbers.
		if v.Major == bv.Major && v.Minor == bv.Minor {
			filtered = append(filtered, b)
		}
	}
	return filtered, nil
}

func GetBundleByName(bundles []Bundle, name string) (Bundle, error) {
	if len(bundles) == 0 {
		return Bundle{}, microerror.Maskf(executionFailedError, "bundles must not be empty")
	}
	if name == "" {
		return Bundle{}, microerror.Maskf(executionFailedError, "name must not be empty")
	}

	for _, b := range bundles {
		if b.Name == name {
			return b, nil
		}
	}

	return Bundle{}, microerror.Maskf(bundleNotFoundError, name)
}

func GetBundleByNameForProvider(bundles []Bundle, name, provider string) (Bundle, error) {
	if len(bundles) == 0 {
		return Bundle{}, microerror.Maskf(executionFailedError, "bundles must not be empty")
	}
	if name == "" {
		return Bundle{}, microerror.Maskf(executionFailedError, "name must not be empty")
	}
	if provider == "" {
		return Bundle{}, microerror.Maskf(executionFailedError, "provider must not be empty")
	}

	for _, b := range bundles {
		if b.Name == name && b.Provider == provider {
			return b, nil
		}
	}

	return Bundle{}, microerror.Maskf(bundleNotFoundError, name)
}

func GetNewestBundle(bundles []Bundle) (Bundle, error) {
	return GetNewestBundleForProvider(bundles, "")
}

func GetNewestBundleForProvider(bundles []Bundle, provider string) (Bundle, error) {
	if len(bundles) == 0 {
		return Bundle{}, microerror.Maskf(executionFailedError, "bundles must not be empty")
	}

	// filter bundles by provider if provider is specificed
	if provider != "" {
		for i := 0; i < len(bundles); i++ {
			if bundles[i].Provider != provider {
				bundles = append(bundles[:i], bundles[i+1:]...)
				i--
			}
		}
	}

	if len(bundles) == 0 {
		return Bundle{}, microerror.Maskf(bundleNotFoundError, "no bundle found for provider %s", provider)
	}

	s := SortBundlesByVersion(bundles)
	sort.Sort(s)

	return s[len(s)-1], nil
}

func validateIncrementsOverTime(bundles []Bundle) error {
	b1 := CopyBundles(bundles)
	b2 := CopyBundles(bundles)
	sort.Sort(SortBundlesByTime(b1))
	sort.Sort(SortBundlesByVersion(b2))
	if !reflect.DeepEqual(b1, b2) {
		return microerror.Maskf(invalidBundlesError, "version bundle versions must always increment")
	}
	return nil
}
