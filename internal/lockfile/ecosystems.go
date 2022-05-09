package lockfile

func KnownEcosystems() []Ecosystem {
	return []Ecosystem{
		NpmEcosystem,
		CargoEcosystem,
		BundlerEcosystem,
		ComposerEcosystem,
		GoEcosystem,
		MavenEcosystem,
		PipEcosystem,
	}
}
