// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fredericlemoine/gotree/io/utils"
	"github.com/fredericlemoine/gotree/tree"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var renameintree string
var renameouttree string
var renamemapfile string

func readMapFile(file string) (map[string]string, error) {
	outmap := make(map[string]string, 0)
	var mapfile *os.File
	var err error
	if mapfile, err = os.Open(file); err != nil {
		return outmap, err
	}

	reader := bufio.NewReader(mapfile)
	line, e := utils.Readln(reader)
	nl := 1
	for e == nil {
		cols := strings.Split(line, "\t")
		if len(cols) != 2 {
			return outmap, errors.New("Map file does not have 2 fields at line: " + fmt.Sprintf("%d", nl))
		}
		outmap[cols[0]] = cols[1]
		line, e = utils.Readln(reader)
		nl++
	}

	if err = mapfile.Close(); err != nil {
		return outmap, err
	}

	return outmap, nil
}

// This function renames nodes of the tree based on the map in argument
// If a name in the map does not exist in the tree, then returns an error
// If a node/tip in the tree does not have a name in the map: OK
// After rename, tip index is updated, as well as bitsets of the edges
func renameTree(t *tree.Tree, namemap map[string]string) error {
	nodeindex := tree.NewNodeIndex(t)
	for name, newname := range namemap {
		node, ok := nodeindex.GetNode(name)
		if !ok {
			return errors.New(fmt.Sprintf("The node %s does not exist in the tree", name))
		}
		node.SetName(newname)
	}
	// After we update bitsets if any, and node indexes
	t.UpdateTipIndex()
	err := t.ClearBitSets()
	if err != nil {
		return err
	}
	t.UpdateBitSet()
	return nil
}

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Renames tips of the input tree, given a map file",
	Long: `Renames tips of the input tree, given a map file.

Map file must be tab separated with columns:
1) Current name of the tip
2) Desired new name of the tip

If a tip name does not appear in the map file, it will not be renamed. 
If a name that does not exist appears in the map file, it will throw an error.
`,
	Run: func(cmd *cobra.Command, args []string) {

		if renamemapfile == "none" {
			panic("map file is not given")
		}

		// Read map file
		namemap, err := readMapFile(renamemapfile)
		if err != nil {
			panic(err)
		}

		// Read Tree
		var tree *tree.Tree
		tree, err = utils.ReadRefTree(renameintree)
		if err != nil {
			panic(err)
		}
		var f *os.File
		if renameouttree != "stdout" {
			f, err = os.Create(renameouttree)
		} else {
			f = os.Stdout
		}
		if err != nil {
			panic(err)
		}

		err = renameTree(tree, namemap)
		if err != nil {
			panic(err)
		}

		f.WriteString(tree.Newick() + "\n")
		f.Close()
	},
}

func init() {
	RootCmd.AddCommand(renameCmd)
	renameCmd.Flags().StringVarP(&renameouttree, "output", "o", "stdout", "Renamed tree output file")
	renameCmd.Flags().StringVarP(&renameintree, "input", "i", "stdin", "Input tree")
	renameCmd.Flags().StringVarP(&renamemapfile, "map", "m", "none", "Tip name map file")
}