package identicon

func parseSymmetric(designSize int, blocks []bool) [][]bool {
	design := make([][]bool, designSize)

	// For each row in the first half, first move its blocks into the design,
	// then copy its blocks backwards into the opposite row.
	for iRow := 0; iRow < designSize / 2; iRow += 1 {
		minBlock := iRow * designSize
		maxBlock := minBlock + designSize - 1
		design[iRow] = blocks[minBlock: maxBlock + 1]

		row := make([]bool, designSize)
		for iBlock := 0; iBlock < designSize; iBlock += 1 {
			row[designSize - 1 - iBlock] = blocks[minBlock + iBlock]
		}
		design[designSize - 1 - iRow] = row
	}
	// If there is a middle row, then for each of the final few blocks,
	// copy it to the middle row of the design, and copy it into the opposite
	// position of the middle row.
	if designSize % 2 == 1 {
		minBlock := (designSize / 2) * designSize
		row := make([]bool, designSize)
		for iBlock := 0; iBlock < designSize / 2 + 1; iBlock += 1 {
			// Note: the centre block will be copied twice
			row[iBlock] = blocks[minBlock + iBlock]
			row[designSize - 1 - iBlock] = blocks[minBlock + iBlock]
		}
		design[designSize / 2] = row
	}
	return design
}
