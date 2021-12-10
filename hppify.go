package main

import (
	"errors"
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
	HasInclude  *regexp.Regexp
	IncludePath *regexp.Regexp
)

// indexWantHFileNames builds a list of potential header files that
// would be associated with discovered CPP files, and therefore should
// be renamed to .hpp if uncovered (not all will exist, some CPP have no H).
func indexWantHFileNames(path string) map[string]bool {
	index := make(map[string]bool)

	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip session_manager sub-directory, as it lacks full path includes
			if strings.Contains(path, "session_manager") {
				return nil
			}

			if filepath.Ext(path) != ".cpp" {
				return nil
			}

			fileName := filepath.Base(path)
			wantHFileName := strings.TrimSuffix(fileName, ".cpp") + ".h"
			if _, contains := index[wantHFileName]; contains {
				fmt.Printf("WARNING: ambiguous header name (multiple CPP with Filename Prefix): %q\n", wantHFileName)
			}

			index[wantHFileName] = true

			return nil
		})
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	return index
}

// renameMatchingHToHPP recursive scans within targetPath and renames
// any header file found which is contained in index (filename.h -> filename.hpp).
// If renamed, and adds the root_path relative path of the renamed file to a
// return map (e.g. map[`lte/gateway/../filename.h`] = true).
func renameMatchingHToHPP(root_path, target_path string, index map[string]bool, dry_run bool) map[string]bool {
	renamed := make(map[string]bool)

	err := filepath.Walk(target_path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip session_manager sub-directory, as it lacks full path includes
			if strings.Contains(path, "session_manager") {
				return nil
			}

			if _, contains := index[filepath.Base(path)]; contains {
				// do rename
				rootRelativePath := strings.TrimPrefix(path, root_path)

				renamed[rootRelativePath] = true

				if dry_run {
					return nil
				}
			}

			return nil
		})
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	return renamed
}

// fixupMatchingIncludes scans through #include statements recursively within
// rootpath, and for any matching include hit in renamedMap (full path) edits
// the include in-place {.h -> .hpp}.
func fixupMatchingIncludes(rootpath, target_path string, rename map[string]bool) {
	err := filepath.Walk(rootpath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			ext := filepath.Ext(path)
			if ext != ".h" && ext != ".c" && ext != ".cpp" {
				return nil
			}

			maybeFixupIncludes(path, rename)

			return nil
		})
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
}

func maybeFixupIncludes(path string, rename map[string]bool) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for i, line := range lines {
		if matched := HasInclude.MatchString(line); !matched {
			continue
		}
		ip := IncludePath.FindStringSubmatch(line)
		if len(ip) != 4 {
			panic(line)
		}
		includePath := ip[2]
		if _, contains := rename[includePath]; !contains {
			continue
		}

		newPath := strings.TrimSuffix(includePath, ".h") + ".hpp"
		updatedLine := IncludePath.ReplaceAllString(line, fmt.Sprintf("%s%s%s", ip[1], newPath, ip[3]))
		lines[i] = updatedLine

	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}

	return
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func findFirstParentWithFile(from_dir, filename string) (string, error) {
	testPath := from_dir + "/" + filename
	if _, err := os.Stat(testPath); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("\t\tos.Stat(%q) == os.ErrNotExist\n", testPath)
		if !strings.Contains(from_dir, "/") {
			return "", errors.New("recursed to / looking for " + filename)
		}
		fmt.Printf("\tgetParentDiretory(%q)\n", from_dir)
		parentDir := getParentDirectory(from_dir)
		fmt.Printf("\t\t=%q\n", parentDir)
		return findFirstParentWithFile(parentDir, filename)
	}

	return testPath, nil
}

func fixupBuildBazel(rename map[string]bool) {
	for path := range rename {
		filename := filepath.Base(path)

		buildBazelPath, err := findFirstParentWithFile(filepath.Dir(path), "BUILD.bazel")
		if err != nil {
			fmt.Printf("Unable to find BUILD.bazel for %q\n", path)
			continue
		}

		input, err := ioutil.ReadFile(buildBazelPath)
		if err != nil {
			log.Fatalln(err)
		}

		lines := strings.Split(string(input), "\n")

		found := false
		for i, line := range lines {
			if strings.Contains(line, filename) {
				// Replace
				found = true
				newFilename := strings.TrimSuffix(filename, ".h") + ".hpp"
				updatedLine := strings.Replace(line, filename, newFilename, 1)
				lines[i] = updatedLine
			}
		}

		if !found {
			fmt.Printf("WARN: Unable to find reference to %q in file %q\n", path, buildBazelPath)
		}

		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile(buildBazelPath, []byte(output), 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func fixupCMakeLists(rename map[string]bool) {
	for path := range rename {
		filename := filepath.Base(path)

		cmakePath, err := findFirstParentWithFile(filepath.Dir(path), "CMakeLists.txt")
		if err != nil {
			fmt.Printf("Unable to find CMakeLists.txt for %q\n", path)
			continue
		}

		input, err := ioutil.ReadFile(cmakePath)
		if err != nil {
			log.Fatalln(err)
		}

		lines := strings.Split(string(input), "\n")

		found := false
		for i, line := range lines {
			if strings.Contains(line, filename) {
				// Replace
				found = true
				newFilename := strings.TrimSuffix(filename, ".h") + ".hpp"
				updatedLine := strings.Replace(line, filename, newFilename, 1)
				lines[i] = updatedLine
			}
		}

		if !found {
			fmt.Printf("WARN: Unable to find reference to %q in file %q\n", path, cmakePath)
		}

		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile(cmakePath, []byte(output), 0644)
		if err != nil {
			log.Fatalln(err)
		}
	}

}

func main() {
	repoRoot := flag.String("root", ".", "path to repository root")
	target := flag.String("target", ".", "path to target directory to recursively hppify")
	dryRun := flag.Bool("dry_run", true, "dry run and output updated files only to stdout")
	flag.Parse()

	fmt.Println("Hppify")
	fmt.Println("Running with:")
	fmt.Printf("\troot: %q\n", *repoRoot)
	fmt.Printf("\ttarget: %q\n", *target)
	fmt.Printf("\tdry_run: %t\n", *dryRun)

	HasInclude = regexp.MustCompile("^#include")
	IncludePath = regexp.MustCompile(`^(#include\s+["|<])\s*(\S+)(\s*["|>].*)`)

	index := indexWantHFileNames(*repoRoot)
	renamed := renameMatchingHToHPP(*repoRoot, *target, index, *dryRun)

	fmt.Println(renamed)
	fmt.Println("Renamed ", len(renamed), " .h to .hpp")

	fixupMatchingIncludes(*repoRoot, *target, renamed)
	fixupCMakeLists(renamed)
	fixupBuildBazel(renamed)

	for filename := range renamed {
		newPath := strings.TrimSuffix(filename, ".h") + ".hpp"
		if err := os.Rename(filename, newPath); err != nil {
			log.Fatalf("os.Rename(%q, %q) = %s", filename, newPath, err)
		}
	}
}

