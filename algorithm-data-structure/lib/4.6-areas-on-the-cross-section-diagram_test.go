package lib

import "testing"

func BenchmarkAreasOnTheCrossSectionDiagramOptimized(b *testing.B) {
	args := `\\///\_/\/\\\\/_/\\///__\\\_\\/_\/_/\`
	for i := 0; i < b.N; i++ {
		AreasOnTheCrossSectionDiagramOptimized(args)
	}
}

func BenchmarkAreasOnTheCrossSectionDiagramOptimized_2(b *testing.B) {
	args := `\\///\_/\/\\\\/_/\\///__\\\_\\/_\/_/\`
	for i := 0; i < b.N; i++ {
		AreasOnTheCrossSectionDiagramOptimized_2(args)
	}
}
