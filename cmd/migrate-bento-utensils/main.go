// One-shot migration: fusionne les lignes legacy cover+eating en une ligne utensils dans recipes/**/*.bentext.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"bentext/internal/recipe"
)

func main() {
	root := "recipes"
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	n := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".bentext") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		text := strings.ReplaceAll(strings.ReplaceAll(string(data), "\r\n", "\n"), "\r", "\n")
		parts := strings.Split(text, "---")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		if len(parts) < 2 {
			return nil
		}
		lastIdx := len(parts) - 1
		lines := nonEmptyLines(parts[lastIdx])
		newLines, ok := recipe.MigrateBentoValueLines(lines)
		if !ok {
			return nil
		}
		parts[lastIdx] = strings.Join(newLines, "\n")
		out := strings.Join(parts, "\n---\n")
		if !strings.HasSuffix(out, "\n") {
			out += "\n"
		}
		if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
			return err
		}
		fmt.Println(path)
		n++
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Migrés : %d fichiers\n", n)
}

func nonEmptyLines(block string) []string {
	var out []string
	for _, line := range strings.Split(block, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			out = append(out, line)
		}
	}
	return out
}
