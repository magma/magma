package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	hasInclude  *regexp.Regexp
	includePath *regexp.Regexp
)

// index_headers finds all header files satisfying filter and returns a map
// from them of form map[header filename]full path.
func index_headers(repoRoot string, filters map[string]bool) map[string]string {
	index := make(map[string]string)
	err := filepath.Walk(repoRoot,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			trimmedPath := strings.TrimPrefix(path, repoRoot)
			if strings.HasPrefix(trimmedPath, "build") {
				return nil
			}

			if _, contains := filters[filepath.Ext(path)]; !contains {
				return nil
			}

			if _, contains := index[filepath.Base(path)]; contains {
				fmt.Printf("WARNING: Found conflicting header filenames:%q\n", path)
			}

			index[filepath.Base(path)] = trimmedPath
			return nil
		})
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	return index
}

func fixupTargets(target string, targetFilters map[string]bool, index map[string]string, dry_run bool) error {
	err := filepath.Walk(target,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if _, contains := targetFilters[filepath.Ext(path)]; !contains {
				return nil
			}

			fmt.Printf("::::fixupIncludes(%q)::::\n", path)
			if err := fixupIncludes(path, index, dry_run); err != nil {
				return err
			}

			return nil
		})
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	return nil
}

func fixupIncludes(path string, index map[string]string, dry_run bool) error {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for i, line := range lines {
		if matched := hasInclude.MatchString(line); !matched {
			continue
		}
		ip := includePath.FindStringSubmatch(line)
		if len(ip) != 4 {
			panic(line)
		}
		fileName := filepath.Base(ip[2])
		if _, contains := index[fileName]; !contains {
			if dry_run {
				fmt.Println(line)
			}
			continue
		}

		if !strings.HasSuffix(index[fileName], ip[2]) {
			// Confirm the match actually holds fully. E.g. if include is <sys/socket.h> we don't want to match on lte/gateway/..../util/socket.h
			continue
		}

		newPath := index[fileName]
		updatedLine := includePath.ReplaceAllString(line, fmt.Sprintf("%s%s%s", ip[1], newPath, ip[3]))
		if dry_run {
			fmt.Println(updatedLine)
		} else {
			lines[i] = updatedLine
		}

	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

	return nil
}

func main() {
	repoRoot := flag.String("root", ".", "path to repository root")
	target := flag.String("target", ".", "path to target directory or file for pathification")
	targetFilter := flag.String("target_filter", ".h", "filename filter for patheification, seperated by `,`")
	recursive := flag.Bool("recursive", true, "recursively evaluate target path subdirectories")
	headerFilter := flag.String("header_filter", ".h", "filename filter for header files, seperated by `,`")
	dryRun := flag.Bool("dry_run", true, "dry run and output updated files only to stdout")
	flag.Parse()

	fmt.Println("Pathifier")
	fmt.Println("Running with:")
	fmt.Printf("\troot: %q\n", *repoRoot)
	fmt.Printf("\ttarget: %q\n", *target)
	fmt.Printf("\ttarget_filter: %q\n", *targetFilter)
	fmt.Printf("\trecursive: %t\n", *recursive)
	fmt.Printf("\theader_filter: %q\n", *headerFilter)
	fmt.Printf("\tdry_run: %t\n", *dryRun)

	hasInclude = regexp.MustCompile("^#include")
	includePath = regexp.MustCompile(`^(#include\s+["|<])\s*(\S+)(\s*["|>].*)`)

	headerFilters := make(map[string]bool)
	{
		f := strings.Split(*headerFilter, ",")
		for _, val := range f {
			headerFilters[val] = true
		}
	}

	targetFilters := make(map[string]bool)
	{
		f := strings.Split(*targetFilter, ",")
		for _, val := range f {
			targetFilters[val] = true
		}
	}

	index := index_headers(*repoRoot, headerFilters)
	//fmt.Println("index:", index)

	if err := fixupTargets(*target, targetFilters, index, *dryRun); err != nil {
		fmt.Println("ERROR: ", err)
	}
}
